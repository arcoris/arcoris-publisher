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
	"path/filepath"
	"testing"

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
