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

func TestModuleSnapshotAccessors(t *testing.T) {
	entry := EntrySnapshot{targetPath: manifest.RelativePath("go.mod")}
	module := ModuleSnapshot{
		name:          manifest.ModuleName("foundation"),
		sourceDir:     "/repo/staging/src/arcoris.dev/foundation",
		moduleRootDir: "/repo/staging/src/arcoris.dev/foundation/pkg",
		entries:       []EntrySnapshot{entry},
		hash:          Hash("sha256:module"),
	}

	if module.Name() != "foundation" {
		t.Fatalf("Name() = %q", module.Name())
	}
	if module.SourceDir() != "/repo/staging/src/arcoris.dev/foundation" {
		t.Fatalf("SourceDir() = %q", module.SourceDir())
	}
	if module.ModuleRootDir() != "/repo/staging/src/arcoris.dev/foundation/pkg" {
		t.Fatalf("ModuleRootDir() = %q", module.ModuleRootDir())
	}
	if module.Hash() != "sha256:module" {
		t.Fatalf("Hash() = %q", module.Hash())
	}
}

func TestModuleSnapshotEntriesAreDetached(t *testing.T) {
	module := ModuleSnapshot{entries: []EntrySnapshot{
		{targetPath: manifest.RelativePath("go.mod")},
		{targetPath: manifest.RelativePath("go.sum")},
	}}

	entries := module.Entries()
	entries[0] = entries[1]

	if module.Entries()[0].TargetPath() != "go.mod" {
		t.Fatal("Entries() returned attached slice")
	}
}
