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

import (
	"errors"
	"fmt"
)

// IssueCollector accumulates validation issues in deterministic order.
type IssueCollector struct {
	issues []Issue
}

// Add appends one issue to the collector.
func (c *IssueCollector) Add(code IssueCode, path string, format string, args ...any) {
	c.issues = append(c.issues, Issue{Code: code, Path: path, Message: fmt.Sprintf(format, args...)})
}

// AddError appends all issues from err when it carries a ValidationError.
//
// Nested validation errors keep their stable issue codes and messages. The
// caller-provided path is prepended so high-level constructors can report both
// the object being validated and the exact field that failed.
func (c *IssueCollector) AddError(path string, err error) {
	if err == nil {
		return
	}
	var validation *ValidationError
	if errors.As(err, &validation) {
		for _, issue := range validation.Issues {
			if issue.Path == "" {
				issue.Path = path
			} else if path != "" {
				issue.Path = path + "." + issue.Path
			}
			c.issues = append(c.issues, issue)
		}
		return
	}
	c.Add(IssueInvalidValue, path, "%s", err.Error())
}

// Issues returns a detached slice of accumulated issues.
func (c *IssueCollector) Issues() []Issue {
	out := make([]Issue, len(c.issues))
	copy(out, c.issues)
	return out
}

// Err returns nil when no issues were collected, otherwise a ValidationError.
func (c *IssueCollector) Err() error {
	if len(c.issues) == 0 {
		return nil
	}
	return &ValidationError{Issues: c.Issues()}
}
