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

func TestAssignmentLookups(t *testing.T) {
	assignments, err := ReleaseTrain(testRegistry(t), MustVersion("v0.4.0"))
	if err != nil {
		t.Fatalf("ReleaseTrain returned error: %v", err)
	}
	path := testModulePath(t, "arcoris.dev/control")

	version, ok := assignments.VersionOfPath(path)
	if !ok || version != "v0.4.0" {
		t.Fatalf("unexpected version of path: %q ok=%v", version, ok)
	}
	item, ok := assignments.ModuleVersion(testModuleName(t, "control"))
	if !ok || item.ModulePath() != path || item.Version() != "v0.4.0" {
		t.Fatalf("unexpected module version: %#v ok=%v", item, ok)
	}
	if !assignments.ContainsPath(path) {
		t.Fatal("expected path to be contained")
	}
}

func TestAssignmentLookupsMiss(t *testing.T) {
	assignments, err := ReleaseTrain(testRegistry(t), MustVersion("v0.4.0"))
	if err != nil {
		t.Fatalf("ReleaseTrain returned error: %v", err)
	}
	missingPath := testModulePath(t, "arcoris.dev/missing")

	if _, ok := assignments.VersionOfPath(missingPath); ok {
		t.Fatal("missing path was found")
	}
	if _, ok := assignments.ModuleVersion(testModuleName(t, "missing")); ok {
		t.Fatal("missing module was found")
	}
	if assignments.ContainsPath(missingPath) {
		t.Fatal("missing path contained")
	}
}
