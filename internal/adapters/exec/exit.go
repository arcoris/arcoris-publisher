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
	"errors"
	"fmt"
	osexec "os/exec"

	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

// finishSuccessfulProcess handles the special case where the OS command
// completed cleanly but the caller explicitly disallowed exit code 0.
func finishSuccessfulProcess(result processport.Result, spec processport.Spec, redactor Redactor) (processport.Result, error) {
	result.ExitCode = 0
	if processport.IsAllowedExitCode(0, spec.AllowedExitCodes) {
		return result, nil
	}
	msg := fmt.Sprintf("process %q exited with code 0", redactor.RedactString(spec.Name))
	return result, processError(processport.CodeFailed, msg, nil, redactedCommandDetails(spec, redactor, 0))
}

// finishFailedProcess classifies all non-successful command outcomes.
//
// The order matters: context timeout/cancel errors are more informative than a
// generic process failure, and ExitError carries the real numeric exit code.
func finishFailedProcess(result processport.Result, err error, parentCtx context.Context, runCtx context.Context, spec processport.Spec, redactor Redactor) (processport.Result, error) {
	result.ExitCode = -1
	if errors.Is(runCtx.Err(), context.DeadlineExceeded) {
		return result, timedOutError(spec, redactor, runCtx.Err(), result.ExitCode)
	}
	if errors.Is(runCtx.Err(), context.Canceled) || errors.Is(parentCtx.Err(), context.Canceled) {
		return result, cancelledError(spec, redactor, runCtx.Err(), result.ExitCode)
	}

	var exitErr *osexec.ExitError
	if errors.As(err, &exitErr) {
		exitCode := exitErr.ExitCode()
		result.ExitCode = exitCode
		if processport.IsAllowedExitCode(exitCode, spec.AllowedExitCodes) {
			return result, nil
		}
		return result, rejectedExitCodeError(spec, redactor, err, exitCode)
	}

	result.ExitCode = 0
	if errors.Is(err, osexec.ErrNotFound) {
		return result, executableNotFoundError(spec, redactor, err)
	}
	return result, startFailedError(spec, redactor, err)
}

// emptyNameError reports validation failure before os/exec is asked to start.
func emptyNameError(spec processport.Spec, redactor Redactor) error {
	return processError(processport.CodeStartFailed, "process name is empty", nil, redactedCommandDetails(spec, redactor, 0))
}

// rejectedExitCodeError reports an exit code that is outside Spec.AllowedExitCodes.
func rejectedExitCodeError(spec processport.Spec, redactor Redactor, cause error, exitCode int) error {
	msg := fmt.Sprintf("process %q exited with code %d", redactor.RedactString(spec.Name), exitCode)
	return processError(processport.CodeFailed, msg, cause, redactedCommandDetails(spec, redactor, exitCode))
}

// timedOutError reports that the derived timeout context expired.
func timedOutError(spec processport.Spec, redactor Redactor, cause error, exitCode int) error {
	msg := fmt.Sprintf("process %q timed out", redactor.RedactString(spec.Name))
	return processError(processport.CodeTimedOut, msg, cause, redactedCommandDetails(spec, redactor, exitCode))
}

// cancelledError reports explicit cancellation from the caller or parent context.
func cancelledError(spec processport.Spec, redactor Redactor, cause error, exitCode int) error {
	msg := fmt.Sprintf("process %q was cancelled", redactor.RedactString(spec.Name))
	return processError(processport.CodeCancelled, msg, cause, redactedCommandDetails(spec, redactor, exitCode))
}

// executableNotFoundError preserves the process_not_found code for missing binaries.
func executableNotFoundError(spec processport.Spec, redactor Redactor, cause error) error {
	msg := fmt.Sprintf("process executable %q was not found", redactor.RedactString(spec.Name))
	return processError(processport.CodeNotFound, msg, cause, redactedCommandDetails(spec, redactor, 0))
}

// startFailedError reports startup failures that are not missing-executable cases.
func startFailedError(spec processport.Spec, redactor Redactor, cause error) error {
	msg := fmt.Sprintf("process %q could not be started", redactor.RedactString(spec.Name))
	return processError(processport.CodeStartFailed, msg, cause, redactedCommandDetails(spec, redactor, 0))
}

// redactedCommandDetails builds safe diagnostics from the original spec.
func redactedCommandDetails(spec processport.Spec, redactor Redactor, exitCode int) map[string]string {
	return commandDetails(redactor.RedactString(spec.Name), redactor.RedactSlice(spec.Args), redactor.RedactString(spec.Dir), exitCode)
}
