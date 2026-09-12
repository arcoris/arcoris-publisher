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
	remove            func(string) error
	syncAcquireParent func(string) error
	syncParent        func(string) error
	beforeRemove      func()
}

type lockRemoveOutcome struct {
	Removed bool
	Synced  bool
}

type lockReleaseOutcome struct {
	Removed bool
	Synced  bool
}

func defaultTransactionLockOps() transactionLockOps {
	return transactionLockOps{
		remove:            os.Remove,
		syncAcquireParent: syncParentDir,
		syncParent:        syncParentDir,
	}
}

func (ops transactionLockOps) withDefaults() transactionLockOps {
	defaults := defaultTransactionLockOps()
	if ops.remove == nil {
		ops.remove = defaults.remove
	}
	if ops.syncAcquireParent == nil {
		ops.syncAcquireParent = defaults.syncAcquireParent
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
		return lockRemoveOutcome{}, fmt.Errorf("%w: %w", errTransactionLockDeleteFailed, err)
	}
	outcome := lockRemoveOutcome{Removed: true}
	if err := ops.syncParent(path); err != nil {
		return outcome, fmt.Errorf("%w: %w", errTransactionLockSyncFailed, err)
	}
	outcome.Synced = true
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

const transactionLockSchemaVersion = "1"

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

	content := fmt.Sprintf("schemaVersion=%s\ntransaction=%s\npid=%d\nstartedAt=%s\ncommand=publish\n", transactionLockSchemaVersion, id, os.Getpid(), now.UTC().Format(time.RFC3339Nano))
	if _, err := file.WriteString(content); err != nil {
		return transactionLock{}, abortTransactionLockAcquire(file, path, ops, err)
	}
	if err := file.Sync(); err != nil {
		return transactionLock{}, abortTransactionLockAcquire(file, path, ops, err)
	}
	if err := file.Close(); err != nil {
		return transactionLock{}, joinTransactionLockAcquireCleanup(path, ops, err)
	}
	if err := ops.syncAcquireParent(path); err != nil {
		return transactionLock{}, joinTransactionLockAcquireCleanup(path, ops, err)
	}
	return transactionLock{path: path, id: id, ops: ops}, nil
}

func abortTransactionLockAcquire(file *os.File, path string, ops transactionLockOps, primary error) error {
	if closeErr := file.Close(); closeErr != nil {
		primary = errors.Join(primary, fmt.Errorf("close partial publish lock failed: %w", closeErr))
	}
	return joinTransactionLockAcquireCleanup(path, ops, primary)
}

func joinTransactionLockAcquireCleanup(path string, ops transactionLockOps, primary error) error {
	outcome, cleanupErr := cleanupCreatedStateFile(path, ops.remove, ops.syncParent)
	if cleanupErr == nil {
		return primary
	}
	if !outcome.Removed {
		cleanupErr = fmt.Errorf("%w during acquire cleanup: %w", errTransactionLockDeleteFailed, cleanupErr)
	} else {
		cleanupErr = fmt.Errorf("%w during acquire cleanup after lock removal: %w", errTransactionLockSyncFailed, cleanupErr)
	}
	return errors.Join(primary, cleanupErr)
}

// Release removes this caller's publish lock. An acquired lock disappearing
// before release is an ownership failure, not a successful no-op; hiding that
// race would make the transaction lifecycle look cleaner than the persisted
// state actually was.
func (l transactionLock) Release() (lockReleaseOutcome, error) {
	if l.path == "" {
		return lockReleaseOutcome{}, nil
	}
	outcome, err := removeTransactionLockIfCurrent(l.path, l.id, l.ops)
	release := lockReleaseOutcome{Removed: outcome.Removed, Synced: outcome.Synced}
	if err != nil {
		return release, err
	}
	return release, nil
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
	data, err := readBoundedStateFile(path, maxTransactionLockBytes)
	if err != nil {
		if isUnsafeStateFileRepresentation(err) {
			return TransactionLockInfo{}, lockCorruptf("publish lock has an unsafe file representation: %v", err)
		}
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
		case "schemaVersion":
			if strings.TrimSpace(value) == "" {
				return TransactionLockInfo{}, lockCorruptf("publish lock schemaVersion is empty")
			}
			if value != transactionLockSchemaVersion {
				return TransactionLockInfo{}, lockCorruptf("unsupported publish lock schemaVersion %q", value)
			}
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
			if value != "publish" {
				return TransactionLockInfo{}, lockCorruptf("unsupported publish lock command %q", value)
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
	if seen["schemaVersion"] {
		if info.PID == "" {
			return TransactionLockInfo{}, lockCorruptf("publish lock pid is missing")
		}
		if info.StartedAt == "" {
			return TransactionLockInfo{}, lockCorruptf("publish lock startedAt is missing")
		}
		if info.Command == "" {
			return TransactionLockInfo{}, lockCorruptf("publish lock command is missing")
		}
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
