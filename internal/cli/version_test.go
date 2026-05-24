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
	"encoding/json"
	"strings"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/buildinfo"
)

func TestRunVersionText(t *testing.T) {
	info := func() buildinfo.Info {
		oldVersion, oldCommit, oldDate, oldDirty := buildinfo.Version, buildinfo.Commit, buildinfo.Date, buildinfo.Dirty
		buildinfo.Version = "v1.2.3"
		buildinfo.Commit = "abc123"
		buildinfo.Date = "2026-05-24T00:00:00Z"
		buildinfo.Dirty = "false"
		out := buildinfo.Current()
		buildinfo.Version, buildinfo.Commit, buildinfo.Date, buildinfo.Dirty = oldVersion, oldCommit, oldDate, oldDirty
		return out
	}

	cli := New(Dependencies{BuildInfo: info}, Options{})
	var stdout, stderr bytes.Buffer
	code := cli.Run(context.Background(), []string{"version"}, &stdout, &stderr)

	if code != ExitOK {
		t.Fatalf("Run(version) code = %d, stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "arcpub v1.2.3") || !strings.Contains(stdout.String(), "  commit: abc123") {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

func TestRunVersionJSON(t *testing.T) {
	t.Parallel()

	cli := New(Dependencies{}, Options{})
	var stdout, stderr bytes.Buffer
	code := cli.Run(context.Background(), []string{"version", "--output", "json"}, &stdout, &stderr)

	if code != ExitOK {
		t.Fatalf("Run(version json) code = %d, stderr = %s", code, stderr.String())
	}
	if !json.Valid(stdout.Bytes()) {
		t.Fatalf("invalid JSON: %s", stdout.String())
	}
	if !strings.HasSuffix(stdout.String(), "\n") {
		t.Fatalf("missing trailing newline: %q", stdout.String())
	}
}

func TestRunVersionDoesNotCallApp(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	cli := New(Dependencies{App: app}, Options{})
	var stdout, stderr bytes.Buffer

	code := cli.Run(context.Background(), []string{"version"}, &stdout, &stderr)

	if code != ExitOK {
		t.Fatalf("Run(version) code = %d stderr = %s", code, stderr.String())
	}
	if app.buildPlanCalled || app.verifyCalled || app.publishCalled {
		t.Fatalf("version command called app: %+v", app)
	}
}

func TestRunVersionInvalidOutput(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	code := New(Dependencies{}, Options{}).Run(
		context.Background(),
		[]string{"version", "--output", "yaml"},
		&stdout,
		&stderr,
	)

	if code != ExitUsage {
		t.Fatalf("Run(version invalid output) code = %d", code)
	}
	if !strings.Contains(stderr.String(), string(CodeInvalidFlags)) {
		t.Fatalf("stderr = %q", stderr.String())
	}
}
