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

package modulefile

import (
	"context"

	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/ports/filesystem"
	"arcoris.dev/arcoris-publisher/internal/workflow/pathutil"
	"arcoris.dev/arcoris-publisher/internal/workflow/target"
)

// Service rewrites target go.mod files using effective plan data.
type Service struct {
	// deps contains infrastructure ports used by go.mod rewriting.
	deps Dependencies

	// opts contains normalized behavior toggles.
	opts Options
}

// New returns a module-file rewrite service.
func New(deps Dependencies, opts Options) Service {
	if opts == (Options{}) {
		opts = DefaultOptions()
	}
	return Service{deps: deps, opts: opts}
}

// Rewrite updates target go.mod files for every planned module.
func (s Service) Rewrite(ctx context.Context, req Request) (Result, error) {
	if err := s.validateRequest(req); err != nil {
		return Result{}, err
	}

	results, err := s.rewriteModules(ctx, req)
	if err != nil {
		return Result{}, err
	}

	return Result{modules: results}, nil
}

// validateRequest rejects incomplete inputs before filesystem mutation starts.
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

// rewriteModules rewrites all planned target go.mod files and aggregates
// validation failures across modules.
func (s Service) rewriteModules(ctx context.Context, req Request) ([]ModuleResult, error) {
	issues := newIssueCollector()
	results := make([]ModuleResult, 0, req.Plan.Len())

	for _, mod := range req.Plan.Modules() {
		result, ok := s.rewriteModule(ctx, mod, req.Targets, &issues)
		if ok {
			results = append(results, result)
		}
	}

	if err := issues.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

// rewriteModule resolves one target go.mod path, reads it, applies managed
// dependency updates, and writes it back only when content changes.
func (s Service) rewriteModule(
	ctx context.Context,
	mod plan.ModulePlan,
	targets target.WorkspaceSet,
	issues *issueCollector,
) (ModuleResult, bool) {
	ws, ok := targets.WorkspaceByModule(mod.Name())
	if !ok {
		issues.Add(IssueMissingTarget, mod.Name(), "target", "target workspace is missing")
		return ModuleResult{}, false
	}

	goMod, ok := resolveGoModPath(mod, ws, issues)
	if !ok {
		return ModuleResult{}, false
	}

	data, ok := s.readGoMod(ctx, mod, goMod, issues)
	if !ok {
		return ModuleResult{}, false
	}

	newData, updates, changed := rewriteGoMod(data, mod, s.opts.RemoveLocalReplaces)
	if changed && !s.writeGoMod(ctx, mod, goMod, newData, issues) {
		return ModuleResult{}, false
	}

	return ModuleResult{
		module:       mod.Name(),
		goModPath:    goMod,
		changed:      changed,
		requirements: updates,
	}, true
}

// resolveGoModPath joins the module root and go.mod path under the worktree and
// rejects any resolved path outside the target worktree.
func resolveGoModPath(
	mod plan.ModulePlan,
	ws target.ModuleWorkspace,
	issues *issueCollector,
) (string, bool) {
	moduleRoot := pathutil.JoinRelative(ws.WorktreeDir(), mod.ModuleRoot())
	goMod := pathutil.JoinRelative(moduleRoot, mod.GoMod())

	if err := pathutil.EnsureInside(ws.WorktreeDir(), goMod); err != nil {
		issues.AddMessage(IssueGoModRewriteFailed, mod.Name(), goMod, err.Error())
		return "", false
	}

	return goMod, true
}

// readGoMod loads the existing target go.mod and maps absence separately from
// read failures so callers get precise diagnostics.
func (s Service) readGoMod(
	ctx context.Context,
	mod plan.ModulePlan,
	goMod string,
	issues *issueCollector,
) ([]byte, bool) {
	exists, err := s.deps.FS.Exists(ctx, goMod)
	if err != nil {
		issues.AddMessage(IssueGoModMissing, mod.Name(), goMod, err.Error())
		return nil, false
	}
	if !exists {
		issues.Add(IssueGoModMissing, mod.Name(), goMod, "go.mod is missing")
		return nil, false
	}

	data, err := s.deps.FS.ReadFile(ctx, goMod)
	if err != nil {
		issues.AddMessage(IssueGoModRewriteFailed, mod.Name(), goMod, err.Error())
		return nil, false
	}

	return data, true
}

// writeGoMod persists changed go.mod content.
func (s Service) writeGoMod(
	ctx context.Context,
	mod plan.ModulePlan,
	goMod string,
	data []byte,
	issues *issueCollector,
) bool {
	err := s.deps.FS.WriteFile(ctx, goMod, data, filesystem.WriteFileOptions{
		CreateDirs: false,
		Overwrite:  true,
	})
	if err != nil {
		issues.AddMessage(IssueGoModRewriteFailed, mod.Name(), goMod, err.Error())
		return false
	}

	return true
}
