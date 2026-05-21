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

import (
	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/versioning"
)

// DependencyRequirement is a direct internal module requirement that later
// modulefile workflow stages can write into go.mod.
//
// Requirements are intentionally direct only. The Go toolchain is responsible
// for maintaining indirect transitive requirements.
type DependencyRequirement struct {
	// module is the resolved module name of the dependency.
	module manifest.ModuleName

	// modulePath is the Go module path to require.
	modulePath manifest.ModulePath

	// version is the assigned version for modulePath.
	version versioning.Version
}

// newDependencyRequirement converts a versioning requirement into plan data.
func newDependencyRequirement(req versioning.Requirement) DependencyRequirement {
	return DependencyRequirement{
		module:     req.Module(),
		modulePath: req.ModulePath(),
		version:    req.Version(),
	}
}

// Module returns the internal dependency module name.
func (r DependencyRequirement) Module() manifest.ModuleName { return r.module }

// ModulePath returns the Go module path to require.
func (r DependencyRequirement) ModulePath() manifest.ModulePath { return r.modulePath }

// Version returns the assigned version required for ModulePath.
func (r DependencyRequirement) Version() versioning.Version { return r.version }

// cloneRequirements detaches dependency requirement slices.
func cloneRequirements(in []DependencyRequirement) []DependencyRequirement {
	out := make([]DependencyRequirement, len(in))
	copy(out, in)
	return out
}
