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

func TestRunVerifyPassesWorkflowRequest(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{
		"verify",
		"--manifest", "staging/arcpub.yaml",
		"--version", "v0.3.0",
		"--source-repo", "/repo",
		"--staging-dir", "/repo/staging",
		"--target-root", "/target",
	}, &stdout, &stderr)

	if code != ExitOK {
		t.Fatalf("Run(verify) code = %d, stderr = %s", code, stderr.String())
	}
	if !app.verifyCalled {
		t.Fatal("Verify was not called")
	}
	if app.verifyRequest.ManifestPath != "staging/arcpub.yaml" || app.verifyRequest.SourceRepositoryDir != "/repo" || app.verifyRequest.TargetRootDir != "/target" {
		t.Fatalf("verify request = %+v", app.verifyRequest)
	}
	if !strings.Contains(stdout.String(), "Workflow") {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

func TestRunVerifyUsesDefaults(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{"verify", "--version", "v0.3.0"}, &stdout, &stderr)
	if code != ExitOK {
		t.Fatalf("Run(verify defaults) code = %d stderr = %s", code, stderr.String())
	}
	if app.verifyRequest.ManifestPath != defaultManifestPath || app.verifyRequest.SourceRepositoryDir != defaultSourceRepositoryDir || app.verifyRequest.TargetRootDir != defaultTargetRootDir {
		t.Fatalf("verify request defaults = %+v", app.verifyRequest)
	}
}

func TestRunVerifyRejectsDryRunFlag(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(
		context.Background(),
		[]string{"verify", "--version", "v0.3.0", "--dry-run"},
		&stdout,
		&stderr,
	)

	if code != ExitUsage {
		t.Fatalf("Run(verify --dry-run) code = %d", code)
	}
	if app.verifyCalled {
		t.Fatal("Verify called for invalid flags")
	}
}

func TestRunVerifyInvalidOutputDoesNotCallApp(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(
		context.Background(),
		[]string{"verify", "--version", "v0.3.0", "--output", "yaml"},
		&stdout,
		&stderr,
	)

	if code != ExitUsage {
		t.Fatalf("Run(verify invalid output) code = %d", code)
	}
	if app.verifyCalled {
		t.Fatal("Verify called for invalid output")
	}
}

func TestRunVerifyAppError(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	app.verifyError = errors.New("boom")
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{"verify", "--version", "v0.3.0"}, &stdout, &stderr)

	if code != ExitError {
		t.Fatalf("Run(verify app error) code = %d", code)
	}
	if !strings.Contains(stderr.String(), string(CodeUseCaseFailed)) {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestRunVerifyFailedReturnsVerificationExit(t *testing.T) {
	t.Parallel()

	cli := New(
		Dependencies{App: newRealApplication(realApplicationOptions{tidyError: errors.New("tidy failed")})},
		Options{},
	)
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), workflowCommandArgs("verify"), &stdout, &stderr)

	if code != ExitVerificationFailed {
		t.Fatalf("Run(verify failed) code = %d stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "verification_failed") {
		t.Fatalf("stdout = %q", stdout.String())
	}
}
