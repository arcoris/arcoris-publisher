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

func TestCodeString(t *testing.T) {
	tests := []struct {
		code Code
		want string
	}{
		{code: Code(""), want: ""},
		{code: Code("git_failed"), want: "git_failed"},
	}

	for _, tt := range tests {
		if got := tt.code.String(); got != tt.want {
			t.Fatalf("%#v.String() = %q, want %q", tt.code, got, tt.want)
		}
	}
}
