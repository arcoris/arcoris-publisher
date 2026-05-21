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

import "testing"

func TestErrorStringUsesMessageFirst(t *testing.T) {
	err := &Error{Message: "explicit message"}
	if got := err.Error(); got != "explicit message" {
		t.Fatalf("Error() = %q, want explicit message", got)
	}
}

func TestErrorStringFallbacks(t *testing.T) {
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
