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

func TestGraphBasicAccessors(t *testing.T) {
	g := mustGraph(t,
		testModule{name: "foundation"},
		testModule{name: "control", dependencies: []string{"foundation"}},
		testModule{name: "scheduler", dependencies: []string{"control"}},
	)
	if g.Len() != 3 {
		t.Fatalf("Len() = %d, want 3", g.Len())
	}
	if g.Empty() {
		t.Fatalf("Empty() = true")
	}
	if !g.Contains(manifest.ModuleName("control")) {
		t.Fatalf("Contains(control) = false")
	}
	node, ok := g.Node(manifest.ModuleName("control"))
	if !ok {
		t.Fatalf("Node(control) missing")
	}
	if node.Name() != manifest.ModuleName("control") ||
		node.ModulePath() != manifest.ModulePath("arcoris.dev/control") ||
		!node.Publishable() {
		t.Fatalf("unexpected node: %#v", node)
	}
	assertNames(t, g.ModuleNames(), "foundation", "control", "scheduler")
	edges := g.Edges()
	if len(edges) != 2 {
		t.Fatalf("len(Edges()) = %d, want 2", len(edges))
	}
	if edges[0].From() != manifest.ModuleName("foundation") ||
		edges[0].To() != manifest.ModuleName("control") {
		t.Fatalf("first edge = %s -> %s", edges[0].From(), edges[0].To())
	}

	nodes := g.Nodes()
	if len(nodes) != 3 || nodes[0].Name() != "foundation" {
		t.Fatalf("unexpected nodes: %#v", nodes)
	}
}

func TestGraphAccessorsReturnDetachedSlices(t *testing.T) {
	g := mustGraph(t,
		testModule{name: "foundation"},
		testModule{name: "control", dependencies: []string{"foundation"}},
	)
	names := g.ModuleNames()
	names[0] = manifest.ModuleName("mutated")
	assertNames(t, g.ModuleNames(), "foundation", "control")
	dependencies, ok := g.DirectDependencies(manifest.ModuleName("control"))
	if !ok {
		t.Fatalf("DirectDependencies(control) missing")
	}
	dependencies[0] = manifest.ModuleName("mutated")
	assertNames(t, mustDirectDependencies(t, g, "control"), "foundation")

	nodes := g.Nodes()
	nodes = nodes[:1]
	if len(g.Nodes()) != 2 {
		t.Fatal("Nodes returned aliased slice")
	}
}

func TestEmptyGraphAccessors(t *testing.T) {
	g := Graph{
		nodes:        map[manifest.ModuleName]Node{},
		dependencies: map[manifest.ModuleName][]manifest.ModuleName{},
		dependents:   map[manifest.ModuleName][]manifest.ModuleName{},
	}
	if !g.Empty() {
		t.Fatal("Empty() = false")
	}
	if g.Len() != 0 {
		t.Fatalf("Len() = %d, want 0", g.Len())
	}
	if _, ok := g.Node("missing"); ok {
		t.Fatal("Node(missing) ok = true")
	}
}

func TestNodeVisibilityPredicates(t *testing.T) {
	g := mustGraph(t,
		testModule{name: "foundation"},
		testModule{name: "tooling", visibility: string(manifest.VisibilityInternal)},
	)

	internal, ok := g.Node("tooling")
	if !ok {
		t.Fatal("tooling node missing")
	}
	if internal.Visibility() != manifest.VisibilityInternal {
		t.Fatalf("unexpected visibility: %s", internal.Visibility())
	}
	if !internal.Internal() {
		t.Fatal("Internal() = false")
	}
	if internal.Disabled() {
		t.Fatal("Disabled() = true")
	}

	disabledNode := Node{visibility: manifest.VisibilityDisabled}
	if !disabledNode.Disabled() {
		t.Fatal("disabled node Disabled() = false")
	}
}

func TestNewSkipsDisabledModules(t *testing.T) {
	reg := mustRegistry(t,
		testModule{name: "foundation"},
		testModule{name: "disabled", visibility: string(manifest.VisibilityDisabled)},
	)

	g, err := New(reg)
	if err != nil {
		t.Fatal(err)
	}
	if g.Contains("disabled") {
		t.Fatal("disabled module was indexed")
	}
	assertNames(t, g.ModuleNames(), "foundation")
}

func mustDirectDependencies(t *testing.T, g Graph, name string) []manifest.ModuleName {
	t.Helper()
	out, ok := g.DirectDependencies(manifest.ModuleName(name))
	if !ok {
		t.Fatalf("DirectDependencies(%s) missing", name)
	}
	return out
}
