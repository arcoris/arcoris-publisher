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

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBuildFilePayloadUsesResolvedRuntimeValues(t *testing.T) {
	input := testInput(t)

	payload := BuildFilePayload(input)

	if payload.SchemaVersion != SchemaVersion {
		t.Fatalf("SchemaVersion = %q", payload.SchemaVersion)
	}
	if payload.Publisher.Version != "v9.8.7" {
		t.Fatalf("Publisher.Version = %q", payload.Publisher.Version)
	}
	if payload.Module.Name != "foundation" {
		t.Fatalf("Module.Name = %q", payload.Module.Name)
	}
	if payload.Module.ModulePath != "arcoris.dev/foundation" {
		t.Fatalf("Module.ModulePath = %q", payload.Module.ModulePath)
	}
	if payload.Source.Repository != "arcoris/arcoris" {
		t.Fatalf("Source.Repository = %q", payload.Source.Repository)
	}
	if payload.Source.Commit != "abcdef1234567890" {
		t.Fatalf("Source.Commit = %q", payload.Source.Commit)
	}
	if payload.Target.Repository != "arcoris/foundation" {
		t.Fatalf("Target.Repository = %q", payload.Target.Repository)
	}
	if payload.Publication.PublishMode != "explicit-projection" {
		t.Fatalf("Publication.PublishMode = %q", payload.Publication.PublishMode)
	}
	if payload.Projection.EntryCount != 2 {
		t.Fatalf("Projection.EntryCount = %d", payload.Projection.EntryCount)
	}
	if !strings.HasPrefix(payload.Projection.ProjectionHash, "sha256:") {
		t.Fatalf("Projection.ProjectionHash = %q", payload.Projection.ProjectionHash)
	}
}

func TestRenderFilePayloadIsDeterministicJSON(t *testing.T) {
	input := testInput(t)

	first, err := RenderFilePayload(input)
	if err != nil {
		t.Fatalf("RenderFilePayload() error = %v", err)
	}
	second, err := RenderFilePayload(input)
	if err != nil {
		t.Fatalf("RenderFilePayload() second error = %v", err)
	}

	if string(first) != string(second) {
		t.Fatalf("payload is not deterministic:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
	if first[len(first)-1] != '\n' {
		t.Fatalf("payload missing trailing newline:\n%s", first)
	}

	var payload FilePayload
	if err := json.Unmarshal(first, &payload); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
}

func TestRenderFilePayloadDoesNotLeakLocalPaths(t *testing.T) {
	input := testInput(t)

	data, err := RenderFilePayload(input)
	if err != nil {
		t.Fatalf("RenderFilePayload() error = %v", err)
	}
	text := string(data)

	for _, localPath := range []string{"/repo", "/target", "/tmp"} {
		if strings.Contains(text, localPath) {
			t.Fatalf("payload leaks local path %q:\n%s", localPath, text)
		}
	}
}
