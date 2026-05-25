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

func TestTargetPrepareClonesMissingWorktrees(t *testing.T) {
	setup := prepareTargetPrepare(t)

	result, decoded := runTargetPrepareJSON(t, setup, 0, "--remote-template", targetPrepareRemoteTemplate(setup))

	if got := stringField(t, decoded, "kind"); got != "target-prepare" {
		t.Fatalf("kind = %q\n%s", got, result.Stdout)
	}
	if got := stringField(t, decoded, "status"); got != "prepared" {
		t.Fatalf("status = %q\n%s", got, result.Stdout)
	}
	for _, repo := range setup.repositories {
		worktree := targetWorktreePath(setup.targetRoot, repo.name)
		assertWorktreeClean(t, worktree)
		if branch := strings.TrimSpace(mustGitOutput(t, worktree, "branch", "--show-current")); branch != "main" {
			t.Fatalf("%s branch = %q, want main", repo.name, branch)
		}
		if remote := strings.TrimSpace(mustGitOutput(t, worktree, "remote", "get-url", "origin")); remote != targetPrepareRemoteURL(setup, repo) {
			t.Fatalf("%s origin = %q, want %q", repo.name, remote, targetPrepareRemoteURL(setup, repo))
		}
	}
	if _, err := os.Stat(filepath.Join(setup.targetRoot, ".arcpub", "state")); !os.IsNotExist(err) {
		t.Fatalf("target prepare created transaction state; stat err = %v", err)
	}

	runPreflightJSON(t, setup, 0)
}

func TestTargetPrepareValidatesExistingWorktrees(t *testing.T) {
	setup := prepareTargetPrepare(t)
	runTargetPrepareJSON(t, setup, 0, "--remote-template", targetPrepareRemoteTemplate(setup))

	before := map[string]string{}
	for _, repo := range setup.repositories {
		worktree := targetWorktreePath(setup.targetRoot, repo.name)
		before[repo.name] = strings.TrimSpace(mustGitOutput(t, worktree, "rev-parse", "HEAD"))
	}

	_, decoded := runTargetPrepareJSON(t, setup, 0, "--remote-template", targetPrepareRemoteTemplate(setup))

	assertTargetPrepareAction(t, decoded, "validate-worktree", "passed", "")
	assertTargetPrepareAction(t, decoded, "fetch", "passed", "")
	assertTargetPrepareAction(t, decoded, "checkout", "passed", "")
	for _, repo := range setup.repositories {
		worktree := targetWorktreePath(setup.targetRoot, repo.name)
		after := strings.TrimSpace(mustGitOutput(t, worktree, "rev-parse", "HEAD"))
		if after != before[repo.name] {
			t.Fatalf("%s HEAD changed: %s -> %s", repo.name, before[repo.name], after)
		}
		assertWorktreeClean(t, worktree)
	}
}

func TestTargetPrepareMissingRemoteTemplateFailsForMissingWorktree(t *testing.T) {
	setup := prepareTargetPrepare(t)

	_, decoded := runTargetPrepareJSON(t, setup, 1)

	assertTargetPrepareAction(t, decoded, "clone", "failed", "missing_remote_template")
}

func TestTargetPrepareUsesManifestRemoteTemplate(t *testing.T) {
	setup := prepareTargetPrepare(t)
	writeManifestRemoteTemplate(t, setup, targetPrepareRemoteTemplate(setup))

	_, decoded := runTargetPrepareJSON(t, setup, 0)

	assertTargetPrepareAction(t, decoded, "clone", "passed", "")
}

func TestTargetPrepareCLIRemoteTemplateOverridesManifest(t *testing.T) {
	setup := prepareTargetPrepare(t)
	writeManifestRemoteTemplate(t, setup, "file:///does/not/exist/{name}.git")

	_, decoded := runTargetPrepareJSON(t, setup, 0, "--remote-template", targetPrepareRemoteTemplate(setup))

	assertTargetPrepareAction(t, decoded, "clone", "passed", "")
}

func TestTargetPrepareRejectsNonGitWorktree(t *testing.T) {
	setup := prepareTargetPrepare(t)
	if err := os.MkdirAll(targetWorktreePath(setup.targetRoot, "arcoris/foundation"), 0o755); err != nil {
		t.Fatalf("create plain worktree: %v", err)
	}

	_, decoded := runTargetPrepareJSON(t, setup, 1, "--remote-template", targetPrepareRemoteTemplate(setup))

	assertTargetPrepareAction(t, decoded, "validate-worktree", "failed", "worktree_status_failed")
}

func TestTargetPrepareRejectsDirtyWorktree(t *testing.T) {
	setup := prepareLocalPublish(t)
	worktree := targetWorktreePath(setup.targetRoot, "arcoris/foundation")
	if err := os.WriteFile(filepath.Join(worktree, ".dirty"), []byte("dirty\n"), 0o644); err != nil {
		t.Fatalf("write dirty file: %v", err)
	}

	_, decoded := runTargetPrepareJSON(t, setup, 1)

	assertTargetPrepareAction(t, decoded, "clean", "failed", "worktree_dirty")
}

func TestTargetPrepareRejectsRemoteMismatch(t *testing.T) {
	setup := prepareLocalPublish(t)
	wrongRemoteRoot := t.TempDir()
	for _, repo := range setup.repositories {
		bare := filepath.Join(wrongRemoteRoot, repo.bareName)
		initBareGitRepo(t, bare)
		seedBareGitMain(t, bare)
	}
	setup.remoteRoot = wrongRemoteRoot

	_, decoded := runTargetPrepareJSON(t, setup, 1, "--remote-template", targetPrepareRemoteTemplate(setup))

	assertTargetPrepareAction(t, decoded, "remote", "failed", "remote_mismatch")
}

func TestTargetPrepareNoPathLeaksByDefault(t *testing.T) {
	setup := prepareTargetPrepare(t)

	result, _ := runTargetPrepareJSON(t, setup, 0, "--remote-template", targetPrepareRemoteTemplate(setup))

	assertNoLocalPathLeak(t, result.Stdout, setup.root, setup.targetRoot, setup.remoteRoot)
}

func TestTargetPrepareIncludeLocalPaths(t *testing.T) {
	setup := prepareTargetPrepare(t)

	result, _ := runTargetPrepareJSON(t, setup, 0,
		"--remote-template", targetPrepareRemoteTemplate(setup),
		"--include-local-paths",
	)

	assertContains(t, result.Stdout, setup.targetRoot)
	assertContains(t, result.Stdout, setup.remoteRoot)
}

func TestTargetPrepareUnsupportedMultiBranchFails(t *testing.T) {
	setup := prepareTargetPrepare(t)
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
	mustRun(t, setup.root, "git", "commit", "-m", "test: add multi-branch target fixture")

	_, decoded := runTargetPrepareJSON(t, setup, 1, "--remote-template", targetPrepareRemoteTemplate(setup))

	assertTargetPrepareAction(t, decoded, "checkout", "failed", "unsupported_multi_branch")
	for _, repo := range setup.repositories {
		assertPathMissing(t, targetWorktreePath(setup.targetRoot, repo.name))
	}
}

func prepareTargetPrepare(t *testing.T) localPublishSetup {
	t.Helper()
	requireGitAndGo(t)

	setup := localPublishSetup{
		root:       copyFixture(t, "local-publish"),
		targetRoot: filepath.Join(t.TempDir(), "targets"),
		remoteRoot: t.TempDir(),
		repositories: []localPublishRepository{
			{name: "arcoris/foundation", bareName: "foundation.git"},
			{name: "arcoris/control", bareName: "control.git"},
		},
	}
	initGitRepo(t, setup.root)
	for _, repo := range setup.repositories {
		bare := setup.bareRepo(repo.name)
		initBareGitRepo(t, bare)
		seedBareGitMain(t, bare)
	}
	return setup
}

func runTargetPrepareJSON(t *testing.T, setup localPublishSetup, wantCode int, extra ...string) (commandResult, map[string]any) {
	t.Helper()
	args := []string{
		"target", "prepare",
		"--manifest", e2eManifest(setup.root),
		"--version", "v0.1.0",
		"--target-root", setup.targetRoot,
		"--output", "json",
	}
	args = append(args, extra...)
	result := runArcpub(t, args...)
	assertExitCode(t, result, wantCode)
	return result, assertJSON(t, result.Stdout)
}

func targetPrepareRemoteTemplate(setup localPublishSetup) string {
	root := filepath.ToSlash(setup.remoteRoot)
	if !strings.HasPrefix(root, "/") {
		root = "/" + root
	}
	return "file://" + root + "/{name}.git"
}

func targetPrepareRemoteURL(setup localPublishSetup, repo localPublishRepository) string {
	root := filepath.ToSlash(setup.remoteRoot)
	if !strings.HasPrefix(root, "/") {
		root = "/" + root
	}
	return "file://" + root + "/" + repo.bareName
}

func writeManifestRemoteTemplate(t *testing.T, setup localPublishSetup, template string) {
	t.Helper()
	data, err := os.ReadFile(e2eManifest(setup.root))
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	updated := strings.Replace(string(data), "publish:\n", "target:\n  remoteTemplate: \""+template+"\"\npublish:\n", 1)
	if updated == string(data) {
		t.Fatal("manifest publish block not found")
	}
	if err := os.WriteFile(e2eManifest(setup.root), []byte(updated), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}
	mustRun(t, setup.root, "git", "add", "arcpub.yaml")
	mustRun(t, setup.root, "git", "commit", "-m", "test: add target remote template")
}

func assertTargetPrepareAction(t *testing.T, decoded map[string]any, name string, status string, code string) {
	t.Helper()
	modules := arrayField(t, decoded, "modules")
	for _, rawModule := range modules {
		module := rawModule.(map[string]any)
		actions := arrayField(t, module, "actions")
		for _, rawAction := range actions {
			action := rawAction.(map[string]any)
			if action["name"] != name || action["status"] != status {
				continue
			}
			if code == "" || action["code"] == code {
				return
			}
		}
	}
	t.Fatalf("target prepare action %s/%s/%s not found: %#v", name, status, code, decoded)
}
