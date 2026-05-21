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

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

func TestNewModuleVersionValidation(t *testing.T) {
	name := testModuleName(t, "foundation")
	path, err := manifest.ParseModulePath("arcoris.dev/foundation")
	if err != nil {
		t.Fatal(err)
	}
	item, err := NewModuleVersion(name, path, MustVersion("v0.1.0"))
	if err != nil {
		t.Fatalf("NewModuleVersion returned error: %v", err)
	}
	if item.Module() != name || item.ModulePath() != path || item.Version() != "v0.1.0" {
		t.Fatalf("unexpected module version: %#v", item)
	}
	if _, err := NewModuleVersion("", path, MustVersion("v0.1.0")); err == nil {
		t.Fatal("NewModuleVersion with empty module succeeded")
	}
	if _, err := NewModuleVersion(name, "", MustVersion("v0.1.0")); err == nil {
		t.Fatal("NewModuleVersion with empty path succeeded")
	}
	if _, err := NewModuleVersion(name, path, ""); err == nil {
		t.Fatal("NewModuleVersion with empty version succeeded")
	}
}

func TestAssignmentsValidateRejectsInvalidState(t *testing.T) {
	path, err := manifest.ParseModulePath("arcoris.dev/foundation")
	if err != nil {
		t.Fatal(err)
	}
	items := []ModuleVersion{
		{module: testModuleName(t, "foundation"), modulePath: path, version: MustVersion("v0.1.0")},
		{module: testModuleName(t, "foundation"), modulePath: path, version: MustVersion("v0.1.0")},
		{},
	}
	assignments := Assignments{policy: manifest.VersionPolicyReleaseTrain, items: items}
	if err := assignments.Validate(); err == nil {
		t.Fatal("Validate succeeded, expected error")
	}
}

func TestAssignmentsValidateRejectsNilReceiver(t *testing.T) {
	var assignments *Assignments
	if err := assignments.Validate(); err == nil {
		t.Fatal("nil Validate succeeded")
	}
}

func TestAssignmentsValidateRejectsMissingPolicy(t *testing.T) {
	path, err := manifest.ParseModulePath("arcoris.dev/foundation")
	if err != nil {
		t.Fatal(err)
	}
	assignments := Assignments{items: []ModuleVersion{{module: testModuleName(t, "foundation"), modulePath: path, version: MustVersion("v0.1.0")}}}
	if err := assignments.Validate(); err == nil {
		t.Fatal("Validate succeeded without policy")
	}
}

func TestReleaseTrainRejectsEmptyVersion(t *testing.T) {
	if _, err := ReleaseTrain(testRegistry(t), ""); err == nil {
		t.Fatal("ReleaseTrain succeeded with empty version")
	}
}
