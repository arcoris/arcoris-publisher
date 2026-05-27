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
		!strings.Contains(jsonBuf.String(), `"reason": "lock_present"`) ||
		strings.Contains(jsonBuf.String(), "/state") {
		t.Fatalf("transaction lock JSON = %s", jsonBuf.String())
	}

	var textBuf bytes.Buffer
	if err := New(Options{Format: FormatText}).TransactionLock(&textBuf, transactionLockFixture()); err != nil {
		t.Fatalf("TransactionLock(text) error = %v", err)
	}
	if !strings.Contains(textBuf.String(), "Transaction lock") ||
		!strings.Contains(textBuf.String(), "Reason: lock_present") ||
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
		!strings.Contains(jsonBuf.String(), `"reason": "cleared"`) ||
		!strings.Contains(jsonBuf.String(), `"postClearState": "ready_for_publish"`) ||
		strings.Contains(jsonBuf.String(), "/state") {
		t.Fatalf("transaction lock clear JSON = %s", jsonBuf.String())
	}

	var textBuf bytes.Buffer
	if err := New(Options{Format: FormatText}).TransactionLockClear(&textBuf, transactionLockClearFixture()); err != nil {
		t.Fatalf("TransactionLockClear(text) error = %v", err)
	}
	if !strings.Contains(textBuf.String(), "Transaction lock clear") ||
		!strings.Contains(textBuf.String(), "Status: cleared") ||
		!strings.Contains(textBuf.String(), "Reason: cleared") ||
		!strings.Contains(textBuf.String(), "Lock cleared: true") ||
		!strings.Contains(textBuf.String(), "Post-clear: ready_for_publish") {
		t.Fatalf("transaction lock clear text = %s", textBuf.String())
	}
}

func TestTransactionLockClearPartialSyncFailureReport(t *testing.T) {
	t.Parallel()

	result := transactionLockClearFixture()
	result.Status = publish.LockClearStatusFailed
	result.Reason = publish.LockClearReasonSyncFailed
	result.Message = "sync publish lock directory failed"
	result.Journal.Status = publish.TransactionStatusCommitted
	result.Journal.Rollback = ""

	var jsonBuf bytes.Buffer
	if err := New(Options{Format: FormatJSON, Pretty: true}).TransactionLockClear(&jsonBuf, result); err != nil {
		t.Fatalf("TransactionLockClear(JSON) error = %v", err)
	}
	if !strings.Contains(jsonBuf.String(), `"status": "failed"`) ||
		!strings.Contains(jsonBuf.String(), `"reason": "sync_failed"`) ||
		!strings.Contains(jsonBuf.String(), `"cleared": true`) ||
		!strings.Contains(jsonBuf.String(), `"postClearState": "ready_for_publish"`) ||
		!strings.Contains(jsonBuf.String(), `"status": "committed"`) ||
		strings.Contains(jsonBuf.String(), "/state") {
		t.Fatalf("transaction lock clear JSON = %s", jsonBuf.String())
	}

	var textBuf bytes.Buffer
	if err := New(Options{Format: FormatText}).TransactionLockClear(&textBuf, result); err != nil {
		t.Fatalf("TransactionLockClear(text) error = %v", err)
	}
	if !strings.Contains(textBuf.String(), "Status: failed") ||
		!strings.Contains(textBuf.String(), "Reason: sync_failed") ||
		!strings.Contains(textBuf.String(), "Lock cleared: true") ||
		!strings.Contains(textBuf.String(), "Post-clear: ready_for_publish") {
		t.Fatalf("transaction lock clear text = %s", textBuf.String())
	}
}

func TestTransactionLockWarningCodesAreStable(t *testing.T) {
	t.Parallel()

	result := transactionLockFixture()
	result.Warnings = []publish.LockWarning{
		{Code: publish.LockWarningJournalMissing, Message: "publish lock references a transaction journal that does not exist"},
		{Code: publish.LockWarningJournalCorrupt, Message: "publish lock references a transaction journal that cannot be read safely"},
		{Code: publish.LockWarningLockCorrupt, Message: "publish lock is not parseable"},
	}

	var jsonBuf bytes.Buffer
	if err := New(Options{Format: FormatJSON, Pretty: true}).TransactionLock(&jsonBuf, result); err != nil {
		t.Fatalf("TransactionLock(JSON) error = %v", err)
	}
	if !strings.Contains(jsonBuf.String(), `"code": "journal_missing"`) ||
		!strings.Contains(jsonBuf.String(), `"code": "journal_corrupt"`) ||
		!strings.Contains(jsonBuf.String(), `"code": "lock_corrupt"`) {
		t.Fatalf("transaction lock JSON = %s", jsonBuf.String())
	}

	var textBuf bytes.Buffer
	if err := New(Options{Format: FormatText}).TransactionLock(&textBuf, result); err != nil {
		t.Fatalf("TransactionLock(text) error = %v", err)
	}
	if !strings.Contains(textBuf.String(), "journal_missing") ||
		!strings.Contains(textBuf.String(), "journal_corrupt") ||
		!strings.Contains(textBuf.String(), "lock_corrupt") {
		t.Fatalf("transaction lock text = %s", textBuf.String())
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

func TestTransactionDiagnosticsReportJSONAndText(t *testing.T) {
	t.Parallel()

	var jsonBuf bytes.Buffer
	if err := New(Options{Format: FormatJSON, Pretty: true}).TransactionDiagnostics(&jsonBuf, transactionDiagnosticsFixture()); err != nil {
		t.Fatalf("TransactionDiagnostics(JSON) error = %v", err)
	}
	for _, want := range []string{
		`"kind": "transactions-diagnostics"`,
		`"publishBlocked": true`,
		`"kind": "corrupt_journal"`,
		`"name": "bad.lock.json"`,
	} {
		if !strings.Contains(jsonBuf.String(), want) {
			t.Fatalf("diagnostics JSON missing %q:\n%s", want, jsonBuf.String())
		}
	}
	if strings.Contains(jsonBuf.String(), "/state") {
		t.Fatalf("diagnostics JSON leaked local path: %s", jsonBuf.String())
	}

	var textBuf bytes.Buffer
	if err := New(Options{Format: FormatText}).TransactionDiagnostics(&textBuf, transactionDiagnosticsFixture()); err != nil {
		t.Fatalf("TransactionDiagnostics(text) error = %v", err)
	}
	for _, want := range []string{
		"Transaction diagnostics",
		"Publish blocked: true",
		"corrupt_journal",
		"bad.lock.json",
	} {
		if !strings.Contains(textBuf.String(), want) {
			t.Fatalf("diagnostics text missing %q:\n%s", want, textBuf.String())
		}
	}
}

func TestTransactionDiagnosticsReportIncludesLocalPathsWhenRequested(t *testing.T) {
	t.Parallel()

	report := BuildTransactionDiagnosticsReport(transactionDiagnosticsFixture(), Options{IncludeLocalPaths: true})

	if len(report.Journals) != 2 {
		t.Fatalf("journals = %#v", report.Journals)
	}
	if report.Journals[0].Path != "/state/transactions/tx-test.json" {
		t.Fatalf("journal path = %q", report.Journals[0].Path)
	}
	if report.Lock.Lock == nil || report.Lock.Lock.Path != "/state/publish.lock" {
		t.Fatalf("lock = %#v", report.Lock.Lock)
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
		Status:  publish.LockShowStatusPresent,
		Reason:  publish.LockShowReasonLockPresent,
		Message: "publish lock is present",
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
		Status:         publish.LockClearStatusCleared,
		TransactionID:  "tx-test",
		Lock:           fixture.Lock,
		LockCleared:    true,
		Journal:        fixture.Journal,
		Reason:         publish.LockClearReasonCleared,
		Message:        "publish lock cleared",
		PostClearState: publish.LockPostClearReadyForPublish,
	}
}

func transactionDiagnosticsFixture() publish.TransactionStateDiagnostics {
	lock := transactionLockFixture()
	return publish.TransactionStateDiagnostics{
		PublishBlocked: true,
		Blockers: []publish.TransactionStateBlocker{
			{
				Kind:          publish.TransactionBlockerPublishLock,
				TransactionID: "tx-test",
				Status:        publish.TransactionStatusRollbackFailed,
				Reason:        publish.TransactionBlockerReasonRecoveryLock,
				Name:          "tx-test.json",
			},
			{
				Kind:   publish.TransactionBlockerCorruptJournal,
				Reason: "corrupt_journal",
				Name:   "bad.lock.json",
			},
		},
		Lock: lock,
		Journals: []publish.JournalDiagnostic{
			{
				ID:               "tx-test",
				Name:             "tx-test.json",
				Status:           publish.TransactionStatusRollbackFailed,
				Rollback:         publish.RollbackStatusFailed,
				Version:          "v0.1.0",
				StartedAt:        time.Unix(1, 2).UTC(),
				UpdatedAt:        time.Unix(1, 2).UTC(),
				BlocksNewPublish: true,
				AllowsLockClear:  true,
				Path:             "/state/transactions/tx-test.json",
			},
			{
				Name:    "bad.lock.json",
				Corrupt: true,
				Message: "transaction journal bad.lock.json contains unsafe transaction id",
				Path:    "/state/transactions/bad.lock.json",
			},
		},
		Warnings: []publish.TransactionDiagnosticWarning{
			{Code: "journal_corrupt", Message: "transaction journal bad.lock.json is corrupt"},
		},
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
