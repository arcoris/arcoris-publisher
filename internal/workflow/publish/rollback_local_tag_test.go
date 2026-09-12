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

func TestRollbackLocalTagRefusesMovedTag(t *testing.T) {
	fakeGit := porttest.NewGit()
	fakeGit.Tags["v1.0.0"] = true
	fakeGit.CommitHash = "bbbbbbbb"

	runner := transactionRunner{
		service: New(Dependencies{Git: fakeGit}, Options{}),
		journal: TransactionJournal{ID: "tx-test"},
	}
	mod := ModuleTransactionState{
		Module:          "foundation",
		Repository:      "arcoris/foundation",
		WorktreeDir:     "/target/arcoris__foundation",
		FinalTagRef:     "refs/tags/v1.0.0",
		LocalTagCreated: true,
		CreatedCommit:   git.CommitHash("aaaaaaaa"),
	}

	runner.rollbackLocalTag(context.Background(), &mod)

	if !mod.LocalTagCreated {
		t.Fatal("LocalTagCreated = false, want true for moved tag")
	}
	if mod.Rollback.LocalTagDeleted {
		t.Fatal("LocalTagDeleted = true, want false for moved tag")
	}
	if len(runner.journal.ManualActions) != 1 {
		t.Fatalf("manual actions = %#v, want one", runner.journal.ManualActions)
	}
	if !fakeGit.Tags["v1.0.0"] {
		t.Fatal("moved local tag was deleted")
	}
}

func TestRollbackLocalTagDeletesOwnedTag(t *testing.T) {
	fakeGit := porttest.NewGit()
	fakeGit.Tags["v1.0.0"] = true
	fakeGit.CommitHash = "aaaaaaaa"

	runner := transactionRunner{
		service: New(Dependencies{Git: fakeGit}, Options{}),
		journal: TransactionJournal{ID: "tx-test"},
	}
	mod := ModuleTransactionState{
		Module:          "foundation",
		Repository:      "arcoris/foundation",
		WorktreeDir:     "/target/arcoris__foundation",
		FinalTagRef:     "refs/tags/v1.0.0",
		LocalTagCreated: true,
		CreatedCommit:   git.CommitHash("aaaaaaaa"),
	}

	runner.rollbackLocalTag(context.Background(), &mod)

	if mod.LocalTagCreated {
		t.Fatal("LocalTagCreated = true, want false after owned tag deletion")
	}
	if !mod.Rollback.LocalTagDeleted {
		t.Fatal("LocalTagDeleted = false, want true")
	}
	if fakeGit.Tags["v1.0.0"] {
		t.Fatal("owned local tag still exists")
	}
	if len(runner.journal.ManualActions) != 0 {
		t.Fatalf("manual actions = %#v, want none", runner.journal.ManualActions)
	}
}
