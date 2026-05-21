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

package versioning

import (
	"errors"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestAssignReleaseTrain(t *testing.T) {
	req := mustRequest(t, "v0.3.0", string(manifest.VersionPolicyReleaseTrain),
		testModule{name: "foundation"},
		testModule{name: "control", dependencies: []string{"foundation"}},
	)
	assignments, err := Assign(req)
	if err != nil {
		t.Fatalf("Assign() error = %v", err)
	}
	assertModuleNames(t, assignments.ModuleNames(), "foundation", "control")
	version, ok := assignments.VersionOf(manifest.ModuleName("control"))
	if !ok || version != Version("v0.3.0") {
		t.Fatalf("VersionOf(control) = %q, %v", version, ok)
	}
	control, ok := assignments.ModuleVersionByPath(manifest.ModulePath("arcoris.dev/control"))
	if !ok {
		t.Fatalf("ModuleVersionByPath(control) = %#v, %v", control, ok)
	}
	assertModuleVersion(t, control, "control", "arcoris.dev/control", "v0.3.0")
}

func TestAssignDirectRequirements(t *testing.T) {
	req := mustRequest(t, "v0.3.0", string(manifest.VersionPolicyReleaseTrain),
		testModule{name: "foundation"},
		testModule{name: "control", dependencies: []string{"foundation"}},
		testModule{name: "scheduler", dependencies: []string{"control"}},
	)
	assignments, err := Assign(req)
	if err != nil {
		t.Fatalf("Assign() error = %v", err)
	}
	reqs, ok := assignments.RequirementsFor(manifest.ModuleName("scheduler"))
	if !ok {
		t.Fatalf("RequirementsFor(scheduler) not found")
	}
	if len(reqs) != 1 {
		t.Fatalf("len(requirements) = %d, want 1", len(reqs))
	}
	assertRequirement(t, reqs[0], "control", "arcoris.dev/control", "v0.3.0")
}

func TestAssignSkipsInternalModules(t *testing.T) {
	req := mustRequest(t, "v0.3.0", string(manifest.VersionPolicyReleaseTrain),
		testModule{name: "foundation"},
		testModule{name: "tooling", visibility: string(manifest.VisibilityInternal)},
	)
	assignments, err := Assign(req)
	if err != nil {
		t.Fatalf("Assign() error = %v", err)
	}
	assertModuleNames(t, assignments.ModuleNames(), "foundation")
	if _, ok := assignments.VersionOf(manifest.ModuleName("tooling")); ok {
		t.Fatalf("internal module received assignment")
	}
}

func TestAssignSnapshot(t *testing.T) {
	req := mustRequest(t,
		"v0.0.0-20260521143000-abcdefabcdef",
		string(manifest.VersionPolicySnapshot),
		testModule{name: "foundation"},
	)
	assignments, err := Assign(req)
	if err != nil {
		t.Fatalf("Assign() error = %v", err)
	}
	version, ok := assignments.VersionOf(manifest.ModuleName("foundation"))
	if !ok || !version.IsPseudo() {
		t.Fatalf("snapshot assignment = %q, %v", version, ok)
	}
}

func TestAssignRejectsWrongVersionForPolicy(t *testing.T) {
	req := mustRequest(t,
		"v0.0.0-20260521143000-abcdefabcdef",
		string(manifest.VersionPolicyReleaseTrain),
		testModule{name: "foundation"},
	)
	_, err := Assign(req)
	var validation *ValidationError
	if !errors.As(err, &validation) || !validation.Has(IssueInvalidVersion) {
		t.Fatalf("Assign() error = %v, want invalid version", err)
	}

	req = mustRequest(t, "v0.3.0", string(manifest.VersionPolicySnapshot),
		testModule{name: "foundation"},
	)
	_, err = Assign(req)
	if !errors.As(err, &validation) || !validation.Has(IssueInvalidVersion) {
		t.Fatalf("Assign() error = %v, want invalid version", err)
	}
}

func TestAssignReportsGraphOrderCycle(t *testing.T) {
	req := mustRequest(t, "v0.3.0", string(manifest.VersionPolicyReleaseTrain),
		testModule{name: "foundation", dependencies: []string{"control"}},
		testModule{name: "control", dependencies: []string{"foundation"}},
	)
	_, err := Assign(req)
	var validation *ValidationError
	if !errors.As(err, &validation) || !validation.Has(IssueGraphOrder) {
		t.Fatalf("Assign() error = %v, want graph order validation", err)
	}
}

func TestAssignmentsReturnDetachedSlices(t *testing.T) {
	req := mustRequest(t, "v0.3.0", string(manifest.VersionPolicyReleaseTrain),
		testModule{name: "foundation"},
		testModule{name: "control", dependencies: []string{"foundation"}},
	)
	assignments, err := Assign(req)
	if err != nil {
		t.Fatalf("Assign() error = %v", err)
	}
	mods := assignments.Modules()
	mods[0] = ModuleVersion{}
	got, ok := assignments.VersionOf(manifest.ModuleName("foundation"))
	if !ok || got != Version("v0.3.0") {
		t.Fatalf("assignment mutated through Modules(): %q %v", got, ok)
	}
	reqs, ok := assignments.RequirementsFor(manifest.ModuleName("control"))
	if !ok || len(reqs) != 1 {
		t.Fatalf("requirements missing")
	}
	reqs[0] = Requirement{}
	reqs2, _ := assignments.RequirementsFor(manifest.ModuleName("control"))
	if reqs2[0].Module() != manifest.ModuleName("foundation") {
		t.Fatalf("requirements mutated through accessor")
	}
}
