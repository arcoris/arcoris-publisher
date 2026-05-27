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

package report

import (
	"io"
	"time"

	"arcoris.dev/arcoris-publisher/internal/workflow/publish"
)

// TransactionListReport is the stable report DTO for publish transaction lists.
type TransactionListReport struct {
	Kind         string                     `json:"kind"`
	Count        int                        `json:"count"`
	Transactions []TransactionSummaryReport `json:"transactions"`
}

// TransactionSummaryReport describes one durable publish transaction.
type TransactionSummaryReport struct {
	ID             string `json:"id"`
	Status         string `json:"status"`
	RollbackStatus string `json:"rollbackStatus,omitempty"`
	Version        string `json:"version,omitempty"`
	StartedAt      string `json:"startedAt,omitempty"`
	UpdatedAt      string `json:"updatedAt,omitempty"`
}

// TransactionReport is the stable report DTO for one publish transaction.
type TransactionReport struct {
	Kind                  string                       `json:"kind,omitempty"`
	ID                    string                       `json:"id"`
	Status                string                       `json:"status"`
	RollbackStatus        string                       `json:"rollbackStatus,omitempty"`
	Version               string                       `json:"version,omitempty"`
	StartedAt             string                       `json:"startedAt,omitempty"`
	UpdatedAt             string                       `json:"updatedAt,omitempty"`
	Modules               []TransactionModuleReport    `json:"modules"`
	Warnings              []string                     `json:"warnings,omitempty"`
	Failure               string                       `json:"failure,omitempty"`
	ManualRecoveryActions []ManualRecoveryActionReport `json:"manualRecoveryActions,omitempty"`
}

// TransactionPruneReport is the stable report DTO for transaction pruning.
type TransactionPruneReport struct {
	Kind         string                  `json:"kind"`
	Status       string                  `json:"status"`
	MatchedCount int                     `json:"matchedCount"`
	DeletedCount int                     `json:"deletedCount"`
	SkippedCount int                     `json:"skippedCount"`
	Matched      []TransactionPruneEntry `json:"matched,omitempty"`
	Deleted      []TransactionPruneEntry `json:"deleted,omitempty"`
	Skipped      []TransactionPruneEntry `json:"skipped,omitempty"`
	Warnings     []string                `json:"warnings,omitempty"`
}

// TransactionPruneEntry describes one transaction considered by prune.
type TransactionPruneEntry struct {
	ID             string `json:"id"`
	Status         string `json:"status"`
	RollbackStatus string `json:"rollbackStatus,omitempty"`
	Version        string `json:"version,omitempty"`
	StartedAt      string `json:"startedAt,omitempty"`
	UpdatedAt      string `json:"updatedAt,omitempty"`
	Path           string `json:"path,omitempty"`
	Reason         string `json:"reason,omitempty"`
}

// TransactionLockReport is the stable DTO for publish lock inspection.
type TransactionLockReport struct {
	Kind     string                   `json:"kind"`
	Status   string                   `json:"status"`
	Reason   string                   `json:"reason,omitempty"`
	Message  string                   `json:"message,omitempty"`
	Lock     *TransactionLockInfo     `json:"lock"`
	Journal  *TransactionLockJournal  `json:"journal"`
	Warnings []TransactionLockWarning `json:"warnings"`
}

// TransactionLockClearReport is the stable DTO for guarded lock clearing.
type TransactionLockClearReport struct {
	Kind           string                   `json:"kind"`
	Status         string                   `json:"status"`
	TransactionID  string                   `json:"transactionId,omitempty"`
	Reason         string                   `json:"reason,omitempty"`
	Message        string                   `json:"message,omitempty"`
	PostClearState string                   `json:"postClearState,omitempty"`
	Lock           TransactionLockClearInfo `json:"lock"`
	Journal        *TransactionLockJournal  `json:"journal"`
	Warnings       []TransactionLockWarning `json:"warnings"`
}

// TransactionDiagnosticsReport is the stable DTO for transaction state diagnostics.
type TransactionDiagnosticsReport struct {
	Kind           string                                `json:"kind"`
	PublishBlocked bool                                  `json:"publishBlocked"`
	Blockers       []TransactionDiagnosticsBlockerReport `json:"blockers"`
	Lock           TransactionDiagnosticsLockReport      `json:"lock"`
	OperationLock  TransactionOperationLockReport        `json:"operationLock"`
	Journals       []TransactionJournalDiagnosticReport  `json:"journals"`
	Warnings       []TransactionDiagnosticWarningReport  `json:"warnings"`
}

// TransactionDiagnosticsBlockerReport describes one publish blocker.
type TransactionDiagnosticsBlockerReport struct {
	Kind          string `json:"kind"`
	TransactionID string `json:"transactionId,omitempty"`
	Status        string `json:"status,omitempty"`
	Reason        string `json:"reason,omitempty"`
	Name          string `json:"name,omitempty"`
}

// TransactionDiagnosticsLockReport describes publish.lock in diagnostics.
type TransactionDiagnosticsLockReport struct {
	Status   string                   `json:"status"`
	Reason   string                   `json:"reason,omitempty"`
	Message  string                   `json:"message,omitempty"`
	Lock     *TransactionLockInfo     `json:"lock"`
	Journal  *TransactionLockJournal  `json:"journal"`
	Warnings []TransactionLockWarning `json:"warnings"`
}

// TransactionOperationLockReport describes operation.lock in diagnostics.
type TransactionOperationLockReport struct {
	Present    bool   `json:"present"`
	Operation  string `json:"operation,omitempty"`
	PID        string `json:"pid,omitempty"`
	StartedAt  string `json:"startedAt,omitempty"`
	Corrupt    bool   `json:"corrupt,omitempty"`
	ReadFailed bool   `json:"readFailed,omitempty"`
	Message    string `json:"message,omitempty"`
	Path       string `json:"path,omitempty"`
}

// TransactionJournalDiagnosticReport describes one journal diagnostic item.
type TransactionJournalDiagnosticReport struct {
	TransactionID    string `json:"transactionId,omitempty"`
	Name             string `json:"name,omitempty"`
	Status           string `json:"status,omitempty"`
	RollbackStatus   string `json:"rollbackStatus,omitempty"`
	Version          string `json:"version,omitempty"`
	StartedAt        string `json:"startedAt,omitempty"`
	UpdatedAt        string `json:"updatedAt,omitempty"`
	Prunable         bool   `json:"prunable"`
	BlocksNewPublish bool   `json:"blocksNewPublish"`
	AllowsLockClear  bool   `json:"allowsLockClear"`
	Corrupt          bool   `json:"corrupt,omitempty"`
	ReadFailed       bool   `json:"readFailed,omitempty"`
	Message          string `json:"message,omitempty"`
	Path             string `json:"path,omitempty"`
}

// TransactionDiagnosticWarningReport describes diagnostic warnings.
type TransactionDiagnosticWarningReport struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// TransactionLockInfo describes publish.lock metadata.
type TransactionLockInfo struct {
	TransactionID string `json:"transactionId,omitempty"`
	PID           string `json:"pid,omitempty"`
	Command       string `json:"command,omitempty"`
	StartedAt     string `json:"startedAt,omitempty"`
	Path          string `json:"path,omitempty"`
}

// TransactionLockJournal describes the journal referenced by publish.lock.
type TransactionLockJournal struct {
	Present        bool   `json:"present"`
	Status         string `json:"status,omitempty"`
	RollbackStatus string `json:"rollbackStatus,omitempty"`
	Version        string `json:"version,omitempty"`
	StartedAt      string `json:"startedAt,omitempty"`
	UpdatedAt      string `json:"updatedAt,omitempty"`
}

// TransactionLockClearInfo describes the lock file outcome.
type TransactionLockClearInfo struct {
	Cleared bool   `json:"cleared"`
	Path    string `json:"path,omitempty"`
}

// TransactionLockWarning describes lock inspection warnings.
type TransactionLockWarning struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// TransactionModuleReport describes one module's transaction refs and state.
type TransactionModuleReport struct {
	Module              string `json:"module"`
	Repository          string `json:"repository"`
	WorktreeDir         string `json:"worktreeDir,omitempty"`
	TargetBranch        string `json:"targetBranch,omitempty"`
	FinalBranchRef      string `json:"finalBranchRef,omitempty"`
	FinalTagRef         string `json:"finalTagRef,omitempty"`
	CreatedCommit       string `json:"createdCommit,omitempty"`
	CandidateBranchRef  string `json:"candidateBranchRef,omitempty"`
	Skipped             bool   `json:"skipped"`
	CandidatePushed     bool   `json:"candidatePushed"`
	FinalBranchPromoted bool   `json:"finalBranchPromoted"`
	RemoteTagPushed     bool   `json:"remoteTagPushed"`
}

// ManualRecoveryActionReport describes an operator recovery action.
type ManualRecoveryActionReport struct {
	Module       string `json:"module,omitempty"`
	Repository   string `json:"repository,omitempty"`
	Ref          string `json:"ref,omitempty"`
	ExpectedHash string `json:"expectedHash,omitempty"`
	DesiredHash  string `json:"desiredHash,omitempty"`
	Message      string `json:"message"`
	Command      string `json:"command,omitempty"`
}

// BuildTransactionListReport converts transaction summaries to a DTO.
func BuildTransactionListReport(summaries []publish.TransactionSummary) TransactionListReport {
	out := TransactionListReport{
		Kind:         "transactions",
		Count:        len(summaries),
		Transactions: make([]TransactionSummaryReport, 0, len(summaries)),
	}
	for _, summary := range summaries {
		out.Transactions = append(out.Transactions, TransactionSummaryReport{
			ID:             summary.ID.String(),
			Status:         string(summary.Status),
			RollbackStatus: string(summary.Rollback),
			Version:        summary.Version,
			StartedAt:      summary.StartedAt.Format(timeFormat),
			UpdatedAt:      summary.UpdatedAt.Format(timeFormat),
		})
	}
	return out
}

// BuildTransactionReport converts a transaction journal to a path-safe DTO.
func BuildTransactionReport(journal publish.TransactionJournal, opts Options) TransactionReport {
	out := TransactionReport{
		Kind:           "transaction",
		ID:             journal.ID.String(),
		Status:         string(journal.Status),
		RollbackStatus: string(journal.Rollback),
		Version:        journal.Version,
		StartedAt:      journal.StartedAt.Format(timeFormat),
		UpdatedAt:      journal.UpdatedAt.Format(timeFormat),
		Warnings:       append([]string(nil), journal.Warnings...),
		Failure:        journal.Failure,
		Modules:        make([]TransactionModuleReport, 0, len(journal.Modules)),
	}
	for _, mod := range journal.Modules {
		out.Modules = append(out.Modules, TransactionModuleReport{
			Module:              mod.Module.String(),
			Repository:          mod.Repository.String(),
			WorktreeDir:         includePath(mod.WorktreeDir, opts),
			TargetBranch:        mod.TargetBranch.String(),
			FinalBranchRef:      mod.FinalBranchRef,
			FinalTagRef:         mod.FinalTagRef,
			CreatedCommit:       mod.CreatedCommit.String(),
			CandidateBranchRef:  mod.CandidateBranchRef,
			Skipped:             mod.Skipped,
			CandidatePushed:     mod.CandidatePushed,
			FinalBranchPromoted: mod.FinalBranchPromoted,
			RemoteTagPushed:     mod.RemoteTagPushed,
		})
	}
	for _, action := range journal.ManualActions {
		out.ManualRecoveryActions = append(out.ManualRecoveryActions, ManualRecoveryActionReport{
			Module:       action.Module.String(),
			Repository:   action.Repository.String(),
			Ref:          action.Ref,
			ExpectedHash: action.ExpectedHash.String(),
			DesiredHash:  action.DesiredHash.String(),
			Message:      action.Message,
			Command:      action.Command,
		})
	}
	return out
}

// BuildTransactionPruneReport converts prune results to a path-safe DTO.
func BuildTransactionPruneReport(result publish.PruneResult, opts Options) TransactionPruneReport {
	return TransactionPruneReport{
		Kind:         "transactions-prune",
		Status:       string(result.Status),
		MatchedCount: len(result.Matched),
		DeletedCount: len(result.Deleted),
		SkippedCount: len(result.Skipped),
		Matched:      buildTransactionPruneEntries(result.Matched, opts),
		Deleted:      buildTransactionPruneEntries(result.Deleted, opts),
		Skipped:      buildTransactionPruneEntries(result.Skipped, opts),
		Warnings:     append([]string(nil), result.Warnings...),
	}
}

func buildTransactionPruneEntries(entries []publish.PruneEntry, opts Options) []TransactionPruneEntry {
	out := make([]TransactionPruneEntry, 0, len(entries))
	for _, entry := range entries {
		out = append(out, TransactionPruneEntry{
			ID:             entry.ID.String(),
			Status:         string(entry.Status),
			RollbackStatus: string(entry.Rollback),
			Version:        entry.Version,
			StartedAt:      entry.StartedAt.Format(timeFormat),
			UpdatedAt:      entry.UpdatedAt.Format(timeFormat),
			Path:           includePath(entry.Path, opts),
			Reason:         entry.Reason,
		})
	}
	return out
}

// BuildTransactionLockReport converts lock inspection to a path-safe DTO.
func BuildTransactionLockReport(result publish.LockShowResult, opts Options) TransactionLockReport {
	var journal *TransactionLockJournal
	if result.Status != publish.LockShowStatusAbsent {
		journal = buildTransactionLockJournal(result.Journal)
	}
	return TransactionLockReport{
		Kind:     "transactions-lock",
		Status:   string(result.Status),
		Reason:   string(result.Reason),
		Message:  result.Message,
		Lock:     buildTransactionLockInfo(result.Lock, opts),
		Journal:  journal,
		Warnings: buildTransactionLockWarnings(result.Warnings),
	}
}

// BuildTransactionLockClearReport converts lock clear results to a path-safe DTO.
func BuildTransactionLockClearReport(result publish.LockClearResult, opts Options) TransactionLockClearReport {
	return TransactionLockClearReport{
		Kind:           "transactions-lock-clear",
		Status:         string(result.Status),
		TransactionID:  result.TransactionID.String(),
		Reason:         string(result.Reason),
		Message:        result.Message,
		PostClearState: string(result.PostClearState),
		Lock: TransactionLockClearInfo{
			Cleared: result.LockCleared,
			Path:    includePath(result.Lock.Path, opts),
		},
		Journal:  buildTransactionLockJournal(result.Journal),
		Warnings: buildTransactionLockWarnings(result.Warnings),
	}
}

// BuildTransactionDiagnosticsReport converts transaction diagnostics to a path-safe DTO.
func BuildTransactionDiagnosticsReport(result publish.TransactionStateDiagnostics, opts Options) TransactionDiagnosticsReport {
	return TransactionDiagnosticsReport{
		Kind:           "transactions-diagnostics",
		PublishBlocked: result.PublishBlocked,
		Blockers:       buildTransactionDiagnosticBlockers(result.Blockers),
		Lock:           buildTransactionDiagnosticsLock(result.Lock, opts),
		OperationLock:  buildTransactionOperationLock(result.OperationLock, opts),
		Journals:       buildTransactionJournalDiagnostics(result.Journals, opts),
		Warnings:       buildTransactionDiagnosticWarnings(result.Warnings),
	}
}

func buildTransactionLockInfo(info publish.TransactionLockInfo, opts Options) *TransactionLockInfo {
	if info.ID == "" {
		return nil
	}
	return &TransactionLockInfo{
		TransactionID: info.ID.String(),
		PID:           info.PID,
		Command:       info.Command,
		StartedAt:     info.StartedAt,
		Path:          includePath(info.Path, opts),
	}
}

func buildTransactionDiagnosticBlockers(blockers []publish.TransactionStateBlocker) []TransactionDiagnosticsBlockerReport {
	out := make([]TransactionDiagnosticsBlockerReport, 0, len(blockers))
	for _, blocker := range blockers {
		out = append(out, TransactionDiagnosticsBlockerReport{
			Kind:          string(blocker.Kind),
			TransactionID: blocker.TransactionID.String(),
			Status:        string(blocker.Status),
			Reason:        string(blocker.Reason),
			Name:          blocker.Name,
		})
	}
	return out
}

func buildTransactionDiagnosticsLock(result publish.LockShowResult, opts Options) TransactionDiagnosticsLockReport {
	var journal *TransactionLockJournal
	if result.Status != publish.LockShowStatusAbsent {
		journal = buildTransactionLockJournal(result.Journal)
	}
	return TransactionDiagnosticsLockReport{
		Status:   string(result.Status),
		Reason:   string(result.Reason),
		Message:  result.Message,
		Lock:     buildTransactionLockInfo(result.Lock, opts),
		Journal:  journal,
		Warnings: buildTransactionLockWarnings(result.Warnings),
	}
}

func buildTransactionOperationLock(lock publish.OperationLockDiagnostic, opts Options) TransactionOperationLockReport {
	return TransactionOperationLockReport{
		Present:    lock.Present,
		Operation:  lock.Operation,
		PID:        lock.PID,
		StartedAt:  lock.StartedAt,
		Corrupt:    lock.Corrupt,
		ReadFailed: lock.ReadFailed,
		Message:    lock.Message,
		Path:       includePath(lock.Path, opts),
	}
}

func buildTransactionJournalDiagnostics(journals []publish.JournalDiagnostic, opts Options) []TransactionJournalDiagnosticReport {
	out := make([]TransactionJournalDiagnosticReport, 0, len(journals))
	for _, journal := range journals {
		out = append(out, TransactionJournalDiagnosticReport{
			TransactionID:    journal.ID.String(),
			Name:             journal.Name,
			Status:           string(journal.Status),
			RollbackStatus:   string(journal.Rollback),
			Version:          journal.Version,
			StartedAt:        formatOptionalTime(journal.StartedAt),
			UpdatedAt:        formatOptionalTime(journal.UpdatedAt),
			Prunable:         journal.Prunable,
			BlocksNewPublish: journal.BlocksNewPublish,
			AllowsLockClear:  journal.AllowsLockClear,
			Corrupt:          journal.Corrupt,
			ReadFailed:       journal.ReadFailed,
			Message:          journal.Message,
			Path:             includePath(journal.Path, opts),
		})
	}
	return out
}

func formatOptionalTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(timeFormat)
}

func buildTransactionDiagnosticWarnings(warnings []publish.TransactionDiagnosticWarning) []TransactionDiagnosticWarningReport {
	out := make([]TransactionDiagnosticWarningReport, 0, len(warnings))
	for _, warning := range warnings {
		out = append(out, TransactionDiagnosticWarningReport{Code: warning.Code, Message: warning.Message})
	}
	return out
}

func buildTransactionLockJournal(journal publish.LockJournalState) *TransactionLockJournal {
	if !journal.Present && journal.Status == "" {
		return &TransactionLockJournal{Present: false}
	}
	return &TransactionLockJournal{
		Present:        journal.Present,
		Status:         string(journal.Status),
		RollbackStatus: string(journal.Rollback),
		Version:        journal.Version,
		StartedAt:      journal.StartedAt.Format(timeFormat),
		UpdatedAt:      journal.UpdatedAt.Format(timeFormat),
	}
}

func buildTransactionLockWarnings(warnings []publish.LockWarning) []TransactionLockWarning {
	out := make([]TransactionLockWarning, 0, len(warnings))
	for _, warning := range warnings {
		out = append(out, TransactionLockWarning{Code: string(warning.Code), Message: warning.Message})
	}
	return out
}

func writeTransactionListText(w io.Writer, report TransactionListReport) error {
	if err := writeLine(w, "Transactions"); err != nil {
		return err
	}
	if err := writeLine(w, "  Count: %d", report.Count); err != nil {
		return err
	}
	for _, tx := range report.Transactions {
		status := tx.Status
		if tx.RollbackStatus != "" {
			status += " rollback=" + tx.RollbackStatus
		}
		if err := writeLine(w, "  %s: %s", tx.ID, status); err != nil {
			return err
		}
	}
	return nil
}

func writeTransactionPruneText(w io.Writer, report TransactionPruneReport) error {
	if err := writeLine(w, "Transactions prune"); err != nil {
		return err
	}
	if err := writeLine(w, "  Status: %s", report.Status); err != nil {
		return err
	}
	if err := writeLine(w, "  Matched: %d", report.MatchedCount); err != nil {
		return err
	}
	if err := writeLine(w, "  Deleted: %d", report.DeletedCount); err != nil {
		return err
	}
	if err := writeLine(w, "  Skipped: %d", report.SkippedCount); err != nil {
		return err
	}
	if len(report.Matched) > 0 {
		if err := writeLine(w, "  Matched:"); err != nil {
			return err
		}
		for _, entry := range report.Matched {
			if err := writeLine(w, "    %s: %s %s", entry.ID, entry.Status, entry.Reason); err != nil {
				return err
			}
		}
	}
	if len(report.Skipped) > 0 {
		if err := writeLine(w, "  Skipped:"); err != nil {
			return err
		}
		for _, entry := range report.Skipped {
			if err := writeLine(w, "    %s: %s %s", entry.ID, entry.Status, entry.Reason); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeTransactionLockText(w io.Writer, report TransactionLockReport) error {
	if err := writeLine(w, "Transaction lock"); err != nil {
		return err
	}
	if err := writeLine(w, "  Status: %s", report.Status); err != nil {
		return err
	}
	if report.Reason != "" {
		if err := writeLine(w, "  Reason: %s", report.Reason); err != nil {
			return err
		}
	}
	if report.Message != "" {
		if err := writeLine(w, "  Message: %s", report.Message); err != nil {
			return err
		}
	}
	if report.Lock != nil {
		if err := writeLine(w, "  Transaction: %s", report.Lock.TransactionID); err != nil {
			return err
		}
		if report.Lock.Command != "" {
			if err := writeLine(w, "  Command: %s", report.Lock.Command); err != nil {
				return err
			}
		}
		if report.Lock.PID != "" {
			if err := writeLine(w, "  PID: %s", report.Lock.PID); err != nil {
				return err
			}
		}
		if report.Lock.StartedAt != "" {
			if err := writeLine(w, "  Started: %s", report.Lock.StartedAt); err != nil {
				return err
			}
		}
	}
	if report.Journal != nil && report.Journal.Present {
		if err := writeLine(w, "  Journal:"); err != nil {
			return err
		}
		if err := writeLine(w, "    Status: %s", report.Journal.Status); err != nil {
			return err
		}
		if report.Journal.RollbackStatus != "" {
			if err := writeLine(w, "    Rollback: %s", report.Journal.RollbackStatus); err != nil {
				return err
			}
		}
	}
	return writeTransactionLockWarnings(w, report.Warnings)
}

func writeTransactionLockClearText(w io.Writer, report TransactionLockClearReport) error {
	if err := writeLine(w, "Transaction lock clear"); err != nil {
		return err
	}
	if err := writeLine(w, "  Status: %s", report.Status); err != nil {
		return err
	}
	if report.TransactionID != "" {
		if err := writeLine(w, "  Transaction: %s", report.TransactionID); err != nil {
			return err
		}
	}
	if report.Reason != "" {
		if err := writeLine(w, "  Reason: %s", report.Reason); err != nil {
			return err
		}
	}
	if err := writeLine(w, "  Lock cleared: %t", report.Lock.Cleared); err != nil {
		return err
	}
	if report.Message != "" {
		if err := writeLine(w, "  Message: %s", report.Message); err != nil {
			return err
		}
	}
	if report.PostClearState != "" {
		if err := writeLine(w, "  Post-clear: %s", report.PostClearState); err != nil {
			return err
		}
	}
	if report.Journal != nil && report.Journal.Present {
		if err := writeLine(w, "  Journal: %s", report.Journal.Status); err != nil {
			return err
		}
	}
	return writeTransactionLockWarnings(w, report.Warnings)
}

func writeTransactionDiagnosticsText(w io.Writer, report TransactionDiagnosticsReport) error {
	if err := writeLine(w, "Transaction diagnostics"); err != nil {
		return err
	}
	if err := writeLine(w, "  Publish blocked: %t", report.PublishBlocked); err != nil {
		return err
	}
	if err := writeLine(w, "  Lock: %s", report.Lock.Status); err != nil {
		return err
	}
	if report.Lock.Reason != "" {
		if err := writeLine(w, "    Reason: %s", report.Lock.Reason); err != nil {
			return err
		}
	}
	if err := writeLine(w, "  Operation lock: %s", operationLockState(report.OperationLock)); err != nil {
		return err
	}
	if len(report.Blockers) > 0 {
		if err := writeLine(w, "  Blockers:"); err != nil {
			return err
		}
		for _, blocker := range report.Blockers {
			if err := writeLine(w, "    %s: %s", blocker.Kind, formatDiagnosticsState(blockerLabel(blocker), blocker.Reason)); err != nil {
				return err
			}
		}
	}
	if len(report.Journals) > 0 {
		if err := writeLine(w, "  Journals:"); err != nil {
			return err
		}
		for _, journal := range report.Journals {
			state := journal.Status
			if journal.Corrupt {
				state = "corrupt"
			} else if journal.ReadFailed {
				state = "read_failed"
			}
			if err := writeLine(w, "    %s: %s", journalLabel(journal), state); err != nil {
				return err
			}
		}
	}
	for _, warning := range report.Warnings {
		if err := writeLine(w, "  Warning: %s: %s", warning.Code, warning.Message); err != nil {
			return err
		}
	}
	return nil
}

func operationLockState(lock TransactionOperationLockReport) string {
	switch {
	case lock.Corrupt:
		return formatDiagnosticsState("corrupt", lock.Message)
	case lock.ReadFailed:
		return formatDiagnosticsState("read_failed", lock.Message)
	case lock.Present:
		return lock.Operation
	default:
		return "absent"
	}
}

func blockerLabel(blocker TransactionDiagnosticsBlockerReport) string {
	if blocker.TransactionID != "" {
		return blocker.TransactionID
	}
	return blocker.Name
}

func journalLabel(journal TransactionJournalDiagnosticReport) string {
	if journal.TransactionID != "" {
		return journal.TransactionID
	}
	if journal.Name != "" {
		return journal.Name
	}
	return "<unknown>"
}

func formatDiagnosticsState(label, reason string) string {
	if label == "" {
		return reason
	}
	if reason == "" {
		return label
	}
	return label + " " + reason
}

func writeTransactionLockWarnings(w io.Writer, warnings []TransactionLockWarning) error {
	if len(warnings) == 0 {
		return nil
	}
	if err := writeLine(w, "  Warnings:"); err != nil {
		return err
	}
	for _, warning := range warnings {
		if err := writeLine(w, "    %s: %s", warning.Code, warning.Message); err != nil {
			return err
		}
	}
	return nil
}

func writeTransactionText(w io.Writer, report TransactionReport) error {
	if err := writeLine(w, "Transaction"); err != nil {
		return err
	}
	if err := writeLine(w, "  ID: %s", report.ID); err != nil {
		return err
	}
	if err := writeLine(w, "  Status: %s", report.Status); err != nil {
		return err
	}
	if report.RollbackStatus != "" {
		if err := writeLine(w, "  Rollback: %s", report.RollbackStatus); err != nil {
			return err
		}
	}
	for _, mod := range report.Modules {
		if err := writeLine(w, "  %s: branch=%s candidate=%s", mod.Module, mod.FinalBranchRef, mod.CandidateBranchRef); err != nil {
			return err
		}
	}
	for _, action := range report.ManualRecoveryActions {
		if err := writeLine(w, "  manual: %s %s", action.Ref, action.Message); err != nil {
			return err
		}
	}
	return nil
}

const timeFormat = "2006-01-02T15:04:05.000000000Z07:00"
