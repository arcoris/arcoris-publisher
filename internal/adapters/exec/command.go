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
	"os"
	osexec "os/exec"

	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

// outputCapture groups the optional stdout/stderr buffers attached to a command.
//
// Keeping the buffers in one value makes it explicit that command preparation
// owns capture wiring, while result finalization owns reading and redaction.
type outputCapture struct {
	stdout *limitedBuffer
	stderr *limitedBuffer
}

// prepareCommand converts a port process spec into an os/exec command.
//
// The method is intentionally limited to mechanical command construction:
// environment merging, working directory, stdin, and bounded capture buffers.
// It does not run the process or interpret the result.
func (r *Runner) prepareCommand(ctx context.Context, spec processport.Spec) (*osexec.Cmd, *outputCapture) {
	cmd := osexec.CommandContext(ctx, spec.Name, spec.Args...)
	cmd.Dir = spec.Dir
	cmd.Env = MergeEnv(os.Environ(), MergeEnv(r.baseEnv, spec.Env))
	cmd.Stdin = spec.Stdin

	capture := &outputCapture{
		stdout: newLimitedBuffer(outputLimit(spec.MaxStdoutBytes)),
		stderr: newLimitedBuffer(outputLimit(spec.MaxStderrBytes)),
	}
	if spec.CaptureStdout {
		cmd.Stdout = capture.stdout
	}
	if spec.CaptureStderr {
		cmd.Stderr = capture.stderr
	}
	return cmd, capture
}

// outputLimit applies the adapter default when the caller leaves a capture
// limit unset.
func outputLimit(limit int64) int64 {
	if limit == 0 {
		return defaultOutputLimit
	}
	return limit
}
