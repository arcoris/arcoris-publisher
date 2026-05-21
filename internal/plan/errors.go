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

package plan

import (
	"fmt"
	"strings"
)

// IssueCode is a stable machine-readable publication-plan issue code.
type IssueCode string

const (
	// IssueInvalidRequest indicates that required planning inputs are absent or
	// inconsistent.
	IssueInvalidRequest IssueCode = "invalid_request"
	// IssueEmptyPlan indicates that no publishable modules are available for
	// planning.
	IssueEmptyPlan IssueCode = "empty_plan"
	// IssueGraphOrder indicates that the dependency graph cannot produce an
	// acyclic publish order.
	IssueGraphOrder IssueCode = "graph_order"
	// IssueUnknownModule indicates that a publish-order module is absent from
	// registry lookup.
	IssueUnknownModule IssueCode = "unknown_module"
	// IssueNonPublishableModule indicates that a planned module is not
	// public/publishable.
	IssueNonPublishableModule IssueCode = "non_publishable_module"
	// IssueMissingAssignment indicates that a planned module or dependency has
	// no assigned version.
	IssueMissingAssignment IssueCode = "missing_assignment"
	// IssueMissingRequirements indicates that assignment data lacks a
	// requirement set for a module.
	IssueMissingRequirements IssueCode = "missing_requirements"
	// IssueEmptyBranches indicates that a planned module has no effective branch mappings.
	IssueEmptyBranches IssueCode = "empty_branches"
	// IssueEmptyPublishEntries indicates that a planned module has no explicit publish entries.
	IssueEmptyPublishEntries IssueCode = "empty_publish_entries"
	// IssueDuplicateModuleName indicates duplicate module names in the final plan index.
	IssueDuplicateModuleName IssueCode = "duplicate_module_name"
	// IssueDuplicateModulePath indicates duplicate module paths in the final plan index.
	IssueDuplicateModulePath IssueCode = "duplicate_module_path"
	// IssueDuplicateRepository indicates duplicate target repositories in the final plan index.
	IssueDuplicateRepository IssueCode = "duplicate_repository"
)

// Issue describes one planning problem.
type Issue struct {
	// Code is the stable machine-readable reason for the issue.
	Code IssueCode

	// Path identifies the request or plan location that failed.
	Path string

	// Message is the human-readable diagnostic text.
	Message string
}

// Error returns a compact issue message.
func (i Issue) Error() string {
	if i.Path == "" {
		return fmt.Sprintf("%s: %s", i.Code, i.Message)
	}
	return fmt.Sprintf("%s: %s: %s", i.Code, i.Path, i.Message)
}

// ValidationError groups planning issues while preserving deterministic order.
type ValidationError struct {
	// Issues contains all collected planning issues in deterministic order.
	Issues []Issue
}

// Error returns a compact summary of all planning issues.
func (e *ValidationError) Error() string {
	if e == nil || len(e.Issues) == 0 {
		return "plan validation failed"
	}
	parts := make([]string, 0, len(e.Issues))
	for _, issue := range e.Issues {
		parts = append(parts, issue.Error())
	}
	return strings.Join(parts, "; ")
}

// Has reports whether the validation error contains code.
func (e *ValidationError) Has(code IssueCode) bool {
	if e == nil {
		return false
	}
	for _, issue := range e.Issues {
		if issue.Code == code {
			return true
		}
	}
	return false
}

// Empty reports whether the validation error has no issues.
func (e *ValidationError) Empty() bool { return e == nil || len(e.Issues) == 0 }

type issueCollector struct {
	// issues buffers validation issues before they are exposed as a detached
	// ValidationError.
	issues []Issue
}

// add records one formatted issue.
func (c *issueCollector) add(code IssueCode, path string, format string, args ...any) {
	c.issues = append(c.issues, Issue{Code: code, Path: path, Message: fmt.Sprintf(format, args...)})
}

// err returns nil for an empty collector or a detached ValidationError.
func (c *issueCollector) err() error {
	if len(c.issues) == 0 {
		return nil
	}
	out := make([]Issue, len(c.issues))
	copy(out, c.issues)
	return &ValidationError{Issues: out}
}
