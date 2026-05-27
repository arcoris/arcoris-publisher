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
	"strings"
	"testing"
	"time"
)

func TestInspectTransactionStateEmptyStateDirFails(t *testing.T) {
	diagnostics, err := InspectTransactionState(context.Background(), "")
	if err == nil {
		t.Fatal("InspectTransactionState() error = nil")
	}
	if !diagnostics.PublishBlocked {
		t.Fatal("PublishBlocked = false")
	}
	assertDiagnosticBlocker(t, diagnostics, TransactionBlockerStateDirUnavailable, "", "")
	if diagnostics.Lock.Reason != LockShowReasonStateDirMissing {
		t.Fatalf("lock reason = %q", diagnostics.Lock.Reason)
	}
}

func TestInspectTransactionStateNoLockNoJournals(t *testing.T) {
	diagnostics, err := InspectTransactionState(context.Background(), t.TempDir())
	if err != nil {
		t.Fatalf("InspectTransactionState() error = %v", err)
	}
	if diagnostics.PublishBlocked || len(diagnostics.Blockers) != 0 || len(diagnostics.Journals) != 0 {
		t.Fatalf("diagnostics = %#v", diagnostics)
	}
	if diagnostics.Lock.Status != LockShowStatusAbsent || diagnostics.Lock.Reason != LockShowReasonLockAbsent {
		t.Fatalf("lock = %#v", diagnostics.Lock)
	}
	if diagnostics.OperationLock.Present {
		t.Fatalf("operation lock = %#v, want absent", diagnostics.OperationLock)
	}
}

func TestInspectTransactionStateReportsOperationLock(t *testing.T) {
	stateDir := t.TempDir()
	writeOperationLockFile(t, stateDir, operationLockPublish, "other-token")

	diagnostics, err := InspectTransactionState(context.Background(), stateDir)
	if err != nil {
		t.Fatalf("InspectTransactionState() error = %v", err)
	}
	if !diagnostics.PublishBlocked {
		t.Fatalf("PublishBlocked = false diagnostics=%#v", diagnostics)
	}
	if !diagnostics.OperationLock.Present || diagnostics.OperationLock.Operation != "publish" || diagnostics.OperationLock.PID == "" || diagnostics.OperationLock.StartedAt == "" {
		t.Fatalf("diagnostics = %#v", diagnostics)
	}
	assertDiagnosticBlockerReason(t, diagnostics, TransactionBlockerOperationLock, "", "", TransactionBlockerReasonOperationLockExists)
	if _, err := os.Stat(operationLockPath(stateDir)); err != nil {
		t.Fatalf("operation lock missing: %v", err)
	}
}

func TestInspectTransactionStateOperationLockFailures(t *testing.T) {
	t.Run("corrupt operation lock", func(t *testing.T) {
		stateDir := t.TempDir()
		if err := os.WriteFile(operationLockPath(stateDir), []byte("operation=publish\n"), 0o600); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		diagnostics, err := InspectTransactionState(context.Background(), stateDir)
		if err != nil {
			t.Fatalf("InspectTransactionState() error = %v", err)
		}
		if !diagnostics.OperationLock.Present || !diagnostics.OperationLock.Corrupt || diagnostics.OperationLock.Message == "" {
			t.Fatalf("operation lock = %#v", diagnostics.OperationLock)
		}
		assertDiagnosticBlockerReason(t, diagnostics, TransactionBlockerCorruptOperationLock, "", "", TransactionBlockerReasonOperationLockCorrupt)
	})

	t.Run("operation lock read failed", func(t *testing.T) {
		stateDir := t.TempDir()
		if err := os.MkdirAll(operationLockPath(stateDir), 0o700); err != nil {
			t.Fatalf("MkdirAll() error = %v", err)
		}

		diagnostics, err := InspectTransactionState(context.Background(), stateDir)
		if err != nil {
			t.Fatalf("InspectTransactionState() error = %v", err)
		}
		if !diagnostics.OperationLock.Present || !diagnostics.OperationLock.ReadFailed || diagnostics.OperationLock.Message == "" {
			t.Fatalf("operation lock = %#v", diagnostics.OperationLock)
		}
		if strings.Contains(diagnostics.OperationLock.Message, stateDir) {
			t.Fatalf("operation lock message leaked path: %q", diagnostics.OperationLock.Message)
		}
		assertDiagnosticBlockerReason(t, diagnostics, TransactionBlockerOperationLockReadFailed, "", "", TransactionBlockerReasonOperationLockReadFailed)
	})
}

func TestInspectTransactionStateJournalPolicyMatrix(t *testing.T) {
	tests := []struct {
		id          TransactionID
		status      TransactionStatus
		wantBlocked bool
		wantKind    TransactionStateBlockerKind
	}{
		{id: "tx-pending", status: TransactionStatusPending, wantBlocked: true, wantKind: TransactionBlockerActiveJournal},
		{id: "tx-preflighted", status: TransactionStatusPreflighted, wantBlocked: true, wantKind: TransactionBlockerActiveJournal},
		{id: "tx-snapshotted", status: TransactionStatusSnapshotted, wantBlocked: true, wantKind: TransactionBlockerActiveJournal},
		{id: "tx-committed-locally", status: TransactionStatusCommittedLocally, wantBlocked: true, wantKind: TransactionBlockerActiveJournal},
		{id: "tx-candidates-pushed", status: TransactionStatusCandidatesPushed, wantBlocked: true, wantKind: TransactionBlockerActiveJournal},
		{id: "tx-promoting", status: TransactionStatusPromoting, wantBlocked: true, wantKind: TransactionBlockerActiveJournal},
		{id: "tx-branches-promoted", status: TransactionStatusBranchesPromoted, wantBlocked: true, wantKind: TransactionBlockerActiveJournal},
		{id: "tx-tagging", status: TransactionStatusTagging, wantBlocked: true, wantKind: TransactionBlockerActiveJournal},
		{id: "tx-committed", status: TransactionStatusCommitted},
		{id: "tx-failed", status: TransactionStatusFailed, wantBlocked: true, wantKind: TransactionBlockerFailedJournal},
		{id: "tx-rolling-back", status: TransactionStatusRollingBack, wantBlocked: true, wantKind: TransactionBlockerActiveJournal},
		{id: "tx-rolled-back", status: TransactionStatusRolledBack},
		{id: "tx-rollback-failed", status: TransactionStatusRollbackFailed, wantBlocked: true, wantKind: TransactionBlockerRollbackFailed},
		{id: "tx-unknown", status: TransactionStatus("unknown"), wantBlocked: true, wantKind: TransactionBlockerActiveJournal},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.id.String(), func(t *testing.T) {
			stateDir := t.TempDir()
			store := NewFileJournalStore(stateDir)
			writePruneJournal(t, store, tt.id, tt.status, time.Unix(1, 0).UTC())

			diagnostics, err := InspectTransactionState(context.Background(), stateDir)
			if err != nil {
				t.Fatalf("InspectTransactionState() error = %v", err)
			}
			if diagnostics.PublishBlocked != tt.wantBlocked {
				t.Fatalf("PublishBlocked = %v, want %v diagnostics=%#v", diagnostics.PublishBlocked, tt.wantBlocked, diagnostics)
			}
			if len(diagnostics.Journals) != 1 {
				t.Fatalf("journals = %#v", diagnostics.Journals)
			}
			journal := diagnostics.Journals[0]
			if journal.Name != tt.id.String()+".json" {
				t.Fatalf("journal name = %q", journal.Name)
			}
			if journal.Prunable != tt.status.Prunable() || journal.BlocksNewPublish != tt.status.BlocksNewPublish() || journal.AllowsLockClear != tt.status.AllowsLockClear() {
				t.Fatalf("journal policies = %#v", journal)
			}
			if tt.wantBlocked {
				assertDiagnosticBlocker(t, diagnostics, tt.wantKind, tt.id, tt.status)
			}
		})
	}
}

func TestInspectTransactionStateLockBlockers(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T, stateDir string)
		wantKind   TransactionStateBlockerKind
		wantID     TransactionID
		wantStatus TransactionStatus
		wantReason TransactionStateBlockerReason
	}{
		{
			name: "active lock",
			setup: func(t *testing.T, stateDir string) {
				store := NewFileJournalStore(stateDir)
				writePruneJournal(t, store, "tx-active", TransactionStatusPending, time.Unix(1, 0).UTC())
				writeLockFile(t, stateDir, "tx-active")
			},
			wantKind:   TransactionBlockerPublishLock,
			wantID:     "tx-active",
			wantStatus: TransactionStatusPending,
			wantReason: TransactionBlockerReasonActiveJournal,
		},
		{
			name: "terminal stale lock",
			setup: func(t *testing.T, stateDir string) {
				store := NewFileJournalStore(stateDir)
				writePruneJournal(t, store, "tx-committed", TransactionStatusCommitted, time.Unix(1, 0).UTC())
				writeLockFile(t, stateDir, "tx-committed")
			},
			wantKind:   TransactionBlockerPublishLock,
			wantID:     "tx-committed",
			wantStatus: TransactionStatusCommitted,
			wantReason: TransactionBlockerReasonStaleTerminalLock,
		},
		{
			name: "failed stale lock",
			setup: func(t *testing.T, stateDir string) {
				store := NewFileJournalStore(stateDir)
				writePruneJournal(t, store, "tx-failed", TransactionStatusFailed, time.Unix(1, 0).UTC())
				writeLockFile(t, stateDir, "tx-failed")
			},
			wantKind:   TransactionBlockerPublishLock,
			wantID:     "tx-failed",
			wantStatus: TransactionStatusFailed,
			wantReason: TransactionBlockerReasonRecoveryLock,
		},
		{
			name: "missing lock journal",
			setup: func(t *testing.T, stateDir string) {
				writeLockFile(t, stateDir, "tx-missing")
			},
			wantKind:   TransactionBlockerMissingLockJournal,
			wantID:     "tx-missing",
			wantReason: TransactionBlockerReasonMissingLockJournal,
		},
		{
			name: "corrupt lock",
			setup: func(t *testing.T, stateDir string) {
				writeRawLockFile(t, stateDir, "pid=1\n")
			},
			wantKind:   TransactionBlockerCorruptLock,
			wantReason: TransactionBlockerReasonCorruptLock,
		},
		{
			name: "corrupt lock journal",
			setup: func(t *testing.T, stateDir string) {
				writeLockFile(t, stateDir, "tx-corrupt")
				txDir := filepath.Join(stateDir, "transactions")
				if err := os.MkdirAll(txDir, 0o700); err != nil {
					t.Fatalf("MkdirAll() error = %v", err)
				}
				if err := os.WriteFile(filepath.Join(txDir, "tx-corrupt.json"), []byte("{"), 0o600); err != nil {
					t.Fatalf("WriteFile() error = %v", err)
				}
			},
			wantKind:   TransactionBlockerCorruptJournal,
			wantID:     "tx-corrupt",
			wantReason: TransactionBlockerReasonCorruptJournal,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			stateDir := t.TempDir()
			tt.setup(t, stateDir)

			diagnostics, err := InspectTransactionState(context.Background(), stateDir)
			if err != nil && diagnostics.Lock.Status == LockShowStatusFailed {
				t.Fatalf("InspectTransactionState() error = %v diagnostics=%#v", err, diagnostics)
			}
			if !diagnostics.PublishBlocked {
				t.Fatalf("PublishBlocked = false diagnostics=%#v", diagnostics)
			}
			assertDiagnosticBlockerReason(t, diagnostics, tt.wantKind, tt.wantID, tt.wantStatus, tt.wantReason)
		})
	}
}

func TestInspectTransactionStateMultipleBlockers(t *testing.T) {
	stateDir := t.TempDir()
	store := NewFileJournalStore(stateDir)
	writePruneJournal(t, store, "tx-failed", TransactionStatusFailed, time.Unix(1, 0).UTC())
	writePruneJournal(t, store, "tx-active", TransactionStatusPending, time.Unix(2, 0).UTC())
	writeLockFile(t, stateDir, "tx-missing")

	diagnostics, err := InspectTransactionState(context.Background(), stateDir)
	if err != nil {
		t.Fatalf("InspectTransactionState() error = %v", err)
	}
	assertDiagnosticBlocker(t, diagnostics, TransactionBlockerMissingLockJournal, "tx-missing", "")
	assertDiagnosticBlocker(t, diagnostics, TransactionBlockerFailedJournal, "tx-failed", TransactionStatusFailed)
	assertDiagnosticBlocker(t, diagnostics, TransactionBlockerActiveJournal, "tx-active", TransactionStatusPending)
	if len(diagnostics.Blockers) != 3 {
		t.Fatalf("blockers = %#v", diagnostics.Blockers)
	}
}

func TestInspectTransactionStateReadFailures(t *testing.T) {
	t.Run("lock read failed", func(t *testing.T) {
		stateDir := t.TempDir()
		if err := os.MkdirAll(filepath.Join(stateDir, "publish.lock"), 0o700); err != nil {
			t.Fatalf("MkdirAll() error = %v", err)
		}

		diagnostics, err := InspectTransactionState(context.Background(), stateDir)
		if err == nil {
			t.Fatal("InspectTransactionState() error = nil")
		}
		assertDiagnosticBlocker(t, diagnostics, TransactionBlockerLockReadFailed, "", "")
	})

	t.Run("journal directory read failed", func(t *testing.T) {
		stateDir := t.TempDir()
		if err := os.WriteFile(filepath.Join(stateDir, "transactions"), []byte("not a dir"), 0o600); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		diagnostics, err := InspectTransactionState(context.Background(), stateDir)
		if err == nil {
			t.Fatal("InspectTransactionState() error = nil")
		}
		assertDiagnosticBlocker(t, diagnostics, TransactionBlockerJournalDirectoryReadFailed, "", "")
		if len(diagnostics.Journals) != 0 {
			t.Fatalf("journals = %#v", diagnostics.Journals)
		}
	})

	t.Run("journal file read failed", func(t *testing.T) {
		stateDir := t.TempDir()
		txDir := filepath.Join(stateDir, "transactions")
		if err := os.MkdirAll(filepath.Join(txDir, "tx-bad.json"), 0o700); err != nil {
			t.Fatalf("MkdirAll() error = %v", err)
		}

		diagnostics, err := InspectTransactionState(context.Background(), stateDir)
		if err != nil {
			t.Fatalf("InspectTransactionState() error = %v", err)
		}
		assertDiagnosticBlocker(t, diagnostics, TransactionBlockerJournalFileReadFailed, "tx-bad", "")
		if len(diagnostics.Journals) != 1 || !diagnostics.Journals[0].ReadFailed || diagnostics.Journals[0].Name != "tx-bad.json" {
			t.Fatalf("journals = %#v", diagnostics.Journals)
		}
	})
}

func TestInspectTransactionStateCorruptJournalNames(t *testing.T) {
	tests := []struct {
		name        string
		filename    string
		content     string
		wantID      TransactionID
		wantName    string
		wantCorrupt bool
	}{
		{
			name:        "corrupt json derives id from filename",
			filename:    "tx-bad.json",
			content:     "{",
			wantID:      "tx-bad",
			wantName:    "tx-bad.json",
			wantCorrupt: true,
		},
		{
			name:        "unsafe filename keeps safe name",
			filename:    "bad.lock.json",
			content:     `{"schemaVersion":1,"id":"tx-good","status":"committed","startedAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z","modules":[]}`,
			wantName:    "bad.lock.json",
			wantCorrupt: true,
		},
		{
			name:        "identity mismatch keeps filename id",
			filename:    "tx-one.json",
			content:     `{"schemaVersion":1,"id":"tx-two","status":"committed","startedAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z","modules":[]}`,
			wantID:      "tx-one",
			wantName:    "tx-one.json",
			wantCorrupt: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			stateDir := t.TempDir()
			txDir := filepath.Join(stateDir, "transactions")
			if err := os.MkdirAll(txDir, 0o700); err != nil {
				t.Fatalf("MkdirAll() error = %v", err)
			}
			if err := os.WriteFile(filepath.Join(txDir, tt.filename), []byte(tt.content), 0o600); err != nil {
				t.Fatalf("WriteFile() error = %v", err)
			}

			diagnostics, err := InspectTransactionState(context.Background(), stateDir)
			if err != nil {
				t.Fatalf("InspectTransactionState() error = %v", err)
			}
			if len(diagnostics.Journals) != 1 {
				t.Fatalf("journals = %#v", diagnostics.Journals)
			}
			journal := diagnostics.Journals[0]
			if journal.ID != tt.wantID || journal.Name != tt.wantName || journal.Corrupt != tt.wantCorrupt {
				t.Fatalf("journal = %#v", journal)
			}
			assertDiagnosticBlockerReason(t, diagnostics, TransactionBlockerCorruptJournal, tt.wantID, "", TransactionBlockerReasonCorruptJournal)
		})
	}
}

func TestTransactionStateDiagnosticsFinalize(t *testing.T) {
	t.Run("warning without blocker does not block", func(t *testing.T) {
		var diagnostics TransactionStateDiagnostics
		diagnostics.addWarning(TransactionDiagnosticWarning{Code: "journal_missing", Message: "journal missing"})
		diagnostics.finalize()
		if diagnostics.PublishBlocked {
			t.Fatalf("PublishBlocked = true diagnostics=%#v", diagnostics)
		}
	})

	t.Run("duplicate blocker is not duplicated but blocks", func(t *testing.T) {
		var diagnostics TransactionStateDiagnostics
		blocker := TransactionStateBlocker{
			Kind:          TransactionBlockerFailedJournal,
			TransactionID: "tx-one",
			Status:        TransactionStatusFailed,
			Reason:        TransactionBlockerReasonFailedJournal,
		}
		diagnostics.addBlocker(blocker)
		diagnostics.addBlocker(blocker)
		diagnostics.finalize()
		if !diagnostics.PublishBlocked || len(diagnostics.Blockers) != 1 {
			t.Fatalf("diagnostics = %#v", diagnostics)
		}
	})
}

func assertDiagnosticBlocker(t *testing.T, diagnostics TransactionStateDiagnostics, kind TransactionStateBlockerKind, id TransactionID, status TransactionStatus) {
	t.Helper()
	assertDiagnosticBlockerReason(t, diagnostics, kind, id, status, "")
}

func assertDiagnosticBlockerReason(t *testing.T, diagnostics TransactionStateDiagnostics, kind TransactionStateBlockerKind, id TransactionID, status TransactionStatus, reason TransactionStateBlockerReason) {
	t.Helper()
	for _, blocker := range diagnostics.Blockers {
		if blocker.Kind == kind &&
			blocker.TransactionID == id &&
			blocker.Status == status &&
			(reason == "" || blocker.Reason == reason) {
			return
		}
	}
	t.Fatalf("blocker kind=%q id=%q status=%q reason=%q not found: %#v", kind, id, status, reason, diagnostics.Blockers)
}
