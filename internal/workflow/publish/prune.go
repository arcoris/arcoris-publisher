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

// PruneStatus describes whether a prune request only previewed or deleted
// terminal transaction journals.
type PruneStatus string

const (
	// PruneStatusDryRun means matching journals were reported but not removed.
	PruneStatusDryRun PruneStatus = "dry_run"
	// PruneStatusCompleted means all selected terminal journals were removed.
	PruneStatusCompleted PruneStatus = "completed"
	// PruneStatusFailed means at least one selected journal could not be removed.
	PruneStatusFailed PruneStatus = "failed"
)

// PruneOptions selects terminal transaction journals for safe deletion.
type PruneOptions struct {
	Statuses  []TransactionStatus
	OlderThan time.Duration
	Now       time.Time
	DryRun    bool
}

// PruneResult reports every transaction considered by a prune request.
type PruneResult struct {
	Status   PruneStatus
	Matched  []PruneEntry
	Deleted  []PruneEntry
	Skipped  []PruneEntry
	Warnings []string
}

// PruneEntry describes one transaction considered by prune. Path is local
// recovery state and should be hidden from default user reports.
type PruneEntry struct {
	ID        TransactionID
	Status    TransactionStatus
	Rollback  RollbackStatus
	Version   string
	StartedAt time.Time
	UpdatedAt time.Time
	Path      string
	Reason    string
}

// PruneTransactions deletes only selected terminal transaction journals. A
// live publish lock blocks non-dry-run prune because prune mutates transaction
// state and must not race an owner that still holds publish.lock.
func (s Service) PruneTransactions(ctx context.Context, stateDir string, opts PruneOptions) (result PruneResult, err error) {
	if !opts.DryRun {
		var operationLock operationLock
		operationLock, err = acquireOperationLock(ctx, stateDir, operationLockPrune, s.operationLockOps)
		if err != nil {
			return PruneResult{Status: PruneStatusFailed}, operationLockAcquireError(operationLockPrune, err)
		}
		defer func() {
			result, err = releaseOperationLockForPrune(operationLock, result, err)
		}()

		if lock, ok, err := currentTransactionLock(stateDir); err != nil {
			return PruneResult{Status: PruneStatusFailed}, &Error{Code: CodeLockFailed, Message: "read publish transaction lock failed", Cause: err}
		} else if ok {
			return PruneResult{Status: PruneStatusFailed}, &Error{
				Code:    CodeLockFailed,
				Message: fmt.Sprintf("publish transaction lock exists for %s; refusing to prune transaction state", lock.ID),
			}
		}
	}
	result, err = NewFileJournalStore(stateDir).Prune(ctx, opts)
	if err != nil {
		return result, &Error{Code: CodePruneFailed, Message: "prune publish transactions failed", Cause: err}
	}
	return result, nil
}

func releaseOperationLockForPrune(lock operationLock, result PruneResult, err error) (PruneResult, error) {
	outcome, releaseErr := lock.Release()
	if releaseErr == nil {
		return result, err
	}
	releaseFailure := operationLockReleaseError(outcome, releaseErr)
	result.Warnings = append(result.Warnings, operationLockReleaseWarning(outcome, releaseErr))
	if err != nil {
		return result, errors.Join(err, releaseFailure)
	}
	result.Status = PruneStatusFailed
	return result, releaseFailure
}

// Prune safely removes selected terminal transaction journals from disk.
func (s FileJournalStore) Prune(ctx context.Context, opts PruneOptions) (PruneResult, error) {
	result := PruneResult{Status: PruneStatusCompleted}
	if opts.DryRun {
		result.Status = PruneStatusDryRun
	}
	if !opts.DryRun && len(opts.Statuses) == 0 && opts.OlderThan <= 0 {
		result.Status = PruneStatusFailed
		return result, fmt.Errorf("refusing to prune without an explicit status or age filter")
	}
	statuses, err := prunableStatusSet(opts.Statuses)
	if err != nil {
		result.Status = PruneStatusFailed
		return result, err
	}
	now := opts.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}

	summaries, err := s.List(ctx)
	if err != nil {
		result.Status = PruneStatusFailed
		return result, err
	}
	for _, summary := range summaries {
		if err := ctx.Err(); err != nil {
			result.Status = PruneStatusFailed
			return result, err
		}
		entry, selected, err := s.pruneCandidate(ctx, summary, statuses, opts.OlderThan, now)
		if err != nil {
			result.Status = PruneStatusFailed
			return result, err
		}
		if !selected {
			result.Skipped = append(result.Skipped, entry)
			continue
		}

		result.Matched = append(result.Matched, entry)
		if opts.DryRun {
			continue
		}
		if err := os.Remove(entry.Path); err != nil {
			entry.Reason = "delete failed: " + err.Error()
			result.Skipped = append(result.Skipped, entry)
			result.Status = PruneStatusFailed
			return result, err
		}
		if err := syncParentDir(entry.Path); err != nil {
			entry.Reason = "delete parent sync failed: " + err.Error()
			result.Skipped = append(result.Skipped, entry)
			result.Status = PruneStatusFailed
			return result, err
		}
		entry.Reason = "deleted " + entry.Reason
		result.Deleted = append(result.Deleted, entry)
	}
	return result, nil
}

func (s FileJournalStore) pruneCandidate(
	ctx context.Context,
	summary TransactionSummary,
	statuses map[TransactionStatus]bool,
	olderThan time.Duration,
	now time.Time,
) (PruneEntry, bool, error) {
	journal, err := s.Load(ctx, summary.ID)
	if err != nil {
		return PruneEntry{}, false, err
	}
	path, err := s.journalPath(journal.ID)
	if err != nil {
		return PruneEntry{}, false, err
	}
	entry := PruneEntry{
		ID:        journal.ID,
		Status:    journal.Status,
		Rollback:  journal.Rollback,
		Version:   journal.Version,
		StartedAt: journal.StartedAt,
		UpdatedAt: journal.UpdatedAt,
		Path:      path,
	}
	if !journal.Status.Prunable() {
		entry.Reason = "status is not prunable"
		return entry, false, nil
	}
	if !statuses[journal.Status] {
		entry.Reason = "status does not match filter"
		return entry, false, nil
	}
	if olderThan > 0 {
		when := journal.UpdatedAt
		if when.IsZero() {
			when = journal.StartedAt
		}
		if when.IsZero() {
			entry.Reason = "transaction timestamp is missing"
			return entry, false, nil
		}
		if when.After(now.Add(-olderThan)) {
			entry.Reason = fmt.Sprintf("newer than %s", olderThan)
			return entry, false, nil
		}
		entry.Reason = fmt.Sprintf("status %s and older than %s", journal.Status, olderThan)
		return entry, true, nil
	}
	entry.Reason = fmt.Sprintf("status %s", journal.Status)
	return entry, true, nil
}

func prunableStatusSet(statuses []TransactionStatus) (map[TransactionStatus]bool, error) {
	out := map[TransactionStatus]bool{}
	if len(statuses) == 0 {
		out[TransactionStatusCommitted] = true
		out[TransactionStatusRolledBack] = true
		return out, nil
	}
	for _, status := range statuses {
		if !status.Prunable() {
			return nil, fmt.Errorf("transaction status %q is not prunable", status)
		}
		out[status] = true
	}
	return out, nil
}
