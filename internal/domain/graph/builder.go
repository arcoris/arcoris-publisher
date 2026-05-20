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

// builder incrementally constructs a Graph while collecting validation issues.
//
// Keeping construction state here prevents the public Graph type from exposing
// partially built or invalid values.
type builder struct {
	modules []manifest.Module
	graph   Graph
	issues  []Issue
}

// newBuilder detaches input modules before construction begins.
func newBuilder(modules []manifest.Module) builder {
	copied := append([]manifest.Module(nil), modules...)
	return builder{
		modules: copied,
		graph: Graph{
			modules:      make(map[manifest.ModuleName]manifest.Module, len(copied)),
			order:        make([]manifest.ModuleName, 0, len(copied)),
			dependencies: make(map[manifest.ModuleName][]manifest.ModuleName, len(copied)),
			dependents:   make(map[manifest.ModuleName][]manifest.ModuleName, len(copied)),
		},
	}
}

// build creates the final graph or returns every issue found along the way.
func (b builder) build() (Graph, error) {
	b.indexModules()
	b.indexDependencies()
	b.sortAdjacency()
	b.detectCycles()
	if len(b.issues) > 0 {
		return Graph{}, &ValidationError{Issues: b.issues}
	}
	return b.graph, nil
}

// indexModules records each module by name and preserves declaration order.
func (b *builder) indexModules() {
	for index, module := range b.modules {
		name := module.Name()
		if _, exists := b.graph.modules[name]; exists {
			b.issues = append(b.issues, duplicateModuleIssue(name, index))
			continue
		}
		b.graph.modules[name] = module
		b.graph.order = append(b.graph.order, name)
		b.graph.dependencies[name] = nil
		b.graph.dependents[name] = nil
	}
}

// indexDependencies validates and records every direct module dependency.
func (b *builder) indexDependencies() {
	for _, module := range b.modules {
		moduleName := module.Name()
		if _, exists := b.graph.modules[moduleName]; !exists {
			continue
		}
		b.indexModuleDependencies(module)
	}
}

// indexModuleDependencies validates dependencies for one module.
func (b *builder) indexModuleDependencies(module manifest.Module) {
	moduleName := module.Name()
	seenDependencies := map[manifest.ModuleName]struct{}{}
	for _, dependency := range module.Dependencies() {
		dependencyName := dependency.Module()
		switch {
		case dependencyName == moduleName:
			b.issues = append(b.issues, selfDependencyIssue(moduleName))
			continue
		case hasSeenDependency(seenDependencies, dependencyName):
			b.issues = append(b.issues, duplicateDependencyIssue(moduleName, dependencyName))
			continue
		case !b.hasModule(dependencyName):
			b.issues = append(b.issues, unknownDependencyIssue(moduleName, dependencyName))
			continue
		}
		seenDependencies[dependencyName] = struct{}{}
		b.graph.dependencies[moduleName] = append(b.graph.dependencies[moduleName], dependencyName)
		b.graph.dependents[dependencyName] = append(b.graph.dependents[dependencyName], moduleName)
	}
}

// sortAdjacency normalizes every adjacency list after all edges are known.
func (b *builder) sortAdjacency() {
	for name := range b.graph.dependencies {
		b.graph.dependencies[name] = sortedModuleNames(b.graph.dependencies[name])
	}
	for name := range b.graph.dependents {
		b.graph.dependents[name] = sortedModuleNames(b.graph.dependents[name])
	}
}

// detectCycles appends graph cycle issues after adjacency lists are stable.
func (b *builder) detectCycles() {
	for _, cycle := range b.graph.Cycles() {
		b.issues = append(b.issues, cycleIssue(cycle))
	}
}

// hasModule reports whether a dependency target was declared successfully.
func (b *builder) hasModule(name manifest.ModuleName) bool {
	_, exists := b.graph.modules[name]
	return exists
}

// hasSeenDependency checks duplicate detection state.
func hasSeenDependency(seen map[manifest.ModuleName]struct{}, name manifest.ModuleName) bool {
	_, exists := seen[name]
	return exists
}
