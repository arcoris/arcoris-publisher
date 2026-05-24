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

package report

import "fmt"

// ErrorCode identifies report rendering failures.
type ErrorCode string

const (
	// CodeUnsupportedFormat indicates an unknown output format.
	CodeUnsupportedFormat ErrorCode = "unsupported_format"

	// CodeRenderFailed indicates a report could not be rendered.
	CodeRenderFailed ErrorCode = "render_failed"

	// CodeWriteFailed indicates a rendered report could not be written.
	CodeWriteFailed ErrorCode = "write_failed"
)

// Error describes a report-layer failure.
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
