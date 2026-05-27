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
	ops  transactionLockOps
}

var (
	errTransactionLockCorrupt      = errors.New("publish lock corrupt")
	errTransactionLockChanged      = errors.New("publish lock changed")
	errTransactionLockDisappeared  = errors.New("publish lock disappeared")
	errTransactionLockDeleteFailed = errors.New("publish lock delete failed")
	errTransactionLockSyncFailed   = errors.New("publish lock sync failed")
)

type transactionLockOps struct {
	remove       func(string) error
	syncParent   func(string) error
	beforeRemove func()
}

type lockRemoveOutcome struct {
	Removed bool
}

func defaultTransactionLockOps() transactionLockOps {
	return transactionLockOps{
		remove:     os.Remove,
		syncParent: syncParentDir,
	}
}

func (ops transactionLockOps) withDefaults() transactionLockOps {
	defaults := defaultTransactionLockOps()
	if ops.remove == nil {
		ops.remove = defaults.remove
	}
	if ops.syncParent == nil {
		ops.syncParent = defaults.syncParent
	}
	return ops
}

func lockCorruptf(format string, args ...any) error {
	return fmt.Errorf("%w: "+format, append([]any{errTransactionLockCorrupt}, args...)...)
}

// removeTransactionLockIfCurrent performs the final identity check immediately
// before deletion. A sync failure is reported after Removed=true because the
// unlink already happened even though durability could not be confirmed.
func removeTransactionLockIfCurrent(path string, expected TransactionID, ops transactionLockOps) (lockRemoveOutcome, error) {
	ops = ops.withDefaults()
	if ops.beforeRemove != nil {
		ops.beforeRemove()
	}
	info, err := readTransactionLock(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return lockRemoveOutcome{}, fmt.Errorf("%w: %v", errTransactionLockDisappeared, err)
		}
		return lockRemoveOutcome{}, err
	}
	if info.ID != expected {
		return lockRemoveOutcome{}, fmt.Errorf("%w: publish lock changed from %s to %s", errTransactionLockChanged, expected, info.ID)
	}
	if err := ops.remove(path); err != nil {
		return lockRemoveOutcome{}, fmt.Errorf("%w: %v", errTransactionLockDeleteFailed, err)
	}
	outcome := lockRemoveOutcome{Removed: true}
	if err := ops.syncParent(path); err != nil {
		return outcome, fmt.Errorf("%w: %v", errTransactionLockSyncFailed, err)
	}
	return outcome, nil
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

func acquireTransactionLock(ctx context.Context, stateDir string, id TransactionID, now time.Time, ops transactionLockOps) (transactionLock, error) {
	ops = ops.withDefaults()
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
	return transactionLock{path: path, id: id, ops: ops}, nil
}

func (l transactionLock) Release() error {
	if l.path == "" {
		return nil
	}
	if _, err := removeTransactionLockIfCurrent(l.path, l.id, l.ops); err != nil {
		if errors.Is(err, errTransactionLockDisappeared) {
			return nil
		}
		return err
	}
	return nil
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
	seen := map[string]bool{}
	for lineNo, rawLine := range strings.Split(string(data), "\n") {
		line := strings.TrimSuffix(rawLine, "\r")
		if strings.TrimSpace(line) == "" {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return TransactionLockInfo{}, lockCorruptf("malformed publish lock line %d", lineNo+1)
		}
		if seen[key] {
			return TransactionLockInfo{}, lockCorruptf("duplicate publish lock key %q", key)
		}
		seen[key] = true
		switch key {
		case "transaction":
			info.ID = TransactionID(value)
		case "pid":
			if strings.TrimSpace(value) == "" {
				return TransactionLockInfo{}, lockCorruptf("publish lock pid is empty")
			}
			if !isASCIIInteger(value) {
				return TransactionLockInfo{}, lockCorruptf("publish lock pid is not numeric")
			}
			info.PID = value
		case "startedAt":
			if strings.TrimSpace(value) == "" {
				return TransactionLockInfo{}, lockCorruptf("publish lock startedAt is empty")
			}
			if _, err := time.Parse(time.RFC3339Nano, value); err != nil {
				return TransactionLockInfo{}, lockCorruptf("publish lock startedAt is invalid")
			}
			info.StartedAt = value
		case "command":
			if strings.TrimSpace(value) == "" {
				return TransactionLockInfo{}, lockCorruptf("publish lock command is empty")
			}
			info.Command = value
		default:
			return TransactionLockInfo{}, lockCorruptf("unknown publish lock key %q", key)
		}
	}
	if info.ID == "" {
		return TransactionLockInfo{}, lockCorruptf("publish lock is missing transaction id")
	}
	if err := validateTransactionID(info.ID); err != nil {
		return TransactionLockInfo{}, lockCorruptf("publish lock has unsafe transaction id %q", info.ID)
	}
	info.Path = path
	return info, nil
}

func isASCIIInteger(value string) bool {
	if value == "" {
		return false
	}
	for _, ch := range value {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}
