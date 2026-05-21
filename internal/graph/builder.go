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

import (
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/manifest/resolved"
	"arcoris.dev/arcoris-publisher/internal/registry"
)

// New builds a dependency graph from a registry.
//
// Graph construction rejects structurally invalid topology, such as duplicate
// nodes, unknown dependencies, dependencies on disabled modules, and direct
// self-dependencies. Cycles are represented by the graph and reported by
// TopologicalOrder, PublishOrder, HasCycle, and Cycles so callers can inspect
// the cycle before failing a higher-level planning operation.
func New(reg registry.Registry) (Graph, error) {
	builder := builder{
		registry: reg,
		issues:   newIssueCollector(),
	}
	return builder.build()
}

// FromPublicationSet builds a registry and then a graph from set.
func FromPublicationSet(set resolved.PublicationSet) (Graph, error) {
	reg, err := registry.New(set)
	if err != nil {
		return Graph{}, err
	}
	return New(reg)
}

// Must builds a graph and panics on invalid input.
//
// Must is intended for tests and static wiring. Runtime code should call New and
// return diagnostics to the caller.
func Must(reg registry.Registry) Graph {
	g, err := New(reg)
	if err != nil {
		panic(err)
	}
	return g
}

type builder struct {
	registry registry.Registry
	issues   issueCollector

	order        []manifest.ModuleName
	nodes        map[manifest.ModuleName]Node
	dependencies map[manifest.ModuleName][]manifest.ModuleName
	dependents   map[manifest.ModuleName][]manifest.ModuleName
}

// build creates graph indexes in two passes so dependency validation can check
// every edge against a complete node set.
func (b *builder) build() (Graph, error) {
	b.nodes = make(map[manifest.ModuleName]Node)
	b.dependencies = make(map[manifest.ModuleName][]manifest.ModuleName)
	b.dependents = make(map[manifest.ModuleName][]manifest.ModuleName)

	b.indexNodes()
	b.indexEdges()
	if err := b.issues.Err(); err != nil {
		return Graph{}, err
	}

	return Graph{
		order:        cloneModuleNames(b.order),
		nodes:        b.cloneNodes(),
		dependencies: cloneAdjacency(b.dependencies),
		dependents:   cloneAdjacency(b.dependents),
	}, nil
}

// indexNodes adds every non-disabled registry module as a graph node while
// preserving registry declaration order.
func (b *builder) indexNodes() {
	modules := b.registry.Modules()
	for i, module := range modules {
		if module.Visibility() == manifest.VisibilityDisabled {
			continue
		}

		name := module.Name()
		path := fmt.Sprintf("modules[%d]", i)
		if _, exists := b.nodes[name]; exists {
			b.addDuplicateNodeIssue(path+".name", name)
			continue
		}

		b.addNode(module)
	}
}

// indexEdges adds dependency-to-dependent edges for every non-disabled module.
func (b *builder) indexEdges() {
	modules := b.registry.Modules()
	for i, module := range modules {
		if module.Visibility() == manifest.VisibilityDisabled {
			continue
		}

		moduleName := module.Name()
		for j, dependency := range module.Dependencies() {
			path := fmt.Sprintf("modules[%d].dependencies[%d]", i, j)
			b.addDependency(path, moduleName, dependency)
		}
	}
}

// addDependency validates one declared dependency and stores both adjacency
// directions used by traversal and topological ordering.
func (b *builder) addDependency(
	path string,
	moduleName manifest.ModuleName,
	dependency manifest.ModuleName,
) {
	if dependency == moduleName {
		b.addSelfDependencyIssue(path, moduleName)
		return
	}

	depModule, ok := b.registry.ModuleByName(dependency)
	if !ok {
		b.addUnknownDependencyIssue(path, dependency)
		return
	}

	if depModule.Visibility() == manifest.VisibilityDisabled {
		b.addDisabledDependencyIssue(path, moduleName, dependency)
		return
	}

	if _, ok := b.nodes[dependency]; !ok {
		b.addMissingDependencyNodeIssue(path, dependency)
		return
	}

	b.addEdge(moduleName, dependency)
}

// addNode stores one non-disabled module in every node index.
func (b *builder) addNode(module resolved.PublicationModule) {
	name := module.Name()

	b.order = append(b.order, name)
	b.nodes[name] = Node{
		name:       name,
		modulePath: module.ModulePath(),
		visibility: module.Visibility(),
	}
	b.dependencies[name] = nil
	b.dependents[name] = nil
}

// addEdge stores one dependency edge in both traversal directions.
func (b *builder) addEdge(moduleName manifest.ModuleName, dependency manifest.ModuleName) {
	b.dependencies[moduleName] = append(b.dependencies[moduleName], dependency)
	b.dependents[dependency] = append(b.dependents[dependency], moduleName)
}

func (b *builder) addDuplicateNodeIssue(path string, name manifest.ModuleName) {
	b.issues.Add(IssueDuplicateNode, path, "duplicate graph node %q", name)
}

// addSelfDependencyIssue records a direct module-to-itself dependency.
func (b *builder) addSelfDependencyIssue(path string, moduleName manifest.ModuleName) {
	b.issues.Add(
		IssueSelfDependency,
		path,
		"module %q cannot depend on itself",
		moduleName,
	)
}

// addUnknownDependencyIssue records a dependency absent from the registry.
func (b *builder) addUnknownDependencyIssue(path string, dependency manifest.ModuleName) {
	b.issues.Add(IssueUnknownDependency, path, "unknown dependency %q", dependency)
}

// addDisabledDependencyIssue records a dependency on a disabled module.
func (b *builder) addDisabledDependencyIssue(
	path string,
	moduleName manifest.ModuleName,
	dependency manifest.ModuleName,
) {
	b.issues.Add(
		IssueDisabledDependency,
		path,
		"module %q depends on disabled module %q",
		moduleName,
		dependency,
	)
}

// addMissingDependencyNodeIssue records a registry module omitted from the graph.
func (b *builder) addMissingDependencyNodeIssue(
	path string,
	dependency manifest.ModuleName,
) {
	b.issues.Add(
		IssueUnknownDependency,
		path,
		"dependency %q is not present in graph",
		dependency,
	)
}

// cloneNodes detaches the graph's node index from the mutable builder.
func (b *builder) cloneNodes() map[manifest.ModuleName]Node {
	out := make(map[manifest.ModuleName]Node, len(b.nodes))
	for key, value := range b.nodes {
		out[key] = value
	}
	return out
}

// cloneAdjacency detaches adjacency lists before exposing the finished graph.
func cloneAdjacency(
	in map[manifest.ModuleName][]manifest.ModuleName,
) map[manifest.ModuleName][]manifest.ModuleName {
	out := make(map[manifest.ModuleName][]manifest.ModuleName, len(in))
	for key, value := range in {
		out[key] = cloneModuleNames(value)
	}
	return out
}
