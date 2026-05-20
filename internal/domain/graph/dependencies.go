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

// DependenciesOf returns direct dependencies of module in deterministic order.
func (g Graph) DependenciesOf(module manifest.ModuleName) []manifest.ModuleName {
	return append([]manifest.ModuleName(nil), g.dependencies[module]...)
}

// DependentsOf returns direct dependents of module in deterministic order.
func (g Graph) DependentsOf(module manifest.ModuleName) []manifest.ModuleName {
	return append([]manifest.ModuleName(nil), g.dependents[module]...)
}

// HasDependency reports whether module directly depends on dependency.
func (g Graph) HasDependency(module manifest.ModuleName, dependency manifest.ModuleName) bool {
	for _, candidate := range g.dependencies[module] {
		if candidate == dependency {
			return true
		}
	}
	return false
}
