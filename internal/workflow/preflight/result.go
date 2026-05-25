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

package preflight

import "arcoris.dev/arcoris-publisher/internal/manifest"

// Status is the stable readiness state for a preflight result.
type Status string

const (
	StatusPassed  Status = "passed"
	StatusFailed  Status = "failed"
	StatusSkipped Status = "skipped"
	StatusWarning Status = "warning"
)

// Severity identifies whether a check blocks publication.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Result contains global and per-module publish readiness checks.
type Result struct {
	status  Status
	version string
	checks  []CheckResult
	modules []ModuleResult
}

// Status returns the aggregate preflight status.
func (r Result) Status() Status { return r.status }

// Version returns the requested publication version.
func (r Result) Version() string { return r.version }

// Checks returns detached global checks.
func (r Result) Checks() []CheckResult {
	out := make([]CheckResult, len(r.checks))
	copy(out, r.checks)
	return out
}

// Modules returns detached module checks.
func (r Result) Modules() []ModuleResult {
	out := make([]ModuleResult, len(r.modules))
	copy(out, r.modules)
	return out
}

// Failed reports whether any blocking check failed.
func (r Result) Failed() bool { return r.status == StatusFailed }

// Empty reports whether no checks have run.
func (r Result) Empty() bool { return len(r.checks) == 0 && len(r.modules) == 0 }

// CheckResult describes one stable readiness check.
type CheckResult struct {
	name     string
	status   Status
	severity Severity
	code     string
	message  string
	path     string
}

// Name returns the stable check identifier.
func (c CheckResult) Name() string { return c.name }

// Status returns this check's state.
func (c CheckResult) Status() Status { return c.status }

// Severity returns this check's severity.
func (c CheckResult) Severity() Severity { return c.severity }

// Code returns a stable machine-readable detail code.
func (c CheckResult) Code() string { return c.code }

// Message returns a human-readable diagnostic.
func (c CheckResult) Message() string { return c.message }

// Path returns the local path related to the check, if any.
func (c CheckResult) Path() string { return c.path }

// ModuleResult contains checks for one planned module.
type ModuleResult struct {
	name        manifest.ModuleName
	repository  manifest.RepositoryRef
	worktreeDir string
	status      Status
	checks      []CheckResult
}

// Name returns the module name.
func (m ModuleResult) Name() manifest.ModuleName { return m.name }

// Repository returns the target repository identity.
func (m ModuleResult) Repository() manifest.RepositoryRef { return m.repository }

// WorktreeDir returns the expected local target worktree path.
func (m ModuleResult) WorktreeDir() string { return m.worktreeDir }

// Status returns the aggregate module status.
func (m ModuleResult) Status() Status { return m.status }

// Checks returns detached module checks.
func (m ModuleResult) Checks() []CheckResult {
	out := make([]CheckResult, len(m.checks))
	copy(out, m.checks)
	return out
}

type resultBuilder struct{ result Result }

func (b *resultBuilder) addGlobal(check CheckResult) {
	b.result.checks = append(b.result.checks, check)
}

func (b *resultBuilder) addModule(module ModuleResult) {
	module.status = aggregate(module.checks)
	b.result.modules = append(b.result.modules, module)
}

func (b *resultBuilder) build() Result {
	b.result.status = aggregate(append(b.result.checks, moduleChecks(b.result.modules)...))
	return b.result
}

func moduleChecks(modules []ModuleResult) []CheckResult {
	out := []CheckResult{}
	for _, mod := range modules {
		out = append(out, mod.checks...)
	}
	return out
}

func aggregate(checks []CheckResult) Status {
	if len(checks) == 0 {
		return StatusSkipped
	}
	status := StatusPassed
	for _, check := range checks {
		switch check.Status() {
		case StatusFailed:
			return StatusFailed
		case StatusWarning:
			status = StatusWarning
		}
	}
	return status
}

func passed(name, message string) CheckResult {
	return CheckResult{name: name, status: StatusPassed, severity: SeverityInfo, message: message}
}

func skipped(name, message string) CheckResult {
	return CheckResult{name: name, status: StatusSkipped, severity: SeverityInfo, message: message}
}

func warning(name, code, message string) CheckResult {
	return CheckResult{name: name, status: StatusWarning, severity: SeverityWarning, code: code, message: message}
}

func failed(name, code, message string) CheckResult {
	return CheckResult{name: name, status: StatusFailed, severity: SeverityError, code: code, message: message}
}

func pathCheck(check CheckResult, path string) CheckResult {
	check.path = path
	return check
}
