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

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestCandidateRefSanitizesTransactionAndModule(t *testing.T) {
	got := candidateRef("tx:bad/value", manifest.ModuleName("control/api"))
	want := "refs/heads/arcpub/tx/tx-bad-value/control-api"
	if got != want {
		t.Fatalf("candidateRef() = %q, want %q", got, want)
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
