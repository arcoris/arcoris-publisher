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

// RepositoryRef identifies a repository hosted by a remote provider.
//
// Host lets adapters distinguish providers or enterprise instances. Owner may
// be empty for providers that support ownerless repositories or when the caller
// has already scoped the adapter to an owner.
type RepositoryRef struct {
	// Host is the provider hostname, such as github.com.
	Host string
	// Owner is the organization or user namespace when the provider has one.
	Owner string
	// Name is the repository name.
	Name string
}

// FullName returns the owner/name repository name where an owner is present.
//
// Host is intentionally excluded because many provider APIs and UI labels use
// owner/name as the repository's display identity.
func (r RepositoryRef) FullName() string {
	if r.Owner == "" {
		return r.Name
	}
	return r.Owner + "/" + r.Name
}

// Repository describes remote repository metadata returned by a provider.
//
// The fields intentionally capture publish-relevant state only. Additional
// provider metadata should stay in adapter-specific types until workflow code
// needs a stable cross-provider contract.
type Repository struct {
	// Ref identifies the repository.
	Ref RepositoryRef
	// CloneURL is the HTTPS clone URL.
	CloneURL string
	// SSHURL is the SSH clone URL.
	SSHURL string
	// WebURL is the browser URL.
	WebURL string
	// DefaultBranch is the provider's default branch name.
	DefaultBranch string
	// Private reports whether the repository is private.
	Private bool
	// Archived reports whether the repository is read-only archived.
	Archived bool
	// Disabled reports whether the provider has disabled the repository.
	Disabled bool
	// Permissions contains effective permissions for the current principal.
	Permissions RepositoryPermissions
}
