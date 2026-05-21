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

package manifest_test

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestNewPublishEntryAppliesDirectoryDefaults(t *testing.T) {
	entry, err := manifest.NewPublishEntry(manifest.PublishEntrySpec{
		Type: "directory",
		From: "runtime",
		To:   "runtime",
	})
	if err != nil {
		t.Fatalf("NewPublishEntry returned error: %v", err)
	}
	if entry.Kind() != manifest.PublishEntryDirectory || !entry.Recursive() || entry.Optional() {
		t.Fatalf("unexpected directory defaults: %#v", entry)
	}
	if entry.From().String() != "runtime" || entry.To().String() != "runtime" {
		t.Fatalf("unexpected entry paths")
	}
}

func TestNewPublishEntryAcceptsExplicitOptionalAndRecursiveFlags(t *testing.T) {
	entry, err := manifest.NewPublishEntry(manifest.PublishEntrySpec{
		Type:      "directory",
		From:      "contracts",
		To:        "contracts",
		Optional:  boolPtr(true),
		Recursive: boolPtr(false),
	})
	if err != nil {
		t.Fatalf("NewPublishEntry returned error: %v", err)
	}
	if !entry.Optional() || entry.Recursive() {
		t.Fatalf("explicit flags were not applied")
	}
	spec := entry.Spec()
	if spec.Optional == nil || !*spec.Optional || spec.Recursive == nil || *spec.Recursive {
		t.Fatalf("unexpected publish entry spec: %#v", spec)
	}
}

func TestNewPublishEntryRejectsInvalidEntries(t *testing.T) {
	for _, spec := range []manifest.PublishEntrySpec{
		{Type: "symlink", From: "go.mod", To: "go.mod"},
		{Type: "file", From: "../go.mod", To: "go.mod"},
		{Type: "file", From: "go.mod", To: "/go.mod"},
		{Type: "file", From: "go.mod", To: "go.mod", Recursive: boolPtr(true)},
	} {
		if _, err := manifest.NewPublishEntry(spec); err == nil {
			t.Fatalf("NewPublishEntry(%#v) returned nil error", spec)
		}
	}
}

func TestNewPublishEntryCollectsInvalidFields(t *testing.T) {
	_, err := manifest.NewPublishEntry(manifest.PublishEntrySpec{
		Type: "symlink",
		From: "../go.mod",
		To:   "/go.mod",
	})

	requireValidationIssuePaths(t, err, "type", "from", "to")
}
