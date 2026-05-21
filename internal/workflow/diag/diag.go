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

package diag

import (
	"fmt"
	"strings"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

// IssueCode constrains workflow issue code types to stable string-like values.
type IssueCode interface {
	~string
}

// Issue describes one workflow validation or state issue.
type Issue[C IssueCode] struct {
	// Code is the stable machine-readable issue reason.
	Code C

	// Module identifies the planned module that owns the issue when applicable.
	Module manifest.ModuleName

	// Path identifies the request field or filesystem path that failed.
	Path string

	// Message is the human-readable diagnostic text.
	Message string
}

// ValidationError aggregates workflow issues in deterministic order.
type ValidationError[C IssueCode] struct {
	// Scope names the workflow stage, such as "source" or "target".
	Scope string

	// Issues contains collected workflow issues in deterministic order.
	Issues []Issue[C]
}

// Error returns a compact validation error summary.
func (e *ValidationError[C]) Error() string {
	scope := e.scope()
	if e == nil || len(e.Issues) == 0 {
		return scope + " validation failed"
	}
	if len(e.Issues) == 1 {
		issue := e.Issues[0]
		if issue.Path == "" {
			return issue.Message
		}
		return issue.Path + ": " + issue.Message
	}

	var b strings.Builder
	fmt.Fprintf(&b, "%s validation failed with %d issues", scope, len(e.Issues))
	for _, issue := range e.Issues {
		if issue.Path == "" {
			fmt.Fprintf(&b, "; %s", issue.Message)
			continue
		}
		fmt.Fprintf(&b, "; %s: %s", issue.Path, issue.Message)
	}

	return b.String()
}

// Has reports whether the validation error contains code.
func (e *ValidationError[C]) Has(code C) bool {
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

// scope returns the configured stage name or a generic fallback for zero-value
// errors built in tests.
func (e *ValidationError[C]) scope() string {
	if e == nil || e.Scope == "" {
		return "workflow"
	}
	return e.Scope
}

// Collector accumulates workflow diagnostics before exposing a detached error.
type Collector[C IssueCode] struct {
	// scope names the workflow stage.
	scope string

	// issues buffers validation issues before Err returns a detached error.
	issues []Issue[C]
}

// NewCollector creates an empty workflow issue collector for scope.
func NewCollector[C IssueCode](scope string) Collector[C] {
	return Collector[C]{scope: scope}
}

// Add records one formatted issue.
func (c *Collector[C]) Add(
	code C,
	module manifest.ModuleName,
	path string,
	format string,
	args ...any,
) {
	c.AddMessage(code, module, path, fmt.Sprintf(format, args...))
}

// AddMessage records one already formatted issue.
func (c *Collector[C]) AddMessage(
	code C,
	module manifest.ModuleName,
	path string,
	message string,
) {
	c.issues = append(c.issues, Issue[C]{
		Code:    code,
		Module:  module,
		Path:    path,
		Message: message,
	})
}

// Append stores detached issues from a nested validation pass.
func (c *Collector[C]) Append(issues []Issue[C]) {
	c.issues = append(c.issues, CloneIssues(issues)...)
}

// Len returns the number of collected issues.
func (c *Collector[C]) Len() int {
	return len(c.issues)
}

// Issues returns detached collected issues.
func (c *Collector[C]) Issues() []Issue[C] {
	return CloneIssues(c.issues)
}

// Err returns nil for an empty collector or a detached ValidationError.
func (c *Collector[C]) Err() error {
	if len(c.issues) == 0 {
		return nil
	}
	return &ValidationError[C]{
		Scope:  c.scope,
		Issues: CloneIssues(c.issues),
	}
}

// CloneIssues returns a detached copy of issues.
func CloneIssues[C IssueCode](issues []Issue[C]) []Issue[C] {
	out := make([]Issue[C], len(issues))
	copy(out, issues)
	return out
}
