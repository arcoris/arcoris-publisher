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

	"arcoris.dev/arcoris-publisher/internal/ports/filesystem"
	"arcoris.dev/arcoris-publisher/internal/workflow/pathutil"
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
		ws, ok := req.Targets.WorkspaceByModule(mod.Name())
		if !ok {
			issues.Add(IssueMissingTarget, mod.Name(), "target", "target workspace is missing")
			continue
		}
		moduleRoot := pathutil.JoinRelative(ws.WorktreeDir(), mod.ModuleRoot())
		goMod := pathutil.JoinRelative(moduleRoot, mod.GoMod())
		if err := pathutil.EnsureInside(ws.WorktreeDir(), goMod); err != nil {
			issues.AddMessage(IssueGoModRewriteFailed, mod.Name(), goMod, err.Error())
			continue
		}
		exists, err := s.deps.FS.Exists(ctx, goMod)
		if err != nil {
			issues.AddMessage(IssueGoModMissing, mod.Name(), goMod, err.Error())
			continue
		}
		if !exists {
			issues.Add(IssueGoModMissing, mod.Name(), goMod, "go.mod is missing")
			continue
		}
		data, err := s.deps.FS.ReadFile(ctx, goMod)
		if err != nil {
			issues.AddMessage(IssueGoModRewriteFailed, mod.Name(), goMod, err.Error())
			continue
		}
		newData, updates, changed := rewriteGoMod(data, mod, s.opts.RemoveLocalReplaces)
		if changed {
			if err := s.deps.FS.WriteFile(ctx, goMod, newData, filesystem.WriteFileOptions{
				CreateDirs: false,
				Overwrite:  true,
			}); err != nil {
				issues.AddMessage(IssueGoModRewriteFailed, mod.Name(), goMod, err.Error())
				continue
			}
		}
		results = append(results, ModuleResult{
			module:       mod.Name(),
			goModPath:    goMod,
			changed:      changed,
			requirements: updates,
		})
	}
	if err := issues.Err(); err != nil {
		return Result{}, err
	}
	return Result{modules: results}, nil
}
