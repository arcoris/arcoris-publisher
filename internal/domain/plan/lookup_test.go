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
	planValue := testPlan(t)
	if !planValue.ContainsModule(moduleName(t, "foundation")) {
		t.Fatal("foundation should be contained")
	}
	if !planValue.ContainsPath(modulePath(t, "arcoris.dev/control")) {
		t.Fatal("control path should be contained")
	}
	if !planValue.ContainsRepository(repositoryRef(t, "arcoris/control")) {
		t.Fatal("control repository should be contained")
	}
	if _, ok := planValue.ModulePlanByPath(modulePath(t, "arcoris.dev/control")); !ok {
		t.Fatal("ModulePlanByPath did not find control")
	}
	if _, ok := planValue.ModulePlanByRepository(repositoryRef(t, "arcoris/control")); !ok {
		t.Fatal("ModulePlanByRepository did not find control")
	}
	if _, ok := planValue.ModulePlan(moduleName(t, "missing")); ok {
		t.Fatal("missing module was found")
	}
	if _, ok := planValue.ModulePlanByPath(modulePath(t, "arcoris.dev/missing")); ok {
		t.Fatal("missing module path was found")
	}
	if _, ok := planValue.ModulePlanByRepository(repositoryRef(t, "arcoris/missing")); ok {
		t.Fatal("missing repository was found")
	}
}

func TestPlanIndexAccessors(t *testing.T) {
	planValue := testPlan(t)
	names := planValue.ModuleNames()
	if len(names) != 2 || names[0] != moduleName(t, "foundation") || names[1] != moduleName(t, "control") {
		t.Fatalf("unexpected names: %#v", names)
	}
	paths := planValue.ModulePaths()
	if len(paths) != 2 || paths[0] != modulePath(t, "arcoris.dev/foundation") || paths[1] != modulePath(t, "arcoris.dev/control") {
		t.Fatalf("unexpected paths: %#v", paths)
	}
	repositories := planValue.Repositories()
	if len(repositories) != 2 || repositories[0] != repositoryRef(t, "arcoris/foundation") || repositories[1] != repositoryRef(t, "arcoris/control") {
		t.Fatalf("unexpected repositories: %#v", repositories)
	}
}

func TestPlanAccessorsReturnDetachedSlices(t *testing.T) {
	planValue := testPlan(t)
	modules := planValue.Modules()
	modules[0] = ModulePlan{}
	if planValue.Modules()[0].Name() == "" {
		t.Fatal("Modules result mutated plan state")
	}
	skipped := planValue.SkippedModules()
	skipped[0] = SkippedModule{}
	if planValue.SkippedModules()[0].Name() == "" {
		t.Fatal("SkippedModules result mutated plan state")
	}
	names := planValue.ModuleNames()
	names[0] = moduleName(t, "control")
	if planValue.ModuleNames()[0] != moduleName(t, "foundation") {
		t.Fatal("ModuleNames result mutated plan state")
	}
}
