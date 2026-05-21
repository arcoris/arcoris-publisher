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

// Graph is an immutable-by-convention dependency topology snapshot.
//
// The graph contains all non-disabled modules from the registry. Public modules
// participate in PublishOrder. Internal modules participate in topology and
// diagnostics but are omitted from publication ordering.
type Graph struct {
	order []manifest.ModuleName

	nodes map[manifest.ModuleName]Node

	dependencies map[manifest.ModuleName][]manifest.ModuleName
	dependents   map[manifest.ModuleName][]manifest.ModuleName
}

// Len returns the number of graph nodes.
func (g Graph) Len() int { return len(g.order) }

// Empty reports whether the graph has no nodes.
func (g Graph) Empty() bool { return len(g.order) == 0 }

// Contains reports whether name is present as a graph node.
func (g Graph) Contains(name manifest.ModuleName) bool {
	_, ok := g.nodes[name]
	return ok
}

// Node returns the graph node for name.
func (g Graph) Node(name manifest.ModuleName) (Node, bool) {
	node, ok := g.nodes[name]
	return node, ok
}

// Nodes returns graph nodes in deterministic registry declaration order.
func (g Graph) Nodes() []Node {
	out := make([]Node, 0, len(g.order))
	for _, name := range g.order {
		out = append(out, g.nodes[name])
	}
	return out
}

// ModuleNames returns graph module names in deterministic registry declaration order.
func (g Graph) ModuleNames() []manifest.ModuleName {
	return cloneModuleNames(g.order)
}

// Edges returns dependency-to-dependent edges in deterministic order.
func (g Graph) Edges() []Edge {
	edges := make([]Edge, 0)
	for _, from := range g.order {
		for _, to := range g.dependents[from] {
			edges = append(edges, NewEdge(from, to))
		}
	}
	return edges
}

// cloneModuleNames detaches module-name slices before returning them to callers.
func cloneModuleNames(in []manifest.ModuleName) []manifest.ModuleName {
	out := make([]manifest.ModuleName, len(in))
	copy(out, in)
	return out
}
