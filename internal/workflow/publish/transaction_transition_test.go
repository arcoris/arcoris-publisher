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
	"time"
)

func TestTransactionStatusTransitionGraph(t *testing.T) {
	allowed := [][2]TransactionStatus{
		{TransactionStatusPending, TransactionStatusPreflighted},
		{TransactionStatusPreflighted, TransactionStatusSnapshotted},
		{TransactionStatusSnapshotted, TransactionStatusCommittedLocally},
		{TransactionStatusCommittedLocally, TransactionStatusCandidatesPushed},
		{TransactionStatusCandidatesPushed, TransactionStatusPromoting},
		{TransactionStatusPromoting, TransactionStatusBranchesPromoted},
		{TransactionStatusBranchesPromoted, TransactionStatusTagging},
		{TransactionStatusBranchesPromoted, TransactionStatusCommitted},
		{TransactionStatusTagging, TransactionStatusCommitted},
		{TransactionStatusCandidatesPushed, TransactionStatusFailed},
		{TransactionStatusFailed, TransactionStatusRollingBack},
		{TransactionStatusPromoting, TransactionStatusRollingBack},
		{TransactionStatusRollingBack, TransactionStatusRollbackFailed},
		{TransactionStatusRollbackFailed, TransactionStatusRollingBack},
		{TransactionStatusRollingBack, TransactionStatusRolledBack},
		{TransactionStatusRollingBack, TransactionStatusRollingBack},
	}
	for _, pair := range allowed {
		if err := validateTransactionStatusTransition(pair[0], pair[1]); err != nil {
			t.Errorf("transition %s -> %s rejected: %v", pair[0], pair[1], err)
		}
	}

	rejected := [][2]TransactionStatus{
		{TransactionStatusPending, TransactionStatusCommitted},
		{TransactionStatusCommitted, TransactionStatusRollingBack},
		{TransactionStatusCommitted, TransactionStatusFailed},
		{TransactionStatusRolledBack, TransactionStatusRollingBack},
		{TransactionStatusFailed, TransactionStatusCommitted},
		{TransactionStatusRollingBack, TransactionStatusCommitted},
		{TransactionStatus("future"), TransactionStatusFailed},
		{TransactionStatusPending, TransactionStatus("future")},
	}
	for _, pair := range rejected {
		if err := validateTransactionStatusTransition(pair[0], pair[1]); !errors.Is(err, errInvalidTransactionStatusTransition) {
			t.Errorf("transition %s -> %s error = %v, want invalid transition", pair[0], pair[1], err)
		}
	}
}

func TestSetStatusPersistsBeforePublishingInMemoryState(t *testing.T) {
	persistErr := errors.New("journal write denied")
	store := &recordingTransitionStore{updateErr: persistErr}
	started := time.Unix(1, 0).UTC()
	runner := transactionRunner{
		store: store,
		journal: TransactionJournal{
			ID:        "tx-test",
			Status:    TransactionStatusPending,
			UpdatedAt: started,
		},
	}

	err := runner.setStatus(context.Background(), TransactionStatusPreflighted)
	if !errors.Is(err, persistErr) {
		t.Fatalf("setStatus() error = %v, want persistence failure", err)
	}
	if runner.journal.Status != TransactionStatusPending {
		t.Fatalf("status = %q, want pending after failed persistence", runner.journal.Status)
	}
	if !runner.journal.UpdatedAt.Equal(started) {
		t.Fatalf("UpdatedAt changed after failed persistence: %s", runner.journal.UpdatedAt)
	}
	if store.last.Status != TransactionStatusPreflighted {
		t.Fatalf("persisted candidate status = %q", store.last.Status)
	}
}

func TestSetStatusRejectsIllegalTransitionBeforeStoreWrite(t *testing.T) {
	store := &recordingTransitionStore{}
	runner := transactionRunner{store: store, journal: TransactionJournal{ID: "tx-test", Status: TransactionStatusPending}}

	err := runner.setStatus(context.Background(), TransactionStatusCommitted)
	if !errors.Is(err, errInvalidTransactionStatusTransition) {
		t.Fatalf("setStatus() error = %v, want invalid transition", err)
	}
	if store.updates != 0 {
		t.Fatalf("Update() calls = %d, want 0", store.updates)
	}
	if runner.journal.Status != TransactionStatusPending {
		t.Fatalf("status = %q, want pending", runner.journal.Status)
	}
}

type recordingTransitionStore struct {
	updates   int
	last      TransactionJournal
	updateErr error
}

func (*recordingTransitionStore) Create(context.Context, TransactionJournal) error { return nil }
func (s *recordingTransitionStore) Update(_ context.Context, journal TransactionJournal) error {
	s.updates++
	s.last = journal
	return s.updateErr
}
func (*recordingTransitionStore) Load(context.Context, TransactionID) (TransactionJournal, error) {
	return TransactionJournal{}, nil
}
func (*recordingTransitionStore) List(context.Context) ([]TransactionSummary, error) { return nil, nil }
func (*recordingTransitionStore) HasPending(context.Context) (TransactionSummary, bool, error) {
	return TransactionSummary{}, false, nil
}
