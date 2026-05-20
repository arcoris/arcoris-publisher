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

func TestCyclesReturnsNormalizedCycle(t *testing.T) {
	graph := graphWithCycle()

	cycles := graph.Cycles()
	if len(cycles) != 1 {
		t.Fatalf("Cycles() len = %d, want 1: %#v", len(cycles), cycles)
	}
	if got := cycles[0].String(); got != "control -> foundation -> control" {
		t.Fatalf("cycle = %q, want normalized path", got)
	}
	if graph.Acyclic() {
		t.Fatalf("Acyclic() = true, want false")
	}
	if err := graph.Validate(); err == nil {
		t.Fatalf("Validate() error = nil, want cycle error")
	}
}

func TestCycleNodesReturnsDetachedCopy(t *testing.T) {
	cycle := NewCycle([]manifest.ModuleName{name("foundation"), name("control"), name("foundation")})
	nodes := cycle.Nodes()
	nodes[0] = name("other")
	assertNames(t, cycle.Nodes(), "foundation", "control", "foundation")
}

func TestValidateSucceedsForAcyclicGraph(t *testing.T) {
	graph := mustGraph(t, []manifest.ModuleSpec{
		moduleSpec("foundation"),
		moduleSpec("control", "foundation"),
	})
	if err := graph.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if !graph.Acyclic() {
		t.Fatalf("Acyclic() = false, want true")
	}
}

func graphWithCycle() Graph {
	foundation := mustModule(moduleSpec("foundation", "control"))
	control := mustModule(moduleSpec("control", "foundation"))
	return Graph{
		modules: map[manifest.ModuleName]manifest.Module{
			name("foundation"): foundation,
			name("control"):    control,
		},
		order: []manifest.ModuleName{name("foundation"), name("control")},
		dependencies: map[manifest.ModuleName][]manifest.ModuleName{
			name("foundation"): {name("control")},
			name("control"):    {name("foundation")},
		},
		dependents: map[manifest.ModuleName][]manifest.ModuleName{
			name("foundation"): {name("control")},
			name("control"):    {name("foundation")},
		},
	}
}

func mustModule(spec manifest.ModuleSpec) manifest.Module {
	module, err := manifest.NewModule(spec)
	if err != nil {
		panic(err)
	}
	return module
}
