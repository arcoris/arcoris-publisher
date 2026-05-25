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

	"arcoris.dev/arcoris-publisher/internal/app"
)

func TestRunPreflightPassesWorkflowRequest(t *testing.T) {
	t.Parallel()

	fake := newFakeApplication(t)
	cli := New(Dependencies{App: fake}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{
		"preflight",
		"--manifest", "staging/arcpub.yaml",
		"--version", "v0.3.0",
		"--source-repo", "/repo",
		"--staging-dir", "/repo/staging",
		"--target-root", "/target",
	}, &stdout, &stderr)

	if code != ExitOK {
		t.Fatalf("Run(preflight) code = %d stderr = %s", code, stderr.String())
	}
	if !fake.preflightCalled {
		t.Fatal("Preflight was not called")
	}
	if fake.preflightRequest.ManifestPath != "staging/arcpub.yaml" ||
		fake.preflightRequest.SourceRepositoryDir != "/repo" ||
		fake.preflightRequest.TargetRootDir != "/target" {
		t.Fatalf("preflight request = %+v", fake.preflightRequest)
	}
}

func TestRunPreflightPassesStateDirToFactory(t *testing.T) {
	t.Parallel()

	fake := newFakeApplication(t)
	var got bool
	cli := New(Dependencies{
		AppFactory: func(opts app.Options) (Application, error) {
			got = opts.Workflow.Preflight.StateDir == "/state"
			return fake, nil
		},
	}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{
		"preflight",
		"--version", "v0.3.0",
		"--state-dir", "/state",
	}, &stdout, &stderr)

	if code != ExitOK {
		t.Fatalf("Run(preflight) code = %d stderr = %s", code, stderr.String())
	}
	if !got {
		t.Fatal("preflight state dir was not passed to app factory")
	}
}

func TestRunPreflightMissingVersionIsUsage(t *testing.T) {
	t.Parallel()

	fake := newFakeApplication(t)
	cli := New(Dependencies{App: fake}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{"preflight"}, &stdout, &stderr)

	if code != ExitUsage {
		t.Fatalf("Run(preflight) code = %d stderr = %s", code, stderr.String())
	}
	if fake.preflightCalled {
		t.Fatal("Preflight called for missing version")
	}
}

func TestRunPreflightJSON(t *testing.T) {
	t.Parallel()

	fake := newFakeApplication(t)
	cli := New(Dependencies{App: fake}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{
		"preflight",
		"--version", "v0.3.0",
		"--output", "json",
	}, &stdout, &stderr)

	if code != ExitOK {
		t.Fatalf("Run(preflight json) code = %d stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), `"kind": "preflight"`) {
		t.Fatalf("preflight JSON output = %q", stdout.String())
	}
}
