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
	"path/filepath"
	"strings"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
)

func TestValidateRollbackJournalAcceptsCanonicalState(t *testing.T) {
	journal := rollbackValidationFixture(t)

	if err := validateRollbackJournal(journal); err != nil {
		t.Fatalf("validateRollbackJournal() error = %v", err)
	}
}

func TestValidateRollbackJournalRejectsUntrustedRecoveryInputs(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*TransactionJournal)
		want   string
	}{
		{
			name: "unknown status",
			mutate: func(journal *TransactionJournal) {
				journal.Status = "future_status"
			},
			want: "unsupported transaction status",
		},
		{
			name: "option-like remote",
			mutate: func(journal *TransactionJournal) {
				journal.Remote = "--upload-pack=helper"
			},
			want: "option prefix",
		},
		{
			name: "relative worktree",
			mutate: func(journal *TransactionJournal) {
				journal.Modules[0].WorktreeDir = "relative/worktree"
			},
			want: "not absolute",
		},
		{
			name: "branch namespace mismatch",
			mutate: func(journal *TransactionJournal) {
				journal.Modules[0].FinalBranchRef = "refs/tags/main"
			},
			want: "outside refs/heads",
		},
		{
			name: "candidate namespace mismatch",
			mutate: func(journal *TransactionJournal) {
				journal.Modules[0].CandidateBranchRef = "refs/heads/arcpub/tx/other/foundation"
			},
			want: "does not match transaction namespace",
		},
		{
			name: "non hexadecimal object",
			mutate: func(journal *TransactionJournal) {
				journal.Modules[0].CreatedCommit = "HEAD~1"
			},
			want: "non-hexadecimal",
		},
		{
			name: "duplicate module",
			mutate: func(journal *TransactionJournal) {
				journal.Modules = append(journal.Modules, journal.Modules[0])
			},
			want: "duplicates module",
		},
		{
			name: "published state without commit",
			mutate: func(journal *TransactionJournal) {
				journal.Modules[0].CandidatePushed = true
				journal.Modules[0].CreatedCommit = ""
			},
			want: "without its created commit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			journal := rollbackValidationFixture(t)
			tt.mutate(&journal)

			err := validateRollbackJournal(journal)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("validateRollbackJournal() error = %v, want containing %q", err, tt.want)
			}
		})
	}
}

func TestRollbackTransactionRejectsInvalidJournalBeforeGitOperations(t *testing.T) {
	stateDir := t.TempDir()
	journal := rollbackValidationFixture(t)
	journal.WorktreeForTest("relative/worktree")
	store := NewFileJournalStore(stateDir)
	if err := store.Create(context.Background(), journal); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	fakeGit := porttest.NewGit()
	service := New(Dependencies{Git: fakeGit}, Options{})
	_, err := service.RollbackTransaction(context.Background(), stateDir, journal.ID)
	if err == nil || !strings.Contains(err.Error(), "recovery state is invalid") {
		t.Fatalf("RollbackTransaction() error = %v", err)
	}
	if len(fakeGit.Calls) != 0 {
		t.Fatalf("RollbackTransaction() Git calls = %#v, want none", fakeGit.Calls)
	}
}

func rollbackValidationFixture(t *testing.T) TransactionJournal {
	t.Helper()
	id := TransactionID("tx-recovery-validation")
	return TransactionJournal{
		SchemaVersion: transactionSchemaVersion,
		ID:            id,
		Status:        TransactionStatusFailed,
		Remote:        "origin",
		Modules: []ModuleTransactionState{{
			Module:             "foundation",
			Repository:         "arcoris/foundation",
			WorktreeDir:        filepath.Join(t.TempDir(), "foundation"),
			TargetBranch:       "main",
			FinalBranchRef:     "refs/heads/main",
			CandidateBranchRef: candidateRef(id, "foundation"),
		}},
	}
}

func (j *TransactionJournal) WorktreeForTest(path string) {
	j.Modules[0].WorktreeDir = path
}
