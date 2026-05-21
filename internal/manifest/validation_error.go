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

import "strings"

// ValidationError groups multiple manifest issues while preserving stable order.
type ValidationError struct {
	Issues []Issue
}

// Error returns a compact summary of all validation issues.
func (e *ValidationError) Error() string {
	if e == nil || len(e.Issues) == 0 {
		return "manifest validation failed"
	}
	parts := make([]string, 0, len(e.Issues))
	for _, issue := range e.Issues {
		parts = append(parts, issue.Error())
	}
	return strings.Join(parts, "; ")
}
