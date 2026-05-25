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

package publish

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
)

func TestPruneDryRunCommittedDoesNotDelete(t *testing.T) {
	ctx := context.Background()
	store := NewFileJournalStore(t.TempDir())
	writePruneJournal(t, store, "tx-committed", TransactionStatusCommitted, time.Unix(1, 0).UTC())

	result, err := store.Prune(ctx, PruneOptions{Statuses: []TransactionStatus{TransactionStatusCommitted}, DryRun: true})
	if err != nil {
		t.Fatalf("Prune() error = %v", err)
	}
	if result.Status != PruneStatusDryRun || len(result.Matched) != 1 || len(result.Deleted) != 0 {
		t.Fatalf("dry-run result = %#v", result)
	}
	assertJournalExists(t, store, "tx-committed")
}

func TestPruneDeletesCommittedAndRolledBack(t *testing.T) {
	ctx := context.Background()
	store := NewFileJournalStore(t.TempDir())
	writePruneJournal(t, store, "tx-committed", TransactionStatusCommitted, time.Unix(1, 0).UTC())
	writePruneJournal(t, store, "tx-rolled-back", TransactionStatusRolledBack, time.Unix(1, 0).UTC())

	result, err := store.Prune(ctx, PruneOptions{
		Statuses: []TransactionStatus{TransactionStatusCommitted, TransactionStatusRolledBack},
	})
	if err != nil {
		t.Fatalf("Prune() error = %v", err)
	}
	if result.Status != PruneStatusCompleted || len(result.Deleted) != 2 {
		t.Fatalf("prune result = %#v", result)
	}
	assertJournalMissing(t, store, "tx-committed")
	assertJournalMissing(t, store, "tx-rolled-back")
}

func TestPrunePreservesNonTerminalAndRollbackFailed(t *testing.T) {
	ctx := context.Background()
	store := NewFileJournalStore(t.TempDir())
	for _, tt := range []struct {
		id     TransactionID
		status TransactionStatus
	}{
		{id: "tx-pending", status: TransactionStatusPending},
		{id: "tx-failed", status: TransactionStatusFailed},
		{id: "tx-rollback-failed", status: TransactionStatusRollbackFailed},
	} {
		writePruneJournal(t, store, tt.id, tt.status, time.Unix(1, 0).UTC())
	}

	result, err := store.Prune(ctx, PruneOptions{OlderThan: time.Hour, Now: time.Unix(10000, 0).UTC()})
	if err != nil {
		t.Fatalf("Prune() error = %v", err)
	}
	if len(result.Deleted) != 0 || len(result.Skipped) != 3 {
		t.Fatalf("prune result = %#v", result)
	}
	assertJournalExists(t, store, "tx-pending")
	assertJournalExists(t, store, "tx-failed")
	assertJournalExists(t, store, "tx-rollback-failed")
}

func TestPruneCorruptedJournalFailsClosed(t *testing.T) {
	ctx := context.Background()
	store := NewFileJournalStore(t.TempDir())
	writePruneJournal(t, store, "tx-committed", TransactionStatusCommitted, time.Unix(1, 0).UTC())
	if err := os.WriteFile(filepath.Join(store.transactionsDir(), "tx-bad.json"), []byte("{"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := store.Prune(ctx, PruneOptions{Statuses: []TransactionStatus{TransactionStatusCommitted}}); err == nil {
		t.Fatal("Prune() error = nil")
	}
	assertJournalExists(t, store, "tx-committed")
}

func TestPruneUnsafeJournalIDFailsClosed(t *testing.T) {
	ctx := context.Background()
	store := NewFileJournalStore(t.TempDir())
	if err := os.MkdirAll(store.transactionsDir(), 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	data := []byte(`{"schemaVersion":1,"id":"../escape","status":"committed","startedAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}`)
	if err := os.WriteFile(filepath.Join(store.transactionsDir(), "tx-unsafe.json"), data, 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := store.Prune(ctx, PruneOptions{Statuses: []TransactionStatus{TransactionStatusCommitted}}); err == nil {
		t.Fatal("Prune() error = nil")
	}
}

func TestPruneOlderThanUsesUpdatedAtThenStartedAt(t *testing.T) {
	ctx := context.Background()
	store := NewFileJournalStore(t.TempDir())
	now := time.Unix(10000, 0).UTC()
	writePruneJournal(t, store, "tx-old", TransactionStatusCommitted, now.Add(-48*time.Hour))
	writePruneJournal(t, store, "tx-fresh", TransactionStatusCommitted, now.Add(-time.Hour))
	startOnly := pruneJournal("tx-start-only", TransactionStatusCommitted, now.Add(-48*time.Hour))
	startOnly.UpdatedAt = time.Time{}
	if err := store.Create(ctx, startOnly); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	result, err := store.Prune(ctx, PruneOptions{OlderThan: 24 * time.Hour, Now: now})
	if err != nil {
		t.Fatalf("Prune() error = %v", err)
	}
	if len(result.Deleted) != 2 {
		t.Fatalf("deleted = %#v", result.Deleted)
	}
	assertJournalMissing(t, store, "tx-old")
	assertJournalExists(t, store, "tx-fresh")
	assertJournalMissing(t, store, "tx-start-only")
}

func TestPruneRejectsUnsafeStatusAndNoFilter(t *testing.T) {
	ctx := context.Background()
	store := NewFileJournalStore(t.TempDir())
	writePruneJournal(t, store, "tx-committed", TransactionStatusCommitted, time.Unix(1, 0).UTC())

	if _, err := store.Prune(ctx, PruneOptions{Statuses: []TransactionStatus{TransactionStatusRollbackFailed}}); err == nil {
		t.Fatal("Prune() rollback_failed error = nil")
	}
	if _, err := store.Prune(ctx, PruneOptions{}); err == nil {
		t.Fatal("Prune() no-filter error = nil")
	}
}

func TestPruneServiceLockPolicy(t *testing.T) {
	ctx := context.Background()
	stateDir := t.TempDir()
	store := NewFileJournalStore(stateDir)
	writePruneJournal(t, store, "tx-committed", TransactionStatusCommitted, time.Unix(1, 0).UTC())
	lock, err := acquireTransactionLock(ctx, stateDir, "tx-active", time.Unix(2, 0).UTC())
	if err != nil {
		t.Fatalf("acquireTransactionLock() error = %v", err)
	}
	defer lock.Release()

	service := New(Dependencies{Git: porttest.NewGit()}, Options{})
	if _, err := service.PruneTransactions(ctx, stateDir, PruneOptions{Statuses: []TransactionStatus{TransactionStatusCommitted}}); err == nil {
		t.Fatal("PruneTransactions() with lock error = nil")
	}
	result, err := service.PruneTransactions(ctx, stateDir, PruneOptions{Statuses: []TransactionStatus{TransactionStatusCommitted}, DryRun: true})
	if err != nil {
		t.Fatalf("PruneTransactions() dry-run error = %v", err)
	}
	if result.Status != PruneStatusDryRun || len(result.Matched) != 1 {
		t.Fatalf("dry-run result = %#v", result)
	}
	assertJournalExists(t, store, "tx-committed")
}

func writePruneJournal(t *testing.T, store FileJournalStore, id TransactionID, status TransactionStatus, at time.Time) {
	t.Helper()
	if err := store.Create(context.Background(), pruneJournal(id, status, at)); err != nil {
		t.Fatalf("Create(%s) error = %v", id, err)
	}
}

func pruneJournal(id TransactionID, status TransactionStatus, at time.Time) TransactionJournal {
	return TransactionJournal{
		SchemaVersion: transactionSchemaVersion,
		ID:            id,
		Status:        status,
		Version:       "v0.1.0",
		StartedAt:     at,
		UpdatedAt:     at,
	}
}

func assertJournalExists(t *testing.T, store FileJournalStore, id TransactionID) {
	t.Helper()
	path, err := store.journalPath(id)
	if err != nil {
		t.Fatalf("journalPath() error = %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("journal %s missing: %v", id, err)
	}
}

func assertJournalMissing(t *testing.T, store FileJournalStore, id TransactionID) {
	t.Helper()
	path, err := store.journalPath(id)
	if err != nil {
		t.Fatalf("journalPath() error = %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("journal %s exists or stat failed: %v", id, err)
	}
}
