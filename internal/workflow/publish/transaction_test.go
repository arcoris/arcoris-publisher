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
	"runtime"
	"strings"
	"testing"
	"time"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestCandidateRefSanitizesTransactionAndModule(t *testing.T) {
	got := candidateRef("tx:bad/value", manifest.ModuleName("control/api"))
	want := "refs/heads/arcpub/tx/tx-bad-value/control-api"
	if got != want {
		t.Fatalf("candidateRef() = %q, want %q", got, want)
	}
	if err := validateGitRef(got); err != nil {
		t.Fatalf("candidateRef() invalid Git ref: %v", err)
	}
}

func TestCandidateRefAvoidsUnsafeGitRefComponents(t *testing.T) {
	got := candidateRef("tx..bad.lock", manifest.ModuleName("@{/.lock"))
	if strings.Contains(got, "..") || strings.Contains(got, "@{") || strings.Contains(got, ".lock") {
		t.Fatalf("candidateRef() kept unsafe syntax: %q", got)
	}
	if err := validateGitRef(got); err != nil {
		t.Fatalf("candidateRef() invalid Git ref: %v", err)
	}
}

func TestValidateGitRefRejectsUnsafeRefs(t *testing.T) {
	for _, ref := range []string{
		"refs/heads/a..b",
		"refs/heads/main.lock",
		"refs/heads/@{bad",
		"refs/heads/trailing.",
		"refs/heads/has space",
		"refs//heads/main",
	} {
		if err := validateGitRef(ref); err == nil {
			t.Fatalf("validateGitRef(%q) error = nil", ref)
		}
	}
}

func TestFileJournalStoreWriteLoadListPending(t *testing.T) {
	ctx := context.Background()
	store := NewFileJournalStore(t.TempDir())
	journal := TransactionJournal{
		SchemaVersion: transactionSchemaVersion,
		ID:            "tx-test",
		Status:        TransactionStatusPending,
		Version:       "v0.1.0",
		StartedAt:     time.Unix(10, 0).UTC(),
		UpdatedAt:     time.Unix(10, 0).UTC(),
	}

	if err := store.Create(ctx, journal); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if runtime.GOOS != "windows" {
		path, err := store.journalPath("tx-test")
		if err != nil {
			t.Fatalf("journalPath() error = %v", err)
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("Stat() error = %v", err)
		}
		if got := info.Mode().Perm(); got != 0o600 {
			t.Fatalf("journal permissions = %o, want 0600", got)
		}
	}
	loaded, err := store.Load(ctx, "tx-test")
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.ID != journal.ID || loaded.Status != journal.Status {
		t.Fatalf("loaded journal = %+v", loaded)
	}
	pending, ok, err := store.HasPending(ctx)
	if err != nil {
		t.Fatalf("HasPending() error = %v", err)
	}
	if !ok || pending.ID != "tx-test" {
		t.Fatalf("pending = %+v ok=%v", pending, ok)
	}

	journal.Status = TransactionStatusRolledBack
	if err := store.Update(ctx, journal); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	_, ok, err = store.HasPending(ctx)
	if err != nil {
		t.Fatalf("HasPending() after terminal error = %v", err)
	}
	if ok {
		t.Fatal("terminal transaction reported as pending")
	}
}

func TestTransactionStatusPolicies(t *testing.T) {
	tests := []struct {
		status          TransactionStatus
		blocksPublish   bool
		prunable        bool
		allowsLockClear bool
	}{
		{status: TransactionStatusPending, blocksPublish: true},
		{status: TransactionStatusPreflighted, blocksPublish: true},
		{status: TransactionStatusSnapshotted, blocksPublish: true},
		{status: TransactionStatusCommittedLocally, blocksPublish: true},
		{status: TransactionStatusCandidatesPushed, blocksPublish: true},
		{status: TransactionStatusPromoting, blocksPublish: true},
		{status: TransactionStatusBranchesPromoted, blocksPublish: true},
		{status: TransactionStatusTagging, blocksPublish: true},
		{status: TransactionStatusCommitted, prunable: true, allowsLockClear: true},
		{status: TransactionStatusFailed, blocksPublish: true, allowsLockClear: true},
		{status: TransactionStatusRollingBack, blocksPublish: true},
		{status: TransactionStatusRolledBack, prunable: true, allowsLockClear: true},
		{status: TransactionStatusRollbackFailed, blocksPublish: true, allowsLockClear: true},
		{status: TransactionStatus("unknown"), blocksPublish: true},
		{status: TransactionStatus(""), blocksPublish: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.BlocksNewPublish(); got != tt.blocksPublish {
				t.Fatalf("BlocksNewPublish() = %v, want %v", got, tt.blocksPublish)
			}
			if got := tt.status.Prunable(); got != tt.prunable {
				t.Fatalf("Prunable() = %v, want %v", got, tt.prunable)
			}
			if got := tt.status.AllowsLockClear(); got != tt.allowsLockClear {
				t.Fatalf("AllowsLockClear() = %v, want %v", got, tt.allowsLockClear)
			}
			if got := tt.status.Terminal(); got != !tt.blocksPublish {
				t.Fatalf("Terminal() = %v, want %v", got, !tt.blocksPublish)
			}
		})
	}
}

func TestFileJournalStoreRejectsUnsafeTransactionID(t *testing.T) {
	store := NewFileJournalStore(t.TempDir())
	if err := store.Create(context.Background(), TransactionJournal{ID: "../escape"}); err == nil {
		t.Fatal("Create() error = nil")
	}
}

func TestFileJournalStoreCorruptJournalBlocksListAndPending(t *testing.T) {
	ctx := context.Background()
	store := NewFileJournalStore(t.TempDir())
	dir := store.transactionsDir()
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "tx-bad.json"), []byte("{"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := store.List(ctx); err == nil {
		t.Fatal("List() error = nil")
	}
	if _, _, err := store.HasPending(ctx); err == nil {
		t.Fatal("HasPending() error = nil")
	}
	if _, err := store.Load(ctx, "tx-bad"); err == nil {
		t.Fatal("Load() error = nil")
	}
}

func TestFileJournalStoreMismatchedJournalIDFailsClosed(t *testing.T) {
	ctx := context.Background()
	store := NewFileJournalStore(t.TempDir())
	dir := store.transactionsDir()
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	data := []byte(`{"schemaVersion":1,"id":"tx-other","status":"committed","startedAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}`)
	if err := os.WriteFile(filepath.Join(dir, "tx-mismatch.json"), data, 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if _, err := store.List(ctx); err == nil {
		t.Fatal("List() error = nil")
	}
	if _, err := store.Load(ctx, "tx-mismatch"); err == nil {
		t.Fatal("Load() error = nil")
	}
}

func TestTransactionLockConflict(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	first, err := acquireTransactionLock(ctx, dir, "tx-one", time.Unix(1, 0).UTC(), transactionLockOps{})
	if err != nil {
		t.Fatalf("first lock error = %v", err)
	}
	if _, err := acquireTransactionLock(ctx, dir, "tx-two", time.Unix(2, 0).UTC(), transactionLockOps{}); err == nil {
		t.Fatal("second lock error = nil")
	}
	if err := first.Release(); err != nil {
		t.Fatalf("Release() error = %v", err)
	}
	if second, err := acquireTransactionLock(ctx, dir, "tx-two", time.Unix(2, 0).UTC(), transactionLockOps{}); err != nil {
		t.Fatalf("lock after release error = %v", err)
	} else {
		_ = second.Release()
	}
}

func TestTransactionLockReleaseSyncsParentAfterRemove(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	var synced bool
	lock, err := acquireTransactionLock(ctx, dir, "tx-one", time.Unix(1, 0).UTC(), transactionLockOps{
		syncParent: func(path string) error {
			synced = filepath.Base(path) == "publish.lock"
			return nil
		},
	})
	if err != nil {
		t.Fatalf("lock error = %v", err)
	}

	if err := lock.Release(); err != nil {
		t.Fatalf("Release() error = %v", err)
	}
	if !synced {
		t.Fatal("Release() did not sync parent directory")
	}
	if _, err := os.Stat(lock.path); !os.IsNotExist(err) {
		t.Fatalf("lock exists or stat failed: %v", err)
	}
}

func TestTransactionLockReleaseReportsRemoveFailures(t *testing.T) {
	ctx := context.Background()

	t.Run("delete failed preserves lock", func(t *testing.T) {
		dir := t.TempDir()
		lock, err := acquireTransactionLock(ctx, dir, "tx-one", time.Unix(1, 0).UTC(), transactionLockOps{
			remove: func(string) error {
				return errors.New("delete refused")
			},
		})
		if err != nil {
			t.Fatalf("lock error = %v", err)
		}

		if err := lock.Release(); err == nil {
			t.Fatal("Release() error = nil")
		}
		if _, err := os.Stat(lock.path); err != nil {
			t.Fatalf("lock missing after failed delete: %v", err)
		}
	})

	t.Run("sync failed after delete removes lock", func(t *testing.T) {
		dir := t.TempDir()
		lock, err := acquireTransactionLock(ctx, dir, "tx-one", time.Unix(1, 0).UTC(), transactionLockOps{
			syncParent: func(string) error {
				return errors.New("sync refused")
			},
		})
		if err != nil {
			t.Fatalf("lock error = %v", err)
		}

		if err := lock.Release(); err == nil {
			t.Fatal("Release() error = nil")
		}
		if _, err := os.Stat(lock.path); !os.IsNotExist(err) {
			t.Fatalf("lock exists or stat failed after sync failure: %v", err)
		}
	})
}

func TestTransactionLockReleaseMissingIsNoop(t *testing.T) {
	lock := transactionLock{path: filepath.Join(t.TempDir(), "publish.lock"), id: "tx-missing"}
	if err := lock.Release(); err != nil {
		t.Fatalf("Release() error = %v", err)
	}
}

func TestTransactionLockReleaseRefusesMismatchedLock(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	first, err := acquireTransactionLock(ctx, dir, "tx-one", time.Unix(1, 0).UTC(), transactionLockOps{})
	if err != nil {
		t.Fatalf("first lock error = %v", err)
	}
	path := first.path
	if err := os.WriteFile(path, []byte("transaction=tx-two\npid=1\nstartedAt=2026-01-01T00:00:00Z\ncommand=publish\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := first.Release(); err == nil {
		t.Fatal("Release() error = nil")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("mismatched lock was removed: %v", err)
	}
}

func TestTransactionLockReleasePreservesCorruptLock(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	lock, err := acquireTransactionLock(ctx, dir, "tx-one", time.Unix(1, 0).UTC(), transactionLockOps{})
	if err != nil {
		t.Fatalf("lock error = %v", err)
	}
	if err := os.WriteFile(lock.path, []byte("pid=1\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := lock.Release(); err == nil {
		t.Fatal("Release() error = nil")
	}
	if _, err := os.Stat(lock.path); err != nil {
		t.Fatalf("corrupt lock was removed: %v", err)
	}
}

func TestReadTransactionLockStrictParser(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "valid current format",
			content: "transaction=tx-one\npid=123\nstartedAt=2026-01-01T00:00:00Z\ncommand=publish\n",
		},
		{
			name:    "valid CRLF format",
			content: "transaction=tx-one\r\npid=123\r\nstartedAt=2026-01-01T00:00:00Z\r\ncommand=publish\r\n",
		},
		{
			name:    "valid trailing empty lines",
			content: "transaction=tx-one\npid=123\nstartedAt=2026-01-01T00:00:00Z\ncommand=publish\n\n  \n",
		},
		{
			name:    "valid required transaction only",
			content: "transaction=tx-one\n",
		},
		{name: "missing transaction", content: "pid=1\nstartedAt=2026-01-01T00:00:00Z\ncommand=publish\n", wantErr: true},
		{name: "invalid transaction", content: "transaction=../bad\npid=1\nstartedAt=2026-01-01T00:00:00Z\ncommand=publish\n", wantErr: true},
		{name: "duplicate transaction", content: "transaction=tx-one\ntransaction=tx-two\npid=1\nstartedAt=2026-01-01T00:00:00Z\ncommand=publish\n", wantErr: true},
		{name: "malformed line", content: "transaction=tx-one\nnot-a-pair\npid=1\nstartedAt=2026-01-01T00:00:00Z\ncommand=publish\n", wantErr: true},
		{name: "duplicate pid", content: "transaction=tx-one\npid=1\npid=2\nstartedAt=2026-01-01T00:00:00Z\ncommand=publish\n", wantErr: true},
		{name: "duplicate startedAt", content: "transaction=tx-one\npid=1\nstartedAt=2026-01-01T00:00:00Z\nstartedAt=2026-01-01T00:00:01Z\ncommand=publish\n", wantErr: true},
		{name: "duplicate command", content: "transaction=tx-one\npid=1\nstartedAt=2026-01-01T00:00:00Z\ncommand=publish\ncommand=publish\n", wantErr: true},
		{name: "invalid pid", content: "transaction=tx-one\npid=-1\nstartedAt=2026-01-01T00:00:00Z\ncommand=publish\n", wantErr: true},
		{name: "invalid startedAt", content: "transaction=tx-one\npid=1\nstartedAt=now\ncommand=publish\n", wantErr: true},
		{name: "empty command", content: "transaction=tx-one\npid=1\nstartedAt=2026-01-01T00:00:00Z\ncommand=\n", wantErr: true},
		{name: "unknown key", content: "transaction=tx-one\npid=1\nstartedAt=2026-01-01T00:00:00Z\ncommand=publish\nowner=operator\n", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "publish.lock")
			if err := os.WriteFile(path, []byte(tt.content), 0o600); err != nil {
				t.Fatalf("WriteFile() error = %v", err)
			}

			info, err := readTransactionLock(path)
			if tt.wantErr {
				if err == nil {
					t.Fatal("readTransactionLock() error = nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("readTransactionLock() error = %v", err)
			}
			if info.ID != "tx-one" {
				t.Fatalf("transaction id = %q", info.ID)
			}
			if strings.HasPrefix(tt.name, "valid ") && tt.name != "valid required transaction only" && (info.PID != "123" || info.Command != "publish" || info.StartedAt != "2026-01-01T00:00:00Z") {
				t.Fatalf("lock info = %#v", info)
			}
		})
	}
}
