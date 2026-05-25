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
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestTransactionsPruneDryRunDoesNotDelete(t *testing.T) {
	stateDir := t.TempDir()
	writeE2ETransactionJournal(t, stateDir, "tx-committed", "committed", time.Now().Add(-48*time.Hour))
	writeE2ETransactionJournal(t, stateDir, "tx-rolled-back", "rolled_back", time.Now().Add(-48*time.Hour))

	result, decoded := runTransactionsPruneJSON(t, stateDir, 0, "--dry-run", "--status", "committed")

	if got := stringField(t, decoded, "status"); got != "dry_run" {
		t.Fatalf("status = %q\n%s", got, result.Stdout)
	}
	if got := floatField(t, decoded, "matchedCount"); got != 1 {
		t.Fatalf("matchedCount = %v\n%s", got, result.Stdout)
	}
	assertFileExists(t, transactionJournalPath(stateDir, "tx-committed"))
	assertFileExists(t, transactionJournalPath(stateDir, "tx-rolled-back"))
}

func TestTransactionsPruneDeletesCommitted(t *testing.T) {
	stateDir := t.TempDir()
	writeE2ETransactionJournal(t, stateDir, "tx-committed", "committed", time.Now().Add(-48*time.Hour))

	_, decoded := runTransactionsPruneJSON(t, stateDir, 0, "--status", "committed")

	if got := stringField(t, decoded, "status"); got != "completed" {
		t.Fatalf("status = %q", got)
	}
	assertPathMissing(t, transactionJournalPath(stateDir, "tx-committed"))
}

func TestTransactionsPruneDeletesRolledBack(t *testing.T) {
	stateDir := t.TempDir()
	writeE2ETransactionJournal(t, stateDir, "tx-rolled-back", "rolled_back", time.Now().Add(-48*time.Hour))

	runTransactionsPruneJSON(t, stateDir, 0, "--status", "rolled_back")

	assertPathMissing(t, transactionJournalPath(stateDir, "tx-rolled-back"))
}

func TestTransactionsPrunePreservesRollbackFailed(t *testing.T) {
	stateDir := t.TempDir()
	writeE2ETransactionJournal(t, stateDir, "tx-rollback-failed", "rollback_failed", time.Now().Add(-48*time.Hour))

	_, decoded := runTransactionsPruneJSON(t, stateDir, 0, "--status", "committed")

	if got := floatField(t, decoded, "skippedCount"); got != 1 {
		t.Fatalf("skippedCount = %v", got)
	}
	assertFileExists(t, transactionJournalPath(stateDir, "tx-rollback-failed"))
}

func TestTransactionsPrunePreservesPending(t *testing.T) {
	stateDir := t.TempDir()
	writeE2ETransactionJournal(t, stateDir, "tx-pending", "pending", time.Now().Add(-48*time.Hour))

	runTransactionsPruneJSON(t, stateDir, 0, "--status", "committed")

	assertFileExists(t, transactionJournalPath(stateDir, "tx-pending"))
}

func TestTransactionsPruneCorruptedJournalFailsClosed(t *testing.T) {
	stateDir := t.TempDir()
	txDir := filepath.Join(stateDir, "transactions")
	writeE2ETransactionJournal(t, stateDir, "tx-committed", "committed", time.Now().Add(-48*time.Hour))
	if err := os.WriteFile(filepath.Join(txDir, "tx-bad.json"), []byte("{"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	result := runTransactionsPrune(t, stateDir, 1, "--status", "committed")

	assertContains(t, result.Stderr, "corrupt")
	assertFileExists(t, transactionJournalPath(stateDir, "tx-committed"))
}

func TestTransactionsPruneOlderThan(t *testing.T) {
	stateDir := t.TempDir()
	now := time.Now().UTC()
	writeE2ETransactionJournal(t, stateDir, "tx-old", "committed", now.Add(-48*time.Hour))
	writeE2ETransactionJournal(t, stateDir, "tx-fresh", "committed", now)

	_, decoded := runTransactionsPruneJSON(t, stateDir, 0, "--older-than", "24h")

	if got := floatField(t, decoded, "deletedCount"); got != 1 {
		t.Fatalf("deletedCount = %v", got)
	}
	assertPathMissing(t, transactionJournalPath(stateDir, "tx-old"))
	assertFileExists(t, transactionJournalPath(stateDir, "tx-fresh"))
}

func TestTransactionsPruneNoPathLeaksByDefault(t *testing.T) {
	stateDir := t.TempDir()
	writeE2ETransactionJournal(t, stateDir, "tx-committed", "committed", time.Now().Add(-48*time.Hour))

	result, _ := runTransactionsPruneJSON(t, stateDir, 0, "--dry-run", "--status", "committed")

	assertNoLocalPathLeak(t, result.Stdout, stateDir)
}

func TestTransactionsPruneBlockedByLock(t *testing.T) {
	stateDir := t.TempDir()
	writeE2ETransactionJournal(t, stateDir, "tx-committed", "committed", time.Now().Add(-48*time.Hour))
	if err := os.WriteFile(filepath.Join(stateDir, "publish.lock"), []byte("transaction=tx-active\npid=1\nstartedAt=now\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	actual := runTransactionsPrune(t, stateDir, 1, "--status", "committed")
	assertContains(t, actual.Stderr, "lock")
	assertFileExists(t, transactionJournalPath(stateDir, "tx-committed"))

	runTransactionsPruneJSON(t, stateDir, 0, "--dry-run", "--status", "committed")
	assertFileExists(t, transactionJournalPath(stateDir, "tx-committed"))
}

func runTransactionsPruneJSON(t *testing.T, stateDir string, wantCode int, extra ...string) (commandResult, map[string]any) {
	t.Helper()
	result := runTransactionsPrune(t, stateDir, wantCode, append(extra, "--output", "json")...)
	return result, assertJSON(t, result.Stdout)
}

func runTransactionsPrune(t *testing.T, stateDir string, wantCode int, extra ...string) commandResult {
	t.Helper()
	args := []string{"transactions", "prune", "--state-dir", stateDir}
	args = append(args, extra...)
	result := runArcpub(t, args...)
	assertExitCode(t, result, wantCode)
	return result
}

func writeE2ETransactionJournal(t *testing.T, stateDir string, id string, status string, updatedAt time.Time) {
	t.Helper()
	txDir := filepath.Join(stateDir, "transactions")
	if err := os.MkdirAll(txDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	data := fmt.Sprintf(
		`{"schemaVersion":1,"id":%q,"status":%q,"version":"v0.1.0","startedAt":%q,"updatedAt":%q,"modules":[]}`+"\n",
		id,
		status,
		updatedAt.UTC().Format(time.RFC3339Nano),
		updatedAt.UTC().Format(time.RFC3339Nano),
	)
	if err := os.WriteFile(transactionJournalPath(stateDir, id), []byte(data), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}

func transactionJournalPath(stateDir string, id string) string {
	return filepath.Join(stateDir, "transactions", id+".json")
}
