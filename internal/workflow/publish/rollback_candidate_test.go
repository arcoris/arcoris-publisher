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
	"testing"

	"arcoris.dev/arcoris-publisher/internal/ports/git"
	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
)

func TestRollbackCandidateRefusesMovedRef(t *testing.T) {
	fakeGit := porttest.NewGit()
	worktree := "/target/arcoris__foundation"
	candidate := "refs/heads/arcpub/tx/tx-test/foundation"
	fakeGit.RemoteRefHashes[porttest.RemoteRefKeyForRepo(worktree, "origin", candidate)] = "bbbbbbbb"

	runner := transactionRunner{
		service: New(Dependencies{Git: fakeGit}, Options{}),
		journal: TransactionJournal{ID: "tx-test"},
	}
	mod := ModuleTransactionState{
		Module:             "foundation",
		Repository:         "arcoris/foundation",
		WorktreeDir:        worktree,
		CandidateBranchRef: candidate,
		CandidatePushed:    true,
		CreatedCommit:      git.CommitHash("aaaaaaaa"),
	}

	runner.rollbackCandidate(context.Background(), &mod)

	if !mod.CandidatePushed {
		t.Fatal("CandidatePushed = false, want true for moved ref")
	}
	if mod.Rollback.CandidateDeleted {
		t.Fatal("CandidateDeleted = true, want false for moved ref")
	}
	if len(runner.journal.ManualActions) != 1 {
		t.Fatalf("manual actions = %#v, want one", runner.journal.ManualActions)
	}
	for _, call := range fakeGit.Calls {
		if call.Op == "delete-remote-ref" {
			t.Fatalf("rollback deleted moved candidate ref: %#v", call)
		}
	}
}

func TestRollbackCandidateDeletesOwnedRef(t *testing.T) {
	fakeGit := porttest.NewGit()
	worktree := "/target/arcoris__foundation"
	candidate := "refs/heads/arcpub/tx/tx-test/foundation"
	created := git.CommitHash("aaaaaaaa")
	fakeGit.RemoteRefHashes[porttest.RemoteRefKeyForRepo(worktree, "origin", candidate)] = created

	runner := transactionRunner{
		service: New(Dependencies{Git: fakeGit}, Options{}),
		journal: TransactionJournal{ID: "tx-test"},
	}
	mod := ModuleTransactionState{
		Module:             "foundation",
		Repository:         "arcoris/foundation",
		WorktreeDir:        worktree,
		CandidateBranchRef: candidate,
		CandidatePushed:    true,
		CreatedCommit:      created,
	}

	runner.rollbackCandidate(context.Background(), &mod)

	if mod.CandidatePushed {
		t.Fatal("CandidatePushed = true, want false after owned ref deletion")
	}
	if !mod.Rollback.CandidateDeleted {
		t.Fatal("CandidateDeleted = false, want true")
	}
	if len(runner.journal.ManualActions) != 0 {
		t.Fatalf("manual actions = %#v, want none", runner.journal.ManualActions)
	}
}
