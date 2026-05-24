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

package e2e_test

import "testing"

func TestPlanTextMinimalFixture(t *testing.T) {
	root := copyFixture(t, "minimal")
	result := runArcpub(t, "plan", "--manifest", e2eManifest(root), "--version", "v0.1.0")
	assertExitCode(t, result, 0)
	for _, want := range []string{
		"foundation",
		"control",
		"v0.1.0",
		"arcoris.dev/foundation",
		"arcoris.dev/control",
	} {
		assertContains(t, result.Stdout, want)
	}
}

func TestPlanJSONMinimalFixture(t *testing.T) {
	root := copyFixture(t, "minimal")
	result, decoded := runArcpubJSON(t, 0, "plan", "--manifest", e2eManifest(root), "--version", "v0.1.0", "--output", "json")

	if decoded["kind"] != "plan" {
		t.Fatalf("kind = %#v, want plan", decoded["kind"])
	}
	if got := floatField(t, decoded, "moduleCount"); got != 2 {
		t.Fatalf("moduleCount = %v, want 2", got)
	}

	modules := arrayField(t, decoded, "modules")
	if len(modules) != 2 {
		t.Fatalf("module count = %d, want 2", len(modules))
	}
	foundation, _ := modules[0].(map[string]any)
	control, _ := modules[1].(map[string]any)
	assertPlanModule(t, foundation, "foundation", "arcoris.dev/foundation", "arcoris/foundation")
	assertPlanModule(t, control, "control", "arcoris.dev/control", "arcoris/control")
	assertPlanEntryTargets(t, foundation, "go.mod", "README.md", "contracts")
	assertPlanEntryTargets(t, control, "go.mod", "README.md", "contracts")
	assertControlRequirement(t, control)
	assertNoLocalPathLeak(t, result.Stdout, root)
}

func TestPlanDoesNotLeakLocalPathsByDefault(t *testing.T) {
	root := copyFixture(t, "minimal")
	result, _ := runArcpubJSON(t, 0, "plan", "--manifest", e2eManifest(root), "--version", "v0.1.0", "--output", "json")
	assertNoLocalPathLeak(t, result.Stdout, root)
}

func assertPlanModule(
	t *testing.T,
	module map[string]any,
	name string,
	modulePath string,
	repository string,
) {
	t.Helper()
	if stringField(t, module, "name") != name ||
		stringField(t, module, "modulePath") != modulePath ||
		stringField(t, module, "repository") != repository {
		t.Fatalf("unexpected module report: %#v", module)
	}
	sourceDir := stringField(t, module, "sourceDir")
	if sourceDir == "" || sourceDir[0] == '/' {
		t.Fatalf("sourceDir = %q, want relative path", sourceDir)
	}
}

func assertPlanEntryTargets(t *testing.T, module map[string]any, targets ...string) {
	t.Helper()
	entries := arrayField(t, module, "publishEntries")
	got := map[string]bool{}
	for _, raw := range entries {
		entry, ok := raw.(map[string]any)
		if !ok {
			t.Fatalf("publish entry = %#v, want object", raw)
		}
		got[stringField(t, entry, "to")] = true
	}
	for _, target := range targets {
		if !got[target] {
			t.Fatalf("module %s missing publish entry target %q in %#v", module["name"], target, entries)
		}
	}
}

func assertControlRequirement(t *testing.T, module map[string]any) {
	t.Helper()
	requirements := arrayField(t, module, "requirements")
	if len(requirements) != 1 {
		t.Fatalf("control requirements = %#v, want one requirement", requirements)
	}
	requirement, ok := requirements[0].(map[string]any)
	if !ok {
		t.Fatalf("requirement = %#v, want object", requirements[0])
	}
	if stringField(t, requirement, "module") != "foundation" ||
		stringField(t, requirement, "modulePath") != "arcoris.dev/foundation" ||
		stringField(t, requirement, "version") != "v0.1.0" {
		t.Fatalf("unexpected control requirement: %#v", requirement)
	}
}
