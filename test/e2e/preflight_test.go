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

func TestPreflightMinimalFixturePasses(t *testing.T) {
	setup := prepareLocalPublish(t)

	result, decoded := runPreflightJSON(t, setup, 0)

	if got := stringField(t, decoded, "kind"); got != "preflight" {
		t.Fatalf("kind = %q\n%s", got, result.Stdout)
	}
	if got := stringField(t, decoded, "status"); got != "passed" {
		t.Fatalf("status = %q\n%s", got, result.Stdout)
	}
	if got := floatField(t, decoded, "moduleCount"); got != 2 {
		t.Fatalf("moduleCount = %v\n%s", got, result.Stdout)
	}
	assertContains(t, result.Stdout, "foundation")
	assertContains(t, result.Stdout, "control")
}

func TestPreflightDoesNotMutateTargetWorktrees(t *testing.T) {
	setup := prepareLocalPublish(t)
	before := map[string]string{}
	for _, repo := range setup.repositories {
		worktree := targetWorktreePath(setup.targetRoot, repo.name)
		before[repo.name] = strings.TrimSpace(mustGitOutput(t, worktree, "rev-parse", "HEAD"))
		assertWorktreeClean(t, worktree)
	}

	runPreflightJSON(t, setup, 0)

	for _, repo := range setup.repositories {
		worktree := targetWorktreePath(setup.targetRoot, repo.name)
		after := strings.TrimSpace(mustGitOutput(t, worktree, "rev-parse", "HEAD"))
		if after != before[repo.name] {
			t.Fatalf("%s HEAD changed: %s -> %s", repo.name, before[repo.name], after)
		}
		assertWorktreeClean(t, worktree)
	}
	if _, err := os.Stat(filepath.Join(setup.targetRoot, ".arcpub", "state")); !os.IsNotExist(err) {
		t.Fatalf("preflight created state dir; stat err = %v", err)
	}
}

func TestPreflightFailsOnNonGitTarget(t *testing.T) {
	requireGitAndGo(t)
	root := copyFixture(t, "local-publish")
	initGitRepo(t, root)
	targetRoot := t.TempDir()
	for _, repository := range []string{"arcoris/foundation", "arcoris/control"} {
		if err := os.MkdirAll(targetWorktreePath(targetRoot, repository), 0o755); err != nil {
			t.Fatalf("create target worktree: %v", err)
		}
	}

	result := runPreflight(t, localPublishSetup{root: root, targetRoot: targetRoot}, 1)
	decoded := assertJSON(t, result.Stdout)
	assertPreflightCheckFailed(t, decoded, "target-status")
}

func TestPreflightFailsOnDirtyTarget(t *testing.T) {
	setup := prepareLocalPublish(t)
	worktree := targetWorktreePath(setup.targetRoot, "arcoris/foundation")
	if err := os.WriteFile(filepath.Join(worktree, ".dirty"), []byte("dirty\n"), 0o644); err != nil {
		t.Fatalf("write dirty file: %v", err)
	}

	result, decoded := runPreflightJSON(t, setup, 1)

	assertContains(t, result.Stdout, "target_worktree_dirty")
	assertPreflightCheckFailed(t, decoded, "target-status")
}

func TestPreflightFailsOnExistingLocalTag(t *testing.T) {
	setup := prepareLocalPublish(t)
	worktree := targetWorktreePath(setup.targetRoot, "arcoris/foundation")
	mustRun(t, worktree, "git", "tag", "-a", "v0.1.0", "-m", "existing")

	_, decoded := runPreflightJSON(t, setup, 1)

	assertPreflightCheckFailed(t, decoded, "local-tag")
}

func TestPreflightFailsOnExistingRemoteTag(t *testing.T) {
	setup := prepareLocalPublish(t)
	worktree := targetWorktreePath(setup.targetRoot, "arcoris/foundation")
	mustRun(t, worktree, "git", "tag", "-a", "v0.1.0", "-m", "existing")
	mustRun(t, worktree, "git", "push", "origin", "refs/tags/v0.1.0")
	mustRun(t, worktree, "git", "tag", "-d", "v0.1.0")

	_, decoded := runPreflightJSON(t, setup, 1)

	assertPreflightCheckFailed(t, decoded, "remote-tag")
}

func TestPreflightFailsOnPendingTransaction(t *testing.T) {
	setup := prepareLocalPublish(t)
	writePendingJournal(t, setup, `"status":"pending"`)

	_, decoded := runPreflightJSON(t, setup, 1)

	assertPreflightGlobalFailed(t, decoded, "pending-transactions")
}

func TestPreflightFailsOnCorruptedJournal(t *testing.T) {
	setup := prepareLocalPublish(t)
	txDir := filepath.Join(setup.targetRoot, ".arcpub", "state", "transactions")
	if err := os.MkdirAll(txDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(txDir, "tx-corrupt.json"), []byte("{"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, decoded := runPreflightJSON(t, setup, 1)

	assertPreflightGlobalFailed(t, decoded, "pending-transactions")
}

func TestPreflightFailsOnMissingJournalLock(t *testing.T) {
	setup := prepareLocalPublish(t)
	stateDir := filepath.Join(setup.targetRoot, ".arcpub", "state")
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(stateDir, "publish.lock"), []byte("transaction=tx-other\npid=1\nstartedAt=2026-01-01T00:00:00Z\ncommand=publish\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, decoded := runPreflightJSON(t, setup, 1)

	assertPreflightGlobalFailedCode(t, decoded, "publish-lock", "stale_publish_lock_journal_missing")
}

func TestPreflightFailsOnCorruptLock(t *testing.T) {
	setup := prepareLocalPublish(t)
	stateDir := filepath.Join(setup.targetRoot, ".arcpub", "state")
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(stateDir, "publish.lock"), []byte("pid=1\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, decoded := runPreflightJSON(t, setup, 1)

	assertPreflightGlobalFailedCode(t, decoded, "publish-lock", "publish_lock_corrupt")
}

func TestPreflightFailsOnOperationLockWithoutMutatingIt(t *testing.T) {
	setup := prepareLocalPublish(t)
	stateDir := filepath.Join(setup.targetRoot, ".arcpub", "state")
	writeE2EOperationLock(t, stateDir, "publish")

	result, decoded := runPreflightJSON(t, setup, 1)

	assertPreflightGlobalFailedCode(t, decoded, "publish-lock", "operation_lock_exists")
	if strings.Contains(result.Stdout, "token-one") {
		t.Fatalf("preflight leaked operation lock token:\n%s", result.Stdout)
	}
	assertNoLocalPathLeak(t, result.Stdout, stateDir)
	assertFileExists(t, transactionOperationLockPath(stateDir))
}

func TestPreflightRejectsMultiBranch(t *testing.T) {
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

	_, decoded := runPreflightJSON(t, setup, 1)

	assertPreflightCheckFailed(t, decoded, "multi-branch")
}

func TestPreflightNoPathLeaksByDefault(t *testing.T) {
	setup := prepareLocalPublish(t)

	result, _ := runPreflightJSON(t, setup, 0)

	assertNoLocalPathLeak(t, result.Stdout, setup.root, setup.targetRoot, setup.remoteRoot)
}

func TestPreflightIncludeLocalPaths(t *testing.T) {
	setup := prepareLocalPublish(t)

	result := runPreflightWithArgs(t, setup, 0, "--include-local-paths")

	assertContains(t, result.Stdout, setup.targetRoot)
}

func runPreflightJSON(t *testing.T, setup localPublishSetup, wantCode int) (commandResult, map[string]any) {
	t.Helper()
	result := runPreflight(t, setup, wantCode)
	return result, assertJSON(t, result.Stdout)
}

func runPreflight(t *testing.T, setup localPublishSetup, wantCode int) commandResult {
	t.Helper()
	return runPreflightWithArgs(t, setup, wantCode)
}

func runPreflightWithArgs(t *testing.T, setup localPublishSetup, wantCode int, extra ...string) commandResult {
	t.Helper()
	args := []string{
		"preflight",
		"--manifest", e2eManifest(setup.root),
		"--version", "v0.1.0",
		"--source-repo", setup.root,
		"--staging-dir", setup.root,
		"--target-root", setup.targetRoot,
		"--output", "json",
	}
	args = append(args, extra...)
	result := runArcpub(t, args...)
	assertExitCode(t, result, wantCode)
	return result
}

func mustGitOutput(t *testing.T, dir string, args ...string) string {
	t.Helper()
	result := runCommand(t, dir, nil, "git", args...)
	if result.Code != 0 {
		t.Fatalf("git %v failed\nstdout:\n%s\nstderr:\n%s", args, result.Stdout, result.Stderr)
	}
	return result.Stdout
}

func writePendingJournal(t *testing.T, setup localPublishSetup, statusFragment string) {
	t.Helper()
	txDir := filepath.Join(setup.targetRoot, ".arcpub", "state", "transactions")
	if err := os.MkdirAll(txDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	data := `{"schemaVersion":1,"id":"tx-pending",` + statusFragment + `,"version":"v0.1.0","startedAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z","modules":[]}`
	if err := os.WriteFile(filepath.Join(txDir, "tx-pending.json"), []byte(data), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}

func assertPreflightGlobalFailed(t *testing.T, decoded map[string]any, name string) {
	t.Helper()
	assertPreflightGlobalFailedCode(t, decoded, name, "")
}

func assertPreflightGlobalFailedCode(t *testing.T, decoded map[string]any, name string, code string) {
	t.Helper()
	checks, ok := decoded["checks"].([]any)
	if !ok {
		t.Fatalf("checks missing: %#v", decoded)
	}
	for _, raw := range checks {
		check := raw.(map[string]any)
		if check["name"] == name && check["status"] == "failed" && (code == "" || check["code"] == code) {
			return
		}
	}
	t.Fatalf("failed global check %q code %q not found: %#v", name, code, checks)
}

func assertPreflightCheckFailed(t *testing.T, decoded map[string]any, name string) {
	t.Helper()
	modules, ok := decoded["modules"].([]any)
	if !ok {
		t.Fatalf("modules missing: %#v", decoded)
	}
	for _, rawModule := range modules {
		module := rawModule.(map[string]any)
		checks, _ := module["checks"].([]any)
		for _, raw := range checks {
			check := raw.(map[string]any)
			if check["name"] == name && check["status"] == "failed" {
				return
			}
		}
	}
	t.Fatalf("failed module check %q not found: %#v", name, modules)
}
