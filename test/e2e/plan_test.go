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
	result := runArcpub(t, "plan", "--manifest", e2eManifest(root), "--version", "v0.1.0", "--output", "json")
	assertExitCode(t, result, 0)

	decoded := assertJSON(t, result.Stdout)
	if decoded["kind"] != "plan" {
		t.Fatalf("kind = %#v, want plan", decoded["kind"])
	}
	modules, ok := decoded["modules"].([]any)
	if !ok || len(modules) != 2 {
		t.Fatalf("modules = %#v, want two modules", decoded["modules"])
	}
	first, _ := modules[0].(map[string]any)
	second, _ := modules[1].(map[string]any)
	if first["name"] != "foundation" || second["name"] != "control" {
		t.Fatalf("module order = %#v, %#v; want foundation, control", first["name"], second["name"])
	}
	assertContains(t, result.Stdout, "requirements")
}

func TestPlanDoesNotLeakLocalPathsByDefault(t *testing.T) {
	root := copyFixture(t, "minimal")
	result := runArcpub(t, "plan", "--manifest", e2eManifest(root), "--version", "v0.1.0", "--output", "json")
	assertExitCode(t, result, 0)
	assertNoLocalPathLeak(t, result.Stdout, root)
}
