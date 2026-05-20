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

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

// IssueCode identifies a graph validation issue.
type IssueCode string

const (
	// IssueDuplicateModule reports duplicate module names in graph input.
	IssueDuplicateModule IssueCode = "duplicate_module"
	// IssueUnknownDependency reports a dependency that does not reference a known module.
	IssueUnknownDependency IssueCode = "unknown_dependency"
	// IssueSelfDependency reports a module depending on itself.
	IssueSelfDependency IssueCode = "self_dependency"
	// IssueDuplicateDependency reports a repeated direct dependency declaration.
	IssueDuplicateDependency IssueCode = "duplicate_dependency"
	// IssueCycle reports one dependency cycle.
	IssueCycle IssueCode = "dependency_cycle"
)

// Issue describes one graph validation problem.
type Issue struct {
	Code IssueCode

	Module     manifest.ModuleName
	Dependency manifest.ModuleName
	Cycle      Cycle

	Message string
}

// ValidationError contains one or more graph validation issues.
type ValidationError struct {
	Issues []Issue
}

// Error returns a concise validation failure summary.
func (e *ValidationError) Error() string {
	if e == nil || len(e.Issues) == 0 {
		return "dependency graph validation failed"
	}
	if len(e.Issues) == 1 {
		return e.Issues[0].Message
	}
	return fmt.Sprintf("dependency graph validation failed with %d issues", len(e.Issues))
}
