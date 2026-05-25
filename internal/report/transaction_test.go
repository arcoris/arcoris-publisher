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
	"bytes"
	"strings"
	"testing"
	"time"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/ports/git"
	"arcoris.dev/arcoris-publisher/internal/workflow/publish"
)

func TestBuildTransactionReportHidesLocalPathsByDefault(t *testing.T) {
	t.Parallel()

	report := BuildTransactionReport(transactionFixture(), Options{})

	if report.Modules[0].WorktreeDir != "" {
		t.Fatalf("WorktreeDir leaked = %q", report.Modules[0].WorktreeDir)
	}
}

func TestRendererTransactionJSON(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := New(Options{Format: FormatJSON, Pretty: true}).Transaction(&buf, transactionFixture()); err != nil {
		t.Fatalf("Transaction(JSON) error = %v", err)
	}
	if !strings.Contains(buf.String(), `"kind": "transaction"`) ||
		!strings.Contains(buf.String(), `"status": "rollback_failed"`) {
		t.Fatalf("transaction JSON = %s", buf.String())
	}
}

func TestBuildTransactionListReportIncludesRollbackStatus(t *testing.T) {
	t.Parallel()

	report := BuildTransactionListReport([]publish.TransactionSummary{{
		ID:       "tx-test",
		Status:   publish.TransactionStatusRollbackFailed,
		Rollback: publish.RollbackStatusFailed,
	}})

	if report.Transactions[0].RollbackStatus != "failed" {
		t.Fatalf("RollbackStatus = %q", report.Transactions[0].RollbackStatus)
	}
}

func TestTransactionPruneReportHidesLocalPathsByDefault(t *testing.T) {
	t.Parallel()

	report := BuildTransactionPruneReport(transactionPruneFixture(), Options{})

	if report.Matched[0].Path != "" {
		t.Fatalf("Path leaked = %q", report.Matched[0].Path)
	}
}

func TestRendererTransactionPruneJSONAndText(t *testing.T) {
	t.Parallel()

	var jsonBuf bytes.Buffer
	if err := New(Options{Format: FormatJSON, Pretty: true}).TransactionPrune(&jsonBuf, transactionPruneFixture()); err != nil {
		t.Fatalf("TransactionPrune(JSON) error = %v", err)
	}
	if !strings.Contains(jsonBuf.String(), `"kind": "transactions-prune"`) ||
		!strings.Contains(jsonBuf.String(), `"status": "dry_run"`) {
		t.Fatalf("transaction prune JSON = %s", jsonBuf.String())
	}
	if strings.Contains(jsonBuf.String(), "/state") {
		t.Fatalf("transaction prune JSON leaked local path: %s", jsonBuf.String())
	}

	var textBuf bytes.Buffer
	if err := New(Options{Format: FormatText}).TransactionPrune(&textBuf, transactionPruneFixture()); err != nil {
		t.Fatalf("TransactionPrune(text) error = %v", err)
	}
	if !strings.Contains(textBuf.String(), "Transactions prune") ||
		!strings.Contains(textBuf.String(), "Matched: 1") {
		t.Fatalf("transaction prune text = %s", textBuf.String())
	}
}

func TestTransactionPruneReportIncludesLocalPathsWhenRequested(t *testing.T) {
	t.Parallel()

	report := BuildTransactionPruneReport(transactionPruneFixture(), Options{IncludeLocalPaths: true})

	if report.Matched[0].Path != "/state/transactions/tx-test.json" {
		t.Fatalf("Path = %q", report.Matched[0].Path)
	}
}

func TestTransactionLockReportJSONAndText(t *testing.T) {
	t.Parallel()

	var jsonBuf bytes.Buffer
	if err := New(Options{Format: FormatJSON, Pretty: true}).TransactionLock(&jsonBuf, transactionLockFixture()); err != nil {
		t.Fatalf("TransactionLock(JSON) error = %v", err)
	}
	if !strings.Contains(jsonBuf.String(), `"kind": "transactions-lock"`) ||
		!strings.Contains(jsonBuf.String(), `"status": "present"`) ||
		strings.Contains(jsonBuf.String(), "/state") {
		t.Fatalf("transaction lock JSON = %s", jsonBuf.String())
	}

	var textBuf bytes.Buffer
	if err := New(Options{Format: FormatText}).TransactionLock(&textBuf, transactionLockFixture()); err != nil {
		t.Fatalf("TransactionLock(text) error = %v", err)
	}
	if !strings.Contains(textBuf.String(), "Transaction lock") ||
		!strings.Contains(textBuf.String(), "Transaction: tx-test") {
		t.Fatalf("transaction lock text = %s", textBuf.String())
	}
}

func TestTransactionLockAbsentJSONUsesNullLockAndJournal(t *testing.T) {
	t.Parallel()

	report := BuildTransactionLockReport(publish.LockShowResult{Status: publish.LockShowStatusAbsent}, Options{})

	if report.Lock != nil || report.Journal != nil || report.Status != "absent" {
		t.Fatalf("absent report = %#v", report)
	}
}

func TestTransactionLockClearReportJSONAndText(t *testing.T) {
	t.Parallel()

	var jsonBuf bytes.Buffer
	if err := New(Options{Format: FormatJSON, Pretty: true}).TransactionLockClear(&jsonBuf, transactionLockClearFixture()); err != nil {
		t.Fatalf("TransactionLockClear(JSON) error = %v", err)
	}
	if !strings.Contains(jsonBuf.String(), `"kind": "transactions-lock-clear"`) ||
		!strings.Contains(jsonBuf.String(), `"status": "cleared"`) ||
		strings.Contains(jsonBuf.String(), "/state") {
		t.Fatalf("transaction lock clear JSON = %s", jsonBuf.String())
	}

	var textBuf bytes.Buffer
	if err := New(Options{Format: FormatText}).TransactionLockClear(&textBuf, transactionLockClearFixture()); err != nil {
		t.Fatalf("TransactionLockClear(text) error = %v", err)
	}
	if !strings.Contains(textBuf.String(), "Transaction lock clear") ||
		!strings.Contains(textBuf.String(), "Status: cleared") {
		t.Fatalf("transaction lock clear text = %s", textBuf.String())
	}
}

func TestTransactionLockReportsIncludeLocalPathsWhenRequested(t *testing.T) {
	t.Parallel()

	show := BuildTransactionLockReport(transactionLockFixture(), Options{IncludeLocalPaths: true})
	if show.Lock == nil || show.Lock.Path != "/state/publish.lock" {
		t.Fatalf("lock path = %#v", show.Lock)
	}
	clear := BuildTransactionLockClearReport(transactionLockClearFixture(), Options{IncludeLocalPaths: true})
	if clear.Lock.Path != "/state/publish.lock" {
		t.Fatalf("clear path = %q", clear.Lock.Path)
	}
}

func transactionFixture() publish.TransactionJournal {
	now := time.Unix(1, 2).UTC()
	return publish.TransactionJournal{
		ID:        "tx-test",
		Status:    publish.TransactionStatusRollbackFailed,
		Rollback:  publish.RollbackStatusFailed,
		Version:   "v0.1.0",
		StartedAt: now,
		UpdatedAt: now,
		Modules: []publish.ModuleTransactionState{{
			Module:             manifest.ModuleName("foundation"),
			Repository:         manifest.RepositoryRef("arcoris/foundation"),
			WorktreeDir:        "/target/arcoris__foundation",
			FinalBranchRef:     "refs/heads/main",
			CandidateBranchRef: "refs/heads/arcpub/tx/tx-test/foundation",
			CreatedCommit:      git.CommitHash("abc123"),
		}},
		ManualActions: []publish.ManualRecoveryAction{{
			Module:     manifest.ModuleName("foundation"),
			Repository: manifest.RepositoryRef("arcoris/foundation"),
			Ref:        "refs/heads/main",
			Message:    "manual restore required",
		}},
	}
}

func transactionLockFixture() publish.LockShowResult {
	now := time.Unix(1, 2).UTC()
	return publish.LockShowResult{
		Status: publish.LockShowStatusPresent,
		Lock: publish.TransactionLockInfo{
			ID:        "tx-test",
			PID:       "123",
			Command:   "publish",
			StartedAt: now.Format(timeFormat),
			Path:      "/state/publish.lock",
		},
		Journal: publish.LockJournalState{
			Present:   true,
			Status:    publish.TransactionStatusRollbackFailed,
			Rollback:  publish.RollbackStatusFailed,
			Version:   "v0.1.0",
			StartedAt: now,
			UpdatedAt: now,
		},
	}
}

func transactionLockClearFixture() publish.LockClearResult {
	fixture := transactionLockFixture()
	return publish.LockClearResult{
		Status:        publish.LockClearStatusCleared,
		TransactionID: "tx-test",
		Lock:          fixture.Lock,
		LockCleared:   true,
		Journal:       fixture.Journal,
		Reason:        "publish lock cleared",
	}
}

func transactionPruneFixture() publish.PruneResult {
	now := time.Unix(1, 2).UTC()
	return publish.PruneResult{
		Status: publish.PruneStatusDryRun,
		Matched: []publish.PruneEntry{{
			ID:        "tx-test",
			Status:    publish.TransactionStatusCommitted,
			Version:   "v0.1.0",
			StartedAt: now,
			UpdatedAt: now,
			Path:      "/state/transactions/tx-test.json",
			Reason:    "status committed",
		}},
		Skipped: []publish.PruneEntry{{
			ID:        "tx-pending",
			Status:    publish.TransactionStatusPending,
			StartedAt: now,
			UpdatedAt: now,
			Path:      "/state/transactions/tx-pending.json",
			Reason:    "status is not prunable",
		}},
	}
}
