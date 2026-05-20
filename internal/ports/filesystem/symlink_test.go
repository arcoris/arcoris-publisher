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

package filesystem

import "testing"

func TestSymlinkPolicyString(t *testing.T) {
	tests := []struct {
		policy SymlinkPolicy
		want   string
	}{
		{policy: SymlinkReject, want: "reject"},
		{policy: SymlinkPreserve, want: "preserve"},
		{policy: SymlinkFollow, want: "follow"},
	}

	for _, tt := range tests {
		if got := tt.policy.String(); got != tt.want {
			t.Fatalf("%#v.String() = %q, want %q", tt.policy, got, tt.want)
		}
	}
}

func TestSymlinkPolicyValid(t *testing.T) {
	tests := []struct {
		name   string
		policy SymlinkPolicy
		want   bool
	}{
		{name: "reject", policy: SymlinkReject, want: true},
		{name: "preserve", policy: SymlinkPreserve, want: true},
		{name: "follow", policy: SymlinkFollow, want: true},
		{name: "empty", policy: SymlinkPolicy(""), want: false},
		{name: "unknown", policy: SymlinkPolicy("copy"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.policy.Valid(); got != tt.want {
				t.Fatalf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
