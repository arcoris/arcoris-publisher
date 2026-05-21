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

import "testing"

func TestAssignmentsAccessors(t *testing.T) {
	registryValue := testRegistry(t)
	assignments, err := ReleaseTrain(registryValue, MustVersion("v0.4.0"))
	if err != nil {
		t.Fatalf("ReleaseTrain returned error: %v", err)
	}
	if assignments.Empty() {
		t.Fatal("assignments unexpectedly empty")
	}
	modules := assignments.Modules()
	if len(modules) != 2 || modules[0] != testModuleName(t, "foundation") || modules[1] != testModuleName(t, "control") {
		t.Fatalf("unexpected module order: %#v", modules)
	}
	modules[0] = testModuleName(t, "control")
	if assignments.Modules()[0] != testModuleName(t, "foundation") {
		t.Fatal("Modules result mutated assignment state")
	}
}

func TestAssignmentItemsReturnDetachedSlice(t *testing.T) {
	assignments, err := ReleaseTrain(testRegistry(t), MustVersion("v0.4.0"))
	if err != nil {
		t.Fatalf("ReleaseTrain returned error: %v", err)
	}

	items := assignments.Items()
	items[0] = ModuleVersion{}

	version, ok := assignments.VersionOfModule(testModuleName(t, "foundation"))
	if !ok || version != "v0.4.0" {
		t.Fatal("assignments were mutated through Items result")
	}
}
