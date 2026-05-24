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

package construct

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/testutil/publishertest"
)

func TestDetectProjectionCollisions(t *testing.T) {
	tests := []struct {
		name      string
		entries   []projectionEntry
		wantPath  string
		wantIssue string
	}{
		{
			name: "file same path",
			entries: []projectionEntry{
				fileProjectionEntry(0, "runtime/special.go"),
				fileProjectionEntry(1, "runtime/special.go"),
			},
			wantPath:  "runtime/special.go",
			wantIssue: "publish.entries[1].to",
		},
		{
			name: "directory contains file",
			entries: []projectionEntry{
				directoryProjectionEntry(0, "runtime"),
				fileProjectionEntry(1, "runtime/special.go"),
			},
			wantPath:  "runtime/special.go",
			wantIssue: "publish.entries[1].to",
		},
		{
			name: "directory same path",
			entries: []projectionEntry{
				directoryProjectionEntry(0, "runtime"),
				directoryProjectionEntry(1, "runtime"),
			},
			wantPath:  "runtime",
			wantIssue: "publish.entries[1].to",
		},
		{
			name: "later directory contains previous file",
			entries: []projectionEntry{
				fileProjectionEntry(0, "runtime/special.go"),
				directoryProjectionEntry(1, "runtime"),
			},
			wantPath:  "runtime",
			wantIssue: "publish.entries[1].to",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collisions := detectProjectionCollisions(tt.entries)

			if len(collisions) != 1 {
				t.Fatalf("len(collisions) = %d, want 1: %#v", len(collisions), collisions)
			}
			if collisions[0].current.path != tt.wantPath {
				t.Fatalf("current.path = %q, want %q", collisions[0].current.path, tt.wantPath)
			}
			if collisions[0].current.issuePath() != tt.wantIssue {
				t.Fatalf("issuePath() = %q, want %q", collisions[0].current.issuePath(), tt.wantIssue)
			}
		})
	}
}

func TestDetectProjectionCollisionsOrderingIsDeterministic(t *testing.T) {
	entries := []projectionEntry{
		directoryProjectionEntry(0, "runtime"),
		fileProjectionEntry(1, "runtime/a.go"),
		fileProjectionEntry(2, "runtime/b.go"),
	}

	collisions := detectProjectionCollisions(entries)

	if len(collisions) != 2 {
		t.Fatalf("len(collisions) = %d, want 2", len(collisions))
	}
	if collisions[0].current.path != "runtime/a.go" || collisions[1].current.path != "runtime/b.go" {
		t.Fatalf("collision order = %#v", collisions)
	}
}

func TestProvenanceProjectionCollidesWithExplicitEntry(t *testing.T) {
	provenanceFile := "ARCPUB.json"
	plan, err := publishertest.Plan(
		publishertest.PlanOptions{
			Publish: manifest.PublishSpec{
				Provenance: manifest.ProvenanceSpec{File: &provenanceFile},
			},
		},
		publishertest.Module{
			Name: "foundation",
			Entries: []manifest.PublishEntrySpec{
				publishertest.FileEntry("ARCPUB.json"),
			},
		},
	)
	if err != nil {
		t.Fatalf("publishertest.Plan() error = %v", err)
	}

	mod, _ := plan.ModuleByName("foundation")
	entries := appendProvenanceProjectionEntry(
		projectionEntries(mod),
		Request{Plan: plan},
		Options{GenerateProvenanceFile: true},
	)
	collisions := detectProjectionCollisions(entries)

	if len(collisions) != 1 {
		t.Fatalf("len(collisions) = %d, want 1: %#v", len(collisions), collisions)
	}
	if collisions[0].current.issuePath() != "publish.provenance.file" {
		t.Fatalf("issuePath() = %q", collisions[0].current.issuePath())
	}
}

func fileProjectionEntry(index int, path string) projectionEntry {
	return projectionEntry{
		index: index,
		kind:  manifest.PublishEntryFile,
		path:  cleanProjectionPath(path),
	}
}

func directoryProjectionEntry(index int, path string) projectionEntry {
	return projectionEntry{
		index: index,
		kind:  manifest.PublishEntryDirectory,
		path:  cleanProjectionPath(path),
	}
}
