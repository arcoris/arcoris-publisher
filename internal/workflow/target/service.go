// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package target

import (
	"context"

	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/ports/filesystem"
	"arcoris.dev/arcoris-publisher/internal/ports/git"
	"arcoris.dev/arcoris-publisher/internal/workflow/pathutil"
)

// Service prepares target repository worktrees for a publication plan.
type Service struct {
	// deps contains infrastructure ports used by target preparation.
	deps Dependencies

	// opts contains normalized behavior toggles.
	opts Options
}

// New returns a target preparation service.
func New(deps Dependencies, opts Options) Service {
	defaults := DefaultOptions()
	if opts.RemoteName == "" {
		opts.RemoteName = defaults.RemoteName
	}
	return Service{deps: deps, opts: opts}
}

// Prepare validates and prepares target worktrees for req.
func (s Service) Prepare(ctx context.Context, req Request) (WorkspaceSet, error) {
	root, err := s.validateRequest(req)
	if err != nil {
		return WorkspaceSet{}, err
	}

	issues := newIssueCollector()
	if err := s.ensureRoot(ctx, root, &issues); err != nil {
		return WorkspaceSet{}, err
	}
	if err := issues.Err(); err != nil {
		return WorkspaceSet{}, err
	}

	workspaces := s.prepareModules(ctx, root, req.Plan, &issues)
	if err := issues.Err(); err != nil {
		return WorkspaceSet{}, err
	}

	return WorkspaceSet{workspaces: workspaces}, nil
}

// validateRequest rejects malformed input before any filesystem or Git work.
func (s Service) validateRequest(req Request) (string, error) {
	issues := newIssueCollector()
	if req.Plan.Empty() {
		issues.Add(IssueInvalidRequest, "", "plan", "plan is empty")
	}

	root, err := pathutil.CleanAbs(req.RootDir)
	if err != nil {
		issues.AddMessage(IssueInvalidRequest, "", "rootDir", err.Error())
	}
	if s.deps.FS == nil {
		issues.Add(IssueInvalidRequest, "", "fs", "filesystem dependency is required")
	}
	if s.deps.Git == nil {
		issues.Add(IssueInvalidRequest, "", "git", "git dependency is required")
	}
	if err := issues.Err(); err != nil {
		return "", err
	}

	return root, nil
}

// ensureRoot creates the target root when needed and verifies that the final
// path is a directory.
func (s Service) ensureRoot(
	ctx context.Context,
	root string,
	issues *issueCollector,
) error {
	if ok, err := s.deps.FS.Exists(ctx, root); err != nil {
		return err
	} else if !ok {
		if err := s.deps.FS.MkdirAll(ctx, root, filesystem.MkdirOptions{}); err != nil {
			return err
		}
	}

	if isDir, err := s.deps.FS.IsDir(ctx, root); err != nil {
		return err
	} else if !isDir {
		issues.Add(IssueRootNotDirectory, "", root, "target root is not a directory")
	}

	return nil
}

// prepareModules prepares all planned module worktrees in plan order and keeps
// collecting module diagnostics after individual module failures.
func (s Service) prepareModules(
	ctx context.Context,
	root string,
	p plan.Plan,
	issues *issueCollector,
) []ModuleWorkspace {
	workspaces := make([]ModuleWorkspace, 0, p.Len())
	for _, mod := range p.Modules() {
		ws, ok := s.prepareModule(ctx, root, mod, issues)
		if ok {
			workspaces = append(workspaces, ws)
		}
	}

	return workspaces
}

// prepareModule ensures one module worktree exists, is clean, and has branch
// metadata ready for later workflow stages.
func (s Service) prepareModule(
	ctx context.Context,
	root string,
	mod plan.ModulePlan,
	issues *issueCollector,
) (ModuleWorkspace, bool) {
	worktree := repositoryWorktree(root, mod.Repository())

	if !s.ensureWorktreeExists(ctx, mod, worktree, issues) {
		return ModuleWorkspace{}, false
	}

	if !s.validateWorktreeDirectory(ctx, mod, worktree, issues) {
		return ModuleWorkspace{}, false
	}

	if !s.validateGitWorktree(ctx, mod, worktree, issues) {
		return ModuleWorkspace{}, false
	}

	if !s.fetchWorktree(ctx, mod, worktree, issues) {
		return ModuleWorkspace{}, false
	}
	s.checkCleanWorktree(ctx, mod, worktree, issues)

	branches := branchWorkspaces(mod)
	s.checkoutFirstBranch(ctx, mod, worktree, branches, issues)

	return ModuleWorkspace{
		module:      mod.Name(),
		repository:  mod.Repository(),
		worktreeDir: worktree,
		branches:    branches,
	}, true
}

// ensureWorktreeExists creates, clones, or reports a missing worktree.
func (s Service) ensureWorktreeExists(
	ctx context.Context,
	mod plan.ModulePlan,
	worktree string,
	issues *issueCollector,
) bool {
	exists, err := s.deps.FS.Exists(ctx, worktree)
	if err != nil {
		issues.AddMessage(IssueWorktreeMissing, mod.Name(), worktree, err.Error())
		return false
	}
	if exists {
		return true
	}

	if s.opts.RemoteURL != nil {
		return s.cloneMissingWorktree(ctx, mod, worktree, issues)
	}

	if s.opts.CreateMissing {
		return s.createMissingWorktree(ctx, mod, worktree, issues)
	}

	issues.Add(IssueWorktreeMissing, mod.Name(), worktree, "target worktree is missing")
	return false
}

// cloneMissingWorktree clones a missing worktree when a remote URL resolver is
// configured.
func (s Service) cloneMissingWorktree(
	ctx context.Context,
	mod plan.ModulePlan,
	worktree string,
	issues *issueCollector,
) bool {
	url := s.opts.RemoteURL(mod.Repository())
	if url == "" {
		issues.Add(
			IssueCloneURLMissing,
			mod.Name(),
			worktree,
			"remote URL resolver returned empty URL for %s",
			mod.Repository(),
		)
		return false
	}

	if err := s.deps.Git.Clone(ctx, url, worktree, git.CloneOptions{}); err != nil {
		issues.Add(IssueWorktreeMissing, mod.Name(), worktree, "clone failed: %v", err)
		return false
	}

	return true
}

// createMissingWorktree creates a missing local worktree directory.
func (s Service) createMissingWorktree(
	ctx context.Context,
	mod plan.ModulePlan,
	worktree string,
	issues *issueCollector,
) bool {
	if err := s.deps.FS.MkdirAll(ctx, worktree, filesystem.MkdirOptions{}); err != nil {
		issues.AddMessage(IssueWorktreeMissing, mod.Name(), worktree, err.Error())
		return false
	}

	return true
}

// validateWorktreeDirectory rejects non-directory worktree paths.
func (s Service) validateWorktreeDirectory(
	ctx context.Context,
	mod plan.ModulePlan,
	worktree string,
	issues *issueCollector,
) bool {
	isDir, err := s.deps.FS.IsDir(ctx, worktree)
	if err != nil {
		issues.AddMessage(IssueWorktreeNotDirectory, mod.Name(), worktree, err.Error())
		return false
	}
	if !isDir {
		issues.Add(IssueWorktreeNotDirectory, mod.Name(), worktree, "target worktree is not a directory")
		return false
	}

	return true
}

// validateGitWorktree rejects directories that cannot be inspected as Git
// worktrees before construct or publish stages mutate them.
func (s Service) validateGitWorktree(
	ctx context.Context,
	mod plan.ModulePlan,
	worktree string,
	issues *issueCollector,
) bool {
	if _, err := s.deps.Git.Status(ctx, worktree); err != nil {
		issues.Add(IssueWorktreeStatusFailed, mod.Name(), worktree, "target Git status failed: %v", err)
		return false
	}

	return true
}

// fetchWorktree refreshes remote state when requested. FetchRequired controls
// whether transport failures block construction.
func (s Service) fetchWorktree(
	ctx context.Context,
	mod plan.ModulePlan,
	worktree string,
	issues *issueCollector,
) bool {
	if !s.opts.Fetch {
		return true
	}

	err := s.deps.Git.Fetch(ctx, worktree, s.opts.RemoteName, git.FetchOptions{
		Prune: true,
		Tags:  git.FetchTagsAll,
	})
	if err != nil && s.opts.FetchRequired {
		issues.Add(IssueFetchFailed, mod.Name(), worktree, "git fetch failed: %v", err)
		return false
	}

	return true
}

// checkCleanWorktree records a dirty worktree when policy requires cleanliness.
func (s Service) checkCleanWorktree(
	ctx context.Context,
	mod plan.ModulePlan,
	worktree string,
	issues *issueCollector,
) {
	if !s.opts.RequireClean {
		return
	}

	status, err := s.deps.Git.Status(ctx, worktree)
	if err != nil {
		issues.Add(IssueWorktreeStatusFailed, mod.Name(), worktree, "target Git status failed: %v", err)
		return
	}
	if !status.Clean || status.HasEntries() {
		issues.Add(IssueWorktreeDirty, mod.Name(), worktree, "target worktree is dirty")
	}
}

// branchWorkspaces converts plan branch mappings into target workspace metadata.
func branchWorkspaces(mod plan.ModulePlan) []BranchWorkspace {
	branches := mod.Branches()
	out := make([]BranchWorkspace, 0, len(branches))
	for _, branch := range branches {
		out = append(out, newBranchWorkspace(branch.Source(), branch.Target()))
	}

	return out
}

// checkoutFirstBranch switches to the first target branch when checkout is
// enabled. Later workflow stages use the same prepared worktree.
func (s Service) checkoutFirstBranch(
	ctx context.Context,
	mod plan.ModulePlan,
	worktree string,
	branches []BranchWorkspace,
	issues *issueCollector,
) {
	if !s.opts.CheckoutBranch || len(branches) == 0 {
		return
	}

	ref := branches[0].Target().String()
	if err := s.deps.Git.Checkout(ctx, worktree, ref, git.CheckoutOptions{}); err != nil {
		issues.Add(IssueInvalidRequest, mod.Name(), worktree, "checkout %s failed: %v", ref, err)
	}
}
