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

// TransitiveDependenciesOf returns all transitive dependencies of module in
// dependency-first order.
func (g Graph) TransitiveDependenciesOf(module manifest.ModuleName) []manifest.ModuleName {
	seen := map[manifest.ModuleName]struct{}{}
	result := make([]manifest.ModuleName, 0)
	var visit func(manifest.ModuleName)
	visit = func(current manifest.ModuleName) {
		for _, dependency := range g.dependencies[current] {
			if _, ok := seen[dependency]; ok {
				continue
			}
			visit(dependency)
			seen[dependency] = struct{}{}
			result = append(result, dependency)
		}
	}
	visit(module)
	return result
}

// TransitiveDependentsOf returns all transitive dependents of module in
// deterministic declaration order.
func (g Graph) TransitiveDependentsOf(module manifest.ModuleName) []manifest.ModuleName {
	seen := map[manifest.ModuleName]struct{}{}
	var visit func(manifest.ModuleName)
	visit = func(current manifest.ModuleName) {
		for _, dependent := range g.dependents[current] {
			if _, ok := seen[dependent]; ok {
				continue
			}
			seen[dependent] = struct{}{}
			visit(dependent)
		}
	}
	visit(module)

	result := make([]manifest.ModuleName, 0, len(seen))
	for _, name := range g.order {
		if _, ok := seen[name]; ok {
			result = append(result, name)
		}
	}
	return result
}
