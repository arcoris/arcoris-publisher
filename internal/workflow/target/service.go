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

	"arcoris.dev/arcoris-publisher/internal/manifest"
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
		return WorkspaceSet{}, err
	}
	if ok, err := s.deps.FS.Exists(ctx, root); err != nil {
		return WorkspaceSet{}, err
	} else if !ok {
		if err := s.deps.FS.MkdirAll(ctx, root, filesystem.MkdirOptions{}); err != nil {
			return WorkspaceSet{}, err
		}
	}
	if isDir, err := s.deps.FS.IsDir(ctx, root); err != nil {
		return WorkspaceSet{}, err
	} else if !isDir {
		issues.Add(IssueRootNotDirectory, "", root, "target root is not a directory")
	}
	workspaces := make([]ModuleWorkspace, 0, req.Plan.Len())
	for _, mod := range req.Plan.Modules() {
		ws, ok := s.prepareModule(ctx, root, mod, &issues)
		if ok {
			workspaces = append(workspaces, ws)
		}
	}
	if err := issues.Err(); err != nil {
		return WorkspaceSet{}, err
	}
	return WorkspaceSet{workspaces: workspaces}, nil
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
	exists, err := s.deps.FS.Exists(ctx, worktree)
	if err != nil {
		issues.AddMessage(IssueWorktreeMissing, mod.Name(), worktree, err.Error())
		return ModuleWorkspace{}, false
	}
	if !exists {
		if s.opts.RemoteURL != nil {
			url := s.opts.RemoteURL(mod.Repository())
			if url == "" {
				issues.Add(
					IssueCloneURLMissing,
					mod.Name(),
					worktree,
					"remote URL resolver returned empty URL for %s",
					mod.Repository(),
				)
				return ModuleWorkspace{}, false
			}
			if err := s.deps.Git.Clone(ctx, url, worktree, git.CloneOptions{}); err != nil {
				issues.Add(IssueWorktreeMissing, mod.Name(), worktree, "clone failed: %v", err)
				return ModuleWorkspace{}, false
			}
		} else if s.opts.CreateMissing {
			if err := s.deps.FS.MkdirAll(ctx, worktree, filesystem.MkdirOptions{}); err != nil {
				issues.AddMessage(IssueWorktreeMissing, mod.Name(), worktree, err.Error())
				return ModuleWorkspace{}, false
			}
		} else {
			issues.Add(IssueWorktreeMissing, mod.Name(), worktree, "target worktree is missing")
			return ModuleWorkspace{}, false
		}
	}
	if isDir, err := s.deps.FS.IsDir(ctx, worktree); err != nil {
		issues.AddMessage(IssueWorktreeNotDirectory, mod.Name(), worktree, err.Error())
		return ModuleWorkspace{}, false
	} else if !isDir {
		issues.Add(IssueWorktreeNotDirectory, mod.Name(), worktree, "target worktree is not a directory")
		return ModuleWorkspace{}, false
	}
	if s.opts.Fetch {
		_ = s.deps.Git.Fetch(ctx, worktree, s.opts.RemoteName, git.FetchOptions{
			Prune: true,
			Tags:  git.FetchTagsAll,
		})
	}
	if s.opts.RequireClean {
		status, err := s.deps.Git.Status(ctx, worktree)
		if err == nil && (!status.Clean || status.HasEntries()) {
			issues.Add(IssueWorktreeDirty, mod.Name(), worktree, "target worktree is dirty")
		}
	}
	branchWorkspaces := make([]BranchWorkspace, 0, len(mod.Branches()))
	for _, branch := range mod.Branches() {
		branchWorkspaces = append(branchWorkspaces, newBranchWorkspace(branch.Source(), branch.Target()))
	}
	if s.opts.CheckoutBranch && len(branchWorkspaces) > 0 {
		ref := branchWorkspaces[0].Target().String()
		if err := s.deps.Git.Checkout(ctx, worktree, ref, git.CheckoutOptions{}); err != nil {
			issues.Add(IssueInvalidRequest, mod.Name(), worktree, "checkout %s failed: %v", ref, err)
		}
	}
	return ModuleWorkspace{
		module:      mod.Name(),
		repository:  mod.Repository(),
		worktreeDir: worktree,
		branches:    branchWorkspaces,
	}, true
}

var _ = manifest.ModuleName("")
