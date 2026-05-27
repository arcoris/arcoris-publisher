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
}
