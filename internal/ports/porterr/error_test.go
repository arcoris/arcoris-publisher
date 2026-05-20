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

import (
	"errors"
	"testing"
)

func TestError_Unwrap(t *testing.T) {
	cause := errors.New("root cause")
	err := New(KindGit, Code("git_failed"), "git failed", cause)

	if !errors.Is(err, cause) {
		t.Fatalf("expected wrapped cause to be discoverable")
	}
}

func TestError_WithDetailsClonesInput(t *testing.T) {
	details := Details{"repo": "target"}
	err := New(KindGit, Code("git_failed"), "git failed", nil).WithDetails(details)
	details["repo"] = "mutated"

	if got := err.Details["repo"]; got != "target" {
		t.Fatalf("expected detached details, got %q", got)
	}
}

func TestError_ErrorFallbacks(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{name: "nil", err: nil, want: "<nil>"},
		{name: "kind and code", err: &Error{Kind: KindProcess, Code: Code("process_failed")}, want: "process: process_failed"},
		{name: "kind only", err: &Error{Kind: KindGit}, want: "git"},
		{name: "code only", err: &Error{Code: Code("remote_failed")}, want: "remote_failed"},
		{name: "empty", err: &Error{}, want: "port error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Fatalf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestError_WithTemporaryAndNilReceivers(t *testing.T) {
	if (*Error)(nil).WithTemporary(true) != nil {
		t.Fatalf("expected nil temporary receiver to return nil")
	}
	if (*Error)(nil).WithDetails(Details{"key": "value"}) != nil {
		t.Fatalf("expected nil details receiver to return nil")
	}
	err := New(KindProcess, Code("process_failed"), "failed", nil).WithTemporary(true)
	if !err.Temporary {
		t.Fatalf("expected temporary flag")
	}
}

func TestError_UnwrapNil(t *testing.T) {
	if (*Error)(nil).Unwrap() != nil {
		t.Fatalf("expected nil unwrap for nil receiver")
	}
}

func TestError_ErrorUsesMessage(t *testing.T) {
	err := &Error{Message: "explicit message"}
	if got := err.Error(); got != "explicit message" {
		t.Fatalf("unexpected explicit message %q", got)
	}
}
