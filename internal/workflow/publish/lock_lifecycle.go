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
	// LockShowStatusFailed means inspection failed for reasons other than lock corruption.
	LockShowStatusFailed LockShowStatus = "failed"
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

// LockClearReason is a stable machine-readable lock clear outcome code.
type LockClearReason string

const (
	LockClearReasonCleared              LockClearReason = "cleared"
	LockClearReasonMissingTransactionID LockClearReason = "missing_transaction_id"
	LockClearReasonMissingConfirmation  LockClearReason = "missing_confirmation"
	LockClearReasonConfirmationMismatch LockClearReason = "confirmation_mismatch"
	LockClearReasonInvalidTransactionID LockClearReason = "invalid_transaction_id"
	LockClearReasonStateDirMissing      LockClearReason = "state_dir_missing"
	LockClearReasonLockAbsent           LockClearReason = "lock_absent"
	LockClearReasonLockCorrupt          LockClearReason = "lock_corrupt"
	LockClearReasonLockReadFailed       LockClearReason = "lock_read_failed"
	LockClearReasonTransactionMismatch  LockClearReason = "transaction_mismatch"
	LockClearReasonJournalCorrupt       LockClearReason = "journal_corrupt"
	LockClearReasonJournalReadFailed    LockClearReason = "journal_read_failed"
	LockClearReasonActiveTransaction    LockClearReason = "active_transaction"
	LockClearReasonLockChanged          LockClearReason = "lock_changed"
	LockClearReasonDeleteFailed         LockClearReason = "delete_failed"
	LockClearReasonSyncFailed           LockClearReason = "sync_failed"
)

// LockWarningCode is a stable machine-readable lock warning code.
type LockWarningCode string

const (
	LockWarningJournalMissing LockWarningCode = "journal_missing"
	LockWarningJournalCorrupt LockWarningCode = "journal_corrupt"
	LockWarningLockCorrupt    LockWarningCode = "lock_corrupt"
)

// LockPostClearState tells operators whether clearing publish.lock removed the
// publish blocker or only exposed remaining journal recovery state.
type LockPostClearState string

const (
	LockPostClearReadyForPublish      LockPostClearState = "ready_for_publish"
	LockPostClearTransactionStillOpen LockPostClearState = "transaction_still_blocks_publish"
	LockPostClearUnverifiedNoJournal  LockPostClearState = "unverified_no_journal"
)

// LockWarning describes recoverable lock inspection information.
type LockWarning struct {
	Code    LockWarningCode
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
	Status         LockClearStatus
	TransactionID  TransactionID
	Lock           TransactionLockInfo
	LockCleared    bool
	Journal        LockJournalState
	Reason         LockClearReason
	Message        string
	PostClearState LockPostClearState
	Warnings       []LockWarning
}

// ShowTransactionLock inspects publish.lock and its referenced journal.
func (s Service) ShowTransactionLock(ctx context.Context, stateDir string) (LockShowResult, error) {
	if err := ctx.Err(); err != nil {
		return LockShowResult{Status: LockShowStatusFailed}, err
	}
	lock, ok, err := currentTransactionLock(stateDir)
	if err != nil {
		status := LockShowStatusFailed
		warnings := []LockWarning{}
		if errors.Is(err, errTransactionLockCorrupt) {
			status = LockShowStatusCorrupt
			warnings = append(warnings, LockWarning{Code: LockWarningLockCorrupt, Message: "publish lock is not parseable"})
		}
		return LockShowResult{Status: status, Warnings: warnings}, &Error{Code: CodeLockFailed, Message: "read publish transaction lock failed", Cause: err}
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
		if errors.Is(err, errTransactionJournalCorrupt) {
			result.Status = LockShowStatusJournalCorrupt
			return result, &Error{Code: CodeRecoveryFailed, Message: "load publish transaction journal failed", Cause: err}
		}
		result.Status = LockShowStatusFailed
		return result, &Error{Code: CodeLockFailed, Message: "inspect publish transaction lock failed", Cause: err}
	}
	if !journal.Present {
		result.Status = LockShowStatusJournalMissing
	}
	return result, nil
}

// ClearTransactionLock removes only publish.lock after explicit confirmation.
func (s Service) ClearTransactionLock(ctx context.Context, stateDir string, opts LockClearOptions) (LockClearResult, error) {
	result := LockClearResult{Status: LockClearStatusRefused, TransactionID: opts.TransactionID}
	if err := ctx.Err(); err != nil {
		result.Status = LockClearStatusFailed
		return result, err
	}
	if opts.TransactionID == "" {
		result.Reason = LockClearReasonMissingTransactionID
		result.Message = "transaction id is required"
		return result, &Error{Code: CodeLockFailed, Message: result.Message}
	}
	if opts.Confirm == "" {
		result.Reason = LockClearReasonMissingConfirmation
		result.Message = "confirmation transaction id is required"
		return result, &Error{Code: CodeLockFailed, Message: result.Message}
	}
	if opts.Confirm != opts.TransactionID {
		result.Reason = LockClearReasonConfirmationMismatch
		result.Message = "confirmation transaction id does not match"
		return result, &Error{Code: CodeLockFailed, Message: result.Message}
	}
	if err := validateTransactionID(opts.TransactionID); err != nil {
		result.Reason = LockClearReasonInvalidTransactionID
		result.Message = "invalid transaction id"
		return result, &Error{Code: CodeLockFailed, Message: result.Message, Cause: err}
	}

	path, err := transactionLockPath(stateDir)
	if err != nil {
		result.Status = LockClearStatusFailed
		result.Reason = LockClearReasonStateDirMissing
		result.Message = "resolve publish lock path failed"
		return result, &Error{Code: CodeLockFailed, Message: result.Message, Cause: err}
	}
	lock, err := readTransactionLock(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			result.Reason = LockClearReasonLockAbsent
			result.Message = "publish lock is absent"
			return result, &Error{Code: CodeLockFailed, Message: result.Message}
		}
		result.Status = LockClearStatusFailed
		if errors.Is(err, errTransactionLockCorrupt) {
			result.Reason = LockClearReasonLockCorrupt
			result.Message = "publish lock is not parseable"
			result.Warnings = append(result.Warnings, LockWarning{Code: LockWarningLockCorrupt, Message: result.Message})
		} else {
			result.Reason = LockClearReasonLockReadFailed
			result.Message = "read publish transaction lock failed"
		}
		return result, &Error{Code: CodeLockFailed, Message: "read publish transaction lock failed", Cause: err}
	}
	result.Lock = lock
	if lock.ID != opts.TransactionID {
		result.Reason = LockClearReasonTransactionMismatch
		result.Message = "lock transaction does not match requested transaction"
		return result, &Error{Code: CodeLockFailed, Message: fmt.Sprintf("publish lock belongs to %s, not %s", lock.ID, opts.TransactionID)}
	}

	journal, warning, err := loadLockJournal(ctx, stateDir, lock.ID)
	result.Journal = journal
	if warning.Code != "" {
		result.Warnings = append(result.Warnings, warning)
	}
	if err != nil {
		result.Status = LockClearStatusFailed
		if errors.Is(err, errTransactionJournalCorrupt) {
			result.Reason = LockClearReasonJournalCorrupt
			result.Message = "referenced transaction journal is corrupt"
			return result, &Error{Code: CodeRecoveryFailed, Message: "load publish transaction journal failed", Cause: err}
		}
		result.Reason = LockClearReasonJournalReadFailed
		result.Message = "read referenced transaction journal failed"
		return result, &Error{Code: CodeLockFailed, Message: "read referenced transaction journal failed", Cause: err}
	}
	postState := postClearState(journal)
	if journal.Present && !journal.Status.AllowsLockClear() {
		result.Reason = LockClearReasonActiveTransaction
		result.Message = "referenced transaction is active"
		return result, &Error{Code: CodeLockFailed, Message: fmt.Sprintf("transaction %s is %s; refusing to clear active publish lock", lock.ID, journal.Status)}
	}

	if err := removeTransactionLockIfCurrent(path, lock.ID); err != nil {
		result.Status = LockClearStatusFailed
		result.Reason = lockClearRemoveReason(err)
		result.Message = lockClearRemoveMessage(result.Reason)
		return result, &Error{Code: CodeLockFailed, Message: result.Message, Cause: err}
	}

	result.Status = LockClearStatusCleared
	result.LockCleared = true
	result.Reason = LockClearReasonCleared
	result.Message = "publish lock cleared"
	result.PostClearState = postState
	return result, nil
}

func loadLockJournal(ctx context.Context, stateDir string, id TransactionID) (LockJournalState, LockWarning, error) {
	journal, err := NewFileJournalStore(stateDir).Load(ctx, id)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return LockJournalState{}, LockWarning{
				Code:    LockWarningJournalMissing,
				Message: "publish lock references a transaction journal that does not exist",
			}, nil
		}
		if errors.Is(err, errTransactionJournalCorrupt) {
			return LockJournalState{}, LockWarning{
				Code:    LockWarningJournalCorrupt,
				Message: "publish lock references a transaction journal that cannot be read safely",
			}, err
		}
		return LockJournalState{}, LockWarning{}, err
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

func postClearState(journal LockJournalState) LockPostClearState {
	if !journal.Present {
		return LockPostClearUnverifiedNoJournal
	}
	if journal.Status.BlocksNewPublish() {
		return LockPostClearTransactionStillOpen
	}
	return LockPostClearReadyForPublish
}

func lockClearRemoveReason(err error) LockClearReason {
	switch {
	case errors.Is(err, errTransactionLockChanged), errors.Is(err, os.ErrNotExist), errors.Is(err, errTransactionLockCorrupt):
		return LockClearReasonLockChanged
	case errors.Is(err, errTransactionLockDeleteFailed):
		return LockClearReasonDeleteFailed
	case errors.Is(err, errTransactionLockSyncFailed):
		return LockClearReasonSyncFailed
	default:
		return LockClearReasonDeleteFailed
	}
}

func lockClearRemoveMessage(reason LockClearReason) string {
	switch reason {
	case LockClearReasonLockChanged:
		return "publish lock changed before deletion"
	case LockClearReasonSyncFailed:
		return "sync publish lock directory failed"
	case LockClearReasonLockCorrupt:
		return "publish lock is not parseable"
	default:
		return "delete publish lock failed"
	}
}
