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

package diagnostic

import (
	"fmt"
	"strings"
)

// Code constrains diagnostic codes to stable string-like values.
type Code interface {
	~string
}

// Issue describes one validation or state problem.
type Issue[C Code] struct {
	// Code is the stable machine-readable reason for the issue.
	Code C

	// Path identifies the input, derived object, or query location that failed.
	Path string

	// Message is the human-readable diagnostic text.
	Message string
}

// Error returns one compact issue message.
func (i Issue[C]) Error() string {
	if i.Path == "" {
		return fmt.Sprintf("%s: %s", i.Code, i.Message)
	}

	return fmt.Sprintf("%s: %s: %s", i.Code, i.Path, i.Message)
}

// ValidationError groups issues while preserving deterministic order.
type ValidationError[C Code] struct {
	// Scope names the package or workflow stage that produced the issues.
	Scope string

	// Issues contains all collected validation issues in deterministic order.
	Issues []Issue[C]
}

// Error returns all validation issues as one compact message.
func (e *ValidationError[C]) Error() string {
	if e == nil || len(e.Issues) == 0 {
		return validationFailed(e.scope())
	}

	parts := make([]string, 0, len(e.Issues))
	for _, issue := range e.Issues {
		parts = append(parts, issue.Error())
	}

	return strings.Join(parts, "; ")
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

// Empty reports whether the validation error has no issues.
func (e *ValidationError[C]) Empty() bool {
	return e == nil || len(e.Issues) == 0
}

// scope returns the package-specific scope or a generic fallback for zero-value
// errors built directly by tests.
func (e *ValidationError[C]) scope() string {
	if e == nil || e.Scope == "" {
		return ""
	}
	return e.Scope
}

func validationFailed(scope string) string {
	if scope == "" {
		return "validation failed"
	}
	return scope + " validation failed"
}

// Collector accumulates diagnostics before exposing a detached error.
type Collector[C Code] struct {
	// scope names the package or workflow stage that owns the issues.
	scope string

	// issues buffers diagnostics before Err returns a detached error.
	issues []Issue[C]
}

// NewCollector creates an empty issue collector for scope.
func NewCollector[C Code](scope string) Collector[C] {
	return Collector[C]{scope: scope}
}

// Add records one formatted issue.
func (c *Collector[C]) Add(code C, path string, format string, args ...any) {
	c.AddMessage(code, path, fmt.Sprintf(format, args...))
}

// AddMessage records one already formatted issue.
func (c *Collector[C]) AddMessage(code C, path string, message string) {
	c.issues = append(c.issues, Issue[C]{
		Code:    code,
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
func CloneIssues[C Code](issues []Issue[C]) []Issue[C] {
	out := make([]Issue[C], len(issues))
	copy(out, issues)
	return out
}
