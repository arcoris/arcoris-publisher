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
	"strings"
	"testing"
)

func TestRunPlan(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{
		"plan",
		"--manifest", "staging/arcpub.yaml",
		"--version", "v0.3.0",
	}, &stdout, &stderr)

	if code != ExitOK {
		t.Fatalf("Run(plan) code = %d, stderr = %s", code, stderr.String())
	}
	if !app.buildPlanCalled || app.buildPlanManifest != "staging/arcpub.yaml" || app.buildPlanVersion.String() != "v0.3.0" {
		t.Fatalf("BuildPlan call = manifest %q version %q called %t", app.buildPlanManifest, app.buildPlanVersion, app.buildPlanCalled)
	}
	if !strings.Contains(stdout.String(), "Plan") || !strings.Contains(stdout.String(), "foundation") {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

func TestRunPlanInvalidVersion(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{"plan", "--version", "1.2.3"}, &stdout, &stderr)

	if code != ExitUsage {
		t.Fatalf("Run(plan invalid version) code = %d", code)
	}
	if app.buildPlanCalled {
		t.Fatal("BuildPlan called for invalid version")
	}
	if !strings.Contains(stderr.String(), string(CodeInvalidVersion)) {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestRunPlanMissingVersionDoesNotCallApp(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{"plan"}, &stdout, &stderr)

	if code != ExitUsage {
		t.Fatalf("Run(plan missing version) code = %d", code)
	}
	if app.buildPlanCalled {
		t.Fatal("BuildPlan called without version")
	}
}

func TestRunPlanAppError(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	app.buildPlanError = errors.New("boom")
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{"plan", "--version", "v0.3.0"}, &stdout, &stderr)

	if code != ExitError {
		t.Fatalf("Run(plan app error) code = %d", code)
	}
	if !strings.Contains(stderr.String(), string(CodeUseCaseFailed)) {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestRunPlanRenderError(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	cli := New(Dependencies{App: app}, Options{})
	var stderr bytes.Buffer

	code := cli.Run(
		context.Background(),
		[]string{"plan", "--version", "v0.3.0"},
		failingWriter{},
		&stderr,
	)

	if code != ExitError {
		t.Fatalf("Run(plan render error) code = %d", code)
	}
	if !strings.Contains(stderr.String(), string(CodeReportFailed)) {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestRunPlanJSON(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(
		context.Background(),
		[]string{"plan", "--version", "v0.3.0", "--output", "json"},
		&stdout,
		&stderr,
	)

	if code != ExitOK {
		t.Fatalf("Run(plan json) code = %d stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"kind": "plan"`) {
		t.Fatalf("stdout = %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "\n  ") {
		t.Fatalf("JSON output is not pretty by default: %q", stdout.String())
	}
}

func TestRunPlanCompactJSON(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(
		context.Background(),
		[]string{"plan", "--version", "v0.3.0", "--output", "json", "--compact"},
		&stdout,
		&stderr,
	)

	if code != ExitOK {
		t.Fatalf("Run(plan compact json) code = %d stderr = %s", code, stderr.String())
	}
	if strings.Contains(stdout.String(), "\n  ") {
		t.Fatalf("compact JSON output is indented: %q", stdout.String())
	}
}

func TestRunPlanRejectsPrettyCompactConflict(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(
		context.Background(),
		[]string{"plan", "--version", "v0.3.0", "--output", "json", "--pretty", "--compact"},
		&stdout,
		&stderr,
	)

	if code != ExitUsage {
		t.Fatalf("Run(plan pretty compact) code = %d", code)
	}
	if app.buildPlanCalled {
		t.Fatal("BuildPlan called for conflicting output flags")
	}
}
