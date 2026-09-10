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
	"testing"
)

func TestFileJournalStoreUpdateRequiresExistingJournal(t *testing.T) {
	store := NewFileJournalStore(t.TempDir())
	journal := TransactionJournal{ID: "tx-test", Status: TransactionStatusFailed}

	err := store.Update(context.Background(), journal)
	if !errors.Is(err, errTransactionJournalNotFound) {
		t.Fatalf("Update() error = %v, want errTransactionJournalNotFound", err)
	}
	path, pathErr := store.journalPath(journal.ID)
	if pathErr != nil {
		t.Fatalf("journalPath() error = %v", pathErr)
	}
	if _, statErr := os.Stat(path); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("missing journal was created: stat error = %v", statErr)
	}
}

func TestFileJournalStoreUpdatePreservesCreateUpdateLifecycle(t *testing.T) {
	store := NewFileJournalStore(t.TempDir())
	journal := TransactionJournal{ID: "tx-test", Status: TransactionStatusPending, Version: "v1"}
	if err := store.Create(context.Background(), journal); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	journal.Status = TransactionStatusFailed
	journal.Version = "v2"
	if err := store.Update(context.Background(), journal); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	loaded, err := store.Load(context.Background(), journal.ID)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if loaded.Status != TransactionStatusFailed || loaded.Version != "v2" {
		t.Fatalf("loaded journal = %#v", loaded)
	}
}
