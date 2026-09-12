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
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// operationLock serializes transaction mutation decisions and writes. It does
// not replace publish.lock, and stale locks are reported but never auto-cleared.
type operationLock struct {
	path      string
	operation operationLockOperation
	token     string
	owned     bool
	ops       operationLockOps
}

type operationLockOperation string

const (
	operationLockPublish   operationLockOperation = "publish"
	operationLockRollback  operationLockOperation = "rollback"
	operationLockPrune     operationLockOperation = "prune"
	operationLockLockClear operationLockOperation = "lock-clear"
)

const operationLockSchemaVersion = "1"

var (
	errOperationLockExists        = errors.New("transaction operation lock exists")
	errOperationLockCorrupt       = errors.New("transaction operation lock corrupt")
	errOperationLockChanged       = errors.New("transaction operation lock changed")
	errOperationLockDisappeared   = errors.New("transaction operation lock disappeared")
	errOperationLockDeleteFailed  = errors.New("transaction operation lock delete failed")
	errOperationLockSyncFailed    = errors.New("transaction operation lock sync failed")
)

type operationLockInfo struct {
	Operation operationLockOperation
	Token     string
	PID       string
	StartedAt string
	Path      string
}

type operationLockOutcome struct {
	Released bool
	Synced   bool
}

type operationLockOps struct {
	remove            func(string) error
	syncAcquireParent func(string) error
	syncParent        func(string) error
	now               func() time.Time
	token             func() (string, error)
	beforeRemove      func()
}

func defaultOperationLockOps() operationLockOps {
	return operationLockOps{
		remove:            os.Remove,
		syncAcquireParent: syncParentDir,
		syncParent:        syncParentDir,
		now:               func() time.Time { return time.Now().UTC() },
		token:             randomOperationLockToken,
	}
}

func (ops operationLockOps) withDefaults() operationLockOps {
	defaults := defaultOperationLockOps()
	if ops.remove == nil {
		ops.remove = defaults.remove
	}
	if ops.syncAcquireParent == nil {
		ops.syncAcquireParent = defaults.syncAcquireParent
	}
	if ops.syncParent == nil {
		ops.syncParent = defaults.syncParent
	}
	if ops.now == nil {
		ops.now = defaults.now
	}
	if ops.token == nil {
		ops.token = defaults.token
	}
	return ops
}

func acquireOperationLock(ctx context.Context, stateDir string, operation operationLockOperation, ops operationLockOps) (operationLock, error) {
	ops = ops.withDefaults()
	if err := ctx.Err(); err != nil {
		return operationLock{}, err
	}
	if strings.TrimSpace(stateDir) == "" {
		return operationLock{}, fmt.Errorf("state dir is required")
	}
	if !operation.valid() {
		return operationLock{}, fmt.Errorf("operation %q is not supported", operation)
	}
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		return operationLock{}, err
	}
	token, err := ops.token()
	if err != nil {
		return operationLock{}, err
	}
	if !isOperationLockToken(token) {
		return operationLock{}, fmt.Errorf("operation lock identity is invalid")
	}

	path := operationLockPath(stateDir)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return operationLock{}, fmt.Errorf("%w: transaction state operation lock already exists", errOperationLockExists)
		}
		return operationLock{}, err
	}
	content := fmt.Sprintf(
		"schemaVersion=%s\noperation=%s\ntoken=%s\npid=%d\nstartedAt=%s\n",
		operationLockSchemaVersion,
		operation,
		token,
		os.Getpid(),
		ops.now().UTC().Format(time.RFC3339Nano),
	)
	if _, err := file.WriteString(content); err != nil {
		return operationLock{}, abortOperationLockAcquire(file, path, ops, err)
	}
	if err := file.Sync(); err != nil {
		return operationLock{}, abortOperationLockAcquire(file, path, ops, err)
	}
	if err := file.Close(); err != nil {
		return operationLock{}, joinOperationLockAcquireCleanup(path, ops, err)
	}
	if err := ops.syncAcquireParent(path); err != nil {
		return operationLock{}, joinOperationLockAcquireCleanup(path, ops, err)
	}
	return operationLock{path: path, operation: operation, token: token, owned: true, ops: ops}, nil
}

func abortOperationLockAcquire(file *os.File, path string, ops operationLockOps, primary error) error {
	if closeErr := file.Close(); closeErr != nil {
		primary = errors.Join(primary, fmt.Errorf("close partial transaction operation lock failed: %w", closeErr))
	}
	return joinOperationLockAcquireCleanup(path, ops, primary)
}

func joinOperationLockAcquireCleanup(path string, ops operationLockOps, primary error) error {
	outcome, cleanupErr := cleanupCreatedStateFile(path, ops.remove, ops.syncParent)
	if cleanupErr == nil {
		return primary
	}
	if !outcome.Removed {
		cleanupErr = fmt.Errorf("%w during acquire cleanup: %w", errOperationLockDeleteFailed, cleanupErr)
	} else {
		cleanupErr = fmt.Errorf("%w during acquire cleanup after lock removal: %w", errOperationLockSyncFailed, cleanupErr)
	}
	return errors.Join(primary, cleanupErr)
}

// Release removes the operation lock only when this handle owns a lock created
// by acquireOperationLock and the persisted identity still matches. An acquired
// lock disappearing before release is an observable lifecycle failure rather
// than a successful no-op.
func (l operationLock) Release() (operationLockOutcome, error) {
	if !l.owned || l.path == "" {
		return operationLockOutcome{}, nil
	}
	ops := l.ops.withDefaults()
	if ops.beforeRemove != nil {
		ops.beforeRemove()
	}
	info, err := readOperationLock(l.path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return operationLockOutcome{}, fmt.Errorf("%w: %v", errOperationLockDisappeared, err)
		}
		return operationLockOutcome{}, err
	}
	if info.Operation != l.operation || info.Token != l.token {
		return operationLockOutcome{}, fmt.Errorf("%w: transaction operation lock changed", errOperationLockChanged)
	}
	if err := ops.remove(l.path); err != nil {
		return operationLockOutcome{}, fmt.Errorf("%w: %w", errOperationLockDeleteFailed, err)
	}
	outcome := operationLockOutcome{Released: true}
	if err := ops.syncParent(l.path); err != nil {
		return outcome, fmt.Errorf("%w: %w", errOperationLockSyncFailed, err)
	}
	outcome.Synced = true
	return outcome, nil
}

func readOperationLock(path string) (operationLockInfo, error) {
	data, err := readBoundedStateFile(path, maxOperationLockBytes)
	if err != nil {
		if isUnsafeStateFileRepresentation(err) {
			return operationLockInfo{}, operationLockCorruptf("transaction operation lock has an unsafe file representation: %v", err)
		}
		return operationLockInfo{}, err
	}
	info := operationLockInfo{}
	seen := map[string]bool{}
	for lineNo, rawLine := range strings.Split(string(data), "\n") {
		line := strings.TrimSuffix(rawLine, "\r")
		if strings.TrimSpace(line) == "" {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return operationLockInfo{}, operationLockCorruptf("malformed transaction operation lock line %d", lineNo+1)
		}
		if seen[key] {
			return operationLockInfo{}, operationLockCorruptf("duplicate transaction operation lock key")
		}
		seen[key] = true
		switch key {
		case "schemaVersion":
			if value != operationLockSchemaVersion {
				return operationLockInfo{}, operationLockCorruptf("unsupported transaction operation lock schemaVersion %q", value)
			}
		case "operation":
			info.Operation = operationLockOperation(value)
			if !info.Operation.valid() {
				return operationLockInfo{}, operationLockCorruptf("unsupported transaction operation %q", value)
			}
		case "token":
			if !isOperationLockToken(value) {
				return operationLockInfo{}, operationLockCorruptf("transaction operation lock identity is invalid")
			}
			info.Token = value
		case "pid":
			if strings.TrimSpace(value) == "" {
				return operationLockInfo{}, operationLockCorruptf("transaction operation lock pid is empty")
			}
			if !isASCIIInteger(value) {
				return operationLockInfo{}, operationLockCorruptf("transaction operation lock pid is not numeric")
			}
			info.PID = value
		case "startedAt":
			if strings.TrimSpace(value) == "" {
				return operationLockInfo{}, operationLockCorruptf("transaction operation lock startedAt is empty")
			}
			if _, err := time.Parse(time.RFC3339Nano, value); err != nil {
				return operationLockInfo{}, operationLockCorruptf("transaction operation lock startedAt is invalid")
			}
			info.StartedAt = value
		default:
			return operationLockInfo{}, operationLockCorruptf("transaction operation lock contains unsupported key")
		}
	}
	if !seen["schemaVersion"] {
		return operationLockInfo{}, operationLockCorruptf("transaction operation lock schemaVersion is missing")
	}
	if info.Operation == "" {
		return operationLockInfo{}, operationLockCorruptf("transaction operation lock operation is missing")
	}
	if info.Token == "" {
		return operationLockInfo{}, operationLockCorruptf("transaction operation lock identity is missing")
	}
	if info.PID == "" {
		return operationLockInfo{}, operationLockCorruptf("transaction operation lock pid is missing")
	}
	if info.StartedAt == "" {
		return operationLockInfo{}, operationLockCorruptf("transaction operation lock startedAt is missing")
	}
	info.Path = path
	return info, nil
}

func operationLockPath(stateDir string) string {
	return filepath.Join(stateDir, "operation.lock")
}

func (operation operationLockOperation) valid() bool {
	switch operation {
	case operationLockPublish, operationLockRollback, operationLockPrune, operationLockLockClear:
		return true
	default:
		return false
	}
}

func operationLockCorruptf(format string, args ...any) error {
	return fmt.Errorf("%w: "+format, append([]any{errOperationLockCorrupt}, args...)...)
}

func randomOperationLockToken() (string, error) {
	var buf [16]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf[:]), nil
}

func isOperationLockToken(value string) bool {
	if value == "" {
		return false
	}
	for _, ch := range value {
		if (ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' ||
			ch == '_' {
			continue
		}
		return false
	}
	return true
}

func operationLockAcquireError(operation operationLockOperation, err error) *Error {
	message := fmt.Sprintf("transaction state operation lock failed for %s", operation)
	if errors.Is(err, errOperationLockExists) {
		message = fmt.Sprintf("transaction state operation lock already exists; refusing %s", operation)
	}
	return &Error{Code: CodeLockFailed, Message: message, Cause: err}
}

func operationLockReleaseError(outcome operationLockOutcome, err error) *Error {
	return &Error{Code: CodeLockFailed, Message: operationLockReleaseMessage(outcome, err), Cause: err}
}

func operationLockReleaseWarning(outcome operationLockOutcome, err error) string {
	return operationLockReleaseMessage(outcome, err) + ": " + err.Error()
}

func operationLockReleaseMessage(outcome operationLockOutcome, err error) string {
	switch {
	case errors.Is(err, errOperationLockSyncFailed) && outcome.Released:
		return "transaction state operation lock cleanup sync failed after lock removal"
	case errors.Is(err, errOperationLockDeleteFailed):
		return "transaction state operation lock cleanup delete failed"
	case errors.Is(err, errOperationLockDisappeared):
		return "transaction state operation lock cleanup found the acquired lock missing"
	case errors.Is(err, errOperationLockChanged):
		return "transaction state operation lock cleanup refused changed lock"
	case errors.Is(err, errOperationLockCorrupt):
		return "transaction state operation lock cleanup refused corrupt lock"
	default:
		return "transaction state operation lock cleanup failed"
	}
}
