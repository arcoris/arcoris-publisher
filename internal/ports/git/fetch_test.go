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

package git

import "testing"

func TestFetchTagsModeString(t *testing.T) {
	tests := []struct {
		mode FetchTagsMode
		want string
	}{
		{mode: FetchTagsDefault, want: "default"},
		{mode: FetchTagsAll, want: "all"},
		{mode: FetchTagsNone, want: "none"},
	}

	for _, tt := range tests {
		if got := tt.mode.String(); got != tt.want {
			t.Fatalf("%#v.String() = %q, want %q", tt.mode, got, tt.want)
		}
	}
}

func TestFetchTagsModeValid(t *testing.T) {
	tests := []struct {
		name string
		mode FetchTagsMode
		want bool
	}{
		{name: "default", mode: FetchTagsDefault, want: true},
		{name: "all", mode: FetchTagsAll, want: true},
		{name: "none", mode: FetchTagsNone, want: true},
		{name: "empty", mode: FetchTagsMode(""), want: false},
		{name: "unknown", mode: FetchTagsMode("some"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mode.Valid(); got != tt.want {
				t.Fatalf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
