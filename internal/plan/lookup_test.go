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

package plan

import "testing"

func TestPlanLookup(t *testing.T) {
	p := mustPlan(t, "v0.3.0",
		testModule{name: "foundation"},
		testModule{name: "control", dependencies: []string{"foundation"}},
	)
	if _, ok := p.ModuleByName("control"); !ok {
		t.Fatal("ModuleByName(control) not found")
	}
	if _, ok := p.ModuleByPath("arcoris.dev/control"); !ok {
		t.Fatal("ModuleByPath(arcoris.dev/control) not found")
	}
	if _, ok := p.ModuleByRepository("arcoris/control"); !ok {
		t.Fatal("ModuleByRepository(arcoris/control) not found")
	}
	if !p.ContainsName("foundation") {
		t.Fatal("ContainsName(foundation) = false")
	}
	if !p.ContainsPath("arcoris.dev/foundation") {
		t.Fatal("ContainsPath(arcoris.dev/foundation) = false")
	}
	if !p.ContainsRepository("arcoris/foundation") {
		t.Fatal("ContainsRepository(arcoris/foundation) = false")
	}
	if _, ok := p.ModuleByName("missing"); ok {
		t.Fatal("ModuleByName(missing) found")
	}
	if _, ok := p.ModuleByPath("arcoris.dev/missing"); ok {
		t.Fatal("ModuleByPath(missing) found")
	}
	if _, ok := p.ModuleByRepository("arcoris/missing"); ok {
		t.Fatal("ModuleByRepository(missing) found")
	}
}
