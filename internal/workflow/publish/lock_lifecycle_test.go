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

	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
)

func TestShowTransactionLockStates(t *testing.T) {
	ctx := context.Background()

	t.Run("absent", func(t *testing.T) {
		result, err := InspectTransactionLock(ctx, t.TempDir())
		if err != nil {
			t.Fatalf("InspectTransactionLock() error = %v", err)
		}
		if result.Status != LockShowStatusAbsent {
			t.Fatalf("status = %q", result.Status)
		}
		if result.Reason != LockShowReasonLockAbsent {
			t.Fatalf("reason = %q", result.Reason)
		}
	})

	t.Run("present with journal", func(t *testing.T) {
		stateDir := t.TempDir()
		store := NewFileJournalStore(stateDir)
		writePruneJournal(t, store, "tx-committed", TransactionStatusCommitted, time.Unix(1, 0).UTC())
		writeLockFile(t, stateDir, "tx-committed")

		result, err := InspectTransactionLock(ctx, stateDir)
		if err != nil {
			t.Fatalf("InspectTransactionLock() error = %v", err)
		}
		if result.Status != LockShowStatusPresent || result.Reason != LockShowReasonLockPresent || !result.Journal.Present || result.Journal.Status != TransactionStatusCommitted {
			t.Fatalf("result = %#v", result)
		}
	})

	t.Run("journal missing", func(t *testing.T) {
		stateDir := t.TempDir()
		writeLockFile(t, stateDir, "tx-missing")

		result, err := InspectTransactionLock(ctx, stateDir)
		if err != nil {
			t.Fatalf("InspectTransactionLock() error = %v", err)
		}
		if result.Status != LockShowStatusJournalMissing || result.Reason != LockShowReasonJournalMissing || result.Journal.Present {
			t.Fatalf("result = %#v", result)
		}
	})

	t.Run("corrupt lock", func(t *testing.T) {
		stateDir := t.TempDir()
		writeRawLockFile(t, stateDir, "pid=1\n")

		result, err := InspectTransactionLock(ctx, stateDir)
		if err == nil {
			t.Fatal("InspectTransactionLock() error = nil")
		}
		if result.Status != LockShowStatusCorrupt || result.Reason != LockShowReasonLockCorrupt {
			t.Fatalf("result = %#v", result)
		}
	})

	t.Run("lock read failure", func(t *testing.T) {
		stateDir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(stateDir, "publish.lock"), 0o700); err != nil {
			t.Fatalf("MkdirAll() error = %v", err)
		}

		result, err := InspectTransactionLock(ctx, stateDir)
		if err == nil {
			t.Fatal("InspectTransactionLock() error = nil")
		}
		if result.Status != LockShowStatusFailed || result.Reason != LockShowReasonLockReadFailed {
			t.Fatalf("result = %#v", result)
		}
	})

	t.Run("corrupt journal", func(t *testing.T) {
		stateDir := t.TempDir()
		writeLockFile(t, stateDir, "tx-bad")
		txDir := filepath.Join(stateDir, "transactions")
		if err := os.MkdirAll(txDir, 0o700); err != nil {
			t.Fatalf("MkdirAll() error = %v", err)
		}
		if err := os.WriteFile(filepath.Join(txDir, "tx-bad.json"), []byte("{"), 0o600); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		result, err := InspectTransactionLock(ctx, stateDir)
		if err == nil {
			t.Fatal("InspectTransactionLock() error = nil")
		}
		if result.Status != LockShowStatusJournalCorrupt || result.Reason != LockShowReasonJournalCorrupt {
			t.Fatalf("result = %#v", result)
		}
	})

	t.Run("journal read failure", func(t *testing.T) {
		stateDir := t.TempDir()
		writeLockFile(t, stateDir, "tx-bad")
		txDir := filepath.Join(stateDir, "transactions")
		if err := os.MkdirAll(filepath.Join(txDir, "tx-bad.json"), 0o700); err != nil {
			t.Fatalf("MkdirAll() error = %v", err)
		}

		result, err := InspectTransactionLock(ctx, stateDir)
		if err == nil {
			t.Fatal("InspectTransactionLock() error = nil")
		}
		if result.Status != LockShowStatusFailed || result.Reason != LockShowReasonJournalReadFailed {
			t.Fatalf("result = %#v", result)
		}
	})
}

func TestShowTransactionLockFailedForInspectionErrors(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result, err := InspectTransactionLock(ctx, t.TempDir())
	if err == nil {
		t.Fatal("InspectTransactionLock(canceled) error = nil")
	}
	if result.Status != LockShowStatusFailed || result.Reason != LockShowReasonContextCanceled {
		t.Fatalf("canceled result = %#v", result)
	}

	deadlineCtx, cancelDeadline := context.WithDeadline(context.Background(), time.Unix(1, 0))
	defer cancelDeadline()
	result, err = InspectTransactionLock(deadlineCtx, t.TempDir())
	if err == nil {
		t.Fatal("InspectTransactionLock(deadline) error = nil")
	}
	if result.Status != LockShowStatusFailed || result.Reason != LockShowReasonContextDeadline {
		t.Fatalf("deadline result = %#v", result)
	}

	result, err = InspectTransactionLock(context.Background(), "")
	if err == nil {
		t.Fatal("InspectTransactionLock(empty state dir) error = nil")
	}
	if result.Status != LockShowStatusFailed || result.Reason != LockShowReasonStateDirMissing {
		t.Fatalf("empty state dir result = %#v", result)
	}
}

func TestServiceShowTransactionLockUsesReadOnlyInspection(t *testing.T) {
	result, err := New(Dependencies{Git: porttest.NewGit()}, Options{}).ShowTransactionLock(context.Background(), t.TempDir())
	if err != nil {
		t.Fatalf("ShowTransactionLock() error = %v", err)
	}
	if result.Status != LockShowStatusAbsent {
		t.Fatalf("status = %q", result.Status)
	}
}

func TestClearTransactionLockGuardrails(t *testing.T) {
	ctx := context.Background()
	service := New(Dependencies{Git: porttest.NewGit()}, Options{})

	tests := []struct {
		name       string
		opts       LockClearOptions
		wantReason LockClearReason
	}{
		{name: "missing transaction", opts: LockClearOptions{Confirm: "tx-one"}, wantReason: LockClearReasonMissingTransactionID},
		{name: "missing confirm", opts: LockClearOptions{TransactionID: "tx-one"}, wantReason: LockClearReasonMissingConfirmation},
		{name: "confirm mismatch", opts: LockClearOptions{TransactionID: "tx-one", Confirm: "tx-two"}, wantReason: LockClearReasonConfirmationMismatch},
		{name: "invalid transaction", opts: LockClearOptions{TransactionID: "../bad", Confirm: "../bad"}, wantReason: LockClearReasonInvalidTransactionID},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.ClearTransactionLock(ctx, t.TempDir(), tt.opts)
			if err == nil {
				t.Fatal("ClearTransactionLock() error = nil")
			}
			if result.Status != LockClearStatusRefused {
				t.Fatalf("result = %#v", result)
			}
			if result.Reason != tt.wantReason {
				t.Fatalf("reason = %q, want %q", result.Reason, tt.wantReason)
			}
		})
	}
}

func TestClearTransactionLockReasonCodes(t *testing.T) {
	ctx := context.Background()
	service := New(Dependencies{Git: porttest.NewGit()}, Options{})

	tests := []struct {
		name       string
		stateDir   string
		setup      func(t *testing.T, stateDir string)
		wantStatus LockClearStatus
		wantReason LockClearReason
	}{
		{
			name:       "state dir missing",
			stateDir:   "",
			wantStatus: LockClearStatusFailed,
			wantReason: LockClearReasonStateDirMissing,
		},
		{
			name:       "lock absent",
			wantStatus: LockClearStatusRefused,
			wantReason: LockClearReasonLockAbsent,
		},
		{
			name: "lock corrupt",
			setup: func(t *testing.T, stateDir string) {
				writeRawLockFile(t, stateDir, "pid=1\n")
			},
			wantStatus: LockClearStatusFailed,
			wantReason: LockClearReasonLockCorrupt,
		},
		{
			name: "lock read failed",
			setup: func(t *testing.T, stateDir string) {
				if err := os.MkdirAll(filepath.Join(stateDir, "publish.lock"), 0o700); err != nil {
					t.Fatalf("MkdirAll() error = %v", err)
				}
			},
			wantStatus: LockClearStatusFailed,
			wantReason: LockClearReasonLockReadFailed,
		},
		{
			name: "transaction mismatch",
			setup: func(t *testing.T, stateDir string) {
				writeLockFile(t, stateDir, "tx-other")
			},
			wantStatus: LockClearStatusRefused,
			wantReason: LockClearReasonTransactionMismatch,
		},
		{
			name: "journal corrupt",
			setup: func(t *testing.T, stateDir string) {
				writeLockFile(t, stateDir, "tx-one")
				txDir := filepath.Join(stateDir, "transactions")
				if err := os.MkdirAll(txDir, 0o700); err != nil {
					t.Fatalf("MkdirAll() error = %v", err)
				}
				if err := os.WriteFile(filepath.Join(txDir, "tx-one.json"), []byte("{"), 0o600); err != nil {
					t.Fatalf("WriteFile() error = %v", err)
				}
			},
			wantStatus: LockClearStatusFailed,
			wantReason: LockClearReasonJournalCorrupt,
		},
		{
			name: "journal read failed",
			setup: func(t *testing.T, stateDir string) {
				writeLockFile(t, stateDir, "tx-one")
				txDir := filepath.Join(stateDir, "transactions")
				if err := os.MkdirAll(filepath.Join(txDir, "tx-one.json"), 0o700); err != nil {
					t.Fatalf("MkdirAll() error = %v", err)
				}
			},
			wantStatus: LockClearStatusFailed,
			wantReason: LockClearReasonJournalReadFailed,
		},
		{
			name: "active transaction",
			setup: func(t *testing.T, stateDir string) {
				store := NewFileJournalStore(stateDir)
				writePruneJournal(t, store, "tx-one", TransactionStatusPending, time.Unix(1, 0).UTC())
				writeLockFile(t, stateDir, "tx-one")
			},
			wantStatus: LockClearStatusRefused,
			wantReason: LockClearReasonActiveTransaction,
		},
		{
			name: "cleared",
			setup: func(t *testing.T, stateDir string) {
				store := NewFileJournalStore(stateDir)
				writePruneJournal(t, store, "tx-one", TransactionStatusCommitted, time.Unix(1, 0).UTC())
				writeLockFile(t, stateDir, "tx-one")
			},
			wantStatus: LockClearStatusCleared,
			wantReason: LockClearReasonCleared,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			stateDir := tt.stateDir
			if stateDir == "" && tt.name != "state dir missing" {
				stateDir = t.TempDir()
			}
			if tt.setup != nil {
				tt.setup(t, stateDir)
			}

			result, err := service.ClearTransactionLock(ctx, stateDir, LockClearOptions{TransactionID: "tx-one", Confirm: "tx-one"})
			if err == nil && tt.wantStatus != LockClearStatusCleared {
				t.Fatal("ClearTransactionLock() error = nil")
			}
			if result.Status != tt.wantStatus || result.Reason != tt.wantReason {
				t.Fatalf("result = %#v, want status %q reason %q", result, tt.wantStatus, tt.wantReason)
			}
		})
	}
}

func TestClearTransactionLockPolicy(t *testing.T) {
	ctx := context.Background()
	service := New(Dependencies{Git: porttest.NewGit()}, Options{})

	t.Run("terminal journal", func(t *testing.T) {
		stateDir := t.TempDir()
		store := NewFileJournalStore(stateDir)
		writePruneJournal(t, store, "tx-committed", TransactionStatusCommitted, time.Unix(1, 0).UTC())
		writeLockFile(t, stateDir, "tx-committed")

		result, err := service.ClearTransactionLock(ctx, stateDir, LockClearOptions{TransactionID: "tx-committed", Confirm: "tx-committed"})
		if err != nil {
			t.Fatalf("ClearTransactionLock() error = %v", err)
		}
		if result.Status != LockClearStatusCleared || !result.LockCleared {
			t.Fatalf("result = %#v", result)
		}
		if result.PostClearState != LockPostClearReadyForPublish {
			t.Fatalf("post-clear state = %q", result.PostClearState)
		}
		assertLockMissing(t, stateDir)
		assertJournalExists(t, store, "tx-committed")
		assertHasPending(t, store, false, "")
	})

	t.Run("rolled back journal", func(t *testing.T) {
		stateDir := t.TempDir()
		store := NewFileJournalStore(stateDir)
		writePruneJournal(t, store, "tx-rolled-back", TransactionStatusRolledBack, time.Unix(1, 0).UTC())
		writeLockFile(t, stateDir, "tx-rolled-back")

		result, err := service.ClearTransactionLock(ctx, stateDir, LockClearOptions{TransactionID: "tx-rolled-back", Confirm: "tx-rolled-back"})
		if err != nil {
			t.Fatalf("ClearTransactionLock() error = %v", err)
		}
		if result.Status != LockClearStatusCleared || result.PostClearState != LockPostClearReadyForPublish {
			t.Fatalf("result = %#v", result)
		}
		assertLockMissing(t, stateDir)
		assertJournalExists(t, store, "tx-rolled-back")
		assertHasPending(t, store, false, "")
	})

	t.Run("rollback failed journal", func(t *testing.T) {
		stateDir := t.TempDir()
		store := NewFileJournalStore(stateDir)
		writePruneJournal(t, store, "tx-rollback-failed", TransactionStatusRollbackFailed, time.Unix(1, 0).UTC())
		writeLockFile(t, stateDir, "tx-rollback-failed")

		result, err := service.ClearTransactionLock(ctx, stateDir, LockClearOptions{TransactionID: "tx-rollback-failed", Confirm: "tx-rollback-failed"})
		if err != nil {
			t.Fatalf("ClearTransactionLock() error = %v", err)
		}
		if result.Status != LockClearStatusCleared {
			t.Fatalf("result = %#v", result)
		}
		if result.PostClearState != LockPostClearTransactionStillOpen {
			t.Fatalf("post-clear state = %q", result.PostClearState)
		}
		assertLockMissing(t, stateDir)
		assertJournalExists(t, store, "tx-rollback-failed")
		assertHasPending(t, store, true, "tx-rollback-failed")
	})

	t.Run("failed journal", func(t *testing.T) {
		stateDir := t.TempDir()
		store := NewFileJournalStore(stateDir)
		writePruneJournal(t, store, "tx-failed", TransactionStatusFailed, time.Unix(1, 0).UTC())
		writeLockFile(t, stateDir, "tx-failed")

		result, err := service.ClearTransactionLock(ctx, stateDir, LockClearOptions{TransactionID: "tx-failed", Confirm: "tx-failed"})
		if err != nil {
			t.Fatalf("ClearTransactionLock() error = %v", err)
		}
		if result.Status != LockClearStatusCleared || result.PostClearState != LockPostClearTransactionStillOpen {
			t.Fatalf("result = %#v", result)
		}
		assertLockMissing(t, stateDir)
		assertJournalExists(t, store, "tx-failed")
		assertHasPending(t, store, true, "tx-failed")
	})

	t.Run("missing journal", func(t *testing.T) {
		stateDir := t.TempDir()
		store := NewFileJournalStore(stateDir)
		writeLockFile(t, stateDir, "tx-orphan")

		result, err := service.ClearTransactionLock(ctx, stateDir, LockClearOptions{TransactionID: "tx-orphan", Confirm: "tx-orphan"})
		if err != nil {
			t.Fatalf("ClearTransactionLock() error = %v", err)
		}
		if result.Status != LockClearStatusCleared || len(result.Warnings) != 1 {
			t.Fatalf("result = %#v", result)
		}
		if result.PostClearState != LockPostClearUnverifiedNoJournal {
			t.Fatalf("post-clear state = %q", result.PostClearState)
		}
		assertLockMissing(t, stateDir)
		assertJournalMissing(t, store, "tx-orphan")
	})

	t.Run("active journal refused", func(t *testing.T) {
		stateDir := t.TempDir()
		store := NewFileJournalStore(stateDir)
		writePruneJournal(t, store, "tx-pending", TransactionStatusPending, time.Unix(1, 0).UTC())
		writeLockFile(t, stateDir, "tx-pending")

		result, err := service.ClearTransactionLock(ctx, stateDir, LockClearOptions{TransactionID: "tx-pending", Confirm: "tx-pending"})
		if err == nil {
			t.Fatal("ClearTransactionLock() error = nil")
		}
		if result.Status != LockClearStatusRefused {
			t.Fatalf("result = %#v", result)
		}
		assertLockExists(t, stateDir)
		assertJournalExists(t, store, "tx-pending")
	})

	t.Run("transaction mismatch refused", func(t *testing.T) {
		stateDir := t.TempDir()
		writeLockFile(t, stateDir, "tx-other")

		result, err := service.ClearTransactionLock(ctx, stateDir, LockClearOptions{TransactionID: "tx-want", Confirm: "tx-want"})
		if err == nil {
			t.Fatal("ClearTransactionLock() error = nil")
		}
		if result.Status != LockClearStatusRefused {
			t.Fatalf("result = %#v", result)
		}
		assertLockExists(t, stateDir)
	})

	t.Run("corrupt journal fails", func(t *testing.T) {
		stateDir := t.TempDir()
		writeLockFile(t, stateDir, "tx-bad")
		txDir := filepath.Join(stateDir, "transactions")
		if err := os.MkdirAll(txDir, 0o700); err != nil {
			t.Fatalf("MkdirAll() error = %v", err)
		}
		if err := os.WriteFile(filepath.Join(txDir, "tx-bad.json"), []byte("{"), 0o600); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		result, err := service.ClearTransactionLock(ctx, stateDir, LockClearOptions{TransactionID: "tx-bad", Confirm: "tx-bad"})
		if err == nil {
			t.Fatal("ClearTransactionLock() error = nil")
		}
		if result.Status != LockClearStatusFailed {
			t.Fatalf("result = %#v", result)
		}
		assertLockExists(t, stateDir)
	})
}

func TestClearTransactionLockStatusPolicyMatrix(t *testing.T) {
	ctx := context.Background()
	service := New(Dependencies{Git: porttest.NewGit()}, Options{})

	tests := []struct {
		id            TransactionID
		status        TransactionStatus
		wantCleared   bool
		wantPostClear LockPostClearState
	}{
		{id: "tx-pending", status: TransactionStatusPending},
		{id: "tx-preflighted", status: TransactionStatusPreflighted},
		{id: "tx-snapshotted", status: TransactionStatusSnapshotted},
		{id: "tx-committed-locally", status: TransactionStatusCommittedLocally},
		{id: "tx-candidates-pushed", status: TransactionStatusCandidatesPushed},
		{id: "tx-promoting", status: TransactionStatusPromoting},
		{id: "tx-branches-promoted", status: TransactionStatusBranchesPromoted},
		{id: "tx-tagging", status: TransactionStatusTagging},
		{id: "tx-committed", status: TransactionStatusCommitted, wantCleared: true, wantPostClear: LockPostClearReadyForPublish},
		{id: "tx-failed", status: TransactionStatusFailed, wantCleared: true, wantPostClear: LockPostClearTransactionStillOpen},
		{id: "tx-rolling-back", status: TransactionStatusRollingBack},
		{id: "tx-rolled-back", status: TransactionStatusRolledBack, wantCleared: true, wantPostClear: LockPostClearReadyForPublish},
		{id: "tx-rollback-failed", status: TransactionStatusRollbackFailed, wantCleared: true, wantPostClear: LockPostClearTransactionStillOpen},
		{id: "tx-unknown", status: TransactionStatus("unknown")},
		{id: "tx-empty", status: TransactionStatus("")},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.id.String(), func(t *testing.T) {
			stateDir := t.TempDir()
			store := NewFileJournalStore(stateDir)
			writePruneJournal(t, store, tt.id, tt.status, time.Unix(1, 0).UTC())
			writeLockFile(t, stateDir, tt.id)

			result, err := service.ClearTransactionLock(ctx, stateDir, LockClearOptions{TransactionID: tt.id, Confirm: tt.id})
			if tt.wantCleared {
				if err != nil {
					t.Fatalf("ClearTransactionLock() error = %v", err)
				}
				if result.Status != LockClearStatusCleared || result.Reason != LockClearReasonCleared || result.PostClearState != tt.wantPostClear {
					t.Fatalf("result = %#v", result)
				}
				assertLockMissing(t, stateDir)
				assertJournalExists(t, store, tt.id)
				return
			}
			if err == nil {
				t.Fatal("ClearTransactionLock() error = nil")
			}
			if result.Status != LockClearStatusRefused || result.Reason != LockClearReasonActiveTransaction {
				t.Fatalf("result = %#v", result)
			}
			if result.PostClearState != "" {
				t.Fatalf("post-clear state = %q, want empty for refused clear", result.PostClearState)
			}
			assertLockExists(t, stateDir)
			assertJournalExists(t, store, tt.id)
		})
	}
}

func TestClearTransactionLockDeleteAndSyncFailures(t *testing.T) {
	ctx := context.Background()
	service := New(Dependencies{Git: porttest.NewGit()}, Options{})

	t.Run("delete failed", func(t *testing.T) {
		stateDir := t.TempDir()
		store := NewFileJournalStore(stateDir)
		writePruneJournal(t, store, "tx-one", TransactionStatusCommitted, time.Unix(1, 0).UTC())
		writeLockFile(t, stateDir, "tx-one")
		service := service
		service.lockOps = transactionLockOps{
			remove: func(string) error {
				return errors.New("delete refused")
			},
		}

		result, err := service.ClearTransactionLock(ctx, stateDir, LockClearOptions{TransactionID: "tx-one", Confirm: "tx-one"})
		if err == nil {
			t.Fatal("ClearTransactionLock() error = nil")
		}
		if result.Status != LockClearStatusFailed || result.Reason != LockClearReasonDeleteFailed {
			t.Fatalf("result = %#v", result)
		}
		if result.PostClearState != "" {
			t.Fatalf("post-clear state = %q, want empty for failed clear", result.PostClearState)
		}
		assertLockExists(t, stateDir)
		assertJournalExists(t, store, "tx-one")
	})

	t.Run("sync failed", func(t *testing.T) {
		stateDir := t.TempDir()
		store := NewFileJournalStore(stateDir)
		writePruneJournal(t, store, "tx-one", TransactionStatusCommitted, time.Unix(1, 0).UTC())
		writeLockFile(t, stateDir, "tx-one")
		service := service
		service.lockOps = transactionLockOps{
			syncParent: func(string) error {
				return errors.New("sync refused")
			},
		}

		result, err := service.ClearTransactionLock(ctx, stateDir, LockClearOptions{TransactionID: "tx-one", Confirm: "tx-one"})
		if err == nil {
			t.Fatal("ClearTransactionLock() error = nil")
		}
		if result.Status != LockClearStatusFailed || result.Reason != LockClearReasonSyncFailed {
			t.Fatalf("result = %#v", result)
		}
		if !result.LockCleared {
			t.Fatal("LockCleared = false, want true")
		}
		if result.PostClearState != LockPostClearReadyForPublish {
			t.Fatalf("post-clear state = %q", result.PostClearState)
		}
		assertLockMissing(t, stateDir)
		assertJournalExists(t, store, "tx-one")

		retry, err := New(Dependencies{Git: porttest.NewGit()}, Options{}).ClearTransactionLock(ctx, stateDir, LockClearOptions{TransactionID: "tx-one", Confirm: "tx-one"})
		if err == nil {
			t.Fatal("retry ClearTransactionLock() error = nil")
		}
		if retry.Status != LockClearStatusRefused || retry.Reason != LockClearReasonLockAbsent {
			t.Fatalf("retry result = %#v", retry)
		}
	})
}

func TestClearTransactionLockRefusesChangedLockBeforeDelete(t *testing.T) {
	ctx := context.Background()
	service := New(Dependencies{Git: porttest.NewGit()}, Options{})
	stateDir := t.TempDir()
	store := NewFileJournalStore(stateDir)
	writePruneJournal(t, store, "tx-one", TransactionStatusCommitted, time.Unix(1, 0).UTC())
	writeLockFile(t, stateDir, "tx-one")

	service.lockOps = transactionLockOps{
		beforeRemove: func() {
			writeLockFile(t, stateDir, "tx-two")
		},
	}

	result, err := service.ClearTransactionLock(ctx, stateDir, LockClearOptions{TransactionID: "tx-one", Confirm: "tx-one"})
	if err == nil {
		t.Fatal("ClearTransactionLock() error = nil")
	}
	if result.Status != LockClearStatusFailed || result.Reason != LockClearReasonLockChanged {
		t.Fatalf("result = %#v", result)
	}
	info, err := readTransactionLock(filepath.Join(stateDir, "publish.lock"))
	if err != nil {
		t.Fatalf("readTransactionLock() error = %v", err)
	}
	if info.ID != "tx-two" {
		t.Fatalf("lock transaction = %q, want tx-two", info.ID)
	}
	assertJournalExists(t, store, "tx-one")
}

func TestClearTransactionLockRefusesDisappearedLockBeforeDelete(t *testing.T) {
	ctx := context.Background()
	service := New(Dependencies{Git: porttest.NewGit()}, Options{})
	stateDir := t.TempDir()
	store := NewFileJournalStore(stateDir)
	writePruneJournal(t, store, "tx-one", TransactionStatusCommitted, time.Unix(1, 0).UTC())
	writeLockFile(t, stateDir, "tx-one")

	service.lockOps = transactionLockOps{
		beforeRemove: func() {
			if err := os.Remove(filepath.Join(stateDir, "publish.lock")); err != nil {
				t.Fatalf("Remove() error = %v", err)
			}
		},
	}

	result, err := service.ClearTransactionLock(ctx, stateDir, LockClearOptions{TransactionID: "tx-one", Confirm: "tx-one"})
	if err == nil {
		t.Fatal("ClearTransactionLock() error = nil")
	}
	if result.Status != LockClearStatusFailed || result.Reason != LockClearReasonLockDisappeared {
		t.Fatalf("result = %#v", result)
	}
	if result.LockCleared || result.PostClearState != "" {
		t.Fatalf("partial result = %#v", result)
	}
	assertLockMissing(t, stateDir)
	assertJournalExists(t, store, "tx-one")
}

func TestClearTransactionLockRefusesCorruptLockBeforeDelete(t *testing.T) {
	ctx := context.Background()
	service := New(Dependencies{Git: porttest.NewGit()}, Options{})
	stateDir := t.TempDir()
	store := NewFileJournalStore(stateDir)
	writePruneJournal(t, store, "tx-one", TransactionStatusCommitted, time.Unix(1, 0).UTC())
	writeLockFile(t, stateDir, "tx-one")

	service.lockOps = transactionLockOps{
		beforeRemove: func() {
			writeRawLockFile(t, stateDir, "transaction=../bad\n")
		},
	}

	result, err := service.ClearTransactionLock(ctx, stateDir, LockClearOptions{TransactionID: "tx-one", Confirm: "tx-one"})
	if err == nil {
		t.Fatal("ClearTransactionLock() error = nil")
	}
	if result.Status != LockClearStatusFailed || result.Reason != LockClearReasonLockCorrupt {
		t.Fatalf("result = %#v", result)
	}
	if result.LockCleared || result.PostClearState != "" {
		t.Fatalf("partial result = %#v", result)
	}
	assertLockExists(t, stateDir)
	assertJournalExists(t, store, "tx-one")
}

func writeLockFile(t *testing.T, stateDir string, id TransactionID) {
	t.Helper()
	writeRawLockFile(t, stateDir, "transaction="+id.String()+"\npid=1\nstartedAt=2026-01-01T00:00:00Z\ncommand=publish\n")
}

func writeRawLockFile(t *testing.T, stateDir string, data string) {
	t.Helper()
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	path, err := transactionLockPath(stateDir)
	if err != nil {
		t.Fatalf("transactionLockPath() error = %v", err)
	}
	if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}

func assertLockExists(t *testing.T, stateDir string) {
	t.Helper()
	path, err := transactionLockPath(stateDir)
	if err != nil {
		t.Fatalf("transactionLockPath() error = %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("lock missing: %v", err)
	}
}

func assertLockMissing(t *testing.T, stateDir string) {
	t.Helper()
	path, err := transactionLockPath(stateDir)
	if err != nil {
		t.Fatalf("transactionLockPath() error = %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("lock exists or stat failed: %v", err)
	}
}

func assertHasPending(t *testing.T, store FileJournalStore, want bool, wantID TransactionID) {
	t.Helper()
	summary, ok, err := store.HasPending(context.Background())
	if err != nil {
		t.Fatalf("HasPending() error = %v", err)
	}
	if ok != want {
		t.Fatalf("HasPending() ok = %v, want %v summary=%#v", ok, want, summary)
	}
	if want && summary.ID != wantID {
		t.Fatalf("pending id = %q, want %q", summary.ID, wantID)
	}
}
