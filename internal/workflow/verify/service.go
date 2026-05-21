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

package verify

import (
	"context"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/ports/gotoolchain"
	"arcoris.dev/arcoris-publisher/internal/workflow/pathutil"
)

// Service verifies constructed target repositories before publication.
type Service struct {
	// deps contains infrastructure ports used by verification.
	deps Dependencies

	// opts contains normalized verification options.
	opts Options
}

// New returns a verification service.
func New(deps Dependencies, opts Options) Service {
	if opts == (Options{}) {
		opts = DefaultOptions()
	}
	return Service{deps: deps, opts: opts}
}

// Verify runs verification checks and returns check results.
//
// A non-nil error means verification could not run. Failed checks are reported
// through Result.Failed so callers can distinguish infrastructure errors from
// target state failures.
func (s Service) Verify(ctx context.Context, req Request) (Result, error) {
	if req.Plan.Empty() {
		return Result{}, &Error{Code: CodeInvalidRequest, Message: "plan is empty"}
	}
	if s.deps.FS == nil {
		return Result{}, &Error{
			Code:    CodeDependencyMissing,
			Message: "filesystem dependency is required",
		}
	}
	modules := make([]ModuleResult, 0, req.Plan.Len())
	for _, mod := range req.Plan.Modules() {
		modules = append(modules, s.verifyModule(ctx, req, mod.Name()))
	}
	return Result{modules: modules}, nil
}

// verifyModule runs target, go.mod, Go toolchain, and cleanliness checks for
// one planned module.
func (s Service) verifyModule(
	ctx context.Context,
	req Request,
	name manifest.ModuleName,
) ModuleResult {
	mod, _ := req.Plan.ModuleByName(name)
	ws, ok := req.Targets.WorkspaceByModule(name)
	checks := []CheckResult{}
	if !ok {
		return ModuleResult{
			module: name,
			checks: []CheckResult{
				NewCheckResult("target-worktree", StatusFailed, SeverityError, "target workspace is missing"),
			},
		}
	}
	if ok, err := s.deps.FS.Exists(ctx, ws.WorktreeDir()); err != nil || !ok {
		checks = append(checks, NewCheckResult(
			"target-worktree",
			StatusFailed,
			SeverityError,
			"target worktree is missing",
		).withPath(ws.WorktreeDir()))
	} else if isDir, err := s.deps.FS.IsDir(ctx, ws.WorktreeDir()); err != nil || !isDir {
		checks = append(checks, NewCheckResult(
			"target-worktree",
			StatusFailed,
			SeverityError,
			"target worktree is not a directory",
		).withPath(ws.WorktreeDir()))
	} else {
		checks = append(checks, NewCheckResult(
			"target-worktree",
			StatusPassed,
			SeverityInfo,
			"target worktree exists",
		).withPath(ws.WorktreeDir()))
	}
	moduleRoot := pathutil.JoinRelative(ws.WorktreeDir(), mod.ModuleRoot())
	goModPath := pathutil.JoinRelative(moduleRoot, mod.GoMod())
	data, err := s.deps.FS.ReadFile(ctx, goModPath)
	if err != nil {
		checks = append(checks, NewCheckResult(
			"go-mod",
			StatusFailed,
			SeverityError,
			"go.mod cannot be read",
		).withPath(goModPath))
		return ModuleResult{module: name, checks: checks}
	}
	info := parseGoMod(data)
	if info.module != mod.ModulePath().String() {
		checks = append(checks, NewCheckResult(
			"go-mod-module",
			StatusFailed,
			SeverityError,
			"go.mod module path does not match plan",
		).withPath(goModPath))
	} else {
		checks = append(checks, NewCheckResult(
			"go-mod-module",
			StatusPassed,
			SeverityInfo,
			"go.mod module path matches plan",
		).withPath(goModPath))
	}
	for _, rep := range info.localReplaces {
		if mod.Verification().LocalReplacePolicy() == manifest.LocalReplacePolicyForbid {
			checks = append(checks, NewCheckResult(
				"go-mod-replace",
				StatusFailed,
				SeverityError,
				"local replace directive is forbidden: "+rep,
			).withPath(goModPath))
		} else if mod.Verification().LocalReplacePolicy() == manifest.LocalReplacePolicyWarn {
			checks = append(checks, NewCheckResult(
				"go-mod-replace",
				StatusWarning,
				SeverityWarning,
				"local replace directive is present: "+rep,
			).withPath(goModPath))
		}
	}
	for _, reqr := range mod.Requirements() {
		if got := info.requires[reqr.ModulePath().String()]; got != reqr.Version().String() {
			checks = append(checks, NewCheckResult(
				"go-mod-require",
				StatusFailed,
				SeverityError,
				"requirement mismatch for "+reqr.ModulePath().String(),
			).withPath(goModPath))
		}
	}
	goPolicy := mod.Verification().Go()
	common := gotoolchain.CommonOptions{GoBinary: s.opts.GoBinary, Timeout: s.opts.Timeout}
	if goPolicy.WorkspaceMode() == manifest.GoWorkspaceModeOff {
		common.WorkspaceMode = gotoolchain.WorkspaceOff
	} else {
		common.WorkspaceMode = gotoolchain.WorkspaceDefault
	}
	if s.deps.Go != nil && goPolicy.List() {
		if _, err := s.deps.Go.List(ctx, moduleRoot, gotoolchain.ListOptions{
			CommonOptions: common,
			Patterns:      goPolicy.Patterns(),
			Test:          true,
		}); err != nil {
			checks = append(checks, NewCheckResult("go-list", StatusFailed, SeverityError, err.Error()))
		} else {
			checks = append(checks, NewCheckResult(
				"go-list",
				StatusPassed,
				SeverityInfo,
				"go list succeeded",
			))
		}
	}
	if s.deps.Go != nil && goPolicy.Test() {
		if _, err := s.deps.Go.Test(ctx, moduleRoot, gotoolchain.TestOptions{
			CommonOptions: common,
			Patterns:      goPolicy.Patterns(),
		}); err != nil {
			checks = append(checks, NewCheckResult("go-test", StatusFailed, SeverityError, err.Error()))
		} else {
			checks = append(checks, NewCheckResult(
				"go-test",
				StatusPassed,
				SeverityInfo,
				"go test succeeded",
			))
		}
	}
	if s.deps.Go != nil && goPolicy.Tidy() {
		if _, err := s.deps.Go.ModTidy(ctx, moduleRoot, gotoolchain.ModTidyOptions{
			CommonOptions: common,
		}); err != nil {
			checks = append(checks, NewCheckResult("go-mod-tidy", StatusFailed, SeverityError, err.Error()))
		} else {
			checks = append(checks, NewCheckResult(
				"go-mod-tidy",
				StatusPassed,
				SeverityInfo,
				"go mod tidy succeeded",
			))
		}
	}
	if s.opts.RequireClean && s.deps.Git != nil {
		status, err := s.deps.Git.Status(ctx, ws.WorktreeDir())
		if err == nil && (!status.Clean || status.HasEntries()) {
			checks = append(checks, NewCheckResult(
				"git-clean",
				StatusFailed,
				SeverityError,
				"target worktree changed during verification",
			))
		}
	}
	return ModuleResult{module: name, checks: checks}
}
