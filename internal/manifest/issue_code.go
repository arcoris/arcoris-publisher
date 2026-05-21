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

// IssueCode is a stable machine-readable validation issue code.
type IssueCode string

const (
	// IssueInvalidAPIVersion indicates an unsupported apiVersion value.
	IssueInvalidAPIVersion IssueCode = "invalid_api_version"
	// IssueInvalidKind indicates an unexpected manifest kind.
	IssueInvalidKind IssueCode = "invalid_kind"
	// IssueMissingField indicates that a required field is absent or empty.
	IssueMissingField IssueCode = "missing_field"
	// IssueInvalidValue indicates a malformed scalar value.
	IssueInvalidValue IssueCode = "invalid_value"
	// IssueInvalidPath indicates an unsafe or malformed manifest path.
	IssueInvalidPath IssueCode = "invalid_path"
	// IssueDuplicateValue indicates a duplicate value where uniqueness is required.
	IssueDuplicateValue IssueCode = "duplicate_value"
	// IssueUnknownModule indicates a reference to an undeclared module.
	IssueUnknownModule IssueCode = "unknown_module"
	// IssueModuleNameMismatch indicates that staging and module manifests disagree.
	IssueModuleNameMismatch IssueCode = "module_name_mismatch"
	// IssueMissingModuleManifest indicates that a staging module has no module manifest.
	IssueMissingModuleManifest IssueCode = "missing_module_manifest"
	// IssueInvalidDependency indicates an invalid module dependency relation.
	IssueInvalidDependency IssueCode = "invalid_dependency"
	// IssueInvalidPublishEntry indicates an invalid explicit publish entry.
	IssueInvalidPublishEntry IssueCode = "invalid_publish_entry"
)
