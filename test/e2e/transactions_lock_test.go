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
	"time"
)

func TestTransactionsLockShowAbsent(t *testing.T) {
	stateDir := t.TempDir()

	_, decoded := runTransactionsLockShowJSON(t, stateDir, 0)

	if got := stringField(t, decoded, "status"); got != "absent" {
		t.Fatalf("status = %q", got)
	}
	if got := stringField(t, decoded, "reason"); got != "lock_absent" {
		t.Fatalf("reason = %q", got)
	}
}

func TestTransactionsLockShowPresentWithJournal(t *testing.T) {
	stateDir := t.TempDir()
	writeE2ETransactionJournal(t, stateDir, "tx-committed", "committed", time.Now())
	writeE2ELockFile(t, stateDir, "tx-committed")

	_, decoded := runTransactionsLockShowJSON(t, stateDir, 0)

	if got := stringField(t, decoded, "status"); got != "present" {
		t.Fatalf("status = %q", got)
	}
	if got := stringField(t, decoded, "reason"); got != "lock_present" {
		t.Fatalf("reason = %q", got)
	}
	journal := objectField(t, decoded, "journal")
	if got := stringField(t, journal, "status"); got != "committed" {
		t.Fatalf("journal status = %q", got)
	}
}

func TestTransactionsLockShowJournalMissing(t *testing.T) {
	stateDir := t.TempDir()
	writeE2ELockFile(t, stateDir, "tx-missing")

	_, decoded := runTransactionsLockShowJSON(t, stateDir, 0)

	if got := stringField(t, decoded, "status"); got != "journal_missing" {
		t.Fatalf("status = %q", got)
	}
	if got := stringField(t, decoded, "reason"); got != "journal_missing" {
		t.Fatalf("reason = %q", got)
	}
}

func TestTransactionsLockShowCorruptLock(t *testing.T) {
	stateDir := t.TempDir()
	writeE2ERawLockFile(t, stateDir, "pid=1\n")

	result, decoded := runTransactionsLockShowJSON(t, stateDir, 1)

	if got := stringField(t, decoded, "status"); got != "corrupt" {
		t.Fatalf("status = %q\nstdout=%s\nstderr=%s", got, result.Stdout, result.Stderr)
	}
	if got := stringField(t, decoded, "reason"); got != "lock_corrupt" {
		t.Fatalf("reason = %q\nstdout=%s\nstderr=%s", got, result.Stdout, result.Stderr)
	}
}

func TestTransactionsLockShowCorruptJournal(t *testing.T) {
	stateDir := t.TempDir()
	writeE2ELockFile(t, stateDir, "tx-corrupt")
	txDir := filepath.Join(stateDir, "transactions")
	if err := os.MkdirAll(txDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(txDir, "tx-corrupt.json"), []byte("{"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	result, decoded := runTransactionsLockShowJSON(t, stateDir, 1)

	if got := stringField(t, decoded, "status"); got != "journal_corrupt" {
		t.Fatalf("status = %q\nstdout=%s\nstderr=%s", got, result.Stdout, result.Stderr)
	}
	if got := stringField(t, decoded, "reason"); got != "journal_corrupt" {
		t.Fatalf("reason = %q\nstdout=%s\nstderr=%s", got, result.Stdout, result.Stderr)
	}
}

func TestTransactionsLockClearRequiresTransactionAndConfirm(t *testing.T) {
	stateDir := t.TempDir()
	for _, tt := range []struct {
		name string
		args []string
	}{
		{name: "missing transaction", args: []string{"transactions", "lock", "clear", "--state-dir", stateDir, "--confirm", "tx-one"}},
		{name: "missing confirm", args: []string{"transactions", "lock", "clear", "--state-dir", stateDir, "--transaction", "tx-one"}},
		{name: "confirm mismatch", args: []string{"transactions", "lock", "clear", "--state-dir", stateDir, "--transaction", "tx-one", "--confirm", "tx-two"}},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result := runArcpub(t, tt.args...)
			assertExitCode(t, result, 64)
		})
	}
}

func TestTransactionsLockClearRejectsMismatch(t *testing.T) {
	stateDir := t.TempDir()
	writeE2ELockFile(t, stateDir, "tx-other")

	result := runTransactionsLockClear(t, stateDir, 1, "tx-want")

	assertContains(t, result.Stderr, "lock")
	assertFileExists(t, transactionLockPath(stateDir))
}

func TestTransactionsLockClearDeletesTerminalLockOnly(t *testing.T) {
	stateDir := t.TempDir()
	writeE2ETransactionJournal(t, stateDir, "tx-committed", "committed", time.Now())
	writeE2ELockFile(t, stateDir, "tx-committed")

	_, decoded := runTransactionsLockClearJSON(t, stateDir, 0, "tx-committed")

	if got := stringField(t, decoded, "status"); got != "cleared" {
		t.Fatalf("status = %q", got)
	}
	if got := stringField(t, decoded, "reason"); got != "cleared" {
		t.Fatalf("reason = %q", got)
	}
	if got := stringField(t, decoded, "postClearState"); got != "ready_for_publish" {
		t.Fatalf("postClearState = %q", got)
	}
	assertPathMissing(t, transactionLockPath(stateDir))
	assertFileExists(t, transactionJournalPath(stateDir, "tx-committed"))
}

func TestTransactionsLockClearRejectsActiveTransaction(t *testing.T) {
	stateDir := t.TempDir()
	writeE2ETransactionJournal(t, stateDir, "tx-pending", "pending", time.Now())
	writeE2ELockFile(t, stateDir, "tx-pending")

	result := runTransactionsLockClear(t, stateDir, 1, "tx-pending")

	assertContains(t, result.Stderr, "active")
	assertFileExists(t, transactionLockPath(stateDir))
	assertFileExists(t, transactionJournalPath(stateDir, "tx-pending"))
}

func TestTransactionsLockClearRollbackFailedPolicy(t *testing.T) {
	stateDir := t.TempDir()
	writeE2ETransactionJournal(t, stateDir, "tx-rollback-failed", "rollback_failed", time.Now())
	writeE2ELockFile(t, stateDir, "tx-rollback-failed")

	_, decoded := runTransactionsLockClearJSON(t, stateDir, 0, "tx-rollback-failed")

	if got := stringField(t, decoded, "postClearState"); got != "transaction_still_blocks_publish" {
		t.Fatalf("postClearState = %q", got)
	}
	assertPathMissing(t, transactionLockPath(stateDir))
	assertFileExists(t, transactionJournalPath(stateDir, "tx-rollback-failed"))
}

func TestTransactionsLockClearMissingJournalIsUnverified(t *testing.T) {
	stateDir := t.TempDir()
	writeE2ELockFile(t, stateDir, "tx-orphan")

	_, decoded := runTransactionsLockClearJSON(t, stateDir, 0, "tx-orphan")

	if got := stringField(t, decoded, "status"); got != "cleared" {
		t.Fatalf("status = %q", got)
	}
	if got := stringField(t, decoded, "postClearState"); got != "unverified_no_journal" {
		t.Fatalf("postClearState = %q", got)
	}
	assertPathMissing(t, transactionLockPath(stateDir))
	assertPathMissing(t, transactionJournalPath(stateDir, "tx-orphan"))
}

func TestTransactionsLockClearRefusesOperationLock(t *testing.T) {
	stateDir := t.TempDir()
	writeE2ETransactionJournal(t, stateDir, "tx-committed", "committed", time.Now())
	writeE2ELockFile(t, stateDir, "tx-committed")
	writeE2EOperationLock(t, stateDir, "publish")

	_, decoded := runTransactionsLockClearJSON(t, stateDir, 1, "tx-committed")

	if got := stringField(t, decoded, "reason"); got != "operation_lock_exists" {
		t.Fatalf("reason = %q", got)
	}
	assertFileExists(t, transactionLockPath(stateDir))
	assertFileExists(t, transactionJournalPath(stateDir, "tx-committed"))
	assertFileExists(t, transactionOperationLockPath(stateDir))
}

func TestTransactionsLockNoPathLeaksByDefault(t *testing.T) {
	stateDir := t.TempDir()
	writeE2ETransactionJournal(t, stateDir, "tx-committed", "committed", time.Now())
	writeE2ELockFile(t, stateDir, "tx-committed")

	result, _ := runTransactionsLockShowJSON(t, stateDir, 0)

	assertNoLocalPathLeak(t, result.Stdout, stateDir)
}

func TestTransactionsLockDiagnoseReportsBlockersReadOnly(t *testing.T) {
	stateDir := t.TempDir()
	writeE2ETransactionJournal(t, stateDir, "tx-failed", "failed", time.Now())
	writeE2EOperationLock(t, stateDir, "publish")

	result := runArcpub(t, "transactions", "diagnose", "--state-dir", stateDir, "--output", "json")
	assertExitCode(t, result, 0)
	assertContains(t, result.Stdout, `"kind": "transactions-diagnostics"`)
	assertContains(t, result.Stdout, `"publishBlocked": true`)
	assertContains(t, result.Stdout, `"operationLock":`)
	assertContains(t, result.Stdout, `"operation": "publish"`)
	assertContains(t, result.Stdout, `"kind": "failed_journal"`)
	assertContains(t, result.Stdout, `"kind": "operation_lock"`)
	if strings.Contains(result.Stdout, "token-one") {
		t.Fatalf("diagnostics leaked operation lock token: %s", result.Stdout)
	}
	assertNoLocalPathLeak(t, result.Stdout, stateDir)
	assertFileExists(t, transactionJournalPath(stateDir, "tx-failed"))
	assertFileExists(t, transactionOperationLockPath(stateDir))
}

func TestTransactionsLockDiagnoseRendersPartialErrors(t *testing.T) {
	stateDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(stateDir, "transactions"), []byte("not a directory"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	result := runArcpub(t, "transactions", "diagnose", "--state-dir", stateDir, "--output", "json")
	assertExitCode(t, result, 1)
	assertContains(t, result.Stdout, `"kind": "transactions-diagnostics"`)
	assertContains(t, result.Stdout, `"kind": "journal_directory_read_failed"`)
	assertContains(t, result.Stdout, `"reason": "journal_directory_read_failed"`)
	assertNoLocalPathLeak(t, result.Stdout, stateDir)
}

func runTransactionsLockShowJSON(t *testing.T, stateDir string, wantCode int) (commandResult, map[string]any) {
	t.Helper()
	result := runArcpub(t, "transactions", "lock", "show", "--state-dir", stateDir, "--output", "json")
	assertExitCode(t, result, wantCode)
	return result, assertJSON(t, result.Stdout)
}

func runTransactionsLockClearJSON(t *testing.T, stateDir string, wantCode int, id string) (commandResult, map[string]any) {
	t.Helper()
	result := runTransactionsLockClear(t, stateDir, wantCode, id, "--output", "json")
	return result, assertJSON(t, result.Stdout)
}

func runTransactionsLockClear(t *testing.T, stateDir string, wantCode int, id string, extra ...string) commandResult {
	t.Helper()
	args := []string{
		"transactions", "lock", "clear",
		"--state-dir", stateDir,
		"--transaction", id,
		"--confirm", id,
	}
	args = append(args, extra...)
	result := runArcpub(t, args...)
	assertExitCode(t, result, wantCode)
	return result
}

func writeE2ELockFile(t *testing.T, stateDir string, id string) {
	t.Helper()
	writeE2ERawLockFile(t, stateDir, "transaction="+id+"\npid=1\nstartedAt=2026-01-01T00:00:00Z\ncommand=publish\n")
}

func writeE2ERawLockFile(t *testing.T, stateDir string, data string) {
	t.Helper()
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(transactionLockPath(stateDir), []byte(data), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}

func transactionLockPath(stateDir string) string {
	return filepath.Join(stateDir, "publish.lock")
}

func writeE2EOperationLock(t *testing.T, stateDir string, operation string) {
	t.Helper()
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	content := "schemaVersion=1\noperation=" + operation + "\ntoken=token-one\npid=1\nstartedAt=2026-01-01T00:00:00Z\n"
	if err := os.WriteFile(transactionOperationLockPath(stateDir), []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}

func transactionOperationLockPath(stateDir string) string {
	return filepath.Join(stateDir, "operation.lock")
}
