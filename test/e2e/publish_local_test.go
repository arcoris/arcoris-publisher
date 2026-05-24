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
	"path/filepath"
	"testing"
)

func TestPublishLocalBareRepositories(t *testing.T) {
	requireGitAndGo(t)

	root := copyFixture(t, "local-publish")
	initGitRepo(t, root)
	targetRoot := t.TempDir()
	remoteRoot := t.TempDir()

	repositories := []struct {
		name     string
		bareName string
	}{
		{name: "arcoris/foundation", bareName: "foundation.git"},
		{name: "arcoris/control", bareName: "control.git"},
	}
	for _, repo := range repositories {
		bare := filepath.Join(remoteRoot, repo.bareName)
		initBareGitRepo(t, bare)
		initTargetGitWorktree(t, targetWorktreePath(targetRoot, repo.name), bare)
	}

	result, decoded := runArcpubJSON(t,
		0,
		"publish",
		"--manifest", e2eManifest(root),
		"--version", "v0.1.0",
		"--source-repo", root,
		"--staging-dir", root,
		"--target-root", targetRoot,
		"--output", "json",
	)
	if decoded["status"] != "published" {
		t.Fatalf("workflow status = %#v, want published\n%s", decoded["status"], result.Stdout)
	}
	assertNoLocalPathLeak(t, result.Stdout, root, targetRoot, remoteRoot)

	for _, repo := range repositories {
		bare := filepath.Join(remoteRoot, repo.bareName)
		assertGitRefExists(t, bare, "refs/heads/main")
		assertGitRefExists(t, bare, "refs/tags/v0.1.0")

		message := gitLogMessage(t, bare, "refs/heads/main")
		assertContains(t, message, "Arcoris-Source-Commit:")
		assertContains(t, message, "Arcoris-Projection-Hash:")
		assertContains(t, message, "Arcoris-Publisher-Version:")

		gitTreeContains(t, bare, "refs/heads/main", "go.mod")
		gitTreeContains(t, bare, "refs/heads/main", "README.md")
		gitTreeContains(t, bare, "refs/heads/main", "contracts/doc.go")
		gitTreeContains(t, bare, "refs/heads/main", ".arcoris/provenance.json")
		gitTreeMissing(t, bare, "refs/heads/main", "secret.txt")
		gitTreeMissing(t, bare, "refs/heads/main", "private/secret.go")
		gitTreeMissing(t, bare, "refs/heads/main", "private")
	}
}
