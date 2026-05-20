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

// CreatePullRequestRequest describes a pull request creation request.
//
// Branch names are provider branch names, not fully qualified Git refs. The
// request assumes the branches already exist or will be made visible by the
// caller through the git port before this provider operation runs.
type CreatePullRequestRequest struct {
	// Repository identifies the target repository.
	Repository RepositoryRef
	// Title is the pull request title.
	Title string
	// Body is the pull request description.
	Body string
	// HeadBranch is the source branch name.
	HeadBranch string
	// BaseBranch is the target branch name.
	BaseBranch string
	// Draft creates a draft pull request when supported.
	Draft bool
}

// PullRequest describes a remote pull request.
//
// Number and URL are the stable values workflow code can show to operators or
// store in publish reports. Provider-specific IDs are intentionally omitted until
// a workflow needs them.
type PullRequest struct {
	// Number is the provider-visible pull request number.
	Number int
	// URL is the browser URL for the pull request.
	URL string
}
