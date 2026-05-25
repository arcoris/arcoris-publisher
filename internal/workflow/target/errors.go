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

package target

import "arcoris.dev/arcoris-publisher/internal/workflow/diag"

// IssueCode identifies a target preparation validation issue.
type IssueCode string

const (
	IssueInvalidRequest       IssueCode = "invalid_request"
	IssueRootMissing          IssueCode = "root_missing"
	IssueRootNotDirectory     IssueCode = "root_not_directory"
	IssueWorktreeMissing      IssueCode = "worktree_missing"
	IssueWorktreeNotDirectory IssueCode = "worktree_not_directory"
	IssueWorktreeStatusFailed IssueCode = "worktree_status_failed"
	IssueWorktreeDirty        IssueCode = "worktree_dirty"
	IssueFetchFailed          IssueCode = "fetch_failed"
	IssueCloneURLMissing      IssueCode = "clone_url_missing"
)

// Issue describes one target preparation validation issue.
type Issue = diag.Issue[IssueCode]

// ValidationError aggregates target preparation validation issues.
type ValidationError = diag.ValidationError[IssueCode]

type issueCollector = diag.Collector[IssueCode]

func newIssueCollector() issueCollector {
	return diag.NewCollector[IssueCode]("target")
}
