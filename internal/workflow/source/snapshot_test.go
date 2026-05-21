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

package source

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestSnapshotAccessorsReturnDetachedSlices(t *testing.T) {
	snap := Snapshot{
		modules: []ModuleSnapshot{
			{name: manifest.ModuleName("foundation")},
			{name: manifest.ModuleName("control")},
		},
		warnings: []Issue{{Code: IssueDirtySource}},
	}

	modules := snap.Modules()
	modules[0] = modules[1]
	if snap.Modules()[0].Name() != "foundation" {
		t.Fatal("Modules() returned attached slice")
	}

	warnings := snap.Warnings()
	warnings[0].Code = IssueDetachedHead
	if snap.Warnings()[0].Code != IssueDirtySource {
		t.Fatal("Warnings() returned attached slice")
	}
}

func TestSnapshotModuleNamesPreserveOrder(t *testing.T) {
	snap := Snapshot{modules: []ModuleSnapshot{
		{name: manifest.ModuleName("foundation")},
		{name: manifest.ModuleName("control")},
	}}

	names := snap.ModuleNames()

	if len(names) != 2 || names[0] != "foundation" || names[1] != "control" {
		t.Fatalf("ModuleNames() = %v", names)
	}
}
