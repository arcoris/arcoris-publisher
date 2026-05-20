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
	"testing"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

func TestGraphAccessorsReturnDetachedSlices(t *testing.T) {
	graph := mustGraph(t, []manifest.ModuleSpec{
		moduleSpec("foundation"),
		moduleSpec("control", "foundation"),
	})

	modules := graph.Modules()
	modules[0] = manifest.Module{}
	assertNames(t, graph.ModuleNames(), "foundation", "control")

	deps := graph.DependenciesOf(name("control"))
	deps[0] = name("other")
	assertNames(t, graph.DependenciesOf(name("control")), "foundation")
}

func TestModuleLookup(t *testing.T) {
	graph := mustGraph(t, []manifest.ModuleSpec{moduleSpec("foundation")})

	module, ok := graph.Module(name("foundation"))
	if !ok || module.Name() != name("foundation") {
		t.Fatalf("Module(foundation) = (%v, %v), want foundation true", module.Name(), ok)
	}
	if graph.Contains(name("missing")) {
		t.Fatalf("Contains(missing) = true, want false")
	}
}
