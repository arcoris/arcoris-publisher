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
	"fmt"
)

// ListTransactions returns durable publish transaction summaries.
func (s Service) ListTransactions(ctx context.Context, stateDir string) ([]TransactionSummary, error) {
	summaries, err := NewFileJournalStore(stateDir).List(ctx)
	if err != nil {
		return nil, &Error{Code: CodeRecoveryFailed, Message: "list publish transactions failed", Cause: err}
	}
	return summaries, nil
}

// ShowTransaction loads one durable publish transaction journal.
func (s Service) ShowTransaction(ctx context.Context, stateDir string, id TransactionID) (TransactionJournal, error) {
	journal, err := NewFileJournalStore(stateDir).Load(ctx, id)
	if err != nil {
		return TransactionJournal{}, &Error{Code: CodeRecoveryFailed, Message: "load publish transaction failed", Cause: err}
	}
	return journal, nil
}

// RollbackTransaction attempts compensating rollback from a durable journal.
func (s Service) RollbackTransaction(ctx context.Context, stateDir string, id TransactionID) (TransactionJournal, error) {
	store := NewFileJournalStore(stateDir)
	journal, err := store.Load(ctx, id)
	if err != nil {
		return TransactionJournal{}, &Error{Code: CodeRecoveryFailed, Message: "load publish transaction failed", Cause: err}
	}
	if lock, ok, err := currentTransactionLock(stateDir); err != nil {
		return journal, &Error{Code: CodeLockFailed, Message: "read publish transaction lock failed", Cause: err}
	} else if ok {
		return journal, &Error{
			Code:    CodeLockFailed,
			Message: fmt.Sprintf("publish transaction lock exists for %s; refusing rollback while publish may be active", lock.ID),
		}
	}
	if journal.Status == TransactionStatusCommitted {
		return journal, &Error{
			Code:    CodeRecoveryFailed,
			Message: fmt.Sprintf("transaction %s is committed and cannot be rolled back without force", id),
		}
	}
	if journal.Status == TransactionStatusRolledBack {
		return journal, nil
	}

	if journal.Remote != "" {
		s.opts.RemoteName = journal.Remote
	}
	runner := transactionRunner{service: s, store: store, journal: journal}
	err = runner.rollback(ctx)
	return runner.journal, err
}
