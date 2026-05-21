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

import "arcoris.dev/arcoris-publisher/internal/manifest"

// DirectDependencies returns modules that name directly depends on.
func (g Graph) DirectDependencies(name manifest.ModuleName) ([]manifest.ModuleName, bool) {
	dependencies, ok := g.dependencies[name]
	if !ok {
		return nil, false
	}
	return cloneModuleNames(dependencies), true
}

// DirectDependents returns modules that directly depend on name.
func (g Graph) DirectDependents(name manifest.ModuleName) ([]manifest.ModuleName, bool) {
	dependents, ok := g.dependents[name]
	if !ok {
		return nil, false
	}
	return cloneModuleNames(dependents), true
}

// TransitiveDependencies returns every dependency required by name.
//
// The result is in global topological order, so dependencies appear before the
// modules that depend on them. The requested module is not included.
func (g Graph) TransitiveDependencies(name manifest.ModuleName) ([]manifest.ModuleName, bool) {
	if !g.Contains(name) {
		return nil, false
	}
	seen := make(map[manifest.ModuleName]struct{})
	g.collectDependencies(name, seen)
	return g.filterByTopologicalOrder(seen), true
}

// TransitiveDependents returns every module that depends on name.
//
// The result is in global topological order. The requested module is not included.
func (g Graph) TransitiveDependents(name manifest.ModuleName) ([]manifest.ModuleName, bool) {
	if !g.Contains(name) {
		return nil, false
	}
	seen := make(map[manifest.ModuleName]struct{})
	g.collectDependents(name, seen)
	return g.filterByTopologicalOrder(seen), true
}

func (g Graph) collectDependencies(
	name manifest.ModuleName,
	seen map[manifest.ModuleName]struct{},
) {
	for _, dependency := range g.dependencies[name] {
		if _, ok := seen[dependency]; ok {
			continue
		}
		seen[dependency] = struct{}{}
		g.collectDependencies(dependency, seen)
	}
}

func (g Graph) collectDependents(name manifest.ModuleName, seen map[manifest.ModuleName]struct{}) {
	for _, dependent := range g.dependents[name] {
		if _, ok := seen[dependent]; ok {
			continue
		}
		seen[dependent] = struct{}{}
		g.collectDependents(dependent, seen)
	}
}

func (g Graph) filterByTopologicalOrder(
	set map[manifest.ModuleName]struct{},
) []manifest.ModuleName {
	order, err := g.TopologicalOrder()
	if err != nil {
		order = g.order
	}
	out := make([]manifest.ModuleName, 0, len(set))
	for _, name := range order {
		if _, ok := set[name]; !ok {
			continue
		}
		out = append(out, name)
	}
	return out
}
