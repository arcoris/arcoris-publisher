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

package modulefile

import (
	"bytes"
	"strings"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/testutil/publishertest"
)

func TestParseGoModModuleLine(t *testing.T) {
	path := parseModuleLine([]byte("module arcoris.dev/control\n"))
	if path.String() != "arcoris.dev/control" {
		t.Fatalf("unexpected module path %q", path)
	}
}

func TestRewriteGoModUpdatesSingleLineRequire(t *testing.T) {
	mod := modulePlan(t, "control")
	data := []byte(`module old.example/control

go 1.25

require arcoris.dev/foundation v0.1.0
require example.com/external v1.2.3
`)

	out, updates, changed, err := rewriteGoMod(data, mod, true)

	if err != nil {
		t.Fatalf("rewriteGoMod() error = %v", err)
	}
	if !changed {
		t.Fatal("changed = false")
	}
	assertContains(t, out, "module arcoris.dev/control")
	assertContains(t, out, "arcoris.dev/foundation v0.3.0")
	assertContains(t, out, "example.com/external v1.2.3")
	assertNotContains(t, out, "arcoris.dev/foundation v0.1.0")
	if len(updates) != 1 || updates[0].ModulePath() != "arcoris.dev/foundation" {
		t.Fatalf("updates = %#v", updates)
	}
}

func TestRewriteGoModUpdatesRequireBlockAndPreservesExternalComments(t *testing.T) {
	mod := modulePlan(t, "control")
	data := []byte(`module arcoris.dev/control

go 1.25

require (
	arcoris.dev/foundation v0.1.0 // old managed
	example.com/external v1.2.3 // keep me
)
`)

	out, _, _, err := rewriteGoMod(data, mod, true)

	if err != nil {
		t.Fatalf("rewriteGoMod() error = %v", err)
	}
	assertContains(t, out, "example.com/external v1.2.3 // keep me")
	assertContains(t, out, "arcoris.dev/foundation v0.3.0")
	assertNotContains(t, out, "old managed")
}

func TestRewriteGoModPreservesExternalIndirectRequirements(t *testing.T) {
	mod := modulePlan(t, "control")
	data := []byte(`module arcoris.dev/control

go 1.25

require (
	arcoris.dev/foundation v0.1.0 // indirect
	example.com/external v1.2.3 // indirect
)
`)

	out, _, _, err := rewriteGoMod(data, mod, true)

	if err != nil {
		t.Fatalf("rewriteGoMod() error = %v", err)
	}
	assertContains(t, out, "example.com/external v1.2.3 // indirect")
	assertContains(t, out, "arcoris.dev/foundation v0.3.0")
	assertNotContains(t, out, "arcoris.dev/foundation v0.3.0 // indirect")
}

func TestRewriteGoModRemovesManagedLocalReplaces(t *testing.T) {
	mod := modulePlan(t, "control")
	data := []byte(`module arcoris.dev/control

go 1.25

require arcoris.dev/foundation v0.1.0

replace arcoris.dev/foundation => ../foundation
replace example.com/external => ../external
replace example.com/remote => example.com/fork v1.2.3
`)

	out, _, _, err := rewriteGoMod(data, mod, true)

	if err != nil {
		t.Fatalf("rewriteGoMod() error = %v", err)
	}
	assertNotContains(t, out, "replace arcoris.dev/foundation")
	assertContains(t, out, "replace example.com/external => ../external")
	assertContains(t, out, "replace example.com/remote => example.com/fork v1.2.3")
}

func TestRewriteGoModRemovesVersionedManagedLocalReplace(t *testing.T) {
	mod := modulePlan(t, "control")
	data := []byte(`module arcoris.dev/control

go 1.25

require arcoris.dev/foundation v0.1.0

replace arcoris.dev/foundation v0.1.0 => ../foundation
`)

	out, _, _, err := rewriteGoMod(data, mod, true)

	if err != nil {
		t.Fatalf("rewriteGoMod() error = %v", err)
	}
	assertNotContains(t, out, "replace arcoris.dev/foundation")
}

func TestRewriteGoModPreservesManagedRemoteReplace(t *testing.T) {
	mod := modulePlan(t, "control")
	data := []byte(`module arcoris.dev/control

go 1.25

require arcoris.dev/foundation v0.1.0

replace arcoris.dev/foundation => arcoris.dev/foundation-fork v0.3.0
`)

	out, _, _, err := rewriteGoMod(data, mod, true)

	if err != nil {
		t.Fatalf("rewriteGoMod() error = %v", err)
	}
	assertContains(t, out, "replace arcoris.dev/foundation => arcoris.dev/foundation-fork v0.3.0")
}

func TestRewriteGoModAddsMissingModuleDirective(t *testing.T) {
	mod := modulePlan(t, "control")
	data := []byte("go 1.25\n")

	out, _, changed, err := rewriteGoMod(data, mod, true)

	if err != nil {
		t.Fatalf("rewriteGoMod() error = %v", err)
	}
	if !changed {
		t.Fatal("changed = false")
	}
	assertContains(t, out, "module arcoris.dev/control")
	assertContains(t, out, "go 1.25")
}

func TestRewriteGoModIsIdempotent(t *testing.T) {
	mod := modulePlan(t, "control")
	data := []byte(`module old.example/control

go 1.25

require arcoris.dev/foundation v0.1.0
`)

	once, _, _, err := rewriteGoMod(data, mod, true)
	if err != nil {
		t.Fatalf("first rewriteGoMod() error = %v", err)
	}
	twice, _, changed, err := rewriteGoMod(once, mod, true)
	if err != nil {
		t.Fatalf("second rewriteGoMod() error = %v", err)
	}
	if changed {
		t.Fatal("second changed = true")
	}
	if !bytes.Equal(once, twice) {
		t.Fatalf("second rewrite changed output:\n%s", twice)
	}
}

func TestRewriteGoModReportsNoChangeForCanonicalFile(t *testing.T) {
	mod := modulePlan(t, "control")
	data := []byte(`module arcoris.dev/control

go 1.25

require arcoris.dev/foundation v0.3.0
`)

	out, _, changed, err := rewriteGoMod(data, mod, true)

	if err != nil {
		t.Fatalf("rewriteGoMod() error = %v", err)
	}
	if changed {
		t.Fatalf("changed = true, output:\n%s", out)
	}
}

func modulePlan(t *testing.T, name string) plan.ModulePlan {
	t.Helper()

	p, err := publishertest.Plan(
		publishertest.PlanOptions{},
		publishertest.Module{Name: "foundation"},
		publishertest.Module{Name: "control", Dependencies: []string{"foundation"}},
	)
	if err != nil {
		t.Fatalf("publishertest.Plan() error = %v", err)
	}

	mod, ok := p.ModuleByName("control")
	if name != "control" {
		mod, ok = p.ModuleByName("foundation")
	}
	if !ok {
		t.Fatalf("ModuleByName(%q) not found", name)
	}

	return mod
}

func assertContains(t *testing.T, data []byte, want string) {
	t.Helper()
	if !strings.Contains(string(data), want) {
		t.Fatalf("output does not contain %q:\n%s", want, data)
	}
}

func assertNotContains(t *testing.T, data []byte, unwanted string) {
	t.Helper()
	if strings.Contains(string(data), unwanted) {
		t.Fatalf("output contains %q:\n%s", unwanted, data)
	}
}
