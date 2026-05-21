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

package verify

import "fmt"

// ErrorCode identifies an infrastructure-level verification error.
type ErrorCode string

const (
	// CodeInvalidRequest indicates malformed verification input.
	CodeInvalidRequest ErrorCode = "invalid_request"

	// CodeDependencyMissing indicates that a required infrastructure port is missing.
	CodeDependencyMissing ErrorCode = "dependency_missing"
)

// Error describes why verification could not run.
//
// Verification failures that are discovered by checks are reported in Result
// instead of this error type.
type Error struct {
	// Code is the stable machine-readable error reason.
	Code ErrorCode

	// Message is the human-readable diagnostic text.
	Message string
}

// Error returns a compact verification infrastructure error string.
func (e *Error) Error() string {
	if e == nil {
		return "verify error"
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}
