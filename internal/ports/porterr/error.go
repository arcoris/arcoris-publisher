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

package porterr

// Error is the common structured error type returned by infrastructure ports.
//
// The type is intentionally small and stable. Port implementations may attach a
// wrapped Cause, but workflow code should normally branch on Kind and Code
// rather than parsing error strings.
//
// Message is for humans; Kind and Code are for control flow. Details may be
// logged and displayed, so adapters must keep secrets in Cause or local
// implementation state and redact them before constructing Error values.
type Error struct {
	// Kind identifies the infrastructure boundary that produced the error.
	Kind Kind
	// Code is the stable machine-readable error code.
	Code Code
	// Message is a concise human-readable error summary.
	Message string
	// Temporary reports whether retrying the same operation may succeed.
	Temporary bool
	// Details carries non-secret structured diagnostic context.
	Details Details
	// Cause is the wrapped implementation-specific error.
	Cause error
}

// New creates a structured infrastructure error.
//
// New leaves Details empty and Temporary false. Call WithDetails and
// WithTemporary to enrich the error while preserving copy-on-write behavior.
func New(kind Kind, code Code, message string, cause error) *Error {
	return &Error{
		Kind:    kind,
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// WithDetails returns a copy of the error with detached details attached.
//
// The receiver is not mutated. This allows adapters to define base errors and
// attach operation-specific context without sharing maps between errors.
func (e *Error) WithDetails(details Details) *Error {
	if e == nil {
		return nil
	}
	copy := *e
	copy.Details = details.Clone()
	return &copy
}

// WithTemporary returns a copy of the error with the temporary flag updated.
//
// Temporary means retrying the same logical operation may succeed, for example
// after rate limiting or transient transport failure.
func (e *Error) WithTemporary(temporary bool) *Error {
	if e == nil {
		return nil
	}
	copy := *e
	copy.Temporary = temporary
	return &copy
}

// Error returns the best available human-readable error text.
//
// It prefers the explicit Message and falls back to Kind and Code only to keep
// zero-message errors useful in logs.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Message != "" {
		return e.Message
	}
	if e.Kind != "" && e.Code != "" {
		return e.Kind.String() + ": " + e.Code.String()
	}
	if e.Kind != "" {
		return e.Kind.String()
	}
	if e.Code != "" {
		return e.Code.String()
	}
	return "port error"
}

// Unwrap returns the wrapped implementation-specific error.
//
// This makes errors.Is and errors.As work for adapter-specific causes while
// still exposing a stable porterr.Error to workflow code.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}
