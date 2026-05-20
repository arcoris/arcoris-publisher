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

import "arcoris.dev/arcoris-publisher/internal/ports/porterr"

const (
	// CodeCommandFailed identifies a Git command that exited unsuccessfully.
	CodeCommandFailed porterr.Code = "git_command_failed"
	// CodeRefNotFound identifies a missing local Git reference.
	CodeRefNotFound porterr.Code = "git_ref_not_found"
	// CodeRemoteRefNotFound identifies a missing remote Git reference.
	CodeRemoteRefNotFound porterr.Code = "git_remote_ref_not_found"
	// CodeDirtyWorktree identifies uncommitted local changes.
	CodeDirtyWorktree porterr.Code = "git_dirty_worktree"
	// CodeNoChanges identifies a commit or publish attempt with no changes.
	CodeNoChanges porterr.Code = "git_no_changes"
	// CodePushRejected identifies a remote push rejection.
	CodePushRejected porterr.Code = "git_push_rejected"
	// CodeAuthenticationFailed identifies invalid or missing Git credentials.
	CodeAuthenticationFailed porterr.Code = "git_authentication_failed"
	// CodeTagAlreadyExists identifies an attempt to recreate an existing tag.
	CodeTagAlreadyExists porterr.Code = "git_tag_already_exists"
	// CodeRepositoryNotFound identifies a missing local or remote repository.
	CodeRepositoryNotFound porterr.Code = "git_repository_not_found"
)
