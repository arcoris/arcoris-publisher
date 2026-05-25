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
	"strings"
	"time"
)

type transactionLock struct {
	path string
	id   TransactionID
}

// TransactionLockInfo describes an existing publish lock without exposing the
// lock file path. Preflight and recovery commands use it for diagnostics.
type TransactionLockInfo struct {
	ID        TransactionID
	PID       string
	StartedAt string
	Command   string
	Path      string
}

func acquireTransactionLock(ctx context.Context, stateDir string, id TransactionID, now time.Time) (transactionLock, error) {
	if err := ctx.Err(); err != nil {
		return transactionLock{}, err
	}
	if stateDir == "" {
		return transactionLock{}, fmt.Errorf("state dir is required")
	}
	if err := validateTransactionID(id); err != nil {
		return transactionLock{}, err
	}
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
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
	content := fmt.Sprintf("transaction=%s\npid=%d\nstartedAt=%s\ncommand=publish\n", id, os.Getpid(), now.UTC().Format(time.RFC3339Nano))
	if _, err := file.WriteString(content); err != nil {
		_ = file.Close()
		_ = os.Remove(path)
		return transactionLock{}, err
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		_ = os.Remove(path)
		return transactionLock{}, err
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(path)
		return transactionLock{}, err
	}
	if err := syncParentDir(path); err != nil {
		_ = os.Remove(path)
		return transactionLock{}, err
	}
	return transactionLock{path: path, id: id}, nil
}

func (l transactionLock) Release() error {
	if l.path == "" {
		return nil
	}
	info, err := readTransactionLock(l.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	if info.ID != l.id {
		return fmt.Errorf("publish lock belongs to transaction %s, not %s", info.ID, l.id)
	}
	return os.Remove(l.path)
}

func currentTransactionLock(stateDir string) (TransactionLockInfo, bool, error) {
	path, err := transactionLockPath(stateDir)
	if err != nil {
		return TransactionLockInfo{}, false, err
	}
	info, err := readTransactionLock(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return TransactionLockInfo{}, false, nil
		}
		return TransactionLockInfo{}, false, err
	}
	return info, true, nil
}

// CurrentTransactionLock returns the current lock, if one exists.
func CurrentTransactionLock(stateDir string) (TransactionLockInfo, bool, error) {
	return currentTransactionLock(stateDir)
}

func transactionLockPath(stateDir string) (string, error) {
	if strings.TrimSpace(stateDir) == "" {
		return "", fmt.Errorf("state dir is required")
	}
	return filepath.Join(stateDir, "publish.lock"), nil
}

func readTransactionLock(path string) (TransactionLockInfo, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return TransactionLockInfo{}, err
	}
	info := TransactionLockInfo{}
	for _, line := range strings.Split(string(data), "\n") {
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		switch key {
		case "transaction":
			info.ID = TransactionID(value)
		case "pid":
			info.PID = value
		case "startedAt":
			info.StartedAt = value
		case "command":
			info.Command = value
		}
	}
	if info.ID == "" {
		return TransactionLockInfo{}, fmt.Errorf("publish lock is missing transaction id")
	}
	if err := validateTransactionID(info.ID); err != nil {
		return TransactionLockInfo{}, fmt.Errorf("publish lock has unsafe transaction id %q", info.ID)
	}
	info.Path = path
	return info, nil
}
