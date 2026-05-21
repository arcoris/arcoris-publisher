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

import "arcoris.dev/arcoris-publisher/internal/diagnostic"

// IssueCode is a stable machine-readable registry validation code.
type IssueCode string

const (
	// IssueDuplicateModuleName reports a repeated module name.
	IssueDuplicateModuleName IssueCode = "duplicate_module_name"
	// IssueDuplicateModulePath reports a repeated module path.
	IssueDuplicateModulePath IssueCode = "duplicate_module_path"
	// IssueDuplicateSourceDir reports a repeated source directory.
	IssueDuplicateSourceDir IssueCode = "duplicate_source_dir"
	// IssueDuplicateRepository reports a repeated target repository.
	IssueDuplicateRepository IssueCode = "duplicate_repository"
	// IssueInvalidPublicationSet reports a structurally unusable resolved set.
	IssueInvalidPublicationSet IssueCode = "invalid_publication_set"
)

// Issue describes one registry validation problem.
type Issue = diagnostic.Issue[IssueCode]

// ValidationError groups registry validation issues in stable order.
type ValidationError = diagnostic.ValidationError[IssueCode]

type issueCollector = diagnostic.Collector[IssueCode]

func newIssueCollector() issueCollector {
	return diagnostic.NewCollector[IssueCode]("registry")
}
