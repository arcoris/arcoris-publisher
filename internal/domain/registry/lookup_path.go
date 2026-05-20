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

// ModuleByPath returns the module with public Go module path and whether it was found.
//
// Module paths are unique in a valid registry because generated go.mod files
// must not publish two repositories under the same import path.
func (r Registry) ModuleByPath(modulePath manifest.ModulePath) (manifest.Module, bool) {
	index, ok := r.byModulePath[modulePath]
	if !ok {
		return manifest.Module{}, false
	}
	return r.modules[index], true
}

// ContainsPath reports whether a module with public Go module path exists.
func (r Registry) ContainsPath(modulePath manifest.ModulePath) bool {
	_, ok := r.byModulePath[modulePath]
	return ok
}
