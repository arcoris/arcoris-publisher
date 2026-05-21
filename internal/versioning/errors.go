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

import "arcoris.dev/arcoris-publisher/internal/diagnostic"

// IssueCode is a stable machine-readable versioning issue code.
type IssueCode string

const (
	// IssueInvalidRequest indicates that required versioning inputs are absent or inconsistent.
	IssueInvalidRequest IssueCode = "invalid_request"
	// IssueInvalidVersion indicates that the supplied version is incompatible with policy.
	IssueInvalidVersion IssueCode = "invalid_version"
	// IssueUnknownModule indicates that a graph module cannot be found in the registry.
	IssueUnknownModule IssueCode = "unknown_module"
	// IssueMissingAssignment indicates that a requirement references a module
	// without an assigned version.
	IssueMissingAssignment IssueCode = "missing_assignment"
	// IssueNonPublishableDependency indicates that a publishable module depends
	// on a non-publishable module.
	IssueNonPublishableDependency IssueCode = "non_publishable_dependency"
	// IssueGraphOrder indicates that the dependency graph cannot produce an
	// acyclic publication order.
	IssueGraphOrder IssueCode = "graph_order"
)

// Issue describes one version assignment problem.
type Issue = diagnostic.Issue[IssueCode]

// ValidationError groups versioning issues while preserving deterministic order.
type ValidationError = diagnostic.ValidationError[IssueCode]

type issueCollector = diagnostic.Collector[IssueCode]

func newIssueCollector() issueCollector {
	return diagnostic.NewCollector[IssueCode]("versioning")
}
