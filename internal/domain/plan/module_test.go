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

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

func TestModulePlanAccessors(t *testing.T) {
	planValue := testPlan(t)
	control, ok := planValue.ModulePlan(moduleName(t, "control"))
	if !ok {
		t.Fatal("control plan not found")
	}
	if control.Name() != moduleName(t, "control") {
		t.Fatalf("Name() = %q", control.Name())
	}
	if control.ModulePath() != modulePath(t, "arcoris.dev/control") {
		t.Fatalf("ModulePath() = %q", control.ModulePath())
	}
	if control.SourceDir() == "" || control.Repository() != repositoryRef(t, "arcoris/control") {
		t.Fatalf("unexpected source/repository: %q %q", control.SourceDir(), control.Repository())
	}
	if control.Action() != ActionPublish {
		t.Fatalf("Action() = %q", control.Action())
	}
	if control.OrderIndex() != 1 {
		t.Fatalf("OrderIndex() = %d", control.OrderIndex())
	}
	if control.Version() != "v0.3.0" {
		t.Fatalf("Version() = %q", control.Version())
	}
	if control.Module().Name() != control.Name() {
		t.Fatalf("Module().Name() = %q", control.Module().Name())
	}
}

func TestBranchAndDependencyPlans(t *testing.T) {
	control, ok := testPlan(t).ModulePlan(moduleName(t, "control"))
	if !ok {
		t.Fatal("control plan not found")
	}
	branches := control.Branches()
	if len(branches) != 2 {
		t.Fatalf("len(Branches()) = %d, want 2", len(branches))
	}
	if branches[0].Source() != "main" || branches[0].Target() != "main" {
		t.Fatalf("unexpected first branch: %#v", branches[0])
	}
	dependencies := control.Dependencies()
	if len(dependencies) != 1 {
		t.Fatalf("len(Dependencies()) = %d, want 1", len(dependencies))
	}
	if dependencies[0].Module() != moduleName(t, "foundation") || dependencies[0].ModulePath() != modulePath(t, "arcoris.dev/foundation") || dependencies[0].Version() != "v0.3.0" {
		t.Fatalf("unexpected dependency: %#v", dependencies[0])
	}
	requirements := control.Requirements()
	if requirements[modulePath(t, "arcoris.dev/foundation")] != "v0.3.0" {
		t.Fatalf("unexpected requirements: %#v", requirements)
	}
}

func TestModulePlanAccessorsReturnDetachedValues(t *testing.T) {
	control, ok := testPlan(t).ModulePlan(moduleName(t, "control"))
	if !ok {
		t.Fatal("control plan not found")
	}
	branches := control.Branches()
	branches[0] = BranchPlan{}
	if control.Branches()[0].Source() == "" {
		t.Fatal("Branches result mutated module plan state")
	}
	dependencies := control.Dependencies()
	dependencies[0] = DependencyPlan{}
	if control.Dependencies()[0].Module() == "" {
		t.Fatal("Dependencies result mutated module plan state")
	}
	requirements := control.Requirements()
	requirements[modulePath(t, "arcoris.dev/foundation")] = "v9.9.9"
	if control.Requirements()[modulePath(t, "arcoris.dev/foundation")] != "v0.3.0" {
		t.Fatal("Requirements result mutated module plan state")
	}
}

func TestNewBranchPlanRejectsZeroMapping(t *testing.T) {
	if _, err := NewBranchPlan(manifest.BranchMapping{}); err == nil {
		t.Fatal("NewBranchPlan returned nil error for invalid mapping")
	}
}
