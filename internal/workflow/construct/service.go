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

package construct

import (
	"context"
	"encoding/json"
	"path/filepath"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/ports/filesystem"
	"arcoris.dev/arcoris-publisher/internal/workflow/pathutil"
	"arcoris.dev/arcoris-publisher/internal/workflow/source"
	"arcoris.dev/arcoris-publisher/internal/workflow/target"
)

// Service constructs target repository trees from explicit publish entries.
type Service struct {
	// deps contains infrastructure ports used by construction.
	deps Dependencies

	// opts contains normalized behavior toggles.
	opts Options
}

// New returns a construction service.
func New(deps Dependencies, opts Options) Service {
	if opts == (Options{}) {
		opts = DefaultOptions()
	}
	return Service{deps: deps, opts: opts}
}

// Construct cleans target worktrees and copies only explicit publish entries.
func (s Service) Construct(ctx context.Context, req Request) (Result, error) {
	if err := s.validateRequest(req); err != nil {
		return Result{}, err
	}

	results, err := s.constructModules(ctx, req)
	if err != nil {
		return Result{}, err
	}

	return Result{modules: results}, nil
}

// validateRequest rejects incomplete inputs before target trees are cleaned.
func (s Service) validateRequest(req Request) error {
	issues := newIssueCollector()
	if req.Plan.Empty() {
		issues.Add(IssueInvalidRequest, "", "plan", "plan is empty")
	}
	if s.deps.FS == nil {
		issues.Add(IssueInvalidRequest, "", "fs", "filesystem dependency is required")
	}

	return issues.Err()
}

// constructModules constructs every planned module and aggregates per-module
// failures so callers get the complete invalid target set.
func (s Service) constructModules(ctx context.Context, req Request) ([]ModuleResult, error) {
	issues := newIssueCollector()
	results := make([]ModuleResult, 0, req.Plan.Len())

	for _, mod := range req.Plan.Modules() {
		mr, ok := s.constructModule(ctx, req, mod.Name(), &issues)
		if ok {
			results = append(results, mr)
		}
	}
	if err := issues.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// constructModule cleans one target worktree and copies only inspected explicit
// publish entries for the same module.
func (s Service) constructModule(
	ctx context.Context,
	req Request,
	name manifest.ModuleName,
	issues *issueCollector,
) (ModuleResult, bool) {
	module, ok := resolveModuleContext(req, name, issues)
	if !ok {
		return ModuleResult{}, false
	}

	operations, ok := s.cleanWorktree(ctx, module, issues)
	if !ok {
		return ModuleResult{}, false
	}

	operations = append(operations, s.copyEntries(ctx, module, issues)...)
	if !s.appendProvenanceFile(ctx, req, module, &operations, issues) {
		return ModuleResult{}, false
	}

	return ModuleResult{
		module:      name,
		worktreeDir: filepath.Clean(module.workspace.WorktreeDir()),
		operations:  operations,
		changed:     len(operations) > 0,
	}, true
}

// moduleContext groups the resolved plan, source, and target records needed to
// construct one module.
type moduleContext struct {
	plan      plan.ModulePlan
	source    source.ModuleSnapshot
	workspace target.ModuleWorkspace
}

// resolveModuleContext binds plan, source snapshot, and target workspace for a
// module name and reports missing upstream workflow results.
func resolveModuleContext(
	req Request,
	name manifest.ModuleName,
	issues *issueCollector,
) (moduleContext, bool) {
	mod, _ := req.Plan.ModuleByName(name)
	src, ok := findSourceModule(req.Source, name)
	if !ok {
		issues.Add(IssueMissingSource, name, "source", "source snapshot for module %q is missing", name)
		return moduleContext{}, false
	}

	ws, ok := req.Targets.WorkspaceByModule(name)
	if !ok {
		issues.Add(IssueMissingTarget, name, "target", "target workspace for module %q is missing", name)
		return moduleContext{}, false
	}

	return moduleContext{plan: mod, source: src, workspace: ws}, true
}

// cleanWorktree removes previous target contents while preserving configured
// repository metadata.
func (s Service) cleanWorktree(
	ctx context.Context,
	module moduleContext,
	issues *issueCollector,
) ([]Operation, bool) {
	preserve := []string{}
	if s.opts.PreserveGitDir {
		preserve = append(preserve, ".git")
	}
	worktree := module.workspace.WorktreeDir()
	if err := s.deps.FS.CleanDir(ctx, worktree, filesystem.CleanDirOptions{
		Preserve:      preserve,
		RequireGitDir: false,
		AllowMissing:  false,
		SafetyRoot:    worktree,
	}); err != nil {
		issues.AddMessage(IssueCleanFailed, module.plan.Name(), worktree, err.Error())
		return nil, false
	}

	return []Operation{newOperation(OperationClean, "", worktree)}, true
}

// copyEntries copies each explicit source entry. Individual copy failures are
// collected while later entries are still attempted.
func (s Service) copyEntries(
	ctx context.Context,
	module moduleContext,
	issues *issueCollector,
) []Operation {
	operations := []Operation{}

	for _, entry := range module.source.Entries() {
		operation, ok := s.copyEntry(ctx, module, entry, issues)
		if ok {
			operations = append(operations, operation)
		}
	}

	return operations
}

// copyEntry copies or records one explicit entry.
func (s Service) copyEntry(
	ctx context.Context,
	module moduleContext,
	entry source.EntrySnapshot,
	issues *issueCollector,
) (Operation, bool) {
	worktree := module.workspace.WorktreeDir()
	targetPath := pathutil.JoinRelative(worktree, entry.TargetPath())

	if !entry.Present() {
		return newOperation(OperationSkipOptional, entry.SourcePath(), targetPath), true
	}

	if err := pathutil.EnsureInside(worktree, targetPath); err != nil {
		issues.AddMessage(IssueTargetPathEscape, module.plan.Name(), targetPath, err.Error())
		return Operation{}, false
	}

	switch entry.Kind() {
	case manifest.PublishEntryFile:
		return s.copyFileEntry(ctx, module, entry, targetPath, issues)
	case manifest.PublishEntryDirectory:
		return s.copyDirectoryEntry(ctx, module, entry, targetPath, issues)
	default:
		return Operation{}, false
	}
}

// copyFileEntry copies one explicit file entry into the target worktree.
func (s Service) copyFileEntry(
	ctx context.Context,
	module moduleContext,
	entry source.EntrySnapshot,
	targetPath string,
	issues *issueCollector,
) (Operation, bool) {
	data, err := s.deps.FS.ReadFile(ctx, entry.SourcePath())
	if err != nil {
		issues.AddMessage(IssueEntryCopyFailed, module.plan.Name(), entry.SourcePath(), err.Error())
		return Operation{}, false
	}

	err = s.deps.FS.WriteFile(ctx, targetPath, data, filesystem.WriteFileOptions{
		CreateDirs: true,
		Overwrite:  true,
	})
	if err != nil {
		issues.AddMessage(IssueEntryCopyFailed, module.plan.Name(), targetPath, err.Error())
		return Operation{}, false
	}

	return newOperation(OperationCopyFile, entry.SourcePath(), targetPath), true
}

// copyDirectoryEntry copies one explicit directory entry into the target
// worktree.
func (s Service) copyDirectoryEntry(
	ctx context.Context,
	module moduleContext,
	entry source.EntrySnapshot,
	targetPath string,
	issues *issueCollector,
) (Operation, bool) {
	_, err := s.deps.FS.CopyTree(ctx, entry.SourcePath(), targetPath, filesystem.CopyTreeOptions{
		PreserveMode:  true,
		SymlinkPolicy: filesystem.SymlinkReject,
		SafetyRoot:    module.workspace.WorktreeDir(),
	})
	if err != nil {
		issues.AddMessage(IssueEntryCopyFailed, module.plan.Name(), entry.SourcePath(), err.Error())
		return Operation{}, false
	}

	return newOperation(OperationCopyDirectory, entry.SourcePath(), targetPath), true
}

// appendProvenanceFile writes generated provenance when enabled and appends the
// corresponding operation.
func (s Service) appendProvenanceFile(
	ctx context.Context,
	req Request,
	module moduleContext,
	operations *[]Operation,
	issues *issueCollector,
) bool {
	if !req.Plan.PublishPolicy().Provenance().FileEnabled() || !s.opts.GenerateProvenanceFile {
		return true
	}

	path := pathutil.JoinRelative(
		module.workspace.WorktreeDir(),
		req.Plan.PublishPolicy().Provenance().File(),
	)
	if err := pathutil.EnsureInside(module.workspace.WorktreeDir(), path); err != nil {
		issues.AddMessage(IssueTargetPathEscape, module.plan.Name(), path, err.Error())
		return false
	}

	data, _ := json.MarshalIndent(provenancePayload(req, module), "", "  ")
	err := s.deps.FS.WriteFile(ctx, path, append(data, '\n'), filesystem.WriteFileOptions{
		CreateDirs: true,
		Overwrite:  true,
	})
	if err != nil {
		issues.AddMessage(IssueEntryCopyFailed, module.plan.Name(), path, err.Error())
		return true
	}

	*operations = append(*operations, newOperation(OperationWriteGenerated, "", path))
	return true
}

// provenancePayload renders the stable generated metadata file contents.
func provenancePayload(req Request, module moduleContext) map[string]string {
	return map[string]string{
		"module":       module.plan.Name().String(),
		"modulePath":   module.plan.ModulePath().String(),
		"version":      module.plan.Version().String(),
		"sourceCommit": req.Source.Repository().Head().String(),
		"sourceBranch": req.Source.Repository().Branch().String(),
		"sourceHash":   module.source.Hash().String(),
	}
}

// findSourceModule locates an inspected source module by name.
func findSourceModule(
	snapshot source.Snapshot,
	name manifest.ModuleName,
) (source.ModuleSnapshot, bool) {
	for _, mod := range snapshot.Modules() {
		if mod.Name() == name {
			return mod, true
		}
	}
	return source.ModuleSnapshot{}, false
}
