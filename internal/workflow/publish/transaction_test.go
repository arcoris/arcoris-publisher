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

func TestTransactionLockConflict(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	first, err := acquireTransactionLock(ctx, dir, "tx-one", time.Unix(1, 0).UTC())
	if err != nil {
		t.Fatalf("first lock error = %v", err)
	}
	if _, err := acquireTransactionLock(ctx, dir, "tx-two", time.Unix(2, 0).UTC()); err == nil {
		t.Fatal("second lock error = nil")
	}
	if err := first.Release(); err != nil {
		t.Fatalf("Release() error = %v", err)
	}
	if second, err := acquireTransactionLock(ctx, dir, "tx-two", time.Unix(2, 0).UTC()); err != nil {
		t.Fatalf("lock after release error = %v", err)
	} else {
		_ = second.Release()
	}
}

func TestTransactionLockReleaseRefusesMismatchedLock(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	first, err := acquireTransactionLock(ctx, dir, "tx-one", time.Unix(1, 0).UTC())
	if err != nil {
		t.Fatalf("first lock error = %v", err)
	}
	path := first.path
	if err := os.WriteFile(path, []byte("transaction=tx-two\npid=1\nstartedAt=now\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	if err := first.Release(); err == nil {
		t.Fatal("Release() error = nil")
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("mismatched lock was removed: %v", err)
	}
}
