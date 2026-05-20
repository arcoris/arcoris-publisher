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

// Edge is a dependency-order edge in the module graph.
//
// Dependency must be published before Dependent.
type Edge struct {
	Dependency manifest.ModuleName
	Dependent  manifest.ModuleName
}

// Edges returns dependency-order edges in deterministic order.
//
// Each edge means Dependency must be published before Dependent.
func (g Graph) Edges() []Edge {
	edges := make([]Edge, 0)
	for _, dependency := range g.order {
		for _, dependent := range g.dependents[dependency] {
			edges = append(edges, Edge{Dependency: dependency, Dependent: dependent})
		}
	}
	return edges
}
