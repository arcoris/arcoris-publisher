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

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestCycles(t *testing.T) {
	g := mustGraphWithCycle(t)
	if !g.HasCycle() {
		t.Fatalf("HasCycle() = false")
	}
	cycles := g.Cycles()
	if len(cycles) != 1 {
		t.Fatalf("len(Cycles()) = %d, want 1", len(cycles))
	}
	nodes := cycles[0].Nodes()
	if len(nodes) < 4 {
		t.Fatalf("cycle nodes too short: %v", nodes)
	}
	if nodes[0] != nodes[len(nodes)-1] {
		t.Fatalf("cycle should repeat first node as last: %v", nodes)
	}
	if cycles[0].String() == "" {
		t.Fatalf("Cycle.String() is empty")
	}
}

func TestCyclesReturnsNilForAcyclicGraph(t *testing.T) {
	g := mustGraph(t,
		testModule{name: "foundation"},
		testModule{name: "control", dependencies: []string{"foundation"}},
	)
	if g.HasCycle() {
		t.Fatal("HasCycle() = true")
	}
	if cycles := g.Cycles(); cycles != nil {
		t.Fatalf("Cycles() = %#v, want nil", cycles)
	}
}

func TestValidateAcyclic(t *testing.T) {
	g := mustGraph(t, testModule{name: "foundation"})
	if err := g.ValidateAcyclic(); err != nil {
		t.Fatalf("ValidateAcyclic() error = %v", err)
	}
	g = mustGraphWithCycle(t)
	if err := g.ValidateAcyclic(); err == nil {
		t.Fatalf("ValidateAcyclic() error = nil, want cycle error")
	}
}

func mustGraphWithCycle(t *testing.T) Graph {
	t.Helper()
	// Construct directly to test cycle diagnostics. The resolved publication model
	// rejects cycles earlier in normal end-to-end flows, but graph still supports
	// cycle discovery for defensive diagnostics and tests.
	return Graph{
		order: []manifest.ModuleName{"foundation", "control", "scheduler"},
		nodes: map[manifest.ModuleName]Node{
			"foundation": {
				name:       "foundation",
				modulePath: "arcoris.dev/foundation",
				visibility: manifest.VisibilityPublic,
			},
			"control": {
				name:       "control",
				modulePath: "arcoris.dev/control",
				visibility: manifest.VisibilityPublic,
			},
			"scheduler": {
				name:       "scheduler",
				modulePath: "arcoris.dev/scheduler",
				visibility: manifest.VisibilityPublic,
			},
		},
		dependencies: map[manifest.ModuleName][]manifest.ModuleName{
			"foundation": {"scheduler"},
			"control":    {"foundation"},
			"scheduler":  {"control"},
		},
		dependents: map[manifest.ModuleName][]manifest.ModuleName{
			"foundation": {"control"},
			"control":    {"scheduler"},
			"scheduler":  {"foundation"},
		},
	}
}
