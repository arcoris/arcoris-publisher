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
	"strings"
	"testing"
)

func TestRunHelp(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	code := New(Dependencies{}, Options{}).Run(context.Background(), []string{"help"}, &stdout, &stderr)

	if code != ExitOK {
		t.Fatalf("Run(help) code = %d", code)
	}
	if !strings.Contains(stdout.String(), "arcpub plan") {
		t.Fatalf("help output = %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestRunSubcommandHelp(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	var stdout, stderr bytes.Buffer
	code := New(Dependencies{App: app}, Options{}).Run(
		context.Background(),
		[]string{"plan", "--help"},
		&stdout,
		&stderr,
	)

	if code != ExitOK {
		t.Fatalf("Run(plan --help) code = %d", code)
	}
	if app.buildPlanCalled {
		t.Fatal("BuildPlan called for help")
	}
	if !strings.Contains(stdout.String(), "arcpub plan") {
		t.Fatalf("stdout = %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestRunUnknownCommand(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	code := New(Dependencies{}, Options{}).Run(context.Background(), []string{"wat"}, &stdout, &stderr)

	if code != ExitUsage {
		t.Fatalf("Run(unknown) code = %d", code)
	}
	if !strings.Contains(stderr.String(), "unknown command") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestRunPlanRequiresApplication(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	code := New(Dependencies{}, Options{}).Run(
		context.Background(),
		[]string{"plan", "--version", "v0.1.0"},
		&stdout,
		&stderr,
	)

	if code != ExitError {
		t.Fatalf("Run(plan without app) code = %d", code)
	}
	if !strings.Contains(stderr.String(), string(CodeMissingApplication)) {
		t.Fatalf("stderr = %q", stderr.String())
	}
}
