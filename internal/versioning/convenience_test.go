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

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestAssignPublicationSet(t *testing.T) {
	set := mustPublicationSet(t, string(manifest.VersionPolicyReleaseTrain),
		testModule{name: "foundation"},
		testModule{name: "control", dependencies: []string{"foundation"}},
	)
	assignments, err := AssignPublicationSet(set, Must("v0.3.0"))
	if err != nil {
		t.Fatalf("AssignPublicationSet() error = %v", err)
	}
	assertModuleNames(t, assignments.ModuleNames(), "foundation", "control")
}

func TestRequirementMapFor(t *testing.T) {
	req := mustRequest(t, "v0.3.0", string(manifest.VersionPolicyReleaseTrain),
		testModule{name: "foundation"},
		testModule{name: "control", dependencies: []string{"foundation"}},
	)
	assignments, err := Assign(req)
	if err != nil {
		t.Fatalf("Assign() error = %v", err)
	}
	m, ok := assignments.RequirementMapFor(manifest.ModuleName("control"))
	if !ok {
		t.Fatalf("RequirementMapFor(control) not found")
	}
	if got := m[manifest.ModulePath("arcoris.dev/foundation")]; got != Version("v0.3.0") {
		t.Fatalf("requirement map foundation = %q", got)
	}
	m[manifest.ModulePath("arcoris.dev/foundation")] = Version("v9.9.9")
	m2, _ := assignments.RequirementMapFor(manifest.ModuleName("control"))
	if got := m2[manifest.ModulePath("arcoris.dev/foundation")]; got != Version("v0.3.0") {
		t.Fatalf("RequirementMapFor returned mutable map: %q", got)
	}
}
