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

import "arcoris.dev/arcoris-publisher/internal/domain/manifest"

// Issue describes one publication-plan validation problem.
type Issue struct {
	// Code classifies the failed invariant.
	Code IssueCode
	// Module identifies the planned or referenced module involved in the issue.
	Module manifest.ModuleName
	// Dependency identifies the dependency involved in the issue, when relevant.
	Dependency manifest.ModuleName
	// Path identifies a module path involved in the issue, when relevant.
	Path manifest.ModulePath
	// Repository identifies a target repository involved in the issue.
	Repository manifest.RepositoryRef
	// Index identifies a broken lookup index when Validate detects one.
	Index string
	// Message explains the issue in human-readable form.
	Message string
}
