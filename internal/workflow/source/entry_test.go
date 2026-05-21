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

func TestEntrySnapshotAccessors(t *testing.T) {
	entry := mustPublishEntry(t, manifest.PublishEntrySpec{
		Type: string(manifest.PublishEntryFile),
		From: "go.mod",
		To:   "module/go.mod",
	})
	snap := EntrySnapshot{
		entry:      entry,
		sourcePath: "/repo/staging/src/arcoris.dev/foundation/go.mod",
		targetPath: entry.To(),
		present:    true,
		hash:       Hash("sha256:file"),
	}

	if snap.Entry() != entry {
		t.Fatal("Entry() did not return original entry")
	}
	if snap.Kind() != manifest.PublishEntryFile {
		t.Fatalf("Kind() = %q", snap.Kind())
	}
	if snap.SourcePath() != "/repo/staging/src/arcoris.dev/foundation/go.mod" {
		t.Fatalf("SourcePath() = %q", snap.SourcePath())
	}
	if snap.TargetPath() != "module/go.mod" {
		t.Fatalf("TargetPath() = %q", snap.TargetPath())
	}
	if snap.Optional() {
		t.Fatal("Optional() = true")
	}
	if !snap.Present() {
		t.Fatal("Present() = false")
	}
	if snap.Hash() != "sha256:file" {
		t.Fatalf("Hash() = %q", snap.Hash())
	}
}
