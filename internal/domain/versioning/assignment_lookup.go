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

import "arcoris.dev/arcoris-publisher/internal/domain/manifest"

// VersionOfModule returns the version assigned to module.
func (a Assignments) VersionOfModule(module manifest.ModuleName) (Version, bool) {
	index, ok := a.byModule[module]
	if !ok {
		return "", false
	}
	return a.items[index].Version(), true
}

// ModuleVersion returns the full assignment for module.
func (a Assignments) ModuleVersion(module manifest.ModuleName) (ModuleVersion, bool) {
	index, ok := a.byModule[module]
	if !ok {
		return ModuleVersion{}, false
	}
	return a.items[index], true
}

// ContainsModule reports whether module has an assigned version.
func (a Assignments) ContainsModule(module manifest.ModuleName) bool {
	_, ok := a.byModule[module]
	return ok
}
