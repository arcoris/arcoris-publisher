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

package process

import (
	"io"
	"time"
)

// Spec describes an external process invocation.
//
// Spec is immutable by convention: helper methods return modified copies, and
// adapters should not retain or mutate caller-owned slices after Run returns.
type Spec struct {
	// Name is the executable name or absolute executable path.
	Name string
	// Args are process arguments excluding Name.
	Args []string
	// Dir is the working directory. Empty means the implementation default.
	Dir string
	// Env contains environment assignments in KEY=VALUE form.
	//
	// Env augments or overrides the adapter's base environment; it is not
	// required to be a complete environment snapshot.
	Env []string
	// Stdin is optional process standard input.
	Stdin io.Reader
	// Timeout is an optional process-level timeout. Zero means no explicit timeout.
	//
	// Context cancellation still applies even when Timeout is zero.
	Timeout time.Duration
	// CaptureStdout asks the implementation to capture stdout into Result.Stdout.
	CaptureStdout bool
	// CaptureStderr asks the implementation to capture stderr into Result.Stderr.
	CaptureStderr bool
	// MaxStdoutBytes optionally limits captured stdout size. Zero means implementation default.
	//
	// Adapters should truncate deterministically and make truncation visible in
	// diagnostics when possible.
	MaxStdoutBytes int64
	// MaxStderrBytes optionally limits captured stderr size. Zero means implementation default.
	MaxStderrBytes int64
	// AllowedExitCodes defines process exit codes accepted as successful completion.
	// If empty, only exit code 0 is considered successful.
	AllowedExitCodes []int
	// SensitiveValues are raw values that MUST be redacted from rendered commands,
	// logs, errors, stdout, and stderr by concrete implementations.
	SensitiveValues []string
}
