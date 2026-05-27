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

package cli

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"arcoris.dev/arcoris-publisher/internal/app"
)

func TestRunTransactionsListUsesStateDir(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{
		"transactions", "list",
		"--state-dir", "/state",
	}, &stdout, &stderr)

	if code != ExitOK {
		t.Fatalf("Run(transactions list) code = %d stderr = %s", code, stderr.String())
	}
	if !app.listCalled || app.transactionRequest.StateDir != "/state" {
		t.Fatalf("transaction request = %+v listCalled=%v", app.transactionRequest, app.listCalled)
	}
}

func TestRunTransactionsListDerivesStateDirFromTargetRoot(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{
		"transactions", "list",
		"--target-root", "/targets",
	}, &stdout, &stderr)

	if code != ExitOK {
		t.Fatalf("Run(transactions list) code = %d stderr = %s", code, stderr.String())
	}
	if want := filepath.Join("/targets", ".arcpub", "state"); app.transactionRequest.StateDir != want {
		t.Fatalf("StateDir = %q, want %q", app.transactionRequest.StateDir, want)
	}
}

func TestRunTransactionsShowRequiresID(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{"transactions", "show"}, &stdout, &stderr)

	if code != ExitUsage {
		t.Fatalf("Run(transactions show) code = %d stderr = %s", code, stderr.String())
	}
	if app.showCalled {
		t.Fatal("ShowTransaction called for missing id")
	}
}

func TestRunTransactionsPrunePassesFilters(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{
		"transactions", "prune",
		"--state-dir", "/state",
		"--status", "committed,rolled_back",
		"--older-than", "30d",
		"--dry-run",
		"--output", "json",
	}, &stdout, &stderr)

	if code != ExitOK {
		t.Fatalf("Run(transactions prune) code = %d stderr = %s", code, stderr.String())
	}
	if !app.pruneCalled {
		t.Fatal("PruneTransactions not called")
	}
	if app.transactionPruneReq.StateDir != "/state" ||
		len(app.transactionPruneReq.Statuses) != 2 ||
		app.transactionPruneReq.OlderThan != 30*24*time.Hour ||
		!app.transactionPruneReq.DryRun {
		t.Fatalf("prune request = %+v", app.transactionPruneReq)
	}
}

func TestRunTransactionsPruneUsageErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
	}{
		{name: "invalid status", args: []string{"transactions", "prune", "--status", "pending", "--dry-run"}},
		{name: "invalid duration", args: []string{"transactions", "prune", "--older-than", "later"}},
		{name: "actual no filter", args: []string{"transactions", "prune"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			app := newFakeApplication(t)
			cli := New(Dependencies{App: app}, Options{})
			var stdout, stderr bytes.Buffer

			code := cli.Run(context.Background(), tt.args, &stdout, &stderr)

			if code != ExitUsage {
				t.Fatalf("Run(%v) code = %d stderr = %s", tt.args, code, stderr.String())
			}
			if app.pruneCalled {
				t.Fatalf("PruneTransactions called for %s", tt.name)
			}
		})
	}
}

func TestRunTransactionsPruneFailureIsUseCaseError(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	app.transactionError = errors.New("prune failed")
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{
		"transactions", "prune",
		"--status", "committed",
	}, &stdout, &stderr)

	if code != ExitError {
		t.Fatalf("Run(transactions prune) code = %d stderr = %s", code, stderr.String())
	}
	if !app.pruneCalled {
		t.Fatal("PruneTransactions not called")
	}
	if stdout.Len() == 0 {
		t.Fatal("expected prune report on stdout")
	}
}

func TestRunTransactionsLockShowPassesStateDir(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{
		"transactions", "lock", "show",
		"--state-dir", "/state",
		"--output", "json",
	}, &stdout, &stderr)

	if code != ExitOK {
		t.Fatalf("Run(transactions lock show) code = %d stderr = %s", code, stderr.String())
	}
	if !app.lockShowCalled || app.transactionLockReq.StateDir != "/state" {
		t.Fatalf("lock request = %+v called=%v", app.transactionLockReq, app.lockShowCalled)
	}
	if !bytes.Contains(stdout.Bytes(), []byte(`"kind": "transactions-lock"`)) {
		t.Fatalf("lock show JSON = %s", stdout.String())
	}
}

func TestRunTransactionsLockClearPassesGuardedRequest(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{
		"transactions", "lock", "clear",
		"--state-dir", "/state",
		"--transaction", " tx-one ",
		"--confirm", " tx-one ",
		"--output", "json",
	}, &stdout, &stderr)

	if code != ExitOK {
		t.Fatalf("Run(transactions lock clear) code = %d stderr = %s", code, stderr.String())
	}
	if !app.lockClearCalled ||
		app.transactionLockReq.StateDir != "/state" ||
		app.transactionLockReq.TransactionID != "tx-one" ||
		app.transactionLockReq.Confirm != "tx-one" {
		t.Fatalf("lock clear request = %+v called=%v", app.transactionLockReq, app.lockClearCalled)
	}
	if !bytes.Contains(stdout.Bytes(), []byte(`"kind": "transactions-lock-clear"`)) {
		t.Fatalf("lock clear JSON = %s", stdout.String())
	}
}

func TestRunTransactionsLockClearUsageErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
	}{
		{name: "missing transaction", args: []string{"transactions", "lock", "clear", "--confirm", "tx-one"}},
		{name: "missing confirm", args: []string{"transactions", "lock", "clear", "--transaction", "tx-one"}},
		{name: "confirm mismatch", args: []string{"transactions", "lock", "clear", "--transaction", "tx-one", "--confirm", "tx-two"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			app := newFakeApplication(t)
			cli := New(Dependencies{App: app}, Options{})
			var stdout, stderr bytes.Buffer

			code := cli.Run(context.Background(), tt.args, &stdout, &stderr)

			if code != ExitUsage {
				t.Fatalf("Run(%v) code = %d stderr = %s", tt.args, code, stderr.String())
			}
			if app.lockClearCalled {
				t.Fatalf("ClearTransactionLock called for %s", tt.name)
			}
		})
	}
}

func TestRunTransactionsDiagnosePassesStateDir(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{
		"transactions", "diagnose",
		"--state-dir", "/state",
		"--output", "json",
	}, &stdout, &stderr)

	if code != ExitOK {
		t.Fatalf("Run(transactions diagnose) code = %d stderr = %s", code, stderr.String())
	}
	if !app.diagnoseCalled || app.transactionRequest.StateDir != "/state" {
		t.Fatalf("diagnose request = %+v called=%v", app.transactionRequest, app.diagnoseCalled)
	}
	if !bytes.Contains(stdout.Bytes(), []byte(`"kind": "transactions-diagnostics"`)) {
		t.Fatalf("diagnostics JSON = %s", stdout.String())
	}
}

func TestRunTransactionsDiagnoseExitSemantics(t *testing.T) {
	t.Parallel()

	t.Run("blockers are successful diagnostics", func(t *testing.T) {
		t.Parallel()
		stateDir := t.TempDir()
		writeCLIOperationLock(t, stateDir)
		cli := New(Dependencies{App: app.New(app.Dependencies{}, app.Options{})}, Options{})
		var stdout, stderr bytes.Buffer

		code := cli.Run(context.Background(), []string{"transactions", "diagnose", "--state-dir", stateDir, "--output", "json"}, &stdout, &stderr)

		if code != ExitOK {
			t.Fatalf("Run(transactions diagnose) code = %d stderr = %s", code, stderr.String())
		}
		if !bytes.Contains(stdout.Bytes(), []byte(`"kind": "operation_lock"`)) {
			t.Fatalf("diagnostics JSON = %s", stdout.String())
		}
	})

	t.Run("partial diagnostics and use case error render report then fail", func(t *testing.T) {
		t.Parallel()
		stateDir := t.TempDir()
		if err := os.WriteFile(filepath.Join(stateDir, "transactions"), []byte("not a directory"), 0o600); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
		cli := New(Dependencies{App: app.New(app.Dependencies{}, app.Options{})}, Options{})
		var stdout, stderr bytes.Buffer

		code := cli.Run(context.Background(), []string{"transactions", "diagnose", "--state-dir", stateDir, "--output", "json"}, &stdout, &stderr)

		if code != ExitError {
			t.Fatalf("Run(transactions diagnose) code = %d stderr = %s", code, stderr.String())
		}
		if !bytes.Contains(stdout.Bytes(), []byte(`"kind": "transactions-diagnostics"`)) ||
			!bytes.Contains(stdout.Bytes(), []byte(`"kind": "journal_directory_read_failed"`)) {
			t.Fatalf("diagnostics JSON = %s", stdout.String())
		}
		if !strings.Contains(stderr.String(), string(CodeUseCaseFailed)) {
			t.Fatalf("stderr = %s", stderr.String())
		}
	})

	t.Run("report render failure wins over use case error", func(t *testing.T) {
		t.Parallel()
		stateDir := t.TempDir()
		if err := os.WriteFile(filepath.Join(stateDir, "transactions"), []byte("not a directory"), 0o600); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}
		cli := New(Dependencies{App: app.New(app.Dependencies{}, app.Options{})}, Options{})
		var stderr bytes.Buffer

		code := cli.Run(context.Background(), []string{"transactions", "diagnose", "--state-dir", stateDir}, failingWriter{}, &stderr)

		if code != ExitError {
			t.Fatalf("Run(transactions diagnose) code = %d stderr = %s", code, stderr.String())
		}
		if !strings.Contains(stderr.String(), string(CodeReportFailed)) {
			t.Fatalf("stderr = %s", stderr.String())
		}
	})

	t.Run("context cancellation is command failure with partial report", func(t *testing.T) {
		t.Parallel()
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		cli := New(Dependencies{App: app.New(app.Dependencies{}, app.Options{})}, Options{})
		var stdout, stderr bytes.Buffer

		code := cli.Run(ctx, []string{"transactions", "diagnose", "--state-dir", t.TempDir(), "--output", "json"}, &stdout, &stderr)

		if code != ExitError {
			t.Fatalf("Run(transactions diagnose) code = %d stderr = %s", code, stderr.String())
		}
		if !bytes.Contains(stdout.Bytes(), []byte(`"reason": "context_canceled"`)) {
			t.Fatalf("diagnostics JSON = %s", stdout.String())
		}
		if !strings.Contains(stderr.String(), string(CodeUseCaseFailed)) {
			t.Fatalf("stderr = %s", stderr.String())
		}
	})
}

func writeCLIOperationLock(t *testing.T, stateDir string) {
	t.Helper()
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	content := "schemaVersion=1\noperation=publish\ntoken=token-one\npid=1\nstartedAt=2026-01-01T00:00:00Z\n"
	if err := os.WriteFile(filepath.Join(stateDir, "operation.lock"), []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}

func TestRunTransactionsLockFailuresRenderReport(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	app.transactionError = errors.New("lock failed")
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{
		"transactions", "lock", "clear",
		"--transaction", "tx-one",
		"--confirm", "tx-one",
	}, &stdout, &stderr)

	if code != ExitError {
		t.Fatalf("Run(transactions lock clear) code = %d stderr = %s", code, stderr.String())
	}
	if !app.lockClearCalled {
		t.Fatal("ClearTransactionLock not called")
	}
	if stdout.Len() == 0 {
		t.Fatal("expected lock clear report on stdout")
	}
}

func TestRunRollbackRequiresTransaction(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{"rollback"}, &stdout, &stderr)

	if code != ExitUsage {
		t.Fatalf("Run(rollback) code = %d stderr = %s", code, stderr.String())
	}
	if app.rollbackCalled {
		t.Fatal("RollbackTransaction called for missing id")
	}
}

func TestRunPublishPassesStateDirAndRollbackMode(t *testing.T) {
	t.Parallel()

	fake := newFakeApplication(t)
	var got bool
	cli := New(Dependencies{
		AppFactory: func(opts app.Options) (Application, error) {
			got = opts.Workflow.Publish.StateDir == "/state" &&
				opts.Workflow.Publish.RollbackMode == app.RollbackManual
			return fake, nil
		},
	}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{
		"publish",
		"--version", "v0.3.0",
		"--state-dir", "/state",
		"--rollback", "manual",
	}, &stdout, &stderr)

	if code != ExitOK {
		t.Fatalf("Run(publish) code = %d stderr = %s", code, stderr.String())
	}
	if !got {
		t.Fatal("publish state dir or rollback mode was not passed to app factory")
	}
}

func TestRunPublishInvalidRollbackModeIsUsage(t *testing.T) {
	t.Parallel()

	fake := newFakeApplication(t)
	cli := New(Dependencies{App: fake}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{
		"publish",
		"--version", "v0.3.0",
		"--rollback", "banana",
	}, &stdout, &stderr)

	if code != ExitUsage {
		t.Fatalf("Run(publish) code = %d stderr = %s", code, stderr.String())
	}
	if fake.publishCalled {
		t.Fatal("Publish called for invalid rollback mode")
	}
}
