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

package provenance

import "testing"

func TestProjectionHashIsIndependentOfEntryOrder(t *testing.T) {
	first := ProjectionHash([]Entry{
		{TargetPath: "contracts/doc.go", Hash: "sha256:two", Present: true},
		{TargetPath: "go.mod", Hash: "sha256:one", Present: true},
	})
	second := ProjectionHash([]Entry{
		{TargetPath: "go.mod", Hash: "sha256:one", Present: true},
		{TargetPath: "contracts/doc.go", Hash: "sha256:two", Present: true},
	})

	if first != second {
		t.Fatalf("ProjectionHash() is order-dependent: %q != %q", first, second)
	}
}

func TestProjectionHashDistinguishesAbsentOptionalEntries(t *testing.T) {
	present := ProjectionHash([]Entry{
		{TargetPath: "optional.go", Hash: "", Present: true},
	})
	absent := ProjectionHash([]Entry{
		{TargetPath: "optional.go", Hash: "", Present: false},
	})

	if present == absent {
		t.Fatalf("ProjectionHash() did not distinguish present and absent entries: %q", present)
	}
}

func TestEntriesFromSourceModuleUsesRelativeTargetsOnly(t *testing.T) {
	input := testInput(t)

	entries := EntriesFromSourceModule(input.SourceModule)

	if len(entries) != 2 {
		t.Fatalf("len(entries) = %d", len(entries))
	}
	for _, entry := range entries {
		if entry.TargetPath == "" {
			t.Fatalf("empty target path in entries: %+v", entries)
		}
		if entry.TargetPath[0] == '/' {
			t.Fatalf("entry leaks absolute path: %+v", entry)
		}
		if !entry.Present {
			t.Fatalf("fixture entry unexpectedly absent: %+v", entry)
		}
	}
}
