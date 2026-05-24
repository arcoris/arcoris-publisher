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

package e2e_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestVerifyWritesProvenanceArtifact(t *testing.T) {
	_, targetRoot, provenancePath := runProvenanceVerify(t)

	assertFileExists(t, provenancePath)
	if want := filepath.Join(targetWorktreePath(targetRoot, "arcoris/provenance"), ".arcoris", "provenance.json"); provenancePath != want {
		t.Fatalf("provenance path = %s, want %s", provenancePath, want)
	}
}

func TestProvenanceArtifactIsParseableAndComplete(t *testing.T) {
	_, _, provenancePath := runProvenanceVerify(t)
	payload := readProvenancePayload(t, provenancePath)

	if stringField(t, payload, "schemaVersion") != "arcoris.provenance/v1" {
		t.Fatalf("schemaVersion = %#v", payload["schemaVersion"])
	}
	module := objectField(t, payload, "module")
	if stringField(t, module, "name") != "provenance" {
		t.Fatalf("module.name = %#v", module["name"])
	}
	if stringField(t, module, "modulePath") != "arcoris.dev/provenance" {
		t.Fatalf("module.modulePath = %#v", module["modulePath"])
	}
	if stringField(t, module, "version") != "v0.1.0" {
		t.Fatalf("module.version = %#v", module["version"])
	}

	source := objectField(t, payload, "source")
	if stringField(t, source, "commit") == "" {
		t.Fatalf("source.commit is empty")
	}
	projection := objectField(t, payload, "projection")
	if !strings.HasPrefix(stringField(t, projection, "projectionHash"), "sha256:") {
		t.Fatalf("projection.projectionHash = %#v", projection["projectionHash"])
	}
	publisher := objectField(t, payload, "publisher")
	if stringField(t, publisher, "version") == "" {
		t.Fatalf("publisher.version is empty")
	}
}

func TestProvenanceDoesNotLeakLocalPaths(t *testing.T) {
	root, targetRoot, provenancePath := runProvenanceVerify(t)
	data, err := os.ReadFile(provenancePath)
	if err != nil {
		t.Fatalf("read provenance file: %v", err)
	}

	assertNoLocalPathLeak(t, string(data), root, targetRoot, repoRoot(t))
}

func TestProvenancePathCannotEscapeTargetRoot(t *testing.T) {
	root := copyFixture(t, "provenance")
	data, err := os.ReadFile(e2eManifest(root))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	updated := strings.Replace(string(data), "file: .arcoris/provenance.json", "file: ../provenance.json", 1)
	if err := os.WriteFile(e2eManifest(root), []byte(updated), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	result := runArcpub(t,
		"plan",
		"--manifest", e2eManifest(root),
		"--version", "v0.1.0",
	)
	assertExitCode(t, result, 1)
	assertContains(t, result.Stderr, "invalid_manifest")
}

func runProvenanceVerify(t *testing.T) (string, string, string) {
	t.Helper()
	requireGitAndGo(t)

	root := copyFixture(t, "provenance")
	initGitRepo(t, root)
	targetRoot := t.TempDir()
	prepareTargetWorktrees(t, targetRoot, "arcoris/provenance")

	result, _ := runArcpubJSON(t,
		0,
		"verify",
		"--manifest", e2eManifest(root),
		"--version", "v0.1.0",
		"--source-repo", root,
		"--staging-dir", root,
		"--target-root", targetRoot,
		"--output", "json",
	)
	assertNoLocalPathLeak(t, result.Stdout, root, targetRoot)

	provenancePath := filepath.Join(
		targetWorktreePath(targetRoot, "arcoris/provenance"),
		".arcoris",
		"provenance.json",
	)
	return root, targetRoot, provenancePath
}

func readProvenancePayload(t *testing.T, path string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read provenance file: %v", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		t.Fatalf("provenance is not valid JSON: %v\n%s", err, string(data))
	}
	return payload
}
