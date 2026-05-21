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

import (
	"errors"
	"strings"
	"testing"
)

func TestErrorFormatsAndUnwraps(t *testing.T) {
	cause := errors.New("cause")
	err := newError(CodeReadFailed, "arcpub.yaml", cause, "failed %s", "read")
	if !strings.Contains(err.Error(), "arcpub.yaml") {
		t.Fatalf("error message lacks path: %s", err)
	}
	if !errors.Is(err, cause) {
		t.Fatal("expected cause to unwrap")
	}
}

func TestErrorWithoutPath(t *testing.T) {
	err := (&Error{Code: CodeConfigNotFound, Message: "missing"}).Error()
	if !strings.Contains(err, "missing") {
		t.Fatalf("unexpected error: %s", err)
	}
}

func TestNilErrorString(t *testing.T) {
	var err *Error
	if err.Error() != "<nil>" {
		t.Fatal("unexpected nil error string")
	}
}

func TestErrorHelpersAssignCodes(t *testing.T) {
	cause := errors.New("cause")
	cases := []struct {
		name string
		err  error
		code Code
	}{
		{name: "read", err: readError("path", cause), code: CodeReadFailed},
		{name: "inspect", err: inspectCandidateError("path", cause), code: CodeReadFailed},
		{name: "not found", err: configNotFoundError("path"), code: CodeConfigNotFound},
		{name: "format", err: unsupportedFormatError("path"), code: CodeUnsupportedFormat},
		{name: "decode", err: decodeError("path", cause, "module"), code: CodeDecodeFailed},
		{
			name: "invalid",
			err:  invalidManifestError("path", cause, "module"),
			code: CodeInvalidManifest,
		},
		{name: "root", err: stagingRootError("path", cause), code: CodeUnsafePath},
		{name: "abs", err: pathAbsError("path", cause), code: CodeUnsafePath},
		{name: "source", err: sourceDirEscapeError("path", cause), code: CodeUnsafePath},
		{
			name: "manifest",
			err:  moduleManifestEscapeError("path", cause),
			code: CodeUnsafePath,
		},
		{name: "resolve", err: resolveError("path", cause), code: CodeResolveFailed},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			configErr, ok := tc.err.(*Error)
			if !ok {
				t.Fatalf("expected *Error, got %T", tc.err)
			}
			if configErr.Code != tc.code {
				t.Fatalf("got %s, want %s", configErr.Code, tc.code)
			}
		})
	}
}
