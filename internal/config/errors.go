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

package config

import "fmt"

// Code classifies a configuration loading error.
type Code string

const (
	// CodeConfigNotFound indicates that automatic discovery could not find a
	// supported top-level arcpub manifest file.
	CodeConfigNotFound Code = "config_not_found"
	// CodeReadFailed indicates that a manifest file could not be read.
	CodeReadFailed Code = "read_failed"
	// CodeUnsupportedFormat indicates that a file extension does not map to a
	// supported configuration format.
	CodeUnsupportedFormat Code = "unsupported_format"
	// CodeDecodeFailed indicates that file bytes could not be strictly decoded.
	CodeDecodeFailed Code = "decode_failed"
	// CodeInvalidManifest indicates that decoded manifest data violates manifest
	// domain validation rules.
	CodeInvalidManifest Code = "invalid_manifest"
	// CodeUnsafePath indicates that a computed path escapes the expected loading
	// boundary or otherwise violates config-level path safety rules.
	CodeUnsafePath Code = "unsafe_path"
	// CodeResolveFailed indicates that validated staging and module manifests
	// could not be bound into an effective publication set.
	CodeResolveFailed Code = "resolve_failed"
)

// Error is a typed configuration loading error.
type Error struct {
	Code    Code
	Path    string
	Message string
	Cause   error
}

// Error returns a stable human-readable error message.
func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Path == "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("%s: %s: %s", e.Code, e.Path, e.Message)
}

// Unwrap returns the underlying cause.
func (e *Error) Unwrap() error {
	return e.Cause
}

// newError creates a typed config error.
func newError(code Code, path string, cause error, format string, args ...any) error {
	return &Error{
		Code:    code,
		Path:    path,
		Message: fmt.Sprintf(format, args...),
		Cause:   cause,
	}
}

// readError wraps filesystem read failures with the config error code expected
// by callers.
func readError(path string, cause error) error {
	return newError(CodeReadFailed, path, cause, "could not read manifest")
}

// inspectCandidateError wraps failures from manifest discovery probes.
func inspectCandidateError(path string, cause error) error {
	return newError(
		CodeReadFailed,
		path,
		cause,
		"could not inspect candidate manifest",
	)
}

// configNotFoundError reports that discovery exhausted all supported names.
func configNotFoundError(path string) error {
	return newError(
		CodeConfigNotFound,
		path,
		nil,
		"no arcpub manifest found",
	)
}

// unsupportedFormatError reports paths whose extension cannot select a decoder.
func unsupportedFormatError(path string) error {
	return newError(
		CodeUnsupportedFormat,
		path,
		nil,
		"unsupported manifest file extension",
	)
}

// decodeError records syntax-level decode failures for a specific manifest kind.
func decodeError(path string, cause error, manifestKind string) error {
	return newError(
		CodeDecodeFailed,
		path,
		cause,
		"could not decode %s manifest",
		manifestKind,
	)
}

// invalidManifestError records domain validation failures after decode succeeds.
func invalidManifestError(path string, cause error, manifestKind string) error {
	return newError(
		CodeInvalidManifest,
		path,
		cause,
		"%s manifest is invalid",
		manifestKind,
	)
}

// stagingRootError records failures while resolving source.stagingRoot.
func stagingRootError(path string, cause error) error {
	return newError(
		CodeUnsafePath,
		path,
		cause,
		"could not resolve staging root",
	)
}

// pathAbsError records path normalization failures.
func pathAbsError(path string, cause error) error {
	return newError(
		CodeUnsafePath,
		path,
		cause,
		"could not make path absolute",
	)
}

// sourceDirEscapeError reports a module.sourceDir outside the staging root.
func sourceDirEscapeError(path string, cause error) error {
	return newError(
		CodeUnsafePath,
		path,
		cause,
		"module sourceDir escapes staging root",
	)
}

// moduleManifestEscapeError reports a module manifest path outside sourceDir.
func moduleManifestEscapeError(path string, cause error) error {
	return newError(
		CodeUnsafePath,
		path,
		cause,
		"module manifest path escapes module sourceDir",
	)
}

// resolveError records failures from binding staging and module manifests.
func resolveError(path string, cause error) error {
	return newError(
		CodeResolveFailed,
		path,
		cause,
		"could not resolve publication set",
	)
}
