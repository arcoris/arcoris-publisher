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

package module_test

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	modulemanifest "arcoris.dev/arcoris-publisher/internal/manifest/module"
)

func TestNewPublishAcceptsUniqueEntries(t *testing.T) {
	publish, err := modulemanifest.NewPublish(modulemanifest.PublishSpec{Entries: []manifest.PublishEntrySpec{
		{Type: "file", From: "go.mod", To: "go.mod"},
		{Type: "directory", From: "runtime", To: "runtime"},
	}})
	if err != nil {
		t.Fatalf("NewPublish returned error: %v", err)
	}
	if len(publish.Entries()) != 2 {
		t.Fatalf("unexpected entry count: %d", len(publish.Entries()))
	}
}

func TestNewPublishRejectsEmptyInvalidAndDuplicateEntries(t *testing.T) {
	for _, spec := range []modulemanifest.PublishSpec{
		{},
		{Entries: []manifest.PublishEntrySpec{{Type: "bad", From: "go.mod", To: "go.mod"}}},
		{Entries: []manifest.PublishEntrySpec{
			{Type: "file", From: "go.mod", To: "go.mod"},
			{Type: "file", From: "other.mod", To: "go.mod"},
		}},
	} {
		if _, err := modulemanifest.NewPublish(spec); err == nil {
			t.Fatalf("NewPublish(%#v) returned nil error", spec)
		}
	}
}

func TestPublishEntriesReturnsDetachedSlice(t *testing.T) {
	publish, err := modulemanifest.NewPublish(modulemanifest.PublishSpec{Entries: []manifest.PublishEntrySpec{{Type: "file", From: "go.mod", To: "go.mod"}}})
	if err != nil {
		t.Fatalf("NewPublish returned error: %v", err)
	}
	entries := publish.Entries()
	entries[0], _ = manifest.NewPublishEntry(manifest.PublishEntrySpec{Type: "file", From: "README.md", To: "README.md"})
	if publish.Entries()[0].From() == "README.md" {
		t.Fatalf("Entries accessor leaked internal slice")
	}
}
