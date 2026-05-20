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

// ModuleByName returns the module with name and whether it was found.
//
// Module names are manifest-local stable identifiers. This lookup is the
// preferred path when another domain object already refers to a module by name.
func (r Registry) ModuleByName(name manifest.ModuleName) (manifest.Module, bool) {
	index, ok := r.byName[name]
	if !ok {
		return manifest.Module{}, false
	}
	return r.modules[index], true
}

// ContainsName reports whether a module with name exists.
func (r Registry) ContainsName(name manifest.ModuleName) bool {
	_, ok := r.byName[name]
	return ok
}
