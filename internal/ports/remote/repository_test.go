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

package remote

import "testing"

func TestRepositoryRefFullName(t *testing.T) {
	tests := []struct {
		name string
		ref  RepositoryRef
		want string
	}{
		{name: "owner", ref: RepositoryRef{Owner: "arcoris", Name: "foundation"}, want: "arcoris/foundation"},
		{name: "ownerless", ref: RepositoryRef{Name: "foundation"}, want: "foundation"},
		{name: "empty", ref: RepositoryRef{}, want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ref.FullName(); got != tt.want {
				t.Fatalf("FullName() = %q, want %q", got, tt.want)
			}
		})
	}
}
