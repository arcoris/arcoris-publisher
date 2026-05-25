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
