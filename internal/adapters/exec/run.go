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
	"context"
	"time"

	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

// Run executes spec and returns a structured process result.
//
// The orchestration stays intentionally linear: create the observable result,
// prepare the OS command, run it once, collect redacted output, then classify
// the lifecycle outcome. The detailed mechanics live in smaller files so this
// method reads like the process adapter contract.
func (r *Runner) Run(ctx context.Context, spec processport.Spec) (processport.Result, error) {
	startedAt := time.Now()
	result := newResult(spec, startedAt)
	redactor := NewRedactor(spec.SensitiveValues...)

	if spec.Name == "" {
		finishResult(&result, startedAt, nil, redactor)
		return result, emptyNameError(spec, redactor)
	}

	runCtx, cancel := contextWithOptionalTimeout(ctx, spec.Timeout)
	defer cancel()

	cmd, capture := r.prepareCommand(runCtx, spec)
	err := cmd.Run()
	finishResult(&result, startedAt, capture, redactor)

	if err == nil {
		return finishSuccessfulProcess(result, spec, redactor)
	}
	return finishFailedProcess(result, err, ctx, runCtx, spec, redactor)
}
