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

// IssueCode identifies a plan construction or validation issue.
type IssueCode string

const (
	// IssueInvalidRegistry reports registry data that cannot build a plan.
	IssueInvalidRegistry IssueCode = "invalid_registry"
	// IssueInvalidGraph reports dependency graph data that cannot produce a plan.
	IssueInvalidGraph IssueCode = "invalid_graph"
	// IssueMissingModule reports a referenced module missing from the registry.
	IssueMissingModule IssueCode = "missing_module"
	// IssueMissingVersion reports a publishable module without a version assignment.
	IssueMissingVersion IssueCode = "missing_version"
	// IssueExtraVersion reports a version assignment outside the publish plan.
	IssueExtraVersion IssueCode = "extra_version"
	// IssueInvalidDependency reports a dependency that cannot become a requirement.
	IssueInvalidDependency IssueCode = "invalid_dependency"
	// IssueDuplicateModule reports duplicate module plans in one publish plan.
	IssueDuplicateModule IssueCode = "duplicate_module"
	// IssueDuplicateModulePath reports duplicate module paths in one publish plan.
	IssueDuplicateModulePath IssueCode = "duplicate_module_path"
	// IssueDuplicateRepository reports duplicate repositories in one publish plan.
	IssueDuplicateRepository IssueCode = "duplicate_repository"
	// IssueInvalidBranch reports an invalid branch mapping in a module plan.
	IssueInvalidBranch IssueCode = "invalid_branch"
	// IssueInvalidAction reports a module plan with an unsupported action.
	IssueInvalidAction IssueCode = "invalid_action"
	// IssueInvalidIndex reports lookup indexes that do not match plan modules.
	IssueInvalidIndex IssueCode = "invalid_index"
)
