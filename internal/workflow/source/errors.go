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

package source

import (
	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/workflow/diag"
)

// IssueCode identifies a source inspection validation issue.
type IssueCode string

const (
	// IssueInvalidRequest indicates malformed source-inspection inputs or
	// infrastructure errors that prevent validation.
	IssueInvalidRequest IssueCode = "invalid_request"

	// IssueRepositoryMissing indicates that the configured source repository
	// root does not exist.
	IssueRepositoryMissing IssueCode = "repository_missing"

	// IssueRepositoryNotDirectory indicates that the repository root exists but
	// is not a directory.
	IssueRepositoryNotDirectory IssueCode = "repository_not_directory"

	// IssueStagingMissing indicates that the configured staging root does not
	// exist.
	IssueStagingMissing IssueCode = "staging_missing"

	// IssueStagingNotDirectory indicates that the staging root exists but is not
	// a directory.
	IssueStagingNotDirectory IssueCode = "staging_not_directory"

	// IssueStagingOutsideRepo indicates that the staging root is outside the
	// source repository root.
	IssueStagingOutsideRepo IssueCode = "staging_outside_repository"

	// IssueDetachedHead indicates that Git reports no current branch and the
	// caller did not explicitly allow detached HEAD inspection.
	IssueDetachedHead IssueCode = "detached_head"

	// IssueDirtySource indicates that the source checkout has Git status entries
	// under a fail-or-warn dirty policy.
	IssueDirtySource IssueCode = "dirty_source"

	// IssueModuleSourceMissing indicates that a planned module source directory
	// does not exist.
	IssueModuleSourceMissing IssueCode = "module_source_missing"

	// IssueModuleSourceNotDir indicates that a planned module source path exists
	// but is not a directory.
	IssueModuleSourceNotDir IssueCode = "module_source_not_directory"

	// IssueModuleRootMissing indicates that a planned module root directory does
	// not exist inside the module source directory.
	IssueModuleRootMissing IssueCode = "module_root_missing"

	// IssueModuleRootNotDir indicates that a planned module root path exists but
	// is not a directory.
	IssueModuleRootNotDir IssueCode = "module_root_not_directory"

	// IssueEntryMissing indicates that a required explicit publish entry source
	// path does not exist.
	IssueEntryMissing IssueCode = "entry_missing"

	// IssueEntryTypeMismatch indicates that a file entry points to a directory
	// or a directory entry points to a file.
	IssueEntryTypeMismatch IssueCode = "entry_type_mismatch"

	// IssueEntryPathEscape indicates that a resolved entry, module source, or
	// module root path escapes its allowed root.
	IssueEntryPathEscape IssueCode = "entry_path_escape"

	// IssueEntryHashFailed indicates that hashing a present source entry failed.
	IssueEntryHashFailed IssueCode = "entry_hash_failed"
)

// Issue describes one source inspection validation issue.
type Issue = diag.Issue[IssueCode]

// ValidationError aggregates source inspection validation issues.
type ValidationError = diag.ValidationError[IssueCode]

// singleIssueError builds a validation error for one issue.
func singleIssueError(
	code IssueCode,
	module manifest.ModuleName,
	path string,
	message string,
) *ValidationError {
	return &ValidationError{
		Scope: "source",
		Issues: []Issue{{
			Code:    code,
			Module:  module,
			Path:    path,
			Message: message,
		}},
	}
}

// cloneIssues detaches issue slices before storing or returning diagnostics.
func cloneIssues(in []Issue) []Issue {
	return diag.CloneIssues(in)
}

// issueCollector accumulates validation diagnostics while preserving the order
// in which inspection discovered them.
type issueCollector = diag.Collector[IssueCode]

// newIssueCollector creates a source-scoped collector.
func newIssueCollector() issueCollector {
	return diag.NewCollector[IssueCode]("source")
}
