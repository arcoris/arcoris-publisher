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

// Package exec implements the process port by using the operating-system
// process execution facilities from the Go standard library.
//
// The adapter is deliberately narrow: it starts commands, captures bounded
// output, maps process lifecycle failures to process-port error codes, and
// redacts configured sensitive values from returned diagnostics.
package exec

import (
	"context"
	"errors"
	"fmt"
	"os"
	osexec "os/exec"
	"time"

	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

const defaultOutputLimit int64 = 8 << 20 // 8 MiB

// Runner executes external processes through os/exec.
//
// Runner is safe to reuse across calls. Its base environment is copied during
// construction, and each Run invocation builds a fresh process specification.
type Runner struct {
	baseEnv []string
}

// Options configures Runner.
type Options struct {
	// Env is a base environment overlaid on top of os.Environ for every command.
	//
	// Per-command Spec.Env values are applied after Options.Env, so a workflow can
	// override adapter defaults for one invocation without mutating the runner.
	Env []string
}

// New creates a process runner.
func New(opts Options) *Runner {
	return &Runner{baseEnv: append([]string(nil), opts.Env...)}
}

// Run executes spec and returns a structured process result.
//
// If a process starts and exits with a rejected status, Run returns both the
// populated Result and a process_failed error. Startup failures return a Result
// with command identity and timing data but no synthetic non-zero exit code.
func (r *Runner) Run(ctx context.Context, spec processport.Spec) (processport.Result, error) {
	startedAt := time.Now()
	result := processport.Result{Name: spec.Name, Args: append([]string(nil), spec.Args...), Dir: spec.Dir, StartedAt: startedAt}
	redactor := NewRedactor(spec.SensitiveValues...)

	if spec.Name == "" {
		result.FinishedAt = time.Now()
		result.Duration = result.FinishedAt.Sub(startedAt)
		return result, processError(processport.CodeStartFailed, "process name is empty", nil, commandDetails("", nil, spec.Dir, 0))
	}

	runCtx := ctx
	cancel := func() {}
	if spec.Timeout > 0 {
		runCtx, cancel = context.WithTimeout(ctx, spec.Timeout)
	}
	defer cancel()

	cmd := osexec.CommandContext(runCtx, spec.Name, spec.Args...)
	cmd.Dir = spec.Dir
	cmd.Env = MergeEnv(os.Environ(), MergeEnv(r.baseEnv, spec.Env))
	cmd.Stdin = spec.Stdin

	stdoutLimit := spec.MaxStdoutBytes
	if stdoutLimit == 0 {
		stdoutLimit = defaultOutputLimit
	}
	stderrLimit := spec.MaxStderrBytes
	if stderrLimit == 0 {
		stderrLimit = defaultOutputLimit
	}
	stdout := newLimitedBuffer(stdoutLimit)
	stderr := newLimitedBuffer(stderrLimit)
	if spec.CaptureStdout {
		cmd.Stdout = stdout
	}
	if spec.CaptureStderr {
		cmd.Stderr = stderr
	}

	err := cmd.Run()
	finishedAt := time.Now()
	result.FinishedAt = finishedAt
	result.Duration = finishedAt.Sub(startedAt)
	if spec.CaptureStdout {
		result.Stdout = redactor.RedactBytes(stdout.Bytes())
	}
	if spec.CaptureStderr {
		result.Stderr = redactor.RedactBytes(stderr.Bytes())
	}

	if err == nil {
		result.ExitCode = 0
		if !processport.IsAllowedExitCode(0, spec.AllowedExitCodes) {
			msg := fmt.Sprintf("process %q exited with code 0", redactor.RedactString(spec.Name))
			return result, processError(processport.CodeFailed, msg, nil, commandDetails(redactor.RedactString(spec.Name), redactor.RedactSlice(spec.Args), redactor.RedactString(spec.Dir), 0))
		}
		return result, nil
	}

	exitCode := -1
	result.ExitCode = exitCode
	if errors.Is(runCtx.Err(), context.DeadlineExceeded) {
		return result, processError(processport.CodeTimedOut, fmt.Sprintf("process %q timed out", redactor.RedactString(spec.Name)), runCtx.Err(), commandDetails(redactor.RedactString(spec.Name), redactor.RedactSlice(spec.Args), redactor.RedactString(spec.Dir), exitCode))
	}
	if errors.Is(runCtx.Err(), context.Canceled) || errors.Is(ctx.Err(), context.Canceled) {
		return result, processError(processport.CodeCancelled, fmt.Sprintf("process %q was cancelled", redactor.RedactString(spec.Name)), runCtx.Err(), commandDetails(redactor.RedactString(spec.Name), redactor.RedactSlice(spec.Args), redactor.RedactString(spec.Dir), exitCode))
	}

	var exitErr *osexec.ExitError
	if errors.As(err, &exitErr) {
		exitCode = exitErr.ExitCode()
		result.ExitCode = exitCode
		if processport.IsAllowedExitCode(exitCode, spec.AllowedExitCodes) {
			return result, nil
		}
		msg := fmt.Sprintf("process %q exited with code %d", redactor.RedactString(spec.Name), exitCode)
		return result, processError(processport.CodeFailed, msg, err, commandDetails(redactor.RedactString(spec.Name), redactor.RedactSlice(spec.Args), redactor.RedactString(spec.Dir), exitCode))
	}

	result.ExitCode = 0
	if errors.Is(err, osexec.ErrNotFound) {
		return result, processError(processport.CodeNotFound, fmt.Sprintf("process executable %q was not found", redactor.RedactString(spec.Name)), err, commandDetails(redactor.RedactString(spec.Name), redactor.RedactSlice(spec.Args), redactor.RedactString(spec.Dir), 0))
	}
	return result, processError(processport.CodeStartFailed, fmt.Sprintf("process %q could not be started", redactor.RedactString(spec.Name)), err, commandDetails(redactor.RedactString(spec.Name), redactor.RedactSlice(spec.Args), redactor.RedactString(spec.Dir), 0))
}

var _ processport.Runner = (*Runner)(nil)
