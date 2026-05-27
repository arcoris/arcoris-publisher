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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// TransactionStateDiagnostics is a read-only snapshot of transaction state
// safety. Blockers describe inspected state; only incomplete inspection should
// make diagnostics commands fail.
type TransactionStateDiagnostics struct {
	PublishBlocked bool
	Blockers       []TransactionStateBlocker
	Lock           LockShowResult
	OperationLock  OperationLockDiagnostic
	Journals       []JournalDiagnostic
	Warnings       []TransactionDiagnosticWarning
}

// TransactionStateBlocker describes one reason a new publish must not start.
type TransactionStateBlocker struct {
	Kind          TransactionStateBlockerKind
	TransactionID TransactionID
	Status        TransactionStatus
	Reason        TransactionStateBlockerReason
	Name          string
}

// TransactionStateBlockerKind is a stable machine-readable blocker class.
type TransactionStateBlockerKind string

const (
	TransactionBlockerActiveJournal              TransactionStateBlockerKind = "active_journal"
	TransactionBlockerFailedJournal              TransactionStateBlockerKind = "failed_journal"
	TransactionBlockerRollbackFailed             TransactionStateBlockerKind = "rollback_failed_journal"
	TransactionBlockerPublishLock                TransactionStateBlockerKind = "publish_lock"
	TransactionBlockerCorruptLock                TransactionStateBlockerKind = "corrupt_lock"
	TransactionBlockerLockReadFailed             TransactionStateBlockerKind = "lock_read_failed"
	TransactionBlockerCorruptJournal             TransactionStateBlockerKind = "corrupt_journal"
	TransactionBlockerMissingLockJournal         TransactionStateBlockerKind = "missing_lock_journal"
	TransactionBlockerJournalFileReadFailed      TransactionStateBlockerKind = "journal_file_read_failed"
	TransactionBlockerJournalDirectoryReadFailed TransactionStateBlockerKind = "journal_directory_read_failed"
	TransactionBlockerStateDirUnavailable        TransactionStateBlockerKind = "state_dir_unavailable"
	TransactionBlockerOperationLock              TransactionStateBlockerKind = "operation_lock"
	TransactionBlockerCorruptOperationLock       TransactionStateBlockerKind = "corrupt_operation_lock"
	TransactionBlockerOperationLockReadFailed    TransactionStateBlockerKind = "operation_lock_read_failed"
)

// TransactionStateBlockerReason is a stable machine-readable blocker reason.
type TransactionStateBlockerReason string

const (
	TransactionBlockerReasonActiveJournal              TransactionStateBlockerReason = "active_transaction"
	TransactionBlockerReasonFailedJournal              TransactionStateBlockerReason = "failed_transaction"
	TransactionBlockerReasonRollbackFailed             TransactionStateBlockerReason = "rollback_failed_transaction"
	TransactionBlockerReasonStaleTerminalLock          TransactionStateBlockerReason = "stale_terminal_transaction"
	TransactionBlockerReasonRecoveryLock               TransactionStateBlockerReason = "recovery_transaction"
	TransactionBlockerReasonMissingLockJournal         TransactionStateBlockerReason = "missing_lock_journal"
	TransactionBlockerReasonCorruptLock                TransactionStateBlockerReason = "corrupt_lock"
	TransactionBlockerReasonLockReadFailed             TransactionStateBlockerReason = "lock_read_failed"
	TransactionBlockerReasonCorruptJournal             TransactionStateBlockerReason = "corrupt_journal"
	TransactionBlockerReasonJournalFileReadFailed      TransactionStateBlockerReason = "journal_file_read_failed"
	TransactionBlockerReasonJournalDirectoryReadFailed TransactionStateBlockerReason = "journal_directory_read_failed"
	TransactionBlockerReasonStateDirMissing            TransactionStateBlockerReason = "state_dir_missing"
	TransactionBlockerReasonOperationLockExists        TransactionStateBlockerReason = "operation_lock_exists"
	TransactionBlockerReasonOperationLockCorrupt       TransactionStateBlockerReason = "operation_lock_corrupt"
	TransactionBlockerReasonOperationLockReadFailed    TransactionStateBlockerReason = "operation_lock_read_failed"
)

// JournalDiagnostic describes one transaction journal without mutating it.
type JournalDiagnostic struct {
	ID               TransactionID
	Name             string
	Status           TransactionStatus
	Rollback         RollbackStatus
	Version          string
	StartedAt        time.Time
	UpdatedAt        time.Time
	Prunable         bool
	BlocksNewPublish bool
	AllowsLockClear  bool
	Path             string
	Corrupt          bool
	ReadFailed       bool
	Message          string
}

// OperationLockDiagnostic reports operation.lock without exposing its identity
// token or clearing stale locks.
type OperationLockDiagnostic struct {
	Present    bool
	Operation  string
	PID        string
	StartedAt  string
	Path       string
	Corrupt    bool
	ReadFailed bool
	Message    string
}

// TransactionDiagnosticWarning is non-fatal read-only diagnostic context.
type TransactionDiagnosticWarning struct {
	Code    string
	Message string
}

// InspectTransactionState reports transaction journals and locks without
// clearing, pruning, rolling back, or touching Git state.
func InspectTransactionState(ctx context.Context, stateDir string) (TransactionStateDiagnostics, error) {
	var diagnostics TransactionStateDiagnostics
	if err := ctx.Err(); err != nil {
		diagnostics.Lock = lockShowFailed(lockShowContextReason(err))
		return diagnostics, err
	}
	if strings.TrimSpace(stateDir) == "" {
		diagnostics.Lock = lockShowFailed(LockShowReasonStateDirMissing, "transaction state directory is unavailable")
		diagnostics.addBlocker(TransactionStateBlocker{
			Kind:   TransactionBlockerStateDirUnavailable,
			Reason: TransactionBlockerReasonStateDirMissing,
		})
		diagnostics.finalize()
		return diagnostics, &Error{Code: CodeLockFailed, Message: diagnostics.Lock.Message}
	}

	lockResult, lockErr := InspectTransactionLock(ctx, stateDir)
	diagnostics.Lock = lockResult
	diagnostics.addLockBlockers(lockResult)
	diagnostics.addLockWarnings(lockResult.Warnings)
	diagnostics.OperationLock = inspectOperationLock(stateDir)
	diagnostics.addOperationLockBlocker(diagnostics.OperationLock)
	diagnostics.addOperationLockWarning(diagnostics.OperationLock)

	journals, err := inspectJournalDiagnostics(ctx, stateDir)
	if err != nil {
		diagnostics.addBlocker(TransactionStateBlocker{
			Kind:   TransactionBlockerJournalDirectoryReadFailed,
			Reason: TransactionBlockerReasonJournalDirectoryReadFailed,
		})
		diagnostics.addWarning(TransactionDiagnosticWarning{
			Code:    string(TransactionBlockerReasonJournalDirectoryReadFailed),
			Message: "read transaction journals directory failed",
		})
		diagnostics.finalize()
		return diagnostics, err
	}
	diagnostics.Journals = journals
	for _, journal := range journals {
		diagnostics.addJournalBlocker(journal)
		diagnostics.addJournalWarning(journal)
	}
	diagnostics.finalize()

	if lockErr != nil && lockResult.Status == LockShowStatusFailed {
		return diagnostics, lockErr
	}
	return diagnostics, nil
}

func (d *TransactionStateDiagnostics) addLockBlockers(result LockShowResult) {
	switch result.Status {
	case LockShowStatusPresent:
		reason := TransactionBlockerReasonActiveJournal
		if !result.Journal.Status.BlocksNewPublish() {
			reason = TransactionBlockerReasonStaleTerminalLock
		} else if result.Journal.Status.AllowsLockClear() {
			reason = TransactionBlockerReasonRecoveryLock
		}
		d.addBlocker(TransactionStateBlocker{
			Kind:          TransactionBlockerPublishLock,
			TransactionID: result.Lock.ID,
			Status:        result.Journal.Status,
			Reason:        reason,
		})
	case LockShowStatusJournalMissing:
		d.addBlocker(TransactionStateBlocker{
			Kind:          TransactionBlockerMissingLockJournal,
			TransactionID: result.Lock.ID,
			Reason:        TransactionBlockerReasonMissingLockJournal,
		})
	case LockShowStatusCorrupt:
		d.addBlocker(TransactionStateBlocker{
			Kind:   TransactionBlockerCorruptLock,
			Reason: TransactionBlockerReasonCorruptLock,
		})
	case LockShowStatusJournalCorrupt:
		d.addBlocker(TransactionStateBlocker{
			Kind:          TransactionBlockerCorruptJournal,
			TransactionID: result.Lock.ID,
			Reason:        TransactionBlockerReasonCorruptJournal,
		})
	case LockShowStatusFailed:
		d.addLockFailureBlocker(result)
	}
}

func (d *TransactionStateDiagnostics) addLockFailureBlocker(result LockShowResult) {
	switch result.Reason {
	case LockShowReasonJournalReadFailed:
		d.addBlocker(TransactionStateBlocker{
			Kind:          TransactionBlockerJournalFileReadFailed,
			TransactionID: result.Lock.ID,
			Reason:        TransactionBlockerReasonJournalFileReadFailed,
		})
	case LockShowReasonStateDirMissing:
		d.addBlocker(TransactionStateBlocker{
			Kind:   TransactionBlockerStateDirUnavailable,
			Reason: TransactionBlockerReasonStateDirMissing,
		})
	case LockShowReasonLockReadFailed:
		d.addBlocker(TransactionStateBlocker{
			Kind:   TransactionBlockerLockReadFailed,
			Reason: TransactionBlockerReasonLockReadFailed,
		})
	}
}

func (d *TransactionStateDiagnostics) addOperationLockBlocker(lock OperationLockDiagnostic) {
	switch {
	case lock.Corrupt:
		d.addBlocker(TransactionStateBlocker{
			Kind:   TransactionBlockerCorruptOperationLock,
			Reason: TransactionBlockerReasonOperationLockCorrupt,
		})
	case lock.ReadFailed:
		d.addBlocker(TransactionStateBlocker{
			Kind:   TransactionBlockerOperationLockReadFailed,
			Reason: TransactionBlockerReasonOperationLockReadFailed,
		})
	case lock.Present:
		d.addBlocker(TransactionStateBlocker{
			Kind:   TransactionBlockerOperationLock,
			Reason: TransactionBlockerReasonOperationLockExists,
			Name:   lock.Operation,
		})
	}
}

func (d *TransactionStateDiagnostics) addJournalBlocker(journal JournalDiagnostic) {
	switch {
	case journal.Corrupt:
		d.addBlocker(TransactionStateBlocker{
			Kind:          TransactionBlockerCorruptJournal,
			TransactionID: journal.ID,
			Reason:        TransactionBlockerReasonCorruptJournal,
			Name:          journal.Name,
		})
	case journal.ReadFailed:
		d.addBlocker(TransactionStateBlocker{
			Kind:          TransactionBlockerJournalFileReadFailed,
			TransactionID: journal.ID,
			Reason:        TransactionBlockerReasonJournalFileReadFailed,
			Name:          journal.Name,
		})
	case journal.Status.BlocksNewPublish():
		d.addBlocker(TransactionStateBlocker{
			Kind:          journalBlockerKind(journal.Status),
			TransactionID: journal.ID,
			Status:        journal.Status,
			Reason:        journalBlockerReason(journal.Status),
			Name:          journal.Name,
		})
	}
}

func (d *TransactionStateDiagnostics) addJournalWarning(journal JournalDiagnostic) {
	switch {
	case journal.Corrupt:
		d.addWarning(TransactionDiagnosticWarning{
			Code:    string(LockWarningJournalCorrupt),
			Message: journal.Message,
		})
	case journal.ReadFailed:
		d.addWarning(TransactionDiagnosticWarning{
			Code:    string(TransactionBlockerReasonJournalFileReadFailed),
			Message: journal.Message,
		})
	}
}

func (d *TransactionStateDiagnostics) addBlocker(blocker TransactionStateBlocker) {
	for _, existing := range d.Blockers {
		if existing.Kind == blocker.Kind &&
			existing.TransactionID == blocker.TransactionID &&
			existing.Status == blocker.Status &&
			existing.Reason == blocker.Reason &&
			existing.Name == blocker.Name {
			return
		}
	}
	d.Blockers = append(d.Blockers, blocker)
}

func (d *TransactionStateDiagnostics) finalize() {
	d.PublishBlocked = len(d.Blockers) > 0
}

func journalBlockerKind(status TransactionStatus) TransactionStateBlockerKind {
	switch status {
	case TransactionStatusFailed:
		return TransactionBlockerFailedJournal
	case TransactionStatusRollbackFailed:
		return TransactionBlockerRollbackFailed
	default:
		return TransactionBlockerActiveJournal
	}
}

func journalBlockerReason(status TransactionStatus) TransactionStateBlockerReason {
	switch status {
	case TransactionStatusFailed:
		return TransactionBlockerReasonFailedJournal
	case TransactionStatusRollbackFailed:
		return TransactionBlockerReasonRollbackFailed
	default:
		return TransactionBlockerReasonActiveJournal
	}
}

func (d *TransactionStateDiagnostics) addOperationLockWarning(lock OperationLockDiagnostic) {
	switch {
	case lock.Corrupt:
		d.addWarning(TransactionDiagnosticWarning{
			Code:    string(TransactionBlockerReasonOperationLockCorrupt),
			Message: lock.Message,
		})
	case lock.ReadFailed:
		d.addWarning(TransactionDiagnosticWarning{
			Code:    string(TransactionBlockerReasonOperationLockReadFailed),
			Message: lock.Message,
		})
	}
}

func (d *TransactionStateDiagnostics) addLockWarnings(warnings []LockWarning) {
	for _, warning := range warnings {
		d.addWarning(TransactionDiagnosticWarning{Code: string(warning.Code), Message: warning.Message})
	}
}

func (d *TransactionStateDiagnostics) addWarning(warning TransactionDiagnosticWarning) {
	for _, existing := range d.Warnings {
		if existing.Code == warning.Code && existing.Message == warning.Message {
			return
		}
	}
	d.Warnings = append(d.Warnings, warning)
}

func inspectJournalDiagnostics(ctx context.Context, stateDir string) ([]JournalDiagnostic, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	store := NewFileJournalStore(stateDir)
	entries, err := os.ReadDir(store.transactionsDir())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, &Error{Code: CodeRecoveryFailed, Message: "read transaction journals failed", Cause: err}
	}

	journals := make([]JournalDiagnostic, 0, len(entries))
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		diagnostic := readJournalDiagnostic(ctx, store, entry.Name())
		journals = append(journals, diagnostic)
	}
	sort.Slice(journals, func(i, j int) bool {
		if journals[i].StartedAt.Equal(journals[j].StartedAt) {
			return journals[i].Path < journals[j].Path
		}
		return journals[i].StartedAt.Before(journals[j].StartedAt)
	})
	return journals, nil
}

func readJournalDiagnostic(ctx context.Context, store FileJournalStore, name string) JournalDiagnostic {
	id := diagnosticJournalID(name)
	path := filepath.Join(store.transactionsDir(), name)
	out := JournalDiagnostic{ID: id, Name: name, Path: path}
	if err := ctx.Err(); err != nil {
		out.ReadFailed = true
		out.Message = err.Error()
		return out
	}
	data, err := os.ReadFile(path)
	if err != nil {
		out.ReadFailed = true
		out.Message = fmt.Sprintf("read transaction journal %s failed: %v", name, err)
		return out
	}
	var journal TransactionJournal
	if err := json.Unmarshal(data, &journal); err != nil {
		out.Corrupt = true
		out.Message = fmt.Sprintf("transaction journal %s is corrupt: %v", name, err)
		return out
	}
	if err := validateJournalIdentity(name, "", journal.ID); err != nil {
		out.Corrupt = true
		out.Message = err.Error()
		return out
	}
	return JournalDiagnostic{
		ID:               journal.ID,
		Name:             name,
		Status:           journal.Status,
		Rollback:         journal.Rollback,
		Version:          journal.Version,
		StartedAt:        journal.StartedAt,
		UpdatedAt:        journal.UpdatedAt,
		Prunable:         journal.Status.Prunable(),
		BlocksNewPublish: journal.Status.BlocksNewPublish(),
		AllowsLockClear:  journal.Status.AllowsLockClear(),
		Path:             path,
	}
}

func diagnosticJournalID(name string) TransactionID {
	id := TransactionID(strings.TrimSuffix(name, ".json"))
	if validateTransactionID(id) != nil {
		return ""
	}
	return id
}

func inspectOperationLock(stateDir string) OperationLockDiagnostic {
	path := operationLockPath(stateDir)
	info, err := readOperationLock(path)
	if err == nil {
		return OperationLockDiagnostic{
			Present:   true,
			Operation: string(info.Operation),
			PID:       info.PID,
			StartedAt: info.StartedAt,
			Path:      path,
		}
	}
	if errors.Is(err, os.ErrNotExist) {
		return OperationLockDiagnostic{}
	}
	out := OperationLockDiagnostic{Path: path, Message: "read transaction operation lock failed"}
	if _, statErr := os.Stat(path); statErr == nil {
		out.Present = true
	}
	if errors.Is(err, errOperationLockCorrupt) {
		out.Present = true
		out.Corrupt = true
		out.Message = err.Error()
		return out
	}
	out.ReadFailed = true
	return out
}
