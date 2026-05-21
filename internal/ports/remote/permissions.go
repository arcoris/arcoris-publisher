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

// RepositoryPermissions describes effective provider permissions for a repository.
//
// Permissions are effective for the current principal, not repository policy in
// general. Providers may compute them from role names, team grants, installation
// scopes, or token capabilities.
type RepositoryPermissions struct {
	// CanRead permits reading repository metadata and contents.
	CanRead bool
	// CanWrite permits writing refs, pull requests, or releases.
	CanWrite bool
	// CanAdmin permits administrative repository actions.
	CanAdmin bool
}

// Allows reports whether the permission set satisfies the requested access level.
//
// Capability is hierarchical: admin implies write and read, and write implies
// read. Unknown access levels are rejected even if all booleans are true.
func (p RepositoryPermissions) Allows(level AccessLevel) bool {
	switch level {
	case AccessRead:
		return p.canRead()
	case AccessWrite:
		return p.canWrite()
	case AccessAdmin:
		return p.canAdmin()
	default:
		return false
	}
}

func (p RepositoryPermissions) canRead() bool {
	return p.CanRead || p.canWrite()
}

func (p RepositoryPermissions) canWrite() bool {
	return p.CanWrite || p.canAdmin()
}

func (p RepositoryPermissions) canAdmin() bool {
	return p.CanAdmin
}
