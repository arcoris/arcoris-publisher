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

func TestPreflightFailsWhenTargetCommitIdentityMissing(t *testing.T) {
	setup := prepareLocalPublish(t)
	clearTargetIdentities(t, setup)

	result, decoded := runPreflightWithEnvJSON(t, setup, isolatedGitConfigEnv(t), 1)

	assertPreflightCheckFailedWithCode(t, decoded, "commit-identity", "missing_commit_identity")
	assertNotContains(t, result.Stdout, "arcoris-test@example.invalid")
	assertNoLocalPathLeak(t, result.Stdout, setup.root, setup.targetRoot, setup.remoteRoot)
}

func TestPreflightPassesWithTargetCommitIdentity(t *testing.T) {
	setup := prepareLocalPublish(t)
	for _, repo := range setup.repositories {
		configureGitIdentity(t, targetWorktreePath(setup.targetRoot, repo.name))
	}

	_, decoded := runPreflightWithEnvJSON(t, setup, isolatedGitConfigEnv(t), 0)

	assertPreflightCheckPassed(t, decoded, "commit-identity")
}

func TestPublishFailsBeforeTransactionWhenIdentityMissing(t *testing.T) {
	setup := prepareLocalPublish(t)
	clearTargetIdentities(t, setup)
	beforeHeads := map[string]string{}
	for _, repo := range setup.repositories {
		beforeHeads[repo.name] = strings.TrimSpace(mustGitOutput(t, targetWorktreePath(setup.targetRoot, repo.name), "rev-parse", "HEAD"))
	}

	result := runLocalPublishWithEnv(t, setup, isolatedGitConfigEnv(t), 1)

	assertContains(t, result.Stderr, "target commit identity failed")
	if _, err := os.Stat(filepath.Join(setup.targetRoot, ".arcpub", "state", "transactions")); !os.IsNotExist(err) {
		t.Fatalf("transaction journal created after identity failure: %v", err)
	}
	if _, err := os.Stat(filepath.Join(setup.targetRoot, ".arcpub", "state", "publish.lock")); !os.IsNotExist(err) {
		t.Fatalf("publish lock created after identity failure: %v", err)
	}
	for _, repo := range setup.repositories {
		worktree := targetWorktreePath(setup.targetRoot, repo.name)
		after := strings.TrimSpace(mustGitOutput(t, worktree, "rev-parse", "HEAD"))
		if after != beforeHeads[repo.name] {
			t.Fatalf("%s HEAD changed: %s -> %s", repo.name, beforeHeads[repo.name], after)
		}
		assertGitRefMissing(t, setup.bareRepo(repo.name), "refs/heads/main")
		assertGitRefMissing(t, setup.bareRepo(repo.name), "refs/tags/v0.1.0")
		assertGitRefMissing(t, setup.bareRepo(repo.name), "refs/heads/arcpub/tx/tx-test/foundation")
		assertGitRefMissing(t, setup.bareRepo(repo.name), "refs/heads/arcpub/tx/tx-test/control")
	}
}

func TestPublishSucceedsWithTargetCommitIdentity(t *testing.T) {
	setup := prepareLocalPublish(t)

	runLocalPublish(t, setup, 0)
}

func clearTargetIdentities(t *testing.T, setup localPublishSetup) {
	t.Helper()
	for _, repo := range setup.repositories {
		clearGitIdentity(t, targetWorktreePath(setup.targetRoot, repo.name))
	}
}

func runPreflightWithEnvJSON(t *testing.T, setup localPublishSetup, env []string, wantCode int) (commandResult, map[string]any) {
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
	result := runArcpubWithEnv(t, env, args...)
	assertExitCode(t, result, wantCode)
	return result, assertJSON(t, result.Stdout)
}

func runLocalPublishWithEnv(t *testing.T, setup localPublishSetup, env []string, wantCode int) commandResult {
	t.Helper()
	result := runArcpubWithEnv(t, env,
		"publish",
		"--manifest", e2eManifest(setup.root),
		"--version", "v0.1.0",
		"--source-repo", setup.root,
		"--staging-dir", setup.root,
		"--target-root", setup.targetRoot,
		"--output", "json",
	)
	assertExitCode(t, result, wantCode)
	return result
}

func assertPreflightCheckPassed(t *testing.T, decoded map[string]any, name string) {
	t.Helper()
	assertPreflightCheckStatus(t, decoded, name, "passed", "")
}

func assertPreflightCheckFailedWithCode(t *testing.T, decoded map[string]any, name string, code string) {
	t.Helper()
	assertPreflightCheckStatus(t, decoded, name, "failed", code)
}

func assertPreflightCheckStatus(t *testing.T, decoded map[string]any, name string, status string, code string) {
	t.Helper()
	modules := arrayField(t, decoded, "modules")
	for _, rawModule := range modules {
		module := rawModule.(map[string]any)
		checks := arrayField(t, module, "checks")
		for _, rawCheck := range checks {
			check := rawCheck.(map[string]any)
			if check["name"] != name || check["status"] != status {
				continue
			}
			if code == "" || check["code"] == code {
				return
			}
		}
	}
	t.Fatalf("preflight check %s/%s/%s not found: %#v", name, status, code, decoded)
}
