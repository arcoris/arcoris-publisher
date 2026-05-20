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

import "arcoris.dev/arcoris-publisher/internal/domain/manifest"

// Issue describes one registry validation or consistency problem.
//
// Fields are intentionally broad rather than split into several issue structs:
// callers can inspect Code first, then read the fields relevant to that issue.
// Message is always populated for direct presentation in logs or diagnostics.
type Issue struct {
	// Code classifies the invariant that failed.
	Code IssueCode
	// Module identifies the module being indexed when the issue was found.
	Module manifest.ModuleName
	// Other identifies the earlier module that already claimed the same key.
	Other manifest.ModuleName
	// ModulePath is populated for module-path index issues.
	ModulePath manifest.ModulePath
	// SourceDir is populated for source-directory index issues.
	SourceDir manifest.SourceDir
	// Repository is populated for repository index issues.
	Repository manifest.RepositoryRef
	// Branch is populated for branch-mapping issues.
	Branch manifest.BranchName
	// Index identifies the internal index when Validate finds registry corruption.
	Index string
	// Message explains the issue in a human-readable form.
	Message string
}
