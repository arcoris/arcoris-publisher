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

package app

import (
	"context"
	"time"

	"arcoris.dev/arcoris-publisher/internal/workflow/publish"
)

// TransactionID identifies one durable publish transaction.
type TransactionID = publish.TransactionID

// TransactionStatus identifies one durable publish transaction status.
type TransactionStatus = publish.TransactionStatus

// RollbackMode controls publish transaction rollback behavior.
type RollbackMode = publish.RollbackMode

const (
	// TransactionStatusCommitted means final refs were published successfully.
	TransactionStatusCommitted = publish.TransactionStatusCommitted
	// TransactionStatusRolledBack means rollback completed successfully.
	TransactionStatusRolledBack = publish.TransactionStatusRolledBack

	// RollbackAutomatic attempts rollback immediately after publish failure.
	RollbackAutomatic = publish.RollbackAutomatic
	// RollbackManual leaves rollback to the recovery command.
	RollbackManual = publish.RollbackManual
	// RollbackDisabled leaves side effects in place and is unsafe for production.
	RollbackDisabled = publish.RollbackDisabled
)

// TransactionRequest selects durable publish transaction state.
type TransactionRequest struct {
	// StateDir contains publish transaction journals and locks.
	StateDir string

	// TransactionID selects one transaction for show or rollback.
	TransactionID publish.TransactionID
}

// TransactionPruneRequest selects terminal publish transactions to prune.
type TransactionPruneRequest struct {
	StateDir  string
	Statuses  []publish.TransactionStatus
	OlderThan time.Duration
	DryRun    bool
}

// TransactionListResult contains durable publish transaction summaries.
type TransactionListResult struct{ summaries []publish.TransactionSummary }

// Summaries returns detached transaction summaries.
func (r TransactionListResult) Summaries() []publish.TransactionSummary {
	out := make([]publish.TransactionSummary, len(r.summaries))
	copy(out, r.summaries)
	return out
}

// TransactionResult contains one durable publish transaction journal.
type TransactionResult struct{ journal publish.TransactionJournal }

// Journal returns the transaction journal.
func (r TransactionResult) Journal() publish.TransactionJournal { return r.journal }

// TransactionPruneResult contains transaction prune details.
type TransactionPruneResult struct{ result publish.PruneResult }

// Result returns the workflow prune result.
func (r TransactionPruneResult) Result() publish.PruneResult { return r.result }

// ListTransactions lists durable publish transactions.
func (a App) ListTransactions(ctx context.Context, req TransactionRequest) (TransactionListResult, error) {
	summaries, err := publish.New(a.workflowDeps.Publish, a.workflowOptions.Publish).
		ListTransactions(ctx, req.StateDir)
	if err != nil {
		return TransactionListResult{}, err
	}
	return TransactionListResult{summaries: summaries}, nil
}

// ShowTransaction loads one durable publish transaction.
func (a App) ShowTransaction(ctx context.Context, req TransactionRequest) (TransactionResult, error) {
	journal, err := publish.New(a.workflowDeps.Publish, a.workflowOptions.Publish).
		ShowTransaction(ctx, req.StateDir, req.TransactionID)
	if err != nil {
		return TransactionResult{}, err
	}
	return TransactionResult{journal: journal}, nil
}

// RollbackTransaction rolls back one durable publish transaction.
func (a App) RollbackTransaction(ctx context.Context, req TransactionRequest) (TransactionResult, error) {
	journal, err := publish.New(a.workflowDeps.Publish, a.workflowOptions.Publish).
		RollbackTransaction(ctx, req.StateDir, req.TransactionID)
	if err != nil {
		return TransactionResult{journal: journal}, err
	}
	return TransactionResult{journal: journal}, nil
}

// PruneTransactions safely removes selected terminal transaction journals.
func (a App) PruneTransactions(ctx context.Context, req TransactionPruneRequest) (TransactionPruneResult, error) {
	result, err := publish.New(a.workflowDeps.Publish, a.workflowOptions.Publish).
		PruneTransactions(ctx, req.StateDir, publish.PruneOptions{
			Statuses:  req.Statuses,
			OlderThan: req.OlderThan,
			DryRun:    req.DryRun,
		})
	if err != nil {
		return TransactionPruneResult{result: result}, err
	}
	return TransactionPruneResult{result: result}, nil
}
