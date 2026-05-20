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

// AccessLevel identifies a required provider-level repository access level.
type AccessLevel string

const (
	// AccessRead requires repository read access.
	AccessRead AccessLevel = "read"
	// AccessWrite requires repository write access.
	AccessWrite AccessLevel = "write"
	// AccessAdmin requires repository administrative access.
	AccessAdmin AccessLevel = "admin"
)

// String returns the stable string representation of the access level.
func (a AccessLevel) String() string {
	return string(a)
}

// Valid reports whether the level is one of the supported access levels.
func (a AccessLevel) Valid() bool {
	switch a {
	case AccessRead, AccessWrite, AccessAdmin:
		return true
	default:
		return false
	}
}

// RepositoryPermissions describes effective provider permissions for a repository.
type RepositoryPermissions struct {
	// CanRead permits reading repository metadata and contents.
	CanRead bool
	// CanWrite permits writing refs, pull requests, or releases.
	CanWrite bool
	// CanAdmin permits administrative repository actions.
	CanAdmin bool
}

// Allows reports whether the permission set satisfies the requested access level.
func (p RepositoryPermissions) Allows(level AccessLevel) bool {
	switch level {
	case AccessRead:
		return p.CanRead || p.CanWrite || p.CanAdmin
	case AccessWrite:
		return p.CanWrite || p.CanAdmin
	case AccessAdmin:
		return p.CanAdmin
	default:
		return false
	}
}
