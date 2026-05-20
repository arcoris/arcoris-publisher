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

package exec

import (
	"time"

	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

// newResult copies command identity into the result before execution starts.
//
// Args are detached so later caller mutations cannot rewrite diagnostics for an
// already-finished process.
func newResult(spec processport.Spec, startedAt time.Time) processport.Result {
	return processport.Result{Name: spec.Name, Args: append([]string(nil), spec.Args...), Dir: spec.Dir, StartedAt: startedAt}
}

// finishResult records timing and copies captured output into the result.
//
// Output is redacted at the adapter boundary because result values often travel
// into higher-level logs and structured errors.
func finishResult(result *processport.Result, startedAt time.Time, capture *outputCapture, redactor Redactor) {
	finishedAt := time.Now()
	result.FinishedAt = finishedAt
	result.Duration = finishedAt.Sub(startedAt)
	if capture == nil {
		return
	}
	if capture.stdout != nil {
		result.Stdout = redactor.RedactBytes(capture.stdout.Bytes())
	}
	if capture.stderr != nil {
		result.Stderr = redactor.RedactBytes(capture.stderr.Bytes())
	}
}
