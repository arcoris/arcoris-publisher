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

package report

// Status is a stable report-level state string used in text and JSON reports.
type Status string

const (
	// StatusEmpty means no meaningful stage data is present.
	StatusEmpty Status = "empty"

	// StatusPresent means a stage has meaningful data but no pass/fail state.
	StatusPresent Status = "present"

	// StatusPartial means a workflow reached some stages but not verification.
	StatusPartial Status = "partial"

	// StatusVerified means verification passed and publication did not run.
	StatusVerified Status = "verified"

	// StatusVerificationFailed means verification ran and found failed checks.
	StatusVerificationFailed Status = "verification_failed"

	// StatusPublished means at least one module was published.
	StatusPublished Status = "published"

	// StatusSkipped means all publishable modules were skipped.
	StatusSkipped Status = "skipped"

	// StatusPassed means checks completed without failures.
	StatusPassed Status = "passed"

	// StatusFailed means at least one check failed.
	StatusFailed Status = "failed"

	// StatusPending means a module has no final published/skipped state.
	StatusPending Status = "pending"

	// StatusCommitted means a publish transaction completed final refs.
	StatusCommitted Status = "committed"

	// StatusRolledBack means failed publish side effects were rolled back.
	StatusRolledBack Status = "rolled_back"

	// StatusRollbackFailed means rollback could not safely finish.
	StatusRollbackFailed Status = "rollback_failed"
)

// String returns the stable textual representation of s.
func (s Status) String() string { return string(s) }
