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
	issues := newIssueCollector()
	if req.Plan.Empty() {
		issues.Add(IssueInvalidRequest, "", "plan", "plan is empty")
	}
	if s.deps.FS == nil {
		issues.Add(IssueInvalidRequest, "", "fs", "filesystem dependency is required")
	}
	if err := issues.Err(); err != nil {
		return Result{}, err
	}
	results := make([]ModuleResult, 0, req.Plan.Len())
	for _, mod := range req.Plan.Modules() {
		mr, ok := s.constructModule(ctx, req, mod.Name(), &issues)
		if ok {
			results = append(results, mr)
		}
	}
	if err := issues.Err(); err != nil {
		return Result{}, err
	}
	return Result{modules: results}, nil
}

// constructModule cleans one target worktree and copies only inspected explicit
// publish entries for the same module.
func (s Service) constructModule(
	ctx context.Context,
	req Request,
	name manifest.ModuleName,
	issues *issueCollector,
) (ModuleResult, bool) {
	mod, _ := req.Plan.ModuleByName(name)
	src, ok := findSourceModule(req.Source, name)
	if !ok {
		issues.Add(IssueMissingSource, name, "source", "source snapshot for module %q is missing", name)
		return ModuleResult{}, false
	}
	ws, ok := req.Targets.WorkspaceByModule(name)
	if !ok {
		issues.Add(IssueMissingTarget, name, "target", "target workspace for module %q is missing", name)
		return ModuleResult{}, false
	}
	operations := []Operation{}
	preserve := []string{}
	if s.opts.PreserveGitDir {
		preserve = append(preserve, ".git")
	}
	if err := s.deps.FS.CleanDir(ctx, ws.WorktreeDir(), filesystem.CleanDirOptions{
		Preserve:      preserve,
		RequireGitDir: false,
		AllowMissing:  false,
		SafetyRoot:    ws.WorktreeDir(),
	}); err != nil {
		issues.AddMessage(IssueCleanFailed, name, ws.WorktreeDir(), err.Error())
		return ModuleResult{}, false
	}
	operations = append(operations, newOperation(OperationClean, "", ws.WorktreeDir()))
	for _, entry := range src.Entries() {
		if !entry.Present() {
			operations = append(operations, newOperation(
				OperationSkipOptional,
				entry.SourcePath(),
				pathutil.JoinRelative(ws.WorktreeDir(), entry.TargetPath()),
			))
			continue
		}
		targetPath := pathutil.JoinRelative(ws.WorktreeDir(), entry.TargetPath())
		if err := pathutil.EnsureInside(ws.WorktreeDir(), targetPath); err != nil {
			issues.AddMessage(IssueTargetPathEscape, name, targetPath, err.Error())
			continue
		}
		switch entry.Kind() {
		case manifest.PublishEntryFile:
			data, err := s.deps.FS.ReadFile(ctx, entry.SourcePath())
			if err != nil {
				issues.AddMessage(IssueEntryCopyFailed, name, entry.SourcePath(), err.Error())
				continue
			}
			if err := s.deps.FS.WriteFile(ctx, targetPath, data, filesystem.WriteFileOptions{
				CreateDirs: true,
				Overwrite:  true,
			}); err != nil {
				issues.AddMessage(IssueEntryCopyFailed, name, targetPath, err.Error())
				continue
			}
			operations = append(operations, newOperation(
				OperationCopyFile,
				entry.SourcePath(),
				targetPath,
			))
		case manifest.PublishEntryDirectory:
			_, err := s.deps.FS.CopyTree(ctx, entry.SourcePath(), targetPath, filesystem.CopyTreeOptions{
				PreserveMode:  true,
				SymlinkPolicy: filesystem.SymlinkReject,
				SafetyRoot:    ws.WorktreeDir(),
			})
			if err != nil {
				issues.AddMessage(IssueEntryCopyFailed, name, entry.SourcePath(), err.Error())
				continue
			}
			operations = append(operations, newOperation(
				OperationCopyDirectory,
				entry.SourcePath(),
				targetPath,
			))
		}
	}
	if req.Plan.PublishPolicy().Provenance().FileEnabled() && s.opts.GenerateProvenanceFile {
		path := pathutil.JoinRelative(ws.WorktreeDir(), req.Plan.PublishPolicy().Provenance().File())
		if err := pathutil.EnsureInside(ws.WorktreeDir(), path); err != nil {
			issues.AddMessage(IssueTargetPathEscape, name, path, err.Error())
			return ModuleResult{}, false
		}
		payload := map[string]string{
			"module":       name.String(),
			"modulePath":   mod.ModulePath().String(),
			"version":      mod.Version().String(),
			"sourceCommit": req.Source.Repository().Head().String(),
			"sourceBranch": req.Source.Repository().Branch().String(),
			"sourceHash":   src.Hash().String(),
		}
		data, _ := json.MarshalIndent(payload, "", "  ")
		if err := s.deps.FS.WriteFile(ctx, path, append(data, '\n'), filesystem.WriteFileOptions{
			CreateDirs: true,
			Overwrite:  true,
		}); err != nil {
			issues.AddMessage(IssueEntryCopyFailed, name, path, err.Error())
		} else {
			operations = append(operations, newOperation(OperationWriteGenerated, "", path))
		}
	}
	return ModuleResult{
		module:      name,
		worktreeDir: filepath.Clean(ws.WorktreeDir()),
		operations:  operations,
		changed:     len(operations) > 0,
	}, true
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

var _ = target.WorkspaceSet{}
