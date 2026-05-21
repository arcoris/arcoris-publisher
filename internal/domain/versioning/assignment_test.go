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
	"time"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

func TestNewReleaseTrainAssignments(t *testing.T) {
	registryValue := testRegistry(t)
	assignments, err := New(registryValue, AssignmentSpec{Policy: manifest.VersionPolicyReleaseTrain, Release: "v0.3.0"})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	if assignments.Policy() != manifest.VersionPolicyReleaseTrain {
		t.Fatalf("unexpected policy: %q", assignments.Policy())
	}
	if assignments.Len() != 2 {
		t.Fatalf("unexpected assignment count: %d", assignments.Len())
	}
	foundation := testModuleName(t, "foundation")
	version, ok := assignments.VersionOfModule(foundation)
	if !ok || version != "v0.3.0" {
		t.Fatalf("unexpected foundation version: %q found=%v", version, ok)
	}
	if assignments.ContainsModule(testModuleName(t, "internal-tools")) {
		t.Fatal("internal module received a publish version")
	}
	items := assignments.Items()
	items[0] = ModuleVersion{}
	version, ok = assignments.VersionOfModule(foundation)
	if !ok || version != "v0.3.0" {
		t.Fatal("assignments were mutated through Items result")
	}
}

func TestNewSnapshotAssignments(t *testing.T) {
	registryValue := testRegistry(t)
	assignments, err := New(registryValue, AssignmentSpec{
		Policy: manifest.VersionPolicySnapshot,
		Snapshot: SnapshotSpec{
			Time:   time.Date(2026, 5, 21, 10, 0, 0, 0, time.UTC),
			Commit: "abcdef1234567890",
		},
	})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	version, ok := assignments.VersionOfModule(testModuleName(t, "control"))
	if !ok {
		t.Fatal("control version not found")
	}
	if version != "v0.0.0-20260521100000-abcdef123456" {
		t.Fatalf("unexpected snapshot version: %q", version)
	}
}

func TestNewRejectsInvalidAssignmentInput(t *testing.T) {
	registryValue := testRegistry(t)
	cases := []AssignmentSpec{
		{Policy: manifest.VersionPolicyReleaseTrain, Release: "bad"},
		{Policy: manifest.VersionPolicySnapshot},
		{Policy: manifest.VersionPolicy("unknown"), Release: "v0.1.0"},
	}
	for _, tc := range cases {
		if _, err := New(registryValue, tc); err == nil {
			t.Fatalf("New(%+v) succeeded, expected error", tc)
		}
	}
}

func TestMustPanicsOnInvalidAssignments(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("Must did not panic")
		}
	}()
	_ = Must(testRegistry(t), AssignmentSpec{Policy: manifest.VersionPolicyReleaseTrain, Release: "bad"})
}
