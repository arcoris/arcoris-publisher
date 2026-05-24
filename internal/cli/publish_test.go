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

	"arcoris.dev/arcoris-publisher/internal/app"
)

func TestRunPublishPassesWorkflowRequest(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{
		"publish",
		"--manifest", "staging/arcpub.yaml",
		"--version", "v0.3.0",
		"--source-repo", "/repo",
		"--staging-dir", "/repo/staging",
		"--target-root", "/target",
	}, &stdout, &stderr)

	if code != ExitOK {
		t.Fatalf("Run(publish) code = %d, stderr = %s", code, stderr.String())
	}
	if !app.publishCalled {
		t.Fatal("Publish was not called")
	}
	if app.publishRequest.ManifestPath != "staging/arcpub.yaml" || app.publishRequest.Version.String() != "v0.3.0" {
		t.Fatalf("publish request = %+v", app.publishRequest)
	}
}

func TestRunPublishFactoryReceivesDryRunOption(t *testing.T) {
	t.Parallel()

	fake := newFakeApplication(t)
	var got app.Options
	cli := New(Dependencies{
		AppFactory: func(opts app.Options) (Application, error) {
			got = opts
			return fake, nil
		},
	}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{"publish", "--version", "v0.3.0", "--dry-run"}, &stdout, &stderr)
	if code != ExitOK {
		t.Fatalf("Run(publish --dry-run) code = %d, stderr = %s", code, stderr.String())
	}
	if !got.Workflow.DryRun {
		t.Fatalf("AppFactory options = %+v", got)
	}
}

func TestRunPublishVerificationFailedReturnsVerificationExit(t *testing.T) {
	t.Parallel()

	cli := New(
		Dependencies{App: newRealApplication(realApplicationOptions{tidyError: errors.New("tidy failed")})},
		Options{},
	)
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), workflowCommandArgs("publish"), &stdout, &stderr)

	if code != ExitVerificationFailed {
		t.Fatalf("Run(publish verification failed) code = %d stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "verification_failed") {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

func TestRunPublishSkippedReturnsOK(t *testing.T) {
	t.Parallel()

	cli := New(Dependencies{App: newRealApplication(realApplicationOptions{})}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), workflowCommandArgs("publish"), &stdout, &stderr)

	if code != ExitOK {
		t.Fatalf("Run(publish skipped) code = %d stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Publication: skipped") {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

func TestRunPublishAppError(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	app.publishError = errors.New("boom")
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{"publish", "--version", "v0.3.0"}, &stdout, &stderr)

	if code != ExitError {
		t.Fatalf("Run(publish app error) code = %d", code)
	}
	if !strings.Contains(stderr.String(), string(CodeUseCaseFailed)) {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestRunPublishInvalidFlagsDoNotCallApp(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{"publish", "--bad"}, &stdout, &stderr)

	if code != ExitUsage {
		t.Fatalf("Run(publish bad flag) code = %d", code)
	}
	if app.publishCalled {
		t.Fatal("Publish called for invalid flags")
	}
}
