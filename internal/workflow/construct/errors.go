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

package construct

import "arcoris.dev/arcoris-publisher/internal/workflow/diag"

type IssueCode string

const (
	IssueInvalidRequest   IssueCode = "invalid_request"
	IssueMissingSource    IssueCode = "missing_source_snapshot"
	IssueMissingTarget    IssueCode = "missing_target_workspace"
	IssueTargetPathEscape IssueCode = "target_path_escape"
	IssueEntryCopyFailed  IssueCode = "entry_copy_failed"
	IssueCleanFailed      IssueCode = "target_clean_failed"
)

type Issue = diag.Issue[IssueCode]

type ValidationError = diag.ValidationError[IssueCode]

type issueCollector = diag.Collector[IssueCode]

func newIssueCollector() issueCollector {
	return diag.NewCollector[IssueCode]("construct")
}
