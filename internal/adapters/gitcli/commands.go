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

package gitcli

import (
	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

// command builds the process spec for one Git invocation.
//
// The method detaches all slices because specs may be inspected by tests or
// reused by callers after the original option slices are mutated.
func (c *Client) command(repoDir string, args []string, sensitive []string, captureStdout, captureStderr bool) processport.Spec {
	return processport.Spec{
		Name:            c.gitBin,
		Args:            append([]string(nil), args...),
		Dir:             repoDir,
		Env:             append([]string(nil), c.env...),
		CaptureStdout:   captureStdout,
		CaptureStderr:   captureStderr,
		SensitiveValues: append([]string(nil), sensitive...),
	}
}

// stringsOf converts typed Git string values into ordinary CLI arguments.
func stringsOf[T ~string](values []T) []string {
	out := make([]string, len(values))
	for i, value := range values {
		out[i] = string(value)
	}
	return out
}
