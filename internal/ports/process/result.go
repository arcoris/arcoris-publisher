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

import "time"

// Result contains structured information about a completed process invocation.
type Result struct {
	// Name is the executable that was invoked.
	Name string
	// Args are process arguments excluding Name.
	Args []string
	// Dir is the working directory used for the invocation.
	Dir string
	// ExitCode is the process exit code, when available.
	ExitCode int
	// Stdout contains captured standard output, if requested.
	Stdout []byte
	// Stderr contains captured standard error, if requested.
	Stderr []byte
	// StartedAt is the invocation start time.
	StartedAt time.Time
	// FinishedAt is the invocation finish time.
	FinishedAt time.Time
	// Duration is the elapsed process duration.
	Duration time.Duration
}

// Succeeded reports whether the result exit code is accepted by the provided list.
func (r Result) Succeeded(allowed []int) bool {
	return IsAllowedExitCode(r.ExitCode, allowed)
}
