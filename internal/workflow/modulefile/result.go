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

package modulefile

import "arcoris.dev/arcoris-publisher/internal/manifest"

// Result describes rewritten target module files in plan order.
type Result struct{ modules []ModuleResult }

// Modules returns detached module-file rewrite results.
func (r Result) Modules() []ModuleResult {
	out := make([]ModuleResult, len(r.modules))
	copy(out, r.modules)
	return out
}

// ModuleByName returns the module-file result for name.
func (r Result) ModuleByName(name manifest.ModuleName) (ModuleResult, bool) {
	for _, m := range r.modules {
		if m.Module() == name {
			return m, true
		}
	}
	return ModuleResult{}, false
}

// Changed reports whether any target go.mod was rewritten.
func (r Result) Changed() bool {
	for _, m := range r.modules {
		if m.Changed() {
			return true
		}
	}
	return false
}

// ModuleResult describes one target module file rewrite.
type ModuleResult struct {
	// module is the planned module name.
	module manifest.ModuleName

	// goModPath is the rewritten target go.mod path.
	goModPath string

	// changed reports whether file content changed.
	changed bool

	// requirements lists internal dependency requirements written to go.mod.
	requirements []RequirementUpdate
}

// Module returns the planned module name.
func (m ModuleResult) Module() manifest.ModuleName { return m.module }

// GoModPath returns the rewritten target go.mod path.
func (m ModuleResult) GoModPath() string { return m.goModPath }

// Changed reports whether file content changed.
func (m ModuleResult) Changed() bool { return m.changed }

// Requirements returns detached internal dependency requirements written to go.mod.
func (m ModuleResult) Requirements() []RequirementUpdate {
	out := make([]RequirementUpdate, len(m.requirements))
	copy(out, m.requirements)
	return out
}

// RequirementUpdate describes one managed internal requirement written to go.mod.
type RequirementUpdate struct {
	// modulePath is the required module path.
	modulePath manifest.ModulePath

	// version is the required module version string.
	version string
}

// ModulePath returns the required module path.
func (u RequirementUpdate) ModulePath() manifest.ModulePath { return u.modulePath }

// Version returns the required module version string.
func (u RequirementUpdate) Version() string { return u.version }
