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
	"bytes"
	"context"
	"path/filepath"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/ports/git"
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
	if err := s.validateRequest(req); err != nil {
		return Result{}, err
	}

	return Result{modules: s.verifyModules(ctx, req)}, nil
}

// validateRequest rejects invalid verification inputs. Failed verification
// checks are returned in Result, but invalid wiring is returned as an error.
func (s Service) validateRequest(req Request) error {
	if req.Plan.Empty() {
		return &Error{Code: CodeInvalidRequest, Message: "plan is empty"}
	}
	if s.deps.FS == nil {
		return &Error{
			Code:    CodeDependencyMissing,
			Message: "filesystem dependency is required",
		}
	}

	return nil
}

// verifyModules verifies every planned module in deterministic plan order.
func (s Service) verifyModules(ctx context.Context, req Request) []ModuleResult {
	modules := make([]ModuleResult, 0, req.Plan.Len())
	for _, mod := range req.Plan.Modules() {
		modules = append(modules, s.verifyModule(ctx, req, mod.Name()))
	}

	return modules
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
		return missingTargetResult(name)
	}

	checks = append(checks, s.targetWorktreeCheck(ctx, ws.WorktreeDir()))
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
	checks = append(checks, goModModuleCheck(mod, info, goModPath))
	checks = append(checks, localReplaceChecks(mod, info, goModPath)...)
	checks = append(checks, requirementChecks(mod, info, goModPath)...)
	checks = append(checks, s.goToolchainChecks(ctx, mod, ws.WorktreeDir(), moduleRoot)...)
	checks = append(checks, s.gitCleanChecks(ctx, ws.WorktreeDir())...)

	return ModuleResult{module: name, checks: checks}
}

// missingTargetResult reports a missing upstream target preparation result.
func missingTargetResult(name manifest.ModuleName) ModuleResult {
	return ModuleResult{
		module: name,
		checks: []CheckResult{
			NewCheckResult("target-worktree", StatusFailed, SeverityError, "target workspace is missing"),
		},
	}
}

// targetWorktreeCheck verifies that the prepared target worktree exists and is
// a directory.
func (s Service) targetWorktreeCheck(ctx context.Context, worktree string) CheckResult {
	if ok, err := s.deps.FS.Exists(ctx, worktree); err != nil || !ok {
		return NewCheckResult(
			"target-worktree",
			StatusFailed,
			SeverityError,
			"target worktree is missing",
		).withPath(worktree)
	}

	if isDir, err := s.deps.FS.IsDir(ctx, worktree); err != nil || !isDir {
		return NewCheckResult(
			"target-worktree",
			StatusFailed,
			SeverityError,
			"target worktree is not a directory",
		).withPath(worktree)
	}

	return NewCheckResult(
		"target-worktree",
		StatusPassed,
		SeverityInfo,
		"target worktree exists",
	).withPath(worktree)
}

// goModModuleCheck compares the target go.mod module declaration with the plan.
func goModModuleCheck(
	mod plan.ModulePlan,
	info goModInfo,
	goModPath string,
) CheckResult {
	if info.module != mod.ModulePath().String() {
		return NewCheckResult(
			"go-mod-module",
			StatusFailed,
			SeverityError,
			"go.mod module path does not match plan",
		).withPath(goModPath)
	}

	return NewCheckResult(
		"go-mod-module",
		StatusPassed,
		SeverityInfo,
		"go.mod module path matches plan",
	).withPath(goModPath)
}

// localReplaceChecks applies the effective local replace policy.
func localReplaceChecks(
	mod plan.ModulePlan,
	info goModInfo,
	goModPath string,
) []CheckResult {
	checks := []CheckResult{}
	for _, replacement := range info.localReplaces {
		check, ok := localReplaceCheck(mod, replacement, goModPath)
		if ok {
			checks = append(checks, check)
		}
	}

	return checks
}

// localReplaceCheck converts one local replace directive into a policy-specific
// verification check.
func localReplaceCheck(
	mod plan.ModulePlan,
	replacement string,
	goModPath string,
) (CheckResult, bool) {
	switch mod.Verification().LocalReplacePolicy() {
	case manifest.LocalReplacePolicyForbid:
		return NewCheckResult(
			"go-mod-replace",
			StatusFailed,
			SeverityError,
			"local replace directive is forbidden: "+replacement,
		).withPath(goModPath), true
	case manifest.LocalReplacePolicyWarn:
		return NewCheckResult(
			"go-mod-replace",
			StatusWarning,
			SeverityWarning,
			"local replace directive is present: "+replacement,
		).withPath(goModPath), true
	default:
		return CheckResult{}, false
	}
}

// requirementChecks verifies direct internal dependency requirements.
func requirementChecks(
	mod plan.ModulePlan,
	info goModInfo,
	goModPath string,
) []CheckResult {
	checks := []CheckResult{}
	for _, requirement := range mod.Requirements() {
		if got := info.requires[requirement.ModulePath().String()]; got != requirement.Version().String() {
			checks = append(checks, NewCheckResult(
				"go-mod-require",
				StatusFailed,
				SeverityError,
				"requirement mismatch for "+requirement.ModulePath().String(),
			).withPath(goModPath))
		}
	}

	return checks
}

// goToolchainChecks runs configured Go checks when a Go toolchain port is
// provided.
func (s Service) goToolchainChecks(
	ctx context.Context,
	mod plan.ModulePlan,
	worktree string,
	moduleRoot string,
) []CheckResult {
	if s.deps.Go == nil {
		return nil
	}

	policy := mod.Verification().Go()
	common := commonGoOptions(policy, s.opts)
	checks := []CheckResult{}

	if policy.List() {
		checks = append(checks, s.goListCheck(ctx, moduleRoot, common, policy.Patterns()))
	}
	if policy.Test() {
		checks = append(checks, s.goTestCheck(ctx, moduleRoot, common, policy.Patterns()))
	}
	if policy.Tidy() {
		checks = append(checks, s.goTidyCheck(ctx, worktree, moduleRoot, common))
	}

	return checks
}

func commonGoOptions(
	policy manifest.GoVerificationPolicy,
	opts Options,
) gotoolchain.CommonOptions {
	common := gotoolchain.CommonOptions{
		GoBinary: opts.GoBinary,
		Timeout:  opts.Timeout,
	}
	if policy.WorkspaceMode() == manifest.GoWorkspaceModeOff {
		common.WorkspaceMode = gotoolchain.WorkspaceOff
	} else {
		common.WorkspaceMode = gotoolchain.WorkspaceDefault
	}

	return common
}

// goListCheck runs the configured go list verification.
func (s Service) goListCheck(
	ctx context.Context,
	moduleRoot string,
	common gotoolchain.CommonOptions,
	patterns []string,
) CheckResult {
	_, err := s.deps.Go.List(ctx, moduleRoot, gotoolchain.ListOptions{
		CommonOptions: common,
		Patterns:      patterns,
		Test:          true,
	})
	if err != nil {
		return NewCheckResult("go-list", StatusFailed, SeverityError, err.Error())
	}

	return NewCheckResult("go-list", StatusPassed, SeverityInfo, "go list succeeded")
}

// goTestCheck runs the configured go test verification.
func (s Service) goTestCheck(
	ctx context.Context,
	moduleRoot string,
	common gotoolchain.CommonOptions,
	patterns []string,
) CheckResult {
	_, err := s.deps.Go.Test(ctx, moduleRoot, gotoolchain.TestOptions{
		CommonOptions: common,
		Patterns:      patterns,
	})
	if err != nil {
		return NewCheckResult("go-test", StatusFailed, SeverityError, err.Error())
	}

	return NewCheckResult("go-test", StatusPassed, SeverityInfo, "go test succeeded")
}

// goTidyCheck runs go mod tidy as a verification step.
func (s Service) goTidyCheck(
	ctx context.Context,
	worktree string,
	moduleRoot string,
	common gotoolchain.CommonOptions,
) CheckResult {
	before := s.snapshotGoModuleFiles(ctx, moduleRoot)

	_, err := s.deps.Go.ModTidy(ctx, moduleRoot, gotoolchain.ModTidyOptions{
		CommonOptions: common,
	})
	if err != nil {
		return NewCheckResult("go-mod-tidy", StatusFailed, SeverityError, err.Error())
	}

	if status, ok := s.goModuleGitStatus(ctx, worktree, moduleRoot); ok && status.IsDirty() {
		return NewCheckResult(
			"go-mod-tidy",
			StatusFailed,
			SeverityError,
			"go mod tidy changed go.mod or go.sum",
		)
	}

	after := s.snapshotGoModuleFiles(ctx, moduleRoot)
	if before.changed(after) {
		return NewCheckResult(
			"go-mod-tidy",
			StatusFailed,
			SeverityError,
			"go mod tidy changed go.mod or go.sum",
		)
	}

	return NewCheckResult("go-mod-tidy", StatusPassed, SeverityInfo, "go mod tidy succeeded")
}

// goModuleGitStatus returns a status containing only go.mod/go.sum entries when
// Git status is available.
func (s Service) goModuleGitStatus(
	ctx context.Context,
	worktree string,
	moduleRoot string,
) (git.Status, bool) {
	if s.deps.Git == nil {
		return git.Status{}, false
	}

	status, err := s.deps.Git.Status(ctx, worktree)
	if err != nil {
		return git.Status{}, false
	}

	paths := goModuleStatusPaths(worktree, moduleRoot)
	entries := make([]git.StatusEntry, 0, len(status.Entries))
	for _, entry := range status.Entries {
		if _, ok := paths[filepath.ToSlash(filepath.Clean(entry.Path))]; ok {
			entries = append(entries, entry)
		}
	}

	return git.Status{Clean: len(entries) == 0, Entries: entries}, true
}

// goModuleStatusPaths returns repository-relative paths that go mod tidy may
// mutate.
func goModuleStatusPaths(worktree string, moduleRoot string) map[string]struct{} {
	out := map[string]struct{}{}
	for _, name := range []string{"go.mod", "go.sum"} {
		rel, err := filepath.Rel(worktree, filepath.Join(moduleRoot, name))
		if err != nil {
			continue
		}
		out[filepath.ToSlash(filepath.Clean(rel))] = struct{}{}
	}

	return out
}

type goModuleFileSnapshot map[string]goModuleFile

type goModuleFile struct {
	exists bool
	data   []byte
}

func (s Service) snapshotGoModuleFiles(
	ctx context.Context,
	moduleRoot string,
) goModuleFileSnapshot {
	snapshot := goModuleFileSnapshot{}
	for _, name := range []string{"go.mod", "go.sum"} {
		path := filepath.Join(moduleRoot, name)
		snapshot[name] = s.snapshotGoModuleFile(ctx, path)
	}

	return snapshot
}

func (s Service) snapshotGoModuleFile(ctx context.Context, path string) goModuleFile {
	exists, err := s.deps.FS.Exists(ctx, path)
	if err != nil || !exists {
		return goModuleFile{}
	}

	data, err := s.deps.FS.ReadFile(ctx, path)
	if err != nil {
		return goModuleFile{}
	}

	return goModuleFile{exists: true, data: data}
}

func (s goModuleFileSnapshot) changed(other goModuleFileSnapshot) bool {
	for _, name := range []string{"go.mod", "go.sum"} {
		left := s[name]
		right := other[name]
		if left.exists != right.exists {
			return true
		}
		if !bytes.Equal(left.data, right.data) {
			return true
		}
	}

	return false
}

// gitCleanChecks verifies that target verification did not leave the worktree
// dirty when clean verification is required.
func (s Service) gitCleanChecks(ctx context.Context, worktree string) []CheckResult {
	if !s.opts.RequireClean || s.deps.Git == nil {
		return nil
	}

	status, err := s.deps.Git.Status(ctx, worktree)
	if err == nil && (!status.Clean || status.HasEntries()) {
		return []CheckResult{NewCheckResult(
			"git-clean",
			StatusFailed,
			SeverityError,
			"target worktree changed during verification",
		)}
	}

	return nil
}
