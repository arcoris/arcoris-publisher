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

import "arcoris.dev/arcoris-publisher/internal/ports/porterr"

const (
	// CodeRepositoryNotFound identifies a missing remote repository.
	CodeRepositoryNotFound porterr.Code = "remote_repository_not_found"
	// CodeAccessDenied identifies insufficient provider permissions.
	CodeAccessDenied porterr.Code = "remote_access_denied"
	// CodeAuthenticationFailed identifies invalid or missing provider credentials.
	CodeAuthenticationFailed porterr.Code = "remote_authentication_failed"
	// CodeRateLimited identifies provider API rate limiting.
	CodeRateLimited porterr.Code = "remote_rate_limited"
	// CodeBranchProtected identifies an operation blocked by branch protection.
	CodeBranchProtected porterr.Code = "remote_branch_protected"
	// CodePullRequestFailed identifies a pull request API failure.
	CodePullRequestFailed porterr.Code = "remote_pull_request_failed"
	// CodeReleaseFailed identifies a release API failure.
	CodeReleaseFailed porterr.Code = "remote_release_failed"
)
