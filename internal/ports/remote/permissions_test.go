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

func TestRepositoryPermissionsAllows(t *testing.T) {
	tests := []struct {
		name        string
		permission  RepositoryPermissions
		access      AccessLevel
		wantAllowed bool
	}{
		{name: "read allows read", permission: RepositoryPermissions{CanRead: true}, access: AccessRead, wantAllowed: true},
		{name: "read rejects write", permission: RepositoryPermissions{CanRead: true}, access: AccessWrite, wantAllowed: false},
		{name: "read rejects admin", permission: RepositoryPermissions{CanRead: true}, access: AccessAdmin, wantAllowed: false},
		{name: "write allows read", permission: RepositoryPermissions{CanWrite: true}, access: AccessRead, wantAllowed: true},
		{name: "write allows write", permission: RepositoryPermissions{CanWrite: true}, access: AccessWrite, wantAllowed: true},
		{name: "write rejects admin", permission: RepositoryPermissions{CanWrite: true}, access: AccessAdmin, wantAllowed: false},
		{name: "admin allows read", permission: RepositoryPermissions{CanAdmin: true}, access: AccessRead, wantAllowed: true},
		{name: "admin allows write", permission: RepositoryPermissions{CanAdmin: true}, access: AccessWrite, wantAllowed: true},
		{name: "admin allows admin", permission: RepositoryPermissions{CanAdmin: true}, access: AccessAdmin, wantAllowed: true},
		{name: "empty rejects read", permission: RepositoryPermissions{}, access: AccessRead, wantAllowed: false},
		{name: "unknown rejected", permission: RepositoryPermissions{CanAdmin: true}, access: AccessLevel("unknown"), wantAllowed: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.permission.Allows(tt.access); got != tt.wantAllowed {
				t.Fatalf("Allows(%q) = %v, want %v", tt.access, got, tt.wantAllowed)
			}
		})
	}
}
