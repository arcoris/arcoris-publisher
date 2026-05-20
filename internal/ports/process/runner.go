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

// Package process defines the port for executing external operating-system
// processes such as git, go, or user-defined smoke-test commands.
//
// This package is the lowest-level execution boundary. Higher-level ports such
// as git and gotoolchain can be implemented on top of Runner, but workflow code
// should prefer those domain-specific ports when it needs Git or Go behavior.
//
// Implementations are responsible for process lifecycle, cancellation, timeout
// handling, output capture limits, and redaction. The contract intentionally
// keeps command rendering out of workflow code so secrets and platform quirks
// have one implementation boundary.
package process

import "context"

// Runner executes external processes described by Spec.
//
// Implementations MUST honor context cancellation, MUST redact configured
// sensitive values from rendered errors and logs, and MUST return structured
// Result values when a process starts successfully.
type Runner interface {
	// Run starts the requested process and waits for completion.
	//
	// A process that exits with an unaccepted exit code should return both its
	// Result and an error so callers can inspect captured output. If the process
	// cannot start, adapters should return a zero Result plus a structured error.
	Run(ctx context.Context, spec Spec) (Result, error)
}
