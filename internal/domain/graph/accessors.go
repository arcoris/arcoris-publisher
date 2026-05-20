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

package graph

import "arcoris.dev/arcoris-publisher/internal/domain/manifest"

// Len returns the number of modules in the graph.
func (g Graph) Len() int { return len(g.order) }

// Empty reports whether the graph contains no modules.
func (g Graph) Empty() bool { return len(g.order) == 0 }

// Modules returns modules in declaration order.
func (g Graph) Modules() []manifest.Module {
	modules := make([]manifest.Module, 0, len(g.order))
	for _, name := range g.order {
		modules = append(modules, g.modules[name])
	}
	return modules
}

// ModuleNames returns module names in declaration order.
func (g Graph) ModuleNames() []manifest.ModuleName {
	return append([]manifest.ModuleName(nil), g.order...)
}

// Module returns the module with name and whether it exists.
func (g Graph) Module(name manifest.ModuleName) (manifest.Module, bool) {
	module, ok := g.modules[name]
	return module, ok
}

// Contains reports whether the graph contains module name.
func (g Graph) Contains(name manifest.ModuleName) bool {
	_, ok := g.modules[name]
	return ok
}
