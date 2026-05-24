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

package cli

import "fmt"

// ErrorCode identifies command-line failures.
type ErrorCode string

const (
	// CodeInvalidCommand indicates an unknown or missing command.
	CodeInvalidCommand ErrorCode = "invalid_command"

	// CodeInvalidFlags indicates malformed command flags.
	CodeInvalidFlags ErrorCode = "invalid_flags"

	// CodeInvalidVersion indicates an invalid publication version flag.
	CodeInvalidVersion ErrorCode = "invalid_version"

	// CodeMissingApplication indicates that a command requiring app use cases was
	// invoked without an application dependency.
	CodeMissingApplication ErrorCode = "missing_application"

	// CodeUseCaseFailed indicates an app use case returned an error.
	CodeUseCaseFailed ErrorCode = "use_case_failed"

	// CodeReportFailed indicates report rendering failed.
	CodeReportFailed ErrorCode = "report_failed"

	// CodeVerificationFailed indicates verification completed with failed checks.
	CodeVerificationFailed ErrorCode = "verification_failed"
)

var errVerificationFailed = &Error{
	Code:    CodeVerificationFailed,
	Message: "verification failed",
}

// Error describes a CLI-layer failure.
type Error struct {
	Code    ErrorCode
	Message string
	Cause   error
}

// Error returns a human-readable error string.
func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

func (e *Error) isUsage() bool {
	switch e.Code {
	case CodeInvalidCommand, CodeInvalidFlags, CodeInvalidVersion:
		return true
	default:
		return false
	}
}
