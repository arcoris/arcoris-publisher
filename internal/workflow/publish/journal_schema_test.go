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
)

func TestFileJournalStoreCreateCanonicalizesSchema(t *testing.T) {
	store := NewFileJournalStore(t.TempDir())
	journal := TransactionJournal{
		ID:        "tx-test",
		Status:    TransactionStatusPending,
		StartedAt: time.Unix(1, 0).UTC(),
		UpdatedAt: time.Unix(1, 0).UTC(),
	}
	if err := store.Create(context.Background(), journal); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	loaded, err := store.Load(context.Background(), journal.ID)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.SchemaVersion != transactionSchemaVersion {
		t.Fatalf("schemaVersion = %d, want %d", loaded.SchemaVersion, transactionSchemaVersion)
	}
}

func TestFileJournalStoreCreateRefusesExistingIdentity(t *testing.T) {
	store := NewFileJournalStore(t.TempDir())
	original := TransactionJournal{ID: "tx-test", Status: TransactionStatusPending, Version: "v1"}
	if err := store.Create(context.Background(), original); err != nil {
		t.Fatalf("first Create() error = %v", err)
	}
	replacement := original
	replacement.Status = TransactionStatusCommitted
	replacement.Version = "v2"
	if err := store.Create(context.Background(), replacement); !errors.Is(err, errTransactionJournalExists) {
		t.Fatalf("second Create() error = %v, want errTransactionJournalExists", err)
	}
	loaded, err := store.Load(context.Background(), original.ID)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.Status != original.Status || loaded.Version != original.Version {
		t.Fatalf("existing journal overwritten: %#v", loaded)
	}
}

func TestFileJournalStoreRejectsUnsupportedWriteSchema(t *testing.T) {
	store := NewFileJournalStore(t.TempDir())
	err := store.Create(context.Background(), TransactionJournal{
		SchemaVersion: transactionSchemaVersion + 1,
		ID:            "tx-test",
		Status:        TransactionStatusPending,
	})
	if err == nil {
		t.Fatal("Create() error = nil")
	}
}

func TestFileJournalStoreRejectsUnsupportedWireRepresentations(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{name: "missing schema", json: `{"id":"tx-test","status":"pending"}`},
		{name: "future schema", json: `{"schemaVersion":2,"id":"tx-test","status":"pending"}`},
		{name: "unknown field", json: `{"schemaVersion":1,"id":"tx-test","status":"pending","futureField":true}`},
		{name: "multiple values", json: `{"schemaVersion":1,"id":"tx-test","status":"pending"} {}`},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			stateDir := t.TempDir()
			store := NewFileJournalStore(stateDir)
			if err := os.MkdirAll(store.transactionsDir(), 0o700); err != nil {
				t.Fatalf("MkdirAll() error = %v", err)
			}
			path := filepath.Join(store.transactionsDir(), "tx-test.json")
			if err := os.WriteFile(path, []byte(tt.json), 0o600); err != nil {
				t.Fatalf("WriteFile() error = %v", err)
			}

			if _, err := store.Load(context.Background(), "tx-test"); !errors.Is(err, errTransactionJournalCorrupt) {
				t.Fatalf("Load() error = %v, want corrupt journal", err)
			}
			if _, err := store.List(context.Background()); !errors.Is(err, errTransactionJournalCorrupt) {
				t.Fatalf("List() error = %v, want corrupt journal", err)
			}

			diagnostics, err := InspectTransactionState(context.Background(), stateDir)
			if err != nil {
				t.Fatalf("InspectTransactionState() error = %v", err)
			}
			if !diagnostics.PublishBlocked {
				t.Fatal("PublishBlocked = false")
			}
			found := false
			for _, blocker := range diagnostics.Blockers {
				if blocker.Kind == TransactionBlockerCorruptJournal {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("blockers = %#v, want corrupt journal blocker", diagnostics.Blockers)
			}
		})
	}
}

func TestFileJournalStoreKeepsUnknownStatusOpaqueAndBlocking(t *testing.T) {
	store := NewFileJournalStore(t.TempDir())
	journal := TransactionJournal{
		ID:        "tx-unknown",
		Status:    TransactionStatus("future_state"),
		Rollback:  RollbackStatus("future_state"),
		StartedAt: time.Unix(1, 0).UTC(),
		UpdatedAt: time.Unix(1, 0).UTC(),
	}
	if err := store.Create(context.Background(), journal); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	loaded, err := store.Load(context.Background(), journal.ID)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.Status != journal.Status || loaded.Rollback != journal.Rollback {
		t.Fatalf("loaded opaque values = status %q rollback %q", loaded.Status, loaded.Rollback)
	}
	if !loaded.Status.BlocksNewPublish() || loaded.Status.Prunable() || loaded.Status.AllowsLockClear() {
		t.Fatalf("unknown status policy is not fail-closed: %q", loaded.Status)
	}
	if err := validateTransactionStatusTransition(loaded.Status, TransactionStatusFailed); !errors.Is(err, errInvalidTransactionStatusTransition) {
		t.Fatalf("unknown status transition error = %v, want invalid transition", err)
	}
}
