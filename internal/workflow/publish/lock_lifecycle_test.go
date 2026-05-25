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

func TestShowTransactionLockStates(t *testing.T) {
	ctx := context.Background()
	service := New(Dependencies{Git: porttest.NewGit()}, Options{})

	t.Run("absent", func(t *testing.T) {
		result, err := service.ShowTransactionLock(ctx, t.TempDir())
		if err != nil {
			t.Fatalf("ShowTransactionLock() error = %v", err)
		}
		if result.Status != LockShowStatusAbsent {
			t.Fatalf("status = %q", result.Status)
		}
	})

	t.Run("present with journal", func(t *testing.T) {
		stateDir := t.TempDir()
		store := NewFileJournalStore(stateDir)
		writePruneJournal(t, store, "tx-committed", TransactionStatusCommitted, time.Unix(1, 0).UTC())
		writeLockFile(t, stateDir, "tx-committed")

		result, err := service.ShowTransactionLock(ctx, stateDir)
		if err != nil {
			t.Fatalf("ShowTransactionLock() error = %v", err)
		}
		if result.Status != LockShowStatusPresent || !result.Journal.Present || result.Journal.Status != TransactionStatusCommitted {
			t.Fatalf("result = %#v", result)
		}
	})

	t.Run("journal missing", func(t *testing.T) {
		stateDir := t.TempDir()
		writeLockFile(t, stateDir, "tx-missing")

		result, err := service.ShowTransactionLock(ctx, stateDir)
		if err != nil {
			t.Fatalf("ShowTransactionLock() error = %v", err)
		}
		if result.Status != LockShowStatusJournalMissing || result.Journal.Present {
			t.Fatalf("result = %#v", result)
		}
	})

	t.Run("corrupt lock", func(t *testing.T) {
		stateDir := t.TempDir()
		writeRawLockFile(t, stateDir, "pid=1\n")

		result, err := service.ShowTransactionLock(ctx, stateDir)
		if err == nil {
			t.Fatal("ShowTransactionLock() error = nil")
		}
		if result.Status != LockShowStatusCorrupt {
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

		result, err := service.ShowTransactionLock(ctx, stateDir)
		if err == nil {
			t.Fatal("ShowTransactionLock() error = nil")
		}
		if result.Status != LockShowStatusJournalCorrupt {
			t.Fatalf("result = %#v", result)
		}
	})
}

func TestClearTransactionLockGuardrails(t *testing.T) {
	ctx := context.Background()
	service := New(Dependencies{Git: porttest.NewGit()}, Options{})

	tests := []struct {
		name string
		opts LockClearOptions
	}{
		{name: "missing transaction", opts: LockClearOptions{Confirm: "tx-one"}},
		{name: "missing confirm", opts: LockClearOptions{TransactionID: "tx-one"}},
		{name: "confirm mismatch", opts: LockClearOptions{TransactionID: "tx-one", Confirm: "tx-two"}},
		{name: "invalid transaction", opts: LockClearOptions{TransactionID: "../bad", Confirm: "../bad"}},
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
		assertLockMissing(t, stateDir)
		assertJournalExists(t, store, "tx-committed")
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
		assertLockMissing(t, stateDir)
		assertJournalExists(t, store, "tx-rollback-failed")
	})

	t.Run("missing journal", func(t *testing.T) {
		stateDir := t.TempDir()
		writeLockFile(t, stateDir, "tx-orphan")

		result, err := service.ClearTransactionLock(ctx, stateDir, LockClearOptions{TransactionID: "tx-orphan", Confirm: "tx-orphan"})
		if err != nil {
			t.Fatalf("ClearTransactionLock() error = %v", err)
		}
		if result.Status != LockClearStatusCleared || len(result.Warnings) != 1 {
			t.Fatalf("result = %#v", result)
		}
		assertLockMissing(t, stateDir)
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
