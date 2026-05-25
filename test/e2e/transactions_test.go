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
	"testing"
)

func TestPublishTransactionHappyPath(t *testing.T) {
	setup := prepareLocalPublish(t)

	result, decoded := runLocalPublishJSON(t, setup, 0)
	publish := objectField(t, decoded, "publish")
	transaction := objectField(t, publish, "transaction")
	if got := stringField(t, transaction, "status"); got != "committed" {
		t.Fatalf("transaction status = %q\n%s", got, result.Stdout)
	}
	id := stringField(t, transaction, "id")

	for _, repo := range setup.repositories {
		bare := setup.bareRepo(repo.name)
		assertGitRefExists(t, bare, "refs/heads/main")
		assertGitRefExists(t, bare, "refs/tags/v0.1.0")
		assertGitRefMissing(t, bare, "refs/heads/arcpub/tx/"+id+"/"+moduleRefName(repo.name))
	}

	show := runArcpub(t,
		"transactions", "show", id,
		"--target-root", setup.targetRoot,
		"--output", "json",
	)
	assertExitCode(t, show, 0)
	showJSON := assertJSON(t, show.Stdout)
	if got := stringField(t, showJSON, "status"); got != "committed" {
		t.Fatalf("show status = %q\n%s", got, show.Stdout)
	}
}

func TestPublishCandidatePushFailureRollsBack(t *testing.T) {
	setup := prepareLocalPublish(t)
	installBareHook(t, setup.bareRepo("arcoris/control"), "pre-receive", rejectRefsHook("refs/heads/arcpub/tx/"))

	result, decoded := runLocalPublishJSON(t, setup, 1)
	assertContains(t, result.Stderr, "candidate")
	assertTransactionStatus(t, decoded, "rolled_back")
	assertNoFinalRefs(t, setup)
	assertWorktreesClean(t, setup)
}

func TestPublishPromotionFailureRollsBackPromotedBranches(t *testing.T) {
	setup := prepareLocalPublish(t)
	installBareHook(t, setup.bareRepo("arcoris/control"), "pre-receive", rejectRefsHook("refs/heads/main"))

	result, decoded := runLocalPublishJSON(t, setup, 1)
	assertContains(t, result.Stdout+result.Stderr, "promotion")
	assertTransactionStatus(t, decoded, "rolled_back")
	assertNoFinalRefs(t, setup)
	assertWorktreesClean(t, setup)
}

func TestPublishTagPushFailureRollsBackBranchesAndTags(t *testing.T) {
	setup := prepareLocalPublish(t)
	installBareHook(t, setup.bareRepo("arcoris/foundation"), "pre-receive", rejectRefsHook("refs/tags/"))

	result, decoded := runLocalPublishJSON(t, setup, 1)
	assertContains(t, result.Stdout+result.Stderr, "tag")
	assertTransactionStatus(t, decoded, "rolled_back")
	assertNoFinalRefs(t, setup)
	assertWorktreesClean(t, setup)
}

func TestPendingTransactionBlocksNewPublish(t *testing.T) {
	setup := prepareLocalPublish(t)
	installBareHook(t, setup.bareRepo("arcoris/control"), "pre-receive", rejectRefsHook("refs/heads/arcpub/tx/"))
	first := runLocalPublishWithArgs(t, setup, 1, "--rollback", "manual")
	assertContains(t, first.Stderr, "candidate")

	second := runLocalPublish(t, setup, 1)
	assertContains(t, second.Stderr, "pending publish transaction")
}

func TestRollbackCommandCompletesPendingTransaction(t *testing.T) {
	setup := prepareLocalPublish(t)
	installBareHook(t, setup.bareRepo("arcoris/control"), "pre-receive", rejectRefsHook("refs/heads/arcpub/tx/"))
	result, decoded := runLocalPublishJSONWithArgs(t, setup, 1, "--rollback", "manual")
	id := transactionID(t, decoded)

	rollback := runArcpub(t,
		"rollback",
		"--transaction", id,
		"--target-root", setup.targetRoot,
		"--output", "json",
	)
	assertExitCode(t, rollback, 0)
	rollbackJSON := assertJSON(t, rollback.Stdout)
	if got := stringField(t, rollbackJSON, "status"); got != "rolled_back" {
		t.Fatalf("rollback status = %q\npublish:\n%s\nrollback:\n%s", got, result.Stdout, rollback.Stdout)
	}
	assertNoFinalRefs(t, setup)
	assertWorktreesClean(t, setup)
}

func TestRollbackCommandIsIdempotent(t *testing.T) {
	setup := prepareLocalPublish(t)
	installBareHook(t, setup.bareRepo("arcoris/control"), "pre-receive", rejectRefsHook("refs/heads/arcpub/tx/"))
	_, decoded := runLocalPublishJSONWithArgs(t, setup, 1, "--rollback", "manual")
	id := transactionID(t, decoded)

	first := runArcpub(t,
		"rollback",
		"--transaction", id,
		"--target-root", setup.targetRoot,
		"--output", "json",
	)
	assertExitCode(t, first, 0)
	second := runArcpub(t,
		"rollback",
		"--transaction", id,
		"--target-root", setup.targetRoot,
		"--output", "json",
	)
	assertExitCode(t, second, 0)
	secondJSON := assertJSON(t, second.Stdout)
	if got := stringField(t, secondJSON, "status"); got != "rolled_back" {
		t.Fatalf("second rollback status = %q\n%s", got, second.Stdout)
	}
}

func TestRollbackFailureIsReported(t *testing.T) {
	setup := prepareLocalPublish(t)
	installBareHook(t, setup.bareRepo("arcoris/foundation"), "pre-receive", rejectDeleteRefsHook("refs/heads/main"))
	installBareHook(t, setup.bareRepo("arcoris/control"), "pre-receive", rejectRefsHook("refs/heads/main"))

	result, decoded := runLocalPublishJSON(t, setup, 1)

	assertContains(t, result.Stdout, "manualRecoveryActions")
	assertTransactionStatus(t, decoded, "rollback_failed")
}

func TestTransactionsShowNoPathLeaks(t *testing.T) {
	setup := prepareLocalPublish(t)
	_, decoded := runLocalPublishJSON(t, setup, 0)
	id := transactionID(t, decoded)

	show := runArcpub(t,
		"transactions", "show", id,
		"--target-root", setup.targetRoot,
		"--output", "json",
	)
	assertExitCode(t, show, 0)
	assertNoLocalPathLeak(t, show.Stdout, setup.root, setup.targetRoot, setup.remoteRoot)

	visible := runArcpub(t,
		"transactions", "show", id,
		"--target-root", setup.targetRoot,
		"--include-local-paths",
		"--output", "json",
	)
	assertExitCode(t, visible, 0)
	assertContains(t, visible.Stdout, setup.targetRoot)
}

func TestTransactionsListAndShow(t *testing.T) {
	setup := prepareLocalPublish(t)
	_, decoded := runLocalPublishJSON(t, setup, 0)
	id := transactionID(t, decoded)

	list := runArcpub(t,
		"transactions", "list",
		"--target-root", setup.targetRoot,
		"--output", "json",
	)
	assertExitCode(t, list, 0)
	listJSON := assertJSON(t, list.Stdout)
	if got := floatField(t, listJSON, "count"); got != 1 {
		t.Fatalf("transaction count = %v\n%s", got, list.Stdout)
	}

	show := runArcpub(t,
		"transactions", "show", id,
		"--target-root", setup.targetRoot,
		"--output", "json",
	)
	assertExitCode(t, show, 0)
	showJSON := assertJSON(t, show.Stdout)
	if got := stringField(t, showJSON, "id"); got != id {
		t.Fatalf("transaction id = %q, want %q", got, id)
	}
}

func TestCorruptedJournalBlocksPublish(t *testing.T) {
	setup := prepareLocalPublish(t)
	stateDir := filepath.Join(setup.targetRoot, ".arcpub", "state")
	txDir := filepath.Join(stateDir, "transactions")
	if err := os.MkdirAll(txDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(txDir, "tx-bad.json"), []byte("{"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	result := runLocalPublish(t, setup, 1)
	assertContains(t, result.Stderr, "corrupt")
	assertNoFinalRefs(t, setup)
}

func TestStateDirOverride(t *testing.T) {
	setup := prepareLocalPublish(t)
	stateDir := t.TempDir()

	_, decoded := runLocalPublishJSONWithArgs(t, setup, 0, "--state-dir", stateDir)
	id := transactionID(t, decoded)
	assertFileExists(t, filepath.Join(stateDir, "transactions", id+".json"))

	list := runArcpub(t,
		"transactions", "list",
		"--state-dir", stateDir,
		"--output", "json",
	)
	assertExitCode(t, list, 0)
}

func TestLockConflictBlocksPublish(t *testing.T) {
	setup := prepareLocalPublish(t)
	stateDir := filepath.Join(setup.targetRoot, ".arcpub", "state")
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(stateDir, "publish.lock"), []byte("transaction=tx-other\npid=1\nstartedAt=now\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	result := runLocalPublish(t, setup, 1)
	assertContains(t, result.Stderr, "lock")
	assertNoFinalRefs(t, setup)
}

func runLocalPublishJSON(t *testing.T, setup localPublishSetup, wantCode int) (commandResult, map[string]any) {
	t.Helper()
	return runLocalPublishJSONWithArgs(t, setup, wantCode)
}

func runLocalPublishJSONWithArgs(t *testing.T, setup localPublishSetup, wantCode int, extra ...string) (commandResult, map[string]any) {
	t.Helper()
	result := runLocalPublishWithArgs(t, setup, wantCode, extra...)
	return result, assertJSON(t, result.Stdout)
}

func runLocalPublishWithArgs(t *testing.T, setup localPublishSetup, wantCode int, extra ...string) commandResult {
	t.Helper()
	args := []string{
		"publish",
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

func assertTransactionStatus(t *testing.T, decoded map[string]any, want string) {
	t.Helper()
	publish := objectField(t, decoded, "publish")
	transaction := objectField(t, publish, "transaction")
	if got := stringField(t, transaction, "status"); got != want {
		t.Fatalf("transaction status = %q, want %q\ntransaction: %#v", got, want, transaction)
	}
}

func transactionID(t *testing.T, decoded map[string]any) string {
	t.Helper()
	return stringField(t, objectField(t, objectField(t, decoded, "publish"), "transaction"), "id")
}

func assertNoFinalRefs(t *testing.T, setup localPublishSetup) {
	t.Helper()
	for _, repo := range setup.repositories {
		bare := setup.bareRepo(repo.name)
		assertGitRefMissing(t, bare, "refs/heads/main")
		assertGitRefMissing(t, bare, "refs/tags/v0.1.0")
	}
}

func assertWorktreesClean(t *testing.T, setup localPublishSetup) {
	t.Helper()
	for _, repo := range setup.repositories {
		assertWorktreeClean(t, targetWorktreePath(setup.targetRoot, repo.name))
	}
}

func rejectRefsHook(prefix string) string {
	return `#!/usr/bin/env sh
set -eu
while read old new ref
do
	case "$ref" in
		` + prefix + `*)
			if [ "$new" != "0000000000000000000000000000000000000000" ]; then
				echo "rejecting $ref" >&2
				exit 1
			fi
			;;
	esac
done
exit 0
`
}

func rejectDeleteRefsHook(prefix string) string {
	return `#!/usr/bin/env sh
set -eu
while read old new ref
do
	case "$ref" in
		` + prefix + `*)
			if [ "$new" = "0000000000000000000000000000000000000000" ]; then
				echo "rejecting delete $ref" >&2
				exit 1
			fi
			;;
	esac
done
exit 0
`
}

func moduleRefName(repository string) string {
	for i := len(repository) - 1; i >= 0; i-- {
		if repository[i] == '/' {
			return repository[i+1:]
		}
	}
	return repository
}
