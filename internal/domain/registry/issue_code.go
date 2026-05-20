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

package registry

// IssueCode identifies a registry construction or consistency violation.
type IssueCode string

const (
	// IssueDuplicateModuleName reports two modules with the same manifest-local name.
	IssueDuplicateModuleName IssueCode = "duplicate_module_name"
	// IssueDuplicateModulePath reports two modules with the same public Go module path.
	IssueDuplicateModulePath IssueCode = "duplicate_module_path"
	// IssueDuplicateSourceDir reports two modules with the same source directory.
	IssueDuplicateSourceDir IssueCode = "duplicate_source_dir"
	// IssueDuplicateRepository reports two modules with the same target repository.
	IssueDuplicateRepository IssueCode = "duplicate_repository"
	// IssueDuplicateBranchMapping reports duplicate source branch mappings in one module.
	IssueDuplicateBranchMapping IssueCode = "duplicate_branch_mapping"
	// IssueInvalidIndex reports a Registry whose internal indexes no longer match its modules.
	IssueInvalidIndex IssueCode = "invalid_index"
)
