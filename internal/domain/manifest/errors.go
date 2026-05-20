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

package manifest

import "fmt"

// IssueCode identifies a manifest validation issue.
type IssueCode string

const (
	// IssueUnsupportedVersion reports an unknown or missing manifest version.
	IssueUnsupportedVersion IssueCode = "unsupported_version"
	// IssueInvalidSource reports an invalid authoritative source declaration.
	IssueInvalidSource IssueCode = "invalid_source"
	// IssueInvalidPolicy reports an invalid global publication policy.
	IssueInvalidPolicy IssueCode = "invalid_policy"
	// IssueInvalidModule reports an invalid module declaration.
	IssueInvalidModule IssueCode = "invalid_module"
	// IssueDuplicateModule reports repeated module names.
	IssueDuplicateModule IssueCode = "duplicate_module"
	// IssueDuplicatePath reports repeated module paths, source dirs, or repositories.
	IssueDuplicatePath IssueCode = "duplicate_path"
	// IssueUnknownDependency reports a dependency pointing at no declared module.
	IssueUnknownDependency IssueCode = "unknown_dependency"
	// IssueInvalidDependency reports a semantically invalid dependency edge.
	IssueInvalidDependency IssueCode = "invalid_dependency"
)

// Issue describes one manifest validation problem.
type Issue struct {
	Code    IssueCode
	Path    string
	Message string
}

// ValidationError contains one or more manifest validation issues.
type ValidationError struct {
	Issues []Issue
}

// Error returns a concise validation failure summary.
func (e *ValidationError) Error() string {
	if e == nil || len(e.Issues) == 0 {
		return "manifest validation failed"
	}
	if len(e.Issues) == 1 {
		return renderIssue(e.Issues[0])
	}
	return fmt.Sprintf("manifest validation failed with %d issues", len(e.Issues))
}

// renderIssue includes the issue path when one is available.
func renderIssue(issue Issue) string {
	if issue.Path == "" {
		return issue.Message
	}
	return issue.Path + ": " + issue.Message
}

// validationErrorf creates a one-issue validation error.
func validationErrorf(code IssueCode, path string, format string, args ...any) error {
	return &ValidationError{Issues: []Issue{{Code: code, Path: path, Message: fmt.Sprintf(format, args...)}}}
}

// issuePath formats an indexed manifest collection path.
func issuePath(name string, index int) string {
	return fmt.Sprintf("%s[%d]", name, index)
}
