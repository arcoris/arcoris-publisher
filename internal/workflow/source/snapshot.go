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

package source

import "arcoris.dev/arcoris-publisher/internal/manifest"

// Snapshot captures the inspected source checkout state for a publication plan.
type Snapshot struct {
	// repository captures Git state and inspected root paths.
	repository RepositorySnapshot

	// modules contains source snapshots in plan publication order.
	modules []ModuleSnapshot

	// warnings contains non-fatal diagnostics collected during inspection.
	warnings []Issue
}

// Repository returns the source repository snapshot.
func (s Snapshot) Repository() RepositorySnapshot { return s.repository }

// Modules returns detached module snapshots in plan publication order.
func (s Snapshot) Modules() []ModuleSnapshot { return cloneModuleSnapshots(s.modules) }

// ModuleNames returns source snapshot module names in plan publication order.
func (s Snapshot) ModuleNames() []manifest.ModuleName {
	out := make([]manifest.ModuleName, 0, len(s.modules))
	for _, mod := range s.modules {
		out = append(out, mod.Name())
	}
	return out
}

// Warnings returns non-fatal source diagnostics, such as a dirty source checkout
// when the source dirty policy is warn.
func (s Snapshot) Warnings() []Issue { return cloneIssues(s.warnings) }

// cloneModuleSnapshots detaches module snapshot slices before returning them.
func cloneModuleSnapshots(in []ModuleSnapshot) []ModuleSnapshot {
	out := make([]ModuleSnapshot, len(in))
	copy(out, in)
	return out
}
