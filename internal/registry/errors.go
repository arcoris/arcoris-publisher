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

package registry

import (
	"fmt"
	"strings"
)

// IssueCode is a stable machine-readable registry validation code.
type IssueCode string

const (
	// IssueDuplicateModuleName reports a repeated module name.
	IssueDuplicateModuleName IssueCode = "duplicate_module_name"
	// IssueDuplicateModulePath reports a repeated module path.
	IssueDuplicateModulePath IssueCode = "duplicate_module_path"
	// IssueDuplicateSourceDir reports a repeated source directory.
	IssueDuplicateSourceDir IssueCode = "duplicate_source_dir"
	// IssueDuplicateRepository reports a repeated target repository.
	IssueDuplicateRepository IssueCode = "duplicate_repository"
	// IssueInvalidPublicationSet reports a structurally unusable resolved set.
	IssueInvalidPublicationSet IssueCode = "invalid_publication_set"
)

// Issue describes one registry validation problem.
type Issue struct {
	Code    IssueCode
	Path    string
	Message string
}

// Error returns a compact human-readable issue message.
func (i Issue) Error() string {
	if i.Path == "" {
		return fmt.Sprintf("%s: %s", i.Code, i.Message)
	}

	return fmt.Sprintf("%s: %s: %s", i.Code, i.Path, i.Message)
}

// ValidationError groups registry validation issues in stable order.
type ValidationError struct {
	Issues []Issue
}

// Error returns all registry validation issues as one message.
func (e *ValidationError) Error() string {
	if e == nil || len(e.Issues) == 0 {
		return "registry validation failed"
	}

	parts := make([]string, 0, len(e.Issues))
	for _, issue := range e.Issues {
		parts = append(parts, issue.Error())
	}

	return strings.Join(parts, "; ")
}

type issueCollector struct {
	issues []Issue
}

func (c *issueCollector) Add(
	code IssueCode,
	path string,
	format string,
	args ...any,
) {
	c.issues = append(c.issues, Issue{
		Code:    code,
		Path:    path,
		Message: fmt.Sprintf(format, args...),
	})
}

func (c *issueCollector) Issues() []Issue {
	out := make([]Issue, len(c.issues))
	copy(out, c.issues)
	return out
}

func (c *issueCollector) Err() error {
	if len(c.issues) == 0 {
		return nil
	}

	return &ValidationError{
		Issues: c.Issues(),
	}
}
