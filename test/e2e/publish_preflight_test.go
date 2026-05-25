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
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPublishRejectsNonGitTargetWorktree(t *testing.T) {
	requireGitAndGo(t)
	root := copyFixture(t, "local-publish")
	initGitRepo(t, root)
	targetRoot := t.TempDir()
	for _, repository := range []string{"arcoris/foundation", "arcoris/control"} {
		if err := os.MkdirAll(targetWorktreePath(targetRoot, repository), 0o755); err != nil {
			t.Fatalf("create non-git target worktree: %v", err)
		}
	}

	result := runArcpub(t,
		"publish",
		"--manifest", e2eManifest(root),
		"--version", "v0.1.0",
		"--source-repo", root,
		"--staging-dir", root,
		"--target-root", targetRoot,
	)

	assertExitCode(t, result, 1)
	assertContains(t, result.Stderr, "target Git status failed")
}

func TestPublishRejectsExistingTagBeforeMutation(t *testing.T) {
	setup := prepareLocalPublish(t)
	foundationWorktree := targetWorktreePath(setup.targetRoot, "arcoris/foundation")
	mustRun(t, foundationWorktree, "git", "tag", "-a", "v0.1.0", "-m", "existing")

	result := runLocalPublish(t, setup, 1)

	assertContains(t, result.Stderr, "tag")
	for _, repo := range setup.repositories {
		bare := setup.bareRepo(repo.name)
		assertGitRefMissing(t, bare, "refs/heads/main")
		assertGitRefMissing(t, bare, "refs/tags/v0.1.0")
	}
}

func TestPublishSecondRunSameVersionFailsPredictably(t *testing.T) {
	setup := prepareLocalPublish(t)
	first := runLocalPublish(t, setup, 0)
	assertJSON(t, first.Stdout)

	second := runLocalPublish(t, setup, 1)

	assertContains(t, second.Stderr, "tag")
}

func TestMultiBranchPublicationRejected(t *testing.T) {
	setup := prepareLocalPublish(t)
	data, err := os.ReadFile(e2eManifest(setup.root))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	updated := strings.Replace(string(data),
		"    - source: main\n      target: main",
		"    - source: main\n      target: main\n    - source: release\n      target: release",
		1,
	)
	if err := os.WriteFile(e2eManifest(setup.root), []byte(updated), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	mustRun(t, setup.root, "git", "add", "arcpub.yaml")
	mustRun(t, setup.root, "git", "commit", "-m", "test: add multi-branch fixture")

	result := runLocalPublish(t, setup, 1)

	assertContains(t, result.Stderr, "multi-branch")
	for _, repo := range setup.repositories {
		assertGitRefMissing(t, setup.bareRepo(repo.name), "refs/heads/main")
	}
}

func TestPublishDryRunHelpDescribesMutationBoundary(t *testing.T) {
	result := runArcpub(t, "publish", "--help")

	assertExitCode(t, result, 0)
	assertContains(t, result.Stdout, "construct and verify")
	assertContains(t, result.Stdout, "skip commit, tag, and push")
}

func TestProvenancePathCollisionRejected(t *testing.T) {
	requireGitAndGo(t)
	root := copyFixture(t, "provenance")
	initGitRepo(t, root)
	targetRoot := t.TempDir()
	prepareTargetWorktrees(t, targetRoot, "arcoris/provenance")

	moduleManifest := filepath.Join(root, "staging", "provenance", "arcpub.module.yaml")
	data, err := os.ReadFile(moduleManifest)
	if err != nil {
		t.Fatalf("read module manifest: %v", err)
	}
	updated := string(data) + `
    - type: file
      from: public/doc.go
      to: .arcoris/provenance.json
`
	if err := os.WriteFile(moduleManifest, []byte(updated), 0o644); err != nil {
		t.Fatalf("write module manifest: %v", err)
	}
	mustRun(t, root, "git", "add", filepath.Join("staging", "provenance", "arcpub.module.yaml"))
	mustRun(t, root, "git", "commit", "-m", "test: add provenance collision")

	result := runArcpub(t,
		"verify",
		"--manifest", e2eManifest(root),
		"--version", "v0.1.0",
		"--source-repo", root,
		"--staging-dir", root,
		"--target-root", targetRoot,
	)

	assertExitCode(t, result, 1)
	assertContains(t, result.Stderr, "already reserved")
}
