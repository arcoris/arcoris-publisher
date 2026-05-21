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

import "arcoris.dev/arcoris-publisher/internal/manifest"

// Requirement is a direct internal module requirement that later modulefile
// workflow stages can write into go.mod.
//
// Requirements intentionally represent direct manifest dependencies only. The Go
// toolchain is responsible for maintaining indirect transitive requirements.
type Requirement struct {
	// module is the resolved manifest module name of the dependency.
	module manifest.ModuleName

	// modulePath is the Go module path that downstream go.mod rewriting should
	// require.
	modulePath manifest.ModulePath

	// version is the assigned version that should be written for modulePath.
	version Version
}

// newRequirement creates one direct dependency requirement value.
func newRequirement(
	module manifest.ModuleName,
	modulePath manifest.ModulePath,
	version Version,
) Requirement {
	return Requirement{module: module, modulePath: modulePath, version: version}
}

// Module returns the internal dependency module name.
func (r Requirement) Module() manifest.ModuleName { return r.module }

// ModulePath returns the Go module path required by go.mod.
func (r Requirement) ModulePath() manifest.ModulePath { return r.modulePath }

// Version returns the required module version.
func (r Requirement) Version() Version { return r.version }

// cloneRequirements detaches requirement slices before storing or returning
// them from immutable-by-convention snapshots.
func cloneRequirements(in []Requirement) []Requirement {
	out := make([]Requirement, len(in))
	copy(out, in)
	return out
}
