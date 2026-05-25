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

func TestRunNoArgsPrintsHelp(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	code := New(Dependencies{}, Options{}).Run(context.Background(), nil, &stdout, &stderr)

	if code != ExitOK {
		t.Fatalf("Run(no args) code = %d", code)
	}
	assertRootHelp(t, stdout.String())
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestRunRootHelpFlag(t *testing.T) {
	t.Parallel()

	var stdout, stderr bytes.Buffer
	code := New(Dependencies{}, Options{}).Run(context.Background(), []string{"--help"}, &stdout, &stderr)

	if code != ExitOK {
		t.Fatalf("Run(--help) code = %d", code)
	}
	assertRootHelp(t, stdout.String())
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestRunHelpPlan(t *testing.T) {
	t.Parallel()

	app := newFakeApplication(t)
	var stdout, stderr bytes.Buffer
	code := New(Dependencies{App: app}, Options{}).Run(
		context.Background(),
		[]string{"help", "plan"},
		&stdout,
		&stderr,
	)

	if code != ExitOK {
		t.Fatalf("Run(help plan) code = %d", code)
	}
	if app.buildPlanCalled {
		t.Fatal("BuildPlan called for help plan")
	}
	if !strings.Contains(stdout.String(), "Render the executable publication plan") {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

func assertRootHelp(t *testing.T, output string) {
	t.Helper()
	for _, want := range []string{"plan", "preflight", "target", "verify", "publish", "version", "completion"} {
		if !strings.Contains(output, want) {
			t.Fatalf("help missing %q:\n%s", want, output)
		}
	}
}
