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

// IssueCode identifies a versioning validation issue.
type IssueCode string

const (
	// IssueUnsupportedPolicy reports an unknown or missing version policy.
	IssueUnsupportedPolicy IssueCode = "unsupported_policy"
	// IssueInvalidReleaseVersion reports an invalid release-train version.
	IssueInvalidReleaseVersion IssueCode = "invalid_release_version"
	// IssueInvalidSnapshot reports invalid snapshot versioning input.
	IssueInvalidSnapshot IssueCode = "invalid_snapshot"
	// IssueInvalidAssignment reports an invalid module version assignment.
	IssueInvalidAssignment IssueCode = "invalid_assignment"
	// IssueDuplicateModule reports duplicate assignments for one module name.
	IssueDuplicateModule IssueCode = "duplicate_module"
	// IssueDuplicateModulePath reports duplicate assignments for one module path.
	IssueDuplicateModulePath IssueCode = "duplicate_module_path"
	// IssueUnknownDependency reports a dependency missing from the module registry.
	IssueUnknownDependency IssueCode = "unknown_dependency"
	// IssueUnassignedDependency reports a known dependency without assigned publish version.
	IssueUnassignedDependency IssueCode = "unassigned_dependency"
)
