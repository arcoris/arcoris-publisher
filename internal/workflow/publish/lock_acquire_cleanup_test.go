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

func TestTransactionLockAcquireCleanupAfterParentSyncFailure(t *testing.T) {
	tests := []struct {
		name           string
		remove         func(string) error
		cleanupSyncErr error
		wantCleanupErr error
		wantExists     bool
	}{
		{name: "cleanup succeeds"},
		{
			name: "delete fails",
			remove: func(string) error {
				return errTestAcquireCleanupDelete
			},
			wantCleanupErr: errTransactionLockDeleteFailed,
			wantExists:     true,
		},
		{
			name:           "cleanup sync fails after removal",
			cleanupSyncErr: errTestAcquireCleanupSync,
			wantCleanupErr: errTransactionLockSyncFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateDir := t.TempDir()
			primary := errors.New("initial parent sync denied")
			syncCalls := 0
			ops := transactionLockOps{
				remove: os.Remove,
				syncParent: func(string) error {
					syncCalls++
					if syncCalls == 1 {
						return primary
					}
					return tt.cleanupSyncErr
				},
			}
			if tt.remove != nil {
				ops.remove = tt.remove
			}

			_, err := acquireTransactionLock(
				context.Background(),
				stateDir,
				"tx-test",
				time.Unix(1, 0).UTC(),
				ops,
			)
			if !errors.Is(err, primary) {
				t.Fatalf("acquireTransactionLock() error = %v, want primary sync failure", err)
			}
			if tt.wantCleanupErr != nil && !errors.Is(err, tt.wantCleanupErr) {
				t.Fatalf("acquireTransactionLock() error = %v, want cleanup error %v", err, tt.wantCleanupErr)
			}
			if tt.remove != nil && !errors.Is(err, errTestAcquireCleanupDelete) {
				t.Fatalf("acquireTransactionLock() error = %v, want underlying delete failure", err)
			}
			if tt.cleanupSyncErr != nil && !errors.Is(err, tt.cleanupSyncErr) {
				t.Fatalf("acquireTransactionLock() error = %v, want underlying cleanup sync failure", err)
			}
			if tt.remove == nil && syncCalls != 2 {
				t.Fatalf("syncParent calls = %d, want 2", syncCalls)
			}
			assertStatePathExists(t, filepath.Join(stateDir, "publish.lock"), tt.wantExists)
		})
	}
}

func TestOperationLockAcquireCleanupAfterParentSyncFailure(t *testing.T) {
	tests := []struct {
		name           string
		remove         func(string) error
		cleanupSyncErr error
		wantCleanupErr error
		wantExists     bool
	}{
		{name: "cleanup succeeds"},
		{
			name: "delete fails",
			remove: func(string) error {
				return errTestAcquireCleanupDelete
			},
			wantCleanupErr: errOperationLockDeleteFailed,
			wantExists:     true,
		},
		{
			name:           "cleanup sync fails after removal",
			cleanupSyncErr: errTestAcquireCleanupSync,
			wantCleanupErr: errOperationLockSyncFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stateDir := t.TempDir()
			primary := errors.New("initial parent sync denied")
			syncCalls := 0
			ops := testOperationLockOps()
			ops.syncParent = func(string) error {
				syncCalls++
				if syncCalls == 1 {
					return primary
				}
				return tt.cleanupSyncErr
			}
			if tt.remove != nil {
				ops.remove = tt.remove
			}

			_, err := acquireOperationLock(context.Background(), stateDir, operationLockPublish, ops)
			if !errors.Is(err, primary) {
				t.Fatalf("acquireOperationLock() error = %v, want primary sync failure", err)
			}
			if tt.wantCleanupErr != nil && !errors.Is(err, tt.wantCleanupErr) {
				t.Fatalf("acquireOperationLock() error = %v, want cleanup error %v", err, tt.wantCleanupErr)
			}
			if tt.remove != nil && !errors.Is(err, errTestAcquireCleanupDelete) {
				t.Fatalf("acquireOperationLock() error = %v, want underlying delete failure", err)
			}
			if tt.cleanupSyncErr != nil && !errors.Is(err, tt.cleanupSyncErr) {
				t.Fatalf("acquireOperationLock() error = %v, want underlying cleanup sync failure", err)
			}
			if tt.remove == nil && syncCalls != 2 {
				t.Fatalf("syncParent calls = %d, want 2", syncCalls)
			}
			assertStatePathExists(t, operationLockPath(stateDir), tt.wantExists)
		})
	}
}

func TestLockReadersRejectSymlinkState(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	if err := os.WriteFile(target, []byte("schemaVersion=1\ntransaction=tx-test\npid=1\nstartedAt=2026-01-01T00:00:00Z\ncommand=publish\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	publishLock := filepath.Join(dir, "publish.lock")
	if err := os.Symlink(target, publishLock); err != nil {
		t.Skipf("Symlink() unavailable: %v", err)
	}
	if _, err := readTransactionLock(publishLock); !errors.Is(err, errTransactionLockCorrupt) {
		t.Fatalf("readTransactionLock() error = %v, want corrupt lock", err)
	}
}

var (
	errTestAcquireCleanupDelete = errors.New("test acquire cleanup delete failure")
	errTestAcquireCleanupSync   = errors.New("test acquire cleanup sync failure")
)

func assertStatePathExists(t *testing.T, path string, want bool) {
	t.Helper()
	_, err := os.Lstat(path)
	exists := err == nil
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("Lstat(%s) error = %v", filepath.Base(path), err)
	}
	if exists != want {
		t.Fatalf("state path %s exists = %v, want %v", filepath.Base(path), exists, want)
	}
}
