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
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestStateFileReadersRejectOversizedContent(t *testing.T) {
	t.Run("publish lock", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "publish.lock")
		if err := os.WriteFile(path, bytes.Repeat([]byte{'x'}, maxTransactionLockBytes+1), 0o600); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
		if _, err := readTransactionLock(path); !errors.Is(err, errTransactionLockCorrupt) {
			t.Fatalf("readTransactionLock() error = %v, want corrupt", err)
		}
	})

	t.Run("operation lock", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "operation.lock")
		if err := os.WriteFile(path, bytes.Repeat([]byte{'x'}, maxOperationLockBytes+1), 0o600); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
		if _, err := readOperationLock(path); !errors.Is(err, errOperationLockCorrupt) {
			t.Fatalf("readOperationLock() error = %v, want corrupt", err)
		}
	})

	t.Run("journal", func(t *testing.T) {
		stateDir := t.TempDir()
		store := NewFileJournalStore(stateDir)
		if err := os.MkdirAll(store.transactionsDir(), 0o700); err != nil {
			t.Fatalf("MkdirAll() error = %v", err)
		}
		path := filepath.Join(store.transactionsDir(), "tx-test.json")
		if err := os.WriteFile(path, bytes.Repeat([]byte{'x'}, maxTransactionJournalBytes+1), 0o600); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
		if _, err := store.Load(context.Background(), "tx-test"); !errors.Is(err, errTransactionJournalCorrupt) {
			t.Fatalf("Load() error = %v, want corrupt", err)
		}
	})
}

func TestReadOperationLockRequiresCanonicalV1Fields(t *testing.T) {
	tests := []string{
		"schemaVersion=1\noperation=publish\ntoken=token-one\nstartedAt=2026-01-01T00:00:00Z\n",
		"schemaVersion=1\noperation=publish\ntoken=token-one\npid=1\n",
	}
	for _, content := range tests {
		stateDir := t.TempDir()
		path := operationLockPath(stateDir)
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
		if _, err := readOperationLock(path); !errors.Is(err, errOperationLockCorrupt) {
			t.Fatalf("readOperationLock() error = %v, want corrupt", err)
		}
	}
}

func TestReadTransactionLockRequiresCanonicalSchemaV1ButKeepsLegacy(t *testing.T) {
	legacyDir := t.TempDir()
	legacyPath := filepath.Join(legacyDir, "publish.lock")
	if err := os.WriteFile(legacyPath, []byte("transaction=tx-one\n"), 0o600); err != nil {
		t.Fatalf("WriteFile(legacy) error = %v", err)
	}
	if _, err := readTransactionLock(legacyPath); err != nil {
		t.Fatalf("legacy readTransactionLock() error = %v", err)
	}

	invalid := []string{
		"schemaVersion=1\ntransaction=tx-one\nstartedAt=2026-01-01T00:00:00Z\ncommand=publish\n",
		"schemaVersion=1\ntransaction=tx-one\npid=1\ncommand=publish\n",
		"schemaVersion=1\ntransaction=tx-one\npid=1\nstartedAt=2026-01-01T00:00:00Z\n",
		"schemaVersion=1\ntransaction=tx-one\npid=1\nstartedAt=2026-01-01T00:00:00Z\ncommand=repair\n",
	}
	for _, content := range invalid {
		dir := t.TempDir()
		path := filepath.Join(dir, "publish.lock")
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
		if _, err := readTransactionLock(path); !errors.Is(err, errTransactionLockCorrupt) {
			t.Fatalf("readTransactionLock() error = %v, want corrupt", err)
		}
	}
}
