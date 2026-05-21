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

import "arcoris.dev/arcoris-publisher/internal/manifest"

// Status identifies the outcome of one verification check.
type Status string

const (
	// StatusPassed means the check ran and found no issue.
	StatusPassed Status = "passed"

	// StatusFailed means the check ran and found a blocking issue.
	StatusFailed Status = "failed"

	// StatusSkipped means the check was intentionally not run.
	StatusSkipped Status = "skipped"

	// StatusWarning means the check found a non-blocking issue.
	StatusWarning Status = "warning"
)

// Severity identifies how a check result should be treated.
type Severity string

const (
	// SeverityError marks blocking verification failures.
	SeverityError Severity = "error"

	// SeverityWarning marks non-blocking verification warnings.
	SeverityWarning Severity = "warning"

	// SeverityInfo marks informational successful checks.
	SeverityInfo Severity = "info"
)

// CheckName identifies a verification check.
type CheckName string

// Result describes verification checks in plan module order.
type Result struct{ modules []ModuleResult }

// Modules returns detached module verification results.
func (r Result) Modules() []ModuleResult {
	out := make([]ModuleResult, len(r.modules))
	copy(out, r.modules)
	return out
}

// Passed reports whether no module has failed checks.
func (r Result) Passed() bool { return !r.Failed() }

// Failed reports whether at least one check failed.
//
// Failed is distinct from a non-nil verification error. A failed result means
// verification ran and found target state that should not be published.
func (r Result) Failed() bool {
	for _, m := range r.modules {
		if m.Failed() {
			return true
		}
	}
	return false
}

// FailedChecks returns failed checks across all modules in deterministic order.
func (r Result) FailedChecks() []CheckResult {
	var out []CheckResult
	for _, m := range r.modules {
		for _, c := range m.Checks() {
			if c.Status() == StatusFailed {
				out = append(out, c)
			}
		}
	}
	return out
}

// ModuleResult describes verification checks for one module.
type ModuleResult struct {
	// module is the planned module name.
	module manifest.ModuleName

	// checks contains check results in execution order.
	checks []CheckResult
}

// Module returns the planned module name.
func (m ModuleResult) Module() manifest.ModuleName { return m.module }

// Checks returns detached check results.
func (m ModuleResult) Checks() []CheckResult {
	out := make([]CheckResult, len(m.checks))
	copy(out, m.checks)
	return out
}

// Failed reports whether this module has at least one failed check.
func (m ModuleResult) Failed() bool {
	for _, c := range m.checks {
		if c.status == StatusFailed {
			return true
		}
	}
	return false
}

// CheckResult describes one verification check outcome.
type CheckResult struct {
	// name identifies the check.
	name CheckName

	// status is the check outcome.
	status Status

	// severity describes whether the outcome blocks publication.
	severity Severity

	// message is a human-readable diagnostic.
	message string

	// path optionally identifies the checked file or directory.
	path string
}

// NewCheckResult creates one verification check result.
func NewCheckResult(name CheckName, status Status, severity Severity, message string) CheckResult {
	return CheckResult{name: name, status: status, severity: severity, message: message}
}

// Name returns the check name.
func (c CheckResult) Name() CheckName { return c.name }

// Status returns the check outcome.
func (c CheckResult) Status() Status { return c.status }

// Severity returns the check severity.
func (c CheckResult) Severity() Severity { return c.severity }

// Message returns the check diagnostic.
func (c CheckResult) Message() string { return c.message }

// Path returns the checked path when present.
func (c CheckResult) Path() string { return c.path }

// withPath attaches a filesystem path to a check result.
func (c CheckResult) withPath(path string) CheckResult { c.path = path; return c }
