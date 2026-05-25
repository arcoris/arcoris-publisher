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
