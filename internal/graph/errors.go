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

import "arcoris.dev/arcoris-publisher/internal/diagnostic"

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
type Issue = diagnostic.Issue[IssueCode]

// ValidationError groups graph issues while preserving deterministic order.
type ValidationError = diagnostic.ValidationError[IssueCode]

type issueCollector = diagnostic.Collector[IssueCode]

func newIssueCollector() issueCollector {
	return diagnostic.NewCollector[IssueCode]("graph")
}
