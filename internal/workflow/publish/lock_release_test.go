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
	"time"
)

func TestTransactionLockReleaseReportsDisappearedOwnedLock(t *testing.T) {
	stateDir := t.TempDir()
	lock, err := acquireTransactionLock(
		context.Background(),
		stateDir,
		"tx-release-ownership",
		time.Unix(1, 0).UTC(),
		transactionLockOps{
			remove:            os.Remove,
			syncAcquireParent: func(string) error { return nil },
			syncParent:        func(string) error { return nil },
		},
	)
	if err != nil {
		t.Fatalf("acquireTransactionLock() error = %v", err)
	}
	if err := os.Remove(lock.path); err != nil {
		t.Fatalf("Remove() error = %v", err)
	}

	outcome, err := lock.Release()
	if !errors.Is(err, errTransactionLockDisappeared) {
		t.Fatalf("Release() error = %v, want errTransactionLockDisappeared", err)
	}
	if outcome.Removed || outcome.Synced {
		t.Fatalf("Release() outcome = %#v, want zero outcome", outcome)
	}
}

func TestZeroValueTransactionLockReleaseIsNoOp(t *testing.T) {
	outcome, err := (transactionLock{}).Release()
	if err != nil {
		t.Fatalf("Release() error = %v", err)
	}
	if outcome.Removed || outcome.Synced {
		t.Fatalf("Release() outcome = %#v, want zero outcome", outcome)
	}
}
