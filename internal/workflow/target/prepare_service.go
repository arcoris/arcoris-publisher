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
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/ports/filesystem"
	"arcoris.dev/arcoris-publisher/internal/ports/git"
	"arcoris.dev/arcoris-publisher/internal/workflow/pathutil"
)

// PrepareTargets creates or refreshes target Git worktrees without constructing
// publication files, committing, tagging, pushing, or writing transaction state.
func (s Service) PrepareTargets(ctx context.Context, req PrepareRequest) (PrepareResult, error) {
	root, err := s.validatePrepareRequest(req)
	if err != nil {
		return PrepareResult{}, err
	}

	result := PrepareResult{targetRoot: root}
	if err := s.deps.FS.MkdirAll(ctx, root, filesystem.MkdirOptions{}); err != nil {
		return result, err
	}
	if isDir, err := s.deps.FS.IsDir(ctx, root); err != nil {
		return result, err
	} else if !isDir {
		result.status = PrepareStatusFailed
		return result, nil
	}

	for _, mod := range req.Plan.Modules() {
		result.modules = append(result.modules, s.prepareTargetModule(ctx, root, req, mod))
	}
	result.status = aggregatePrepareModules(result.modules)
	return result, nil
}

func (s Service) validatePrepareRequest(req PrepareRequest) (string, error) {
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

func (s Service) prepareTargetModule(ctx context.Context, root string, req PrepareRequest, mod plan.ModulePlan) PrepareModuleResult {
	remoteName := s.opts.RemoteName
	if remoteName == "" {
		remoteName = DefaultOptions().RemoteName
	}
	worktree := RepositoryWorktree(root, mod.Repository())
	module := PrepareModuleResult{
		module:      mod.Name(),
		repository:  mod.Repository(),
		worktreeDir: worktree,
		remoteName:  remoteName,
	}
	module.actions = append(module.actions, preparePath(preparePassed("worktree-path", "target worktree path resolved"), worktree))

	expectedURL, hasExpected, ok := resolvePrepareRemote(req, mod, &module)
	if !ok {
		module.status = PrepareStatusFailed
		return module
	}
	if hasExpected {
		module.remoteURL = expectedURL
	}

	if !prepareSingleBranch(mod, &module) {
		module.status = PrepareStatusFailed
		return module
	}
	if !s.prepareWorktreePresence(ctx, worktree, expectedURL, hasExpected, &module) {
		module.status = PrepareStatusFailed
		return module
	}
	if !s.prepareWorktreeGit(ctx, worktree, &module) {
		module.status = PrepareStatusFailed
		return module
	}
	if !s.prepareRemote(ctx, worktree, remoteName, expectedURL, hasExpected, &module) {
		module.status = PrepareStatusFailed
		return module
	}
	if !s.prepareFetch(ctx, worktree, remoteName, module.remoteURL, &module) {
		module.status = PrepareStatusFailed
		return module
	}
	if !s.prepareBranch(ctx, mod.Branches()[0].Target().String(), worktree, remoteName, &module) {
		module.status = PrepareStatusFailed
		return module
	}
	if !s.prepareClean(ctx, worktree, &module) {
		module.status = PrepareStatusFailed
		return module
	}

	module.status = aggregatePrepareActions(module.actions)
	return module
}

func resolvePrepareRemote(req PrepareRequest, mod plan.ModulePlan, result *PrepareModuleResult) (string, bool, bool) {
	if !req.HasRemoteTemplate {
		return "", false, true
	}
	url, err := req.RemoteTemplate.Resolve(mod.Repository(), mod.Name())
	if err != nil {
		result.actions = append(result.actions, prepareFailed("remote", "remote_template_failed", "remote template could not be resolved"))
		return "", true, false
	}
	return url, true, true
}

func prepareSingleBranch(mod plan.ModulePlan, result *PrepareModuleResult) bool {
	branches := mod.Branches()
	if len(branches) != 1 {
		result.actions = append(result.actions, prepareFailed("checkout", "unsupported_multi_branch", fmt.Sprintf("module has %d branch mappings; multi-branch target preparation is not supported", len(branches))))
		return false
	}
	return true
}

func (s Service) prepareWorktreePresence(
	ctx context.Context,
	worktree string,
	expectedURL string,
	hasExpected bool,
	result *PrepareModuleResult,
) bool {
	exists, err := s.deps.FS.Exists(ctx, worktree)
	if err != nil {
		result.actions = append(result.actions, preparePath(prepareFailed("validate-worktree", "worktree_lookup_failed", "target worktree lookup failed"), worktree))
		return false
	}
	if !exists {
		if !hasExpected {
			result.actions = append(result.actions, preparePath(prepareFailed("clone", "missing_remote_template", "missing target worktree requires --remote-template or target.remoteTemplate"), worktree))
			return false
		}
		if err := s.deps.Git.Clone(ctx, expectedURL, worktree, git.CloneOptions{SensitiveValues: []string{expectedURL}}); err != nil {
			result.actions = append(result.actions, preparePath(prepareRemote(prepareFailed("clone", "clone_failed", "clone failed"), expectedURL), worktree))
			return false
		}
		_ = s.deps.FS.MkdirAll(ctx, worktree, filesystem.MkdirOptions{})
		result.actions = append(result.actions, preparePath(prepareRemote(preparePassed("clone", "missing worktree cloned"), expectedURL), worktree))
		return true
	}

	isDir, err := s.deps.FS.IsDir(ctx, worktree)
	if err != nil {
		result.actions = append(result.actions, preparePath(prepareFailed("validate-worktree", "worktree_lookup_failed", "target worktree lookup failed"), worktree))
		return false
	}
	if !isDir {
		result.actions = append(result.actions, preparePath(prepareFailed("validate-worktree", "worktree_not_directory", "target worktree is not a directory"), worktree))
		return false
	}
	result.actions = append(result.actions, preparePath(prepareSkipped("clone", "target worktree already exists"), worktree))
	result.actions = append(result.actions, preparePath(preparePassed("validate-worktree", "target worktree exists"), worktree))
	return true
}

func (s Service) prepareWorktreeGit(ctx context.Context, worktree string, result *PrepareModuleResult) bool {
	if _, err := s.deps.Git.Status(ctx, worktree); err != nil {
		result.actions = append(result.actions, preparePath(prepareFailed("validate-worktree", "worktree_status_failed", "target Git status failed"), worktree))
		return false
	}
	return true
}

func (s Service) prepareRemote(
	ctx context.Context,
	worktree string,
	remoteName string,
	expectedURL string,
	hasExpected bool,
	result *PrepareModuleResult,
) bool {
	currentURL, ok, err := s.deps.Git.RemoteURL(ctx, worktree, remoteName)
	if err != nil {
		result.actions = append(result.actions, preparePath(prepareFailed("remote", "remote_lookup_failed", "target remote lookup failed"), worktree))
		return false
	}
	if !ok {
		if !hasExpected {
			result.actions = append(result.actions, preparePath(prepareFailed("remote", "remote_missing", "target remote is missing and no remote template is configured"), worktree))
			return false
		}
		if err := s.deps.Git.AddRemote(ctx, worktree, remoteName, expectedURL); err != nil {
			result.actions = append(result.actions, preparePath(prepareRemote(prepareFailed("remote", "remote_add_failed", "target remote add failed"), expectedURL), worktree))
			return false
		}
		result.remoteURL = expectedURL
		result.actions = append(result.actions, preparePath(prepareRemote(preparePassed("remote", "target remote added"), expectedURL), worktree))
		return true
	}

	result.remoteURL = currentURL
	if hasExpected && currentURL != expectedURL {
		result.actions = append(result.actions, preparePath(prepareRemote(prepareFailed("remote", "remote_mismatch", "target remote does not match expected template"), currentURL), worktree))
		return false
	}
	result.actions = append(result.actions, preparePath(prepareRemote(preparePassed("remote", "target remote is configured"), currentURL), worktree))
	return true
}

func (s Service) prepareFetch(ctx context.Context, worktree string, remoteName string, remoteURL string, result *PrepareModuleResult) bool {
	if !s.opts.Fetch {
		result.actions = append(result.actions, prepareSkipped("fetch", "fetch disabled"))
		return true
	}
	if err := s.deps.Git.Fetch(ctx, worktree, remoteName, git.FetchOptions{
		Prune:           true,
		Tags:            git.FetchTagsAll,
		SensitiveValues: []string{remoteURL},
	}); err != nil {
		result.actions = append(result.actions, preparePath(prepareFailed("fetch", "fetch_failed", "git fetch failed"), worktree))
		return false
	}
	result.actions = append(result.actions, preparePath(preparePassed("fetch", "remote refs fetched"), worktree))
	return true
}

func (s Service) prepareBranch(
	ctx context.Context,
	targetBranch string,
	worktree string,
	remoteName string,
	result *PrepareModuleResult,
) bool {
	branchRef := "refs/heads/" + targetBranch
	local, err := s.deps.Git.RefExists(ctx, worktree, branchRef)
	if err != nil {
		result.actions = append(result.actions, prepareFailed("checkout", "branch_lookup_failed", "target branch lookup failed"))
		return false
	}
	if local {
		if err := s.deps.Git.Checkout(ctx, worktree, targetBranch, git.CheckoutOptions{}); err != nil {
			result.actions = append(result.actions, prepareFailed("checkout", "checkout_failed", "target branch checkout failed"))
			return false
		}
		result.actions = append(result.actions, preparePassed("checkout", "target branch checked out"))
		return true
	}

	remoteExists, err := s.deps.Git.RemoteRefExists(ctx, worktree, remoteName, branchRef)
	if err != nil {
		result.actions = append(result.actions, prepareFailed("checkout", "remote_branch_lookup_failed", "remote branch lookup failed"))
		return false
	}
	if !remoteExists {
		result.actions = append(result.actions, prepareFailed("checkout", "branch_missing", "target branch is missing locally and remotely"))
		return false
	}
	startPoint := remoteName + "/" + targetBranch
	if err := s.deps.Git.CreateBranch(ctx, worktree, git.BranchName(targetBranch), startPoint, git.CreateBranchOptions{}); err != nil {
		result.actions = append(result.actions, prepareFailed("checkout", "checkout_failed", "tracking branch creation failed"))
		return false
	}
	if err := s.deps.Git.Checkout(ctx, worktree, targetBranch, git.CheckoutOptions{}); err != nil {
		result.actions = append(result.actions, prepareFailed("checkout", "checkout_failed", "target branch checkout failed"))
		return false
	}
	result.actions = append(result.actions, preparePassed("checkout", "target branch checked out from remote"))
	return true
}

func (s Service) prepareClean(ctx context.Context, worktree string, result *PrepareModuleResult) bool {
	if !s.opts.RequireClean {
		result.actions = append(result.actions, prepareSkipped("clean", "clean worktree validation disabled"))
		return true
	}
	status, err := s.deps.Git.Status(ctx, worktree)
	if err != nil {
		result.actions = append(result.actions, preparePath(prepareFailed("clean", "worktree_status_failed", "target Git status failed"), worktree))
		return false
	}
	if status.IsDirty() {
		result.actions = append(result.actions, preparePath(prepareFailed("clean", "worktree_dirty", "target worktree is dirty"), worktree))
		return false
	}
	result.actions = append(result.actions, preparePath(preparePassed("clean", "target worktree is clean"), worktree))
	return true
}
