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
	"testing"

	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
)

func TestRollbackStopsBeforeGitMutationWhenIntentCannotPersist(t *testing.T) {
	persistErr := errors.New("journal write denied")
	store := &failingJournalStore{failUpdateAt: 1, updateErr: persistErr}
	fakeGit := porttest.NewGit()
	runner := transactionRunner{
		service: New(Dependencies{Git: fakeGit}, Options{}),
		store:   store,
		journal: TransactionJournal{
			ID:     "tx-test",
			Status: TransactionStatusFailed,
			Modules: []ModuleTransactionState{{
				Module:             "foundation",
				WorktreeDir:        "/worktree",
				CandidateBranchRef: "refs/heads/arcpub/tx/tx-test/foundation",
				CandidatePushed:    true,
			}},
		},
	}

	err := runner.rollback(context.Background())
	if !errors.Is(err, persistErr) {
		t.Fatalf("rollback() error = %v, want journal persistence failure", err)
	}
	if len(fakeGit.Calls) != 0 {
		t.Fatalf("rollback mutated Git before durable rolling_back state: %#v", fakeGit.Calls)
	}
	if store.updateCalls != 1 {
		t.Fatalf("Update() calls = %d, want 1", store.updateCalls)
	}
}

func TestRollbackStopsAfterProgressPersistenceFailure(t *testing.T) {
	persistErr := errors.New("journal write denied")
	store := &failingJournalStore{failUpdateAt: 2, updateErr: persistErr}
	fakeGit := porttest.NewGit()
	runner := transactionRunner{
		service: New(Dependencies{Git: fakeGit}, Options{}),
		store:   store,
		journal: TransactionJournal{
			ID:     "tx-test",
			Status: TransactionStatusFailed,
			Modules: []ModuleTransactionState{{
				Module:             "foundation",
				WorktreeDir:        "/worktree",
				CandidateBranchRef: "refs/heads/arcpub/tx/tx-test/foundation",
				CandidatePushed:    true,
			}},
		},
	}

	err := runner.rollback(context.Background())
	if !errors.Is(err, persistErr) {
		t.Fatalf("rollback() error = %v, want journal persistence failure", err)
	}
	if store.updateCalls != 2 {
		t.Fatalf("Update() calls = %d, want 2", store.updateCalls)
	}
	if runner.journal.Status != TransactionStatusRollingBack {
		t.Fatalf("status = %q, want rolling_back", runner.journal.Status)
	}
}

func TestFailDoesNotStartRollbackWhenFailureStateCannotPersist(t *testing.T) {
	publishErr := errors.New("candidate push failed")
	persistErr := errors.New("journal write denied")
	store := &failingJournalStore{failUpdateAt: 1, updateErr: persistErr}
	fakeGit := porttest.NewGit()
	runner := transactionRunner{
		service: New(Dependencies{Git: fakeGit}, Options{RollbackMode: RollbackAutomatic}),
		store:   store,
		journal: TransactionJournal{
			ID:     "tx-test",
			Status: TransactionStatusCandidatesPushed,
			Modules: []ModuleTransactionState{{
				Module:             "foundation",
				WorktreeDir:        "/worktree",
				CandidateBranchRef: "refs/heads/arcpub/tx/tx-test/foundation",
				CandidatePushed:    true,
			}},
		},
	}

	_, err := runner.fail(context.Background(), publishErr)
	if !errors.Is(err, publishErr) || !errors.Is(err, persistErr) {
		t.Fatalf("fail() error = %v, want publish and journal persistence failures", err)
	}
	if len(fakeGit.Calls) != 0 {
		t.Fatalf("fail() started rollback without durable failed state: %#v", fakeGit.Calls)
	}
}

func TestFailSurfacesManualRollbackPersistenceFailure(t *testing.T) {
	publishErr := errors.New("candidate push failed")
	persistErr := errors.New("journal write denied")
	store := &failingJournalStore{failUpdateAt: 2, updateErr: persistErr}
	runner := transactionRunner{
		service: New(Dependencies{Git: porttest.NewGit()}, Options{RollbackMode: RollbackManual}),
		store:   store,
		journal: TransactionJournal{ID: "tx-test", Status: TransactionStatusCandidatesPushed},
	}

	_, err := runner.fail(context.Background(), publishErr)
	if !errors.Is(err, publishErr) || !errors.Is(err, persistErr) {
		t.Fatalf("fail() error = %v, want publish and manual rollback persistence failures", err)
	}
	if store.updateCalls != 2 {
		t.Fatalf("Update() calls = %d, want 2", store.updateCalls)
	}
}

type failingJournalStore struct {
	updateCalls  int
	failUpdateAt int
	updateErr    error
}

func (s *failingJournalStore) Create(context.Context, TransactionJournal) error { return nil }

func (s *failingJournalStore) Update(context.Context, TransactionJournal) error {
	s.updateCalls++
	if s.updateCalls == s.failUpdateAt {
		return s.updateErr
	}
	return nil
}

func (s *failingJournalStore) Load(context.Context, TransactionID) (TransactionJournal, error) {
	return TransactionJournal{}, nil
}

func (s *failingJournalStore) List(context.Context) ([]TransactionSummary, error) { return nil, nil }

func (s *failingJournalStore) HasPending(context.Context) (TransactionSummary, bool, error) {
	return TransactionSummary{}, false, nil
}
