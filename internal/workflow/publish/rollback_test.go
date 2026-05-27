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
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"arcoris.dev/arcoris-publisher/internal/ports/git"
	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
)

func TestRollbackRefusesMovedFinalBranch(t *testing.T) {
	ctx := context.Background()
	fakeGit := porttest.NewGit()
	ref := "refs/heads/main"
	worktree := "/worktree"
	fakeGit.RemoteRefHashes[porttest.RemoteRefKeyForRepo(worktree, "origin", ref)] = "someone-else"

	runner := transactionRunner{
		service: New(Dependencies{Git: fakeGit}, Options{}),
		store:   NewFileJournalStore(t.TempDir()),
		journal: TransactionJournal{
			ID:     "tx-test",
			Status: TransactionStatusBranchesPromoted,
			Modules: []ModuleTransactionState{{
				Module:              "foundation",
				Repository:          "arcoris/foundation",
				WorktreeDir:         worktree,
				FinalBranchRef:      ref,
				CreatedCommit:       "created",
				RemoteBaseCommit:    "base",
				RemoteBaseExists:    true,
				FinalBranchPromoted: true,
			}},
		},
	}

	if err := runner.rollback(ctx); err == nil {
		t.Fatal("rollback() error = nil")
	}
	if runner.journal.Status != TransactionStatusRollbackFailed {
		t.Fatalf("status = %q", runner.journal.Status)
	}
	if len(runner.journal.ManualActions) != 1 {
		t.Fatalf("manual actions = %+v", runner.journal.ManualActions)
	}
	assertCallAbsent(t, fakeGit.Calls, "push")
}

func TestRollbackTreatsMissingCandidateAsAlreadyDeleted(t *testing.T) {
	ctx := context.Background()
	fakeGit := porttest.NewGit()
	runner := transactionRunner{
		service: New(Dependencies{Git: fakeGit}, Options{}),
		store:   NewFileJournalStore(t.TempDir()),
		journal: TransactionJournal{
			ID:     "tx-test",
			Status: TransactionStatusCandidatesPushed,
			Modules: []ModuleTransactionState{{
				Module:             "foundation",
				Repository:         "arcoris/foundation",
				WorktreeDir:        "/worktree",
				CandidateBranchRef: "refs/heads/arcpub/tx/tx-test/foundation",
				CandidatePushed:    true,
			}},
		},
	}

	if err := runner.rollback(ctx); err != nil {
		t.Fatalf("rollback() error = %v", err)
	}
	if runner.journal.Status != TransactionStatusRolledBack {
		t.Fatalf("status = %q", runner.journal.Status)
	}
	if runner.journal.Modules[0].CandidatePushed {
		t.Fatal("candidate still marked pushed")
	}
}

func TestRollbackRestoresBranchWithExactLease(t *testing.T) {
	ctx := context.Background()
	fakeGit := porttest.NewGit()
	ref := "refs/heads/main"
	worktree := "/worktree"
	fakeGit.RemoteRefHashes[porttest.RemoteRefKeyForRepo(worktree, "origin", ref)] = "created"

	runner := transactionRunner{
		service: New(Dependencies{Git: fakeGit}, Options{}),
		store:   NewFileJournalStore(t.TempDir()),
		journal: TransactionJournal{
			ID:     "tx-test",
			Status: TransactionStatusBranchesPromoted,
			Modules: []ModuleTransactionState{{
				Module:              "foundation",
				Repository:          "arcoris/foundation",
				WorktreeDir:         worktree,
				FinalBranchRef:      ref,
				CreatedCommit:       "created",
				RemoteBaseCommit:    "base",
				RemoteBaseExists:    true,
				FinalBranchPromoted: true,
			}},
		},
	}

	if err := runner.rollback(ctx); err != nil {
		t.Fatalf("rollback() error = %v", err)
	}
	push := findCall(fakeGit.Calls, "push")
	if push.ForceWithLeaseRef != ref || push.ForceWithLeaseExpect != git.CommitHash("created") {
		t.Fatalf("lease = ref %q expect %q", push.ForceWithLeaseRef, push.ForceWithLeaseExpect)
	}
}

func TestRollbackTransactionRefusesCommittedJournal(t *testing.T) {
	ctx := context.Background()
	stateDir := t.TempDir()
	store := NewFileJournalStore(stateDir)
	if err := store.Create(ctx, TransactionJournal{
		ID:        "tx-test",
		Status:    TransactionStatusCommitted,
		StartedAt: time.Unix(1, 0).UTC(),
		UpdatedAt: time.Unix(1, 0).UTC(),
	}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	_, err := New(Dependencies{Git: porttest.NewGit()}, Options{}).RollbackTransaction(ctx, stateDir, "tx-test")
	if err == nil {
		t.Fatal("RollbackTransaction() error = nil")
	}
}

func TestRollbackTransactionRefusesExistingLock(t *testing.T) {
	ctx := context.Background()
	stateDir := t.TempDir()
	store := NewFileJournalStore(stateDir)
	if err := store.Create(ctx, TransactionJournal{
		ID:        "tx-test",
		Status:    TransactionStatusFailed,
		StartedAt: time.Unix(1, 0).UTC(),
		UpdatedAt: time.Unix(1, 0).UTC(),
	}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	lock, err := acquireTransactionLock(ctx, stateDir, "tx-test", time.Unix(1, 0).UTC(), transactionLockOps{})
	if err != nil {
		t.Fatalf("acquireTransactionLock() error = %v", err)
	}
	defer func() { _, _ = lock.Release() }()

	_, err = New(Dependencies{Git: porttest.NewGit()}, Options{}).RollbackTransaction(ctx, stateDir, "tx-test")
	if err == nil {
		t.Fatal("RollbackTransaction() error = nil")
	}
	if _, err := os.Stat(operationLockPath(stateDir)); !os.IsNotExist(err) {
		t.Fatalf("operation lock exists or stat failed: %v", err)
	}
}

func TestRollbackTransactionRefusesExistingOperationLock(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, stateDir string)
	}{
		{
			name: "valid failed journal",
			setup: func(t *testing.T, stateDir string) {
				if err := NewFileJournalStore(stateDir).Create(context.Background(), TransactionJournal{
					ID:        "tx-test",
					Status:    TransactionStatusFailed,
					StartedAt: time.Unix(1, 0).UTC(),
					UpdatedAt: time.Unix(1, 0).UTC(),
				}); err != nil {
					t.Fatalf("Create() error = %v", err)
				}
			},
		},
		{
			name: "missing journal",
		},
		{
			name: "corrupt journal",
			setup: func(t *testing.T, stateDir string) {
				txDir := filepath.Join(stateDir, "transactions")
				if err := os.MkdirAll(txDir, 0o700); err != nil {
					t.Fatalf("MkdirAll() error = %v", err)
				}
				if err := os.WriteFile(filepath.Join(txDir, "tx-test.json"), []byte("{"), 0o600); err != nil {
					t.Fatalf("WriteFile() error = %v", err)
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			stateDir := t.TempDir()
			if tt.setup != nil {
				tt.setup(t, stateDir)
			}
			writeOperationLockFile(t, stateDir, operationLockPublish, "other-token")

			journal, err := New(Dependencies{Git: porttest.NewGit()}, Options{}).RollbackTransaction(context.Background(), stateDir, "tx-test")
			if err == nil {
				t.Fatal("RollbackTransaction() error = nil")
			}
			if !errors.Is(err, errOperationLockExists) {
				t.Fatalf("RollbackTransaction() error = %v, want operation lock exists", err)
			}
			if tt.setup == nil && journal.ID != "" {
				t.Fatalf("journal = %#v, want zero value for missing journal blocked before load", journal)
			}
			if _, err := os.Stat(operationLockPath(stateDir)); err != nil {
				t.Fatalf("operation lock missing: %v", err)
			}
		})
	}
}

func TestReadOnlyTransactionsIgnoreOperationLock(t *testing.T) {
	ctx := context.Background()
	stateDir := t.TempDir()
	store := NewFileJournalStore(stateDir)
	if err := store.Create(ctx, TransactionJournal{
		ID:        "tx-test",
		Status:    TransactionStatusFailed,
		StartedAt: time.Unix(1, 0).UTC(),
		UpdatedAt: time.Unix(1, 0).UTC(),
	}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	writeOperationLockFile(t, stateDir, operationLockPublish, "other-token")
	service := New(Dependencies{Git: porttest.NewGit()}, Options{})

	summaries, err := service.ListTransactions(ctx, stateDir)
	if err != nil {
		t.Fatalf("ListTransactions() error = %v", err)
	}
	if len(summaries) != 1 || summaries[0].ID != "tx-test" {
		t.Fatalf("summaries = %#v", summaries)
	}
	journal, err := service.ShowTransaction(ctx, stateDir, "tx-test")
	if err != nil {
		t.Fatalf("ShowTransaction() error = %v", err)
	}
	if journal.ID != "tx-test" {
		t.Fatalf("journal = %#v", journal)
	}
	if _, err := os.Stat(operationLockPath(stateDir)); err != nil {
		t.Fatalf("operation lock missing: %v", err)
	}
}

func TestRollbackTransactionLoadsJournalAfterOperationLock(t *testing.T) {
	ctx := context.Background()
	stateDir := t.TempDir()
	store := NewFileJournalStore(stateDir)
	if err := store.Create(ctx, TransactionJournal{
		ID:        "tx-test",
		Status:    TransactionStatusFailed,
		StartedAt: time.Unix(1, 0).UTC(),
		UpdatedAt: time.Unix(1, 0).UTC(),
	}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	service := New(Dependencies{Git: porttest.NewGit()}, Options{})
	ops := testOperationLockOps()
	mutated := false
	ops.syncParent = func(string) error {
		if !mutated {
			mutated = true
			if err := store.Update(ctx, TransactionJournal{
				ID:        "tx-test",
				Status:    TransactionStatusRolledBack,
				StartedAt: time.Unix(1, 0).UTC(),
				UpdatedAt: time.Unix(2, 0).UTC(),
			}); err != nil {
				t.Fatalf("Update() error = %v", err)
			}
		}
		return nil
	}
	service.operationLockOps = ops

	journal, err := service.RollbackTransaction(ctx, stateDir, "tx-test")
	if err != nil {
		t.Fatalf("RollbackTransaction() error = %v", err)
	}
	if journal.Status != TransactionStatusRolledBack {
		t.Fatalf("journal status = %q, want post-acquire rolled_back", journal.Status)
	}
}

func TestRollbackTransactionSurfacesOperationLockReleaseFailure(t *testing.T) {
	ctx := context.Background()
	stateDir := t.TempDir()
	if err := NewFileJournalStore(stateDir).Create(ctx, TransactionJournal{
		ID:        "tx-test",
		Status:    TransactionStatusRolledBack,
		StartedAt: time.Unix(1, 0).UTC(),
		UpdatedAt: time.Unix(1, 0).UTC(),
	}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	service := New(Dependencies{Git: porttest.NewGit()}, Options{})
	ops := testOperationLockOps()
	failSync := false
	ops.beforeRemove = func() { failSync = true }
	ops.syncParent = func(string) error {
		if failSync {
			return errors.New("operation sync refused")
		}
		return nil
	}
	service.operationLockOps = ops

	journal, err := service.RollbackTransaction(ctx, stateDir, "tx-test")
	if err == nil {
		t.Fatal("RollbackTransaction() error = nil")
	}
	if !errors.Is(err, errOperationLockSyncFailed) {
		t.Fatalf("RollbackTransaction() error = %v, want operation sync failure", err)
	}
	if journal.Status != TransactionStatusRolledBack {
		t.Fatalf("journal = %#v", journal)
	}
}
