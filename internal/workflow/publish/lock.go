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
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type transactionLock struct {
	path string
	id   TransactionID
}

func acquireTransactionLock(ctx context.Context, stateDir string, id TransactionID, now time.Time) (transactionLock, error) {
	if err := ctx.Err(); err != nil {
		return transactionLock{}, err
	}
	if stateDir == "" {
		return transactionLock{}, fmt.Errorf("state dir is required")
	}
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		return transactionLock{}, err
	}
	path := filepath.Join(stateDir, "publish.lock")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return transactionLock{}, fmt.Errorf("publish transaction lock already exists at %s", path)
		}
		return transactionLock{}, err
	}
	content := fmt.Sprintf("transaction=%s\npid=%d\nstartedAt=%s\n", id, os.Getpid(), now.UTC().Format(time.RFC3339Nano))
	if _, err := file.WriteString(content); err != nil {
		_ = file.Close()
		_ = os.Remove(path)
		return transactionLock{}, err
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(path)
		return transactionLock{}, err
	}
	return transactionLock{path: path, id: id}, nil
}

func (l transactionLock) Release() error {
	if l.path == "" {
		return nil
	}
	return os.Remove(l.path)
}

func releaseTransactionLock(stateDir string) error {
	if stateDir == "" {
		return nil
	}
	path := filepath.Join(stateDir, "publish.lock")
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
