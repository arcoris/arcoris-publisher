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
	"testing"

	"arcoris.dev/arcoris-publisher/internal/graph"
	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/registry"
)

func TestAssignVersionDefensiveBranches(t *testing.T) {
	set := mustPublicationSet(t, string(manifest.VersionPolicyReleaseTrain),
		testModule{name: "tooling", visibility: string(manifest.VisibilityInternal)},
	)
	reg := registry.Must(set)

	assigner := assigner{
		request: Request{
			Registry: reg,
			Version:  Must("v0.3.0"),
		},
		versions:     make(map[manifest.ModuleName]ModuleVersion),
		requirements: make(map[manifest.ModuleName][]Requirement),
	}

	assigner.assignVersion(0, manifest.ModuleName("missing"))
	assigner.assignVersion(1, manifest.ModuleName("tooling"))

	err := assigner.issues.err()
	assertValidationHas(t, err, IssueUnknownModule)
	assertValidationHas(t, err, IssueInvalidRequest)
}

func TestRequirementForDependencyDefensiveBranches(t *testing.T) {
	set := mustPublicationSet(t, string(manifest.VersionPolicyReleaseTrain),
		testModule{name: "foundation"},
		testModule{name: "tooling", visibility: string(manifest.VisibilityInternal)},
		testModule{name: "control"},
	)
	reg := registry.Must(set)

	assigner := assigner{
		request: Request{Registry: reg},
		versions: map[manifest.ModuleName]ModuleVersion{
			manifest.ModuleName("foundation"): newModuleVersion(
				manifest.ModuleName("foundation"),
				manifest.ModulePath("arcoris.dev/foundation"),
				Must("v0.3.0"),
			),
		},
	}

	_, ok := assigner.requirementForDependency(
		"dependencies[0]",
		manifest.ModuleName("missing"),
	)
	if ok {
		t.Fatal("missing dependency returned a requirement")
	}
	_, ok = assigner.requirementForDependency(
		"dependencies[1]",
		manifest.ModuleName("tooling"),
	)
	if ok {
		t.Fatal("internal dependency returned a requirement")
	}
	_, ok = assigner.requirementForDependency(
		"dependencies[2]",
		manifest.ModuleName("control"),
	)
	if ok {
		t.Fatal("unassigned dependency returned a requirement")
	}

	err := assigner.issues.err()
	assertValidationHas(t, err, IssueUnknownModule)
	assertValidationHas(t, err, IssueNonPublishableDependency)
	assertValidationHas(t, err, IssueMissingAssignment)
}

func TestAssignModuleRequirementsSkipsInvalidDependency(t *testing.T) {
	set := mustPublicationSet(t, string(manifest.VersionPolicyReleaseTrain),
		testModule{name: "foundation"},
		testModule{name: "control", dependencies: []string{"foundation"}},
	)
	reg := registry.Must(set)
	g := graph.Must(reg)

	assigner := assigner{
		request: Request{Graph: g},
		versions: map[manifest.ModuleName]ModuleVersion{
			manifest.ModuleName("control"): newModuleVersion(
				manifest.ModuleName("control"),
				manifest.ModulePath("arcoris.dev/control"),
				Must("v0.3.0"),
			),
		},
		requirements: make(map[manifest.ModuleName][]Requirement),
	}

	assigner.assignModuleRequirements(0, manifest.ModuleName("control"))

	if got := assigner.requirements[manifest.ModuleName("control")]; len(got) != 0 {
		t.Fatalf("requirements = %#v, want none", got)
	}
	assertValidationHas(t, assigner.issues.err(), IssueUnknownModule)
}
