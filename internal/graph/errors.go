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

package graph

import (
	"fmt"
	"strings"
)

// IssueCode is a stable machine-readable graph validation issue code.
type IssueCode string

const (
	// IssueInvalidRegistry indicates that the input registry cannot be converted
	// into a dependency graph.
	IssueInvalidRegistry IssueCode = "invalid_registry"
	// IssueDuplicateNode indicates that the same module name was seen more than once.
	IssueDuplicateNode IssueCode = "duplicate_node"
	// IssueUnknownDependency indicates a dependency on a module absent from the graph.
	IssueUnknownDependency IssueCode = "unknown_dependency"
	// IssueDisabledDependency indicates a dependency on a disabled module.
	IssueDisabledDependency IssueCode = "disabled_dependency"
	// IssueSelfDependency indicates that a module directly depends on itself.
	IssueSelfDependency IssueCode = "self_dependency"
	// IssueUnknownNode indicates a query for a module not present in the graph.
	IssueUnknownNode IssueCode = "unknown_node"
	// IssueDependencyCycle indicates that no acyclic order can be produced.
	IssueDependencyCycle IssueCode = "dependency_cycle"
)

// Issue describes one graph validation or traversal problem.
type Issue struct {
	Code    IssueCode
	Path    string
	Message string
}

// Error returns a compact issue message.
func (i Issue) Error() string {
	if i.Path == "" {
		return fmt.Sprintf("%s: %s", i.Code, i.Message)
	}
	return fmt.Sprintf("%s: %s: %s", i.Code, i.Path, i.Message)
}

// ValidationError groups graph issues while preserving deterministic order.
type ValidationError struct {
	Issues []Issue
}

// Error returns a compact summary of all graph validation issues.
func (e *ValidationError) Error() string {
	if e == nil || len(e.Issues) == 0 {
		return "graph validation failed"
	}
	parts := make([]string, 0, len(e.Issues))
	for _, issue := range e.Issues {
		parts = append(parts, issue.Error())
	}
	return strings.Join(parts, "; ")
}

// Has reports whether the error contains issue code.
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
	issues []Issue
}

// add appends one issue while preserving discovery order.
func (c *issueCollector) add(
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

// err returns a detached ValidationError or nil when no issues were collected.
func (c *issueCollector) err() error {
	if len(c.issues) == 0 {
		return nil
	}
	out := make([]Issue, len(c.issues))
	copy(out, c.issues)
	return &ValidationError{Issues: out}
}
