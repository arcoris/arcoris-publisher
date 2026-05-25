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
	"time"
)

// LockShowStatus describes the inspected publish lock state.
type LockShowStatus string

const (
	// LockShowStatusAbsent means no publish lock exists.
	LockShowStatusAbsent LockShowStatus = "absent"
	// LockShowStatusPresent means the lock and referenced journal are readable.
	LockShowStatusPresent LockShowStatus = "present"
	// LockShowStatusCorrupt means the lock exists but cannot be parsed safely.
	LockShowStatusCorrupt LockShowStatus = "corrupt"
	// LockShowStatusJournalMissing means the lock references no journal file.
	LockShowStatusJournalMissing LockShowStatus = "journal_missing"
	// LockShowStatusJournalCorrupt means the referenced journal is unreadable.
	LockShowStatusJournalCorrupt LockShowStatus = "journal_corrupt"
)

// LockClearStatus describes a guarded publish lock clear attempt.
type LockClearStatus string

const (
	// LockClearStatusCleared means publish.lock was removed.
	LockClearStatusCleared LockClearStatus = "cleared"
	// LockClearStatusRefused means policy refused to remove publish.lock.
	LockClearStatusRefused LockClearStatus = "refused"
	// LockClearStatusFailed means lock inspection or deletion failed.
	LockClearStatusFailed LockClearStatus = "failed"
)

// LockWarning describes recoverable lock inspection information.
type LockWarning struct {
	Code    string
	Message string
}

// LockJournalState summarizes the transaction journal referenced by a lock.
type LockJournalState struct {
	Present   bool
	Status    TransactionStatus
	Rollback  RollbackStatus
	Version   string
	StartedAt time.Time
	UpdatedAt time.Time
}

// LockShowResult reports publish lock state without mutating it.
type LockShowResult struct {
	Status   LockShowStatus
	Lock     TransactionLockInfo
	Journal  LockJournalState
	Warnings []LockWarning
}

// LockClearOptions are the explicit operator guardrails for lock clearing.
type LockClearOptions struct {
	TransactionID TransactionID
	Confirm       TransactionID
}

// LockClearResult reports a guarded publish lock clear attempt.
type LockClearResult struct {
	Status        LockClearStatus
	TransactionID TransactionID
	Lock          TransactionLockInfo
	LockCleared   bool
	Journal       LockJournalState
	Reason        string
	Warnings      []LockWarning
}

// ShowTransactionLock inspects publish.lock and its referenced journal.
func (s Service) ShowTransactionLock(ctx context.Context, stateDir string) (LockShowResult, error) {
	if err := ctx.Err(); err != nil {
		return LockShowResult{Status: LockShowStatusCorrupt}, err
	}
	lock, ok, err := currentTransactionLock(stateDir)
	if err != nil {
		return LockShowResult{
			Status:   LockShowStatusCorrupt,
			Warnings: []LockWarning{{Code: "lock_corrupt", Message: "publish lock is not parseable"}},
		}, &Error{Code: CodeLockFailed, Message: "read publish transaction lock failed", Cause: err}
	}
	if !ok {
		return LockShowResult{Status: LockShowStatusAbsent}, nil
	}

	journal, warning, err := loadLockJournal(ctx, stateDir, lock.ID)
	result := LockShowResult{Status: LockShowStatusPresent, Lock: lock, Journal: journal}
	if warning.Code != "" {
		result.Warnings = append(result.Warnings, warning)
	}
	if err != nil {
		result.Status = LockShowStatusJournalCorrupt
		return result, &Error{Code: CodeRecoveryFailed, Message: "load publish transaction journal failed", Cause: err}
	}
	if !journal.Present {
		result.Status = LockShowStatusJournalMissing
	}
	return result, nil
}

// ClearTransactionLock removes only publish.lock after explicit confirmation.
func (s Service) ClearTransactionLock(ctx context.Context, stateDir string, opts LockClearOptions) (LockClearResult, error) {
	result := LockClearResult{
		Status:        LockClearStatusRefused,
		TransactionID: opts.TransactionID,
	}
	if err := ctx.Err(); err != nil {
		result.Status = LockClearStatusFailed
		return result, err
	}
	if opts.TransactionID == "" {
		result.Reason = "missing transaction id"
		return result, &Error{Code: CodeLockFailed, Message: "transaction id is required"}
	}
	if opts.Confirm == "" {
		result.Reason = "missing confirmation"
		return result, &Error{Code: CodeLockFailed, Message: "confirmation transaction id is required"}
	}
	if opts.Confirm != opts.TransactionID {
		result.Reason = "confirmation does not match transaction id"
		return result, &Error{Code: CodeLockFailed, Message: "confirmation transaction id does not match"}
	}
	if err := validateTransactionID(opts.TransactionID); err != nil {
		result.Reason = "invalid transaction id"
		return result, &Error{Code: CodeLockFailed, Message: "invalid transaction id", Cause: err}
	}

	path, err := transactionLockPath(stateDir)
	if err != nil {
		result.Status = LockClearStatusFailed
		result.Reason = "state dir is required"
		return result, &Error{Code: CodeLockFailed, Message: "resolve publish lock path failed", Cause: err}
	}
	lock, err := readTransactionLock(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			result.Reason = "publish lock is absent"
			return result, &Error{Code: CodeLockFailed, Message: "publish lock is absent"}
		}
		result.Status = LockClearStatusFailed
		result.Reason = "publish lock is not parseable"
		result.Warnings = append(result.Warnings, LockWarning{Code: "lock_corrupt", Message: "publish lock is not parseable"})
		return result, &Error{Code: CodeLockFailed, Message: "read publish transaction lock failed", Cause: err}
	}
	result.Lock = lock
	if lock.ID != opts.TransactionID {
		result.Reason = "lock transaction does not match requested transaction"
		return result, &Error{Code: CodeLockFailed, Message: fmt.Sprintf("publish lock belongs to %s, not %s", lock.ID, opts.TransactionID)}
	}

	journal, warning, err := loadLockJournal(ctx, stateDir, lock.ID)
	result.Journal = journal
	if warning.Code != "" {
		result.Warnings = append(result.Warnings, warning)
	}
	if err != nil {
		result.Status = LockClearStatusFailed
		result.Reason = "referenced transaction journal is corrupt"
		return result, &Error{Code: CodeRecoveryFailed, Message: "load publish transaction journal failed", Cause: err}
	}
	if journal.Present && !canClearLockForStatus(journal.Status) {
		result.Reason = "referenced transaction is active"
		return result, &Error{Code: CodeLockFailed, Message: fmt.Sprintf("transaction %s is %s; refusing to clear active publish lock", lock.ID, journal.Status)}
	}

	if err := os.Remove(path); err != nil {
		result.Status = LockClearStatusFailed
		result.Reason = "delete publish lock failed"
		return result, &Error{Code: CodeLockFailed, Message: "delete publish lock failed", Cause: err}
	}
	if err := syncParentDir(path); err != nil {
		result.Status = LockClearStatusFailed
		result.Reason = "sync publish lock directory failed"
		return result, &Error{Code: CodeLockFailed, Message: "sync publish lock directory failed", Cause: err}
	}

	result.Status = LockClearStatusCleared
	result.LockCleared = true
	result.Reason = "publish lock cleared"
	return result, nil
}

func loadLockJournal(ctx context.Context, stateDir string, id TransactionID) (LockJournalState, LockWarning, error) {
	journal, err := NewFileJournalStore(stateDir).Load(ctx, id)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return LockJournalState{}, LockWarning{
				Code:    "journal_missing",
				Message: "publish lock references a transaction journal that does not exist",
			}, nil
		}
		return LockJournalState{}, LockWarning{
			Code:    "journal_corrupt",
			Message: "publish lock references a transaction journal that cannot be read safely",
		}, err
	}
	return LockJournalState{
		Present:   true,
		Status:    journal.Status,
		Rollback:  journal.Rollback,
		Version:   journal.Version,
		StartedAt: journal.StartedAt,
		UpdatedAt: journal.UpdatedAt,
	}, LockWarning{}, nil
}

func canClearLockForStatus(status TransactionStatus) bool {
	switch status {
	case TransactionStatusCommitted,
		TransactionStatusRolledBack,
		TransactionStatusFailed,
		TransactionStatusRollbackFailed:
		return true
	default:
		return false
	}
}
