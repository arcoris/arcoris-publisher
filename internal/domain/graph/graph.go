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

// Graph is an immutable-by-convention module dependency graph.
//
// The graph keeps two orientations:
//   - DependenciesOf(module) returns modules that must be published before module.
//   - DependentsOf(module) returns modules that directly depend on module.
//
// Internally the topological adjacency is dependency -> dependent, so a
// successful TopologicalOrder result is already a valid dependency-first publish
// order.
type Graph struct {
	modules map[manifest.ModuleName]manifest.Module
	order   []manifest.ModuleName

	dependencies map[manifest.ModuleName][]manifest.ModuleName
	dependents   map[manifest.ModuleName][]manifest.ModuleName
}
