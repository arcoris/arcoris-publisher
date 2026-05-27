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
// safety. It is shared by preflight and future operator diagnostics so publish,
// prune, lock inspection, and recovery use the same status policy.
type TransactionStateDiagnostics struct {
	PublishBlocked bool
	Blockers       []TransactionStateBlocker
	Lock           LockShowResult
	Journals       []JournalDiagnostic
	Warnings       []TransactionDiagnosticWarning
}

// TransactionStateBlocker describes one reason a new publish must not start.
type TransactionStateBlocker struct {
	Kind          TransactionStateBlockerKind
	TransactionID TransactionID
	Status        TransactionStatus
	Reason        string
}

// TransactionStateBlockerKind is a stable machine-readable blocker class.
type TransactionStateBlockerKind string

const (
	TransactionBlockerActiveJournal       TransactionStateBlockerKind = "active_journal"
	TransactionBlockerFailedJournal       TransactionStateBlockerKind = "failed_journal"
	TransactionBlockerRollbackFailed      TransactionStateBlockerKind = "rollback_failed_journal"
	TransactionBlockerPublishLock         TransactionStateBlockerKind = "publish_lock"
	TransactionBlockerCorruptLock         TransactionStateBlockerKind = "corrupt_lock"
	TransactionBlockerLockReadFailed      TransactionStateBlockerKind = "lock_read_failed"
	TransactionBlockerCorruptJournal      TransactionStateBlockerKind = "corrupt_journal"
	TransactionBlockerMissingLockJournal  TransactionStateBlockerKind = "missing_lock_journal"
	TransactionBlockerJournalReadFailed   TransactionStateBlockerKind = "journal_read_failed"
	TransactionBlockerStateDirUnavailable TransactionStateBlockerKind = "state_dir_unavailable"
)

const (
	transactionDiagnosticReasonActiveJournal      = "active_transaction"
	transactionDiagnosticReasonFailedJournal      = "failed_transaction"
	transactionDiagnosticReasonRollbackFailed     = "rollback_failed_transaction"
	transactionDiagnosticReasonStaleTerminalLock  = "stale_terminal_transaction"
	transactionDiagnosticReasonRecoveryLock       = "recovery_transaction"
	transactionDiagnosticReasonMissingLockJournal = "missing_lock_journal"
	transactionDiagnosticReasonCorruptLock        = "corrupt_lock"
	transactionDiagnosticReasonLockReadFailed     = "lock_read_failed"
	transactionDiagnosticReasonCorruptJournal     = "corrupt_journal"
	transactionDiagnosticReasonJournalReadFailed  = "journal_read_failed"
	transactionDiagnosticReasonStateDirMissing    = "state_dir_missing"
)

// JournalDiagnostic describes one transaction journal without mutating it.
type JournalDiagnostic struct {
	ID               TransactionID
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

// TransactionDiagnosticWarning is non-fatal read-only diagnostic context.
type TransactionDiagnosticWarning struct {
	Code    string
	Message string
}

// InspectTransactionState reports transaction journals and publish.lock
// blockers without clearing, pruning, rolling back, or touching Git state.
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
			Reason: transactionDiagnosticReasonStateDirMissing,
		})
		return diagnostics, &Error{Code: CodeLockFailed, Message: diagnostics.Lock.Message}
	}

	lockResult, lockErr := InspectTransactionLock(ctx, stateDir)
	diagnostics.Lock = lockResult
	diagnostics.addLockBlockers(lockResult)
	diagnostics.addLockWarnings(lockResult.Warnings)

	journals, err := inspectJournalDiagnostics(ctx, stateDir)
	if err != nil {
		diagnostics.addBlocker(TransactionStateBlocker{
			Kind:   TransactionBlockerJournalReadFailed,
			Reason: transactionDiagnosticReasonJournalReadFailed,
		})
		diagnostics.PublishBlocked = len(diagnostics.Blockers) > 0
		return diagnostics, err
	}
	diagnostics.Journals = journals
	for _, journal := range journals {
		diagnostics.addJournalBlocker(journal)
		diagnostics.addJournalWarning(journal)
	}
	diagnostics.PublishBlocked = len(diagnostics.Blockers) > 0

	if lockErr != nil && lockResult.Status == LockShowStatusFailed {
		return diagnostics, lockErr
	}
	return diagnostics, nil
}

func (d *TransactionStateDiagnostics) addLockBlockers(result LockShowResult) {
	switch result.Status {
	case LockShowStatusPresent:
		reason := transactionDiagnosticReasonActiveJournal
		if !result.Journal.Status.BlocksNewPublish() {
			reason = transactionDiagnosticReasonStaleTerminalLock
		} else if result.Journal.Status.AllowsLockClear() {
			reason = transactionDiagnosticReasonRecoveryLock
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
			Reason:        transactionDiagnosticReasonMissingLockJournal,
		})
	case LockShowStatusCorrupt:
		d.addBlocker(TransactionStateBlocker{
			Kind:   TransactionBlockerCorruptLock,
			Reason: transactionDiagnosticReasonCorruptLock,
		})
	case LockShowStatusJournalCorrupt:
		d.addBlocker(TransactionStateBlocker{
			Kind:          TransactionBlockerCorruptJournal,
			TransactionID: result.Lock.ID,
			Reason:        transactionDiagnosticReasonCorruptJournal,
		})
	case LockShowStatusFailed:
		d.addLockFailureBlocker(result)
	}
}

func (d *TransactionStateDiagnostics) addLockFailureBlocker(result LockShowResult) {
	switch result.Reason {
	case LockShowReasonJournalReadFailed:
		d.addBlocker(TransactionStateBlocker{
			Kind:          TransactionBlockerJournalReadFailed,
			TransactionID: result.Lock.ID,
			Reason:        transactionDiagnosticReasonJournalReadFailed,
		})
	case LockShowReasonStateDirMissing:
		d.addBlocker(TransactionStateBlocker{
			Kind:   TransactionBlockerStateDirUnavailable,
			Reason: transactionDiagnosticReasonStateDirMissing,
		})
	case LockShowReasonLockReadFailed:
		d.addBlocker(TransactionStateBlocker{
			Kind:   TransactionBlockerLockReadFailed,
			Reason: transactionDiagnosticReasonLockReadFailed,
		})
	}
}

func (d *TransactionStateDiagnostics) addJournalBlocker(journal JournalDiagnostic) {
	switch {
	case journal.Corrupt:
		d.addBlocker(TransactionStateBlocker{
			Kind:          TransactionBlockerCorruptJournal,
			TransactionID: journal.ID,
			Reason:        transactionDiagnosticReasonCorruptJournal,
		})
	case journal.ReadFailed:
		d.addBlocker(TransactionStateBlocker{
			Kind:          TransactionBlockerJournalReadFailed,
			TransactionID: journal.ID,
			Reason:        transactionDiagnosticReasonJournalReadFailed,
		})
	case journal.Status.BlocksNewPublish():
		d.addBlocker(TransactionStateBlocker{
			Kind:          journalBlockerKind(journal.Status),
			TransactionID: journal.ID,
			Status:        journal.Status,
			Reason:        journalBlockerReason(journal.Status),
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
			Code:    "journal_read_failed",
			Message: journal.Message,
		})
	}
}

func (d *TransactionStateDiagnostics) addBlocker(blocker TransactionStateBlocker) {
	for _, existing := range d.Blockers {
		if existing.Kind == blocker.Kind &&
			existing.TransactionID == blocker.TransactionID &&
			existing.Status == blocker.Status &&
			existing.Reason == blocker.Reason {
			return
		}
	}
	d.Blockers = append(d.Blockers, blocker)
	d.PublishBlocked = true
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

func journalBlockerReason(status TransactionStatus) string {
	switch status {
	case TransactionStatusFailed:
		return transactionDiagnosticReasonFailedJournal
	case TransactionStatusRollbackFailed:
		return transactionDiagnosticReasonRollbackFailed
	default:
		return transactionDiagnosticReasonActiveJournal
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
	out := JournalDiagnostic{ID: id, Path: path}
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
