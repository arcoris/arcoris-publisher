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

// CreateReleaseRequest describes a release creation request.
type CreateReleaseRequest struct {
	// Repository identifies the target repository.
	Repository RepositoryRef
	// TagName is the Git tag that the release describes.
	TagName string
	// Name is the human-readable release name.
	Name string
	// Body is the release notes body.
	Body string
	// Draft creates an unpublished release draft when supported.
	Draft bool
	// Prerelease marks the release as not yet stable.
	Prerelease bool
}

// Release describes a remote release object.
type Release struct {
	// ID is the provider-visible release identifier.
	ID string
	// URL is the browser URL for the release.
	URL string
}
