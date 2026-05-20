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

// Len returns the number of modules indexed by the registry.
func (r Registry) Len() int { return len(r.modules) }

// Empty reports whether the registry contains no modules.
func (r Registry) Empty() bool { return len(r.modules) == 0 }

// Modules returns all modules in manifest declaration order.
//
// A detached slice is returned so callers can sort, truncate, or replace
// entries in the result without mutating the registry. The module values
// themselves are immutable by convention and expose their own detached slices.
func (r Registry) Modules() []manifest.Module {
	return cloneModules(r.modules)
}

// ModuleNames returns all module names in manifest declaration order.
func (r Registry) ModuleNames() []manifest.ModuleName {
	names := make([]manifest.ModuleName, 0, len(r.modules))
	for _, module := range r.modules {
		names = append(names, module.Name())
	}
	return names
}
