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

// New builds a dependency graph from manifest modules.
//
// The input modules are copied into builder state before validation so caller
// mutations after construction cannot affect the resulting graph.
func New(modules []manifest.Module) (Graph, error) {
	builder := newBuilder(modules)
	return builder.build()
}

// FromManifest builds a dependency graph from all modules declared in manifest.
func FromManifest(manifestValue manifest.Manifest) (Graph, error) {
	return New(manifestValue.Modules())
}

// Must constructs a graph and panics when modules are invalid.
//
// Must is intended for tests and static wiring. Runtime code should call New and
// return diagnostics to the caller.
func Must(modules []manifest.Module) Graph {
	graph, err := New(modules)
	if err != nil {
		panic(err)
	}
	return graph
}
