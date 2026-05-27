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

func TestOperationLockAcquireAndRelease(t *testing.T) {
	stateDir := t.TempDir()
	ops := testOperationLockOps()

	lock, err := acquireOperationLock(context.Background(), stateDir, operationLockPublish, ops)
	if err != nil {
		t.Fatalf("acquireOperationLock() error = %v", err)
	}
	info, err := readOperationLock(operationLockPath(stateDir))
	if err != nil {
		t.Fatalf("readOperationLock() error = %v", err)
	}
	if info.Operation != operationLockPublish || info.Token != "token-one" {
		t.Fatalf("operation lock info = %#v", info)
	}

	outcome, err := lock.Release()
	if err != nil {
		t.Fatalf("Release() error = %v", err)
	}
	if !outcome.Released || !outcome.Synced {
		t.Fatalf("release outcome = %#v", outcome)
	}
	if _, err := os.Stat(operationLockPath(stateDir)); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("operation lock still exists: %v", err)
	}
}

func TestOperationLockRefusesExisting(t *testing.T) {
	stateDir := t.TempDir()
	writeOperationLockFile(t, stateDir, operationLockPrune, "stale-token")

	_, err := acquireOperationLock(context.Background(), stateDir, operationLockPublish, testOperationLockOps())
	if !errors.Is(err, errOperationLockExists) {
		t.Fatalf("acquireOperationLock() error = %v, want errOperationLockExists", err)
	}
	info, err := readOperationLock(operationLockPath(stateDir))
	if err != nil {
		t.Fatalf("readOperationLock() error = %v", err)
	}
	if info.Operation != operationLockPrune || info.Token != "stale-token" {
		t.Fatalf("existing operation lock changed: %#v", info)
	}
}

func TestOperationLockReleaseSafety(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(t *testing.T, stateDir string, lock operationLock)
		ops          func(stateDir string) operationLockOps
		wantErr      error
		wantReleased bool
		wantSynced   bool
		wantExists   bool
	}{
		{
			name: "missing lock is no-op",
			setup: func(t *testing.T, stateDir string, _ operationLock) {
				t.Helper()
				if err := os.Remove(operationLockPath(stateDir)); err != nil {
					t.Fatalf("Remove() error = %v", err)
				}
			},
			wantExists: false,
		},
		{
			name: "changed lock preserved",
			setup: func(t *testing.T, stateDir string, _ operationLock) {
				t.Helper()
				writeOperationLockFile(t, stateDir, operationLockRollback, "other-token")
			},
			wantErr:    errOperationLockChanged,
			wantExists: true,
		},
		{
			name: "corrupt lock preserved",
			setup: func(t *testing.T, stateDir string, _ operationLock) {
				t.Helper()
				if err := os.WriteFile(operationLockPath(stateDir), []byte("operation=publish\n"), 0o600); err != nil {
					t.Fatalf("WriteFile() error = %v", err)
				}
			},
			wantErr:    errOperationLockCorrupt,
			wantExists: true,
		},
		{
			name: "delete failure preserves lock",
			ops: func(string) operationLockOps {
				ops := testOperationLockOps()
				ops.remove = func(string) error { return errors.New("remove denied") }
				return ops
			},
			wantErr:    errOperationLockDeleteFailed,
			wantExists: true,
		},
		{
			name: "sync failure after removal",
			ops: func(string) operationLockOps {
				ops := testOperationLockOps()
				failSync := false
				ops.beforeRemove = func() { failSync = true }
				ops.syncParent = func(string) error {
					if failSync {
						return errors.New("sync denied")
					}
					return nil
				}
				return ops
			},
			wantErr:      errOperationLockSyncFailed,
			wantReleased: true,
			wantExists:   false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			stateDir := t.TempDir()
			ops := testOperationLockOps()
			if tt.ops != nil {
				ops = tt.ops(stateDir)
			}
			lock, err := acquireOperationLock(context.Background(), stateDir, operationLockPublish, ops)
			if err != nil {
				t.Fatalf("acquireOperationLock() error = %v", err)
			}
			if tt.setup != nil {
				tt.setup(t, stateDir, lock)
			}

			outcome, err := lock.Release()
			if tt.wantErr == nil {
				if err != nil {
					t.Fatalf("Release() error = %v", err)
				}
			} else if !errors.Is(err, tt.wantErr) {
				t.Fatalf("Release() error = %v, want %v", err, tt.wantErr)
			}
			if outcome.Released != tt.wantReleased || outcome.Synced != tt.wantSynced {
				t.Fatalf("release outcome = %#v", outcome)
			}
			_, statErr := os.Stat(operationLockPath(stateDir))
			exists := statErr == nil
			if exists != tt.wantExists {
				t.Fatalf("operation lock exists = %v, want %v statErr=%v", exists, tt.wantExists, statErr)
			}
		})
	}
}

func TestReadOperationLockRejectsMalformedContent(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{name: "missing schema", content: "operation=publish\ntoken=token-one\npid=1\nstartedAt=2026-01-01T00:00:00Z\n"},
		{name: "unsupported schema", content: "schemaVersion=2\noperation=publish\ntoken=token-one\npid=1\nstartedAt=2026-01-01T00:00:00Z\n"},
		{name: "duplicate operation", content: "schemaVersion=1\noperation=publish\noperation=prune\ntoken=token-one\npid=1\nstartedAt=2026-01-01T00:00:00Z\n"},
		{name: "unknown key", content: "schemaVersion=1\noperation=publish\ntoken=token-one\npid=1\nstartedAt=2026-01-01T00:00:00Z\nextra=true\n"},
		{name: "invalid operation", content: "schemaVersion=1\noperation=repair\ntoken=token-one\npid=1\nstartedAt=2026-01-01T00:00:00Z\n"},
		{name: "invalid pid", content: "schemaVersion=1\noperation=publish\ntoken=token-one\npid=abc\nstartedAt=2026-01-01T00:00:00Z\n"},
		{name: "invalid startedAt", content: "schemaVersion=1\noperation=publish\ntoken=token-one\npid=1\nstartedAt=nope\n"},
		{name: "invalid token", content: "schemaVersion=1\noperation=publish\ntoken=../bad\npid=1\nstartedAt=2026-01-01T00:00:00Z\n"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			stateDir := t.TempDir()
			if err := os.WriteFile(operationLockPath(stateDir), []byte(tt.content), 0o600); err != nil {
				t.Fatalf("WriteFile() error = %v", err)
			}
			if _, err := readOperationLock(operationLockPath(stateDir)); !errors.Is(err, errOperationLockCorrupt) {
				t.Fatalf("readOperationLock() error = %v, want errOperationLockCorrupt", err)
			}
		})
	}
}

func testOperationLockOps() operationLockOps {
	return operationLockOps{
		remove:     os.Remove,
		syncParent: func(string) error { return nil },
		now:        func() time.Time { return time.Unix(1, 0).UTC() },
		token:      func() (string, error) { return "token-one", nil },
	}
}

func writeOperationLockFile(t *testing.T, stateDir string, operation operationLockOperation, token string) {
	t.Helper()
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	content := "schemaVersion=1\noperation=" + string(operation) + "\ntoken=" + token + "\npid=1\nstartedAt=2026-01-01T00:00:00Z\n"
	if err := os.WriteFile(filepath.Join(stateDir, "operation.lock"), []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}
