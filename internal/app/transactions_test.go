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

package app

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"arcoris.dev/arcoris-publisher/internal/workflow/publish"
)

func TestPruneTransactionsReturnsWorkflowResult(t *testing.T) {
	app, _ := appFixture(t)
	stateDir := t.TempDir()
	store := publish.NewFileJournalStore(stateDir)
	if err := store.Create(context.Background(), publish.TransactionJournal{
		SchemaVersion: 1,
		ID:            "tx-committed",
		Status:        publish.TransactionStatusCommitted,
		StartedAt:     time.Unix(1, 0).UTC(),
		UpdatedAt:     time.Unix(1, 0).UTC(),
	}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	result, err := app.PruneTransactions(context.Background(), TransactionPruneRequest{
		StateDir: stateDir,
		Statuses: []TransactionStatus{TransactionStatusCommitted},
		DryRun:   true,
	})
	if err != nil {
		t.Fatalf("PruneTransactions() error = %v", err)
	}
	if got := result.Result(); got.Status != publish.PruneStatusDryRun || len(got.Matched) != 1 {
		t.Fatalf("prune result = %#v", got)
	}
}

func TestTransactionLockUseCasesReturnWorkflowResults(t *testing.T) {
	app, _ := appFixture(t)
	stateDir := t.TempDir()
	store := publish.NewFileJournalStore(stateDir)
	if err := store.Create(context.Background(), publish.TransactionJournal{
		SchemaVersion: 1,
		ID:            "tx-committed",
		Status:        publish.TransactionStatusCommitted,
		StartedAt:     time.Unix(1, 0).UTC(),
		UpdatedAt:     time.Unix(1, 0).UTC(),
	}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(stateDir, "publish.lock"), []byte("transaction=tx-committed\npid=1\nstartedAt=2026-01-01T00:00:00Z\ncommand=publish\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	show, err := app.ShowTransactionLock(context.Background(), TransactionLockRequest{StateDir: stateDir})
	if err != nil {
		t.Fatalf("ShowTransactionLock() error = %v", err)
	}
	if got := show.Result(); got.Status != publish.LockShowStatusPresent || got.Journal.Status != publish.TransactionStatusCommitted {
		t.Fatalf("show result = %#v", got)
	}

	clear, err := app.ClearTransactionLock(context.Background(), TransactionLockRequest{
		StateDir:      stateDir,
		TransactionID: "tx-committed",
		Confirm:       "tx-committed",
	})
	if err != nil {
		t.Fatalf("ClearTransactionLock() error = %v", err)
	}
	if got := clear.Result(); got.Status != publish.LockClearStatusCleared || !got.LockCleared {
		t.Fatalf("clear result = %#v", got)
	}
}
