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

func TestDirectTraversal(t *testing.T) {
	g := mustGraph(t,
		testModule{name: "foundation"},
		testModule{name: "control", dependencies: []string{"foundation"}},
		testModule{name: "scheduler", dependencies: []string{"control"}},
	)
	dependencies, ok := g.DirectDependencies(manifest.ModuleName("control"))
	if !ok {
		t.Fatalf("DirectDependencies(control) missing")
	}
	assertNames(t, dependencies, "foundation")
	dependents, ok := g.DirectDependents(manifest.ModuleName("control"))
	if !ok {
		t.Fatalf("DirectDependents(control) missing")
	}
	assertNames(t, dependents, "scheduler")
	if _, ok := g.DirectDependencies(manifest.ModuleName("missing")); ok {
		t.Fatalf("DirectDependencies(missing) ok = true")
	}
	if _, ok := g.DirectDependents(manifest.ModuleName("missing")); ok {
		t.Fatalf("DirectDependents(missing) ok = true")
	}
}

func TestTransitiveTraversal(t *testing.T) {
	g := mustGraph(t,
		testModule{name: "foundation"},
		testModule{name: "control", dependencies: []string{"foundation"}},
		testModule{name: "runtime", dependencies: []string{"foundation"}},
		testModule{name: "scheduler", dependencies: []string{"control", "runtime"}},
	)
	dependencies, ok := g.TransitiveDependencies(manifest.ModuleName("scheduler"))
	if !ok {
		t.Fatalf("TransitiveDependencies(scheduler) missing")
	}
	assertNames(t, dependencies, "foundation", "control", "runtime")
	dependents, ok := g.TransitiveDependents(manifest.ModuleName("foundation"))
	if !ok {
		t.Fatalf("TransitiveDependents(foundation) missing")
	}
	assertNames(t, dependents, "control", "runtime", "scheduler")
	if _, ok := g.TransitiveDependents(manifest.ModuleName("missing")); ok {
		t.Fatalf("TransitiveDependents(missing) ok = true")
	}
	if _, ok := g.TransitiveDependencies(manifest.ModuleName("missing")); ok {
		t.Fatalf("TransitiveDependencies(missing) ok = true")
	}
}

func TestTransitiveTraversalFallsBackToDeclarationOrderOnCycle(t *testing.T) {
	g := mustGraphWithCycle(t)
	dependencies, ok := g.TransitiveDependencies("scheduler")
	if !ok {
		t.Fatal("TransitiveDependencies(scheduler) missing")
	}
	assertNames(t, dependencies, "foundation", "control", "scheduler")
}
