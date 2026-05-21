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

package versioning

import "fmt"

// ValidationError contains one or more versioning validation issues.
type ValidationError struct {
	// Issues preserves deterministic discovery order so CLI diagnostics and
	// tests stay stable across runs.
	Issues []Issue
}

// Error returns a concise validation failure summary.
func (e *ValidationError) Error() string {
	if e == nil || len(e.Issues) == 0 {
		return "versioning validation failed"
	}
	if len(e.Issues) == 1 {
		issue := e.Issues[0]
		if issue.Path == "" {
			return issue.Message
		}
		return issue.Path + ": " + issue.Message
	}
	return fmt.Sprintf("versioning validation failed with %d issues", len(e.Issues))
}
