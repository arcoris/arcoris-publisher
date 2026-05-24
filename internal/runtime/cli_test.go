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

package runtime

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/buildinfo"
	"arcoris.dev/arcoris-publisher/internal/cli"
)

func TestCLIRunsHelp(t *testing.T) {
	t.Parallel()

	rt := New(Options{})
	var stdout, stderr bytes.Buffer
	code := rt.CLI().Run(context.Background(), []string{"help"}, &stdout, &stderr)
	if code != cli.ExitOK {
		t.Fatalf("CLI help exit code = %d", code)
	}
	if !strings.Contains(stdout.String(), "plan") ||
		!strings.Contains(stdout.String(), "completion") {
		t.Fatalf("help output = %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestCLIRunsVersion(t *testing.T) {
	t.Parallel()

	rt := New(Options{})
	var stdout, stderr bytes.Buffer
	code := rt.CLI().Run(context.Background(), []string{"version"}, &stdout, &stderr)
	if code != cli.ExitOK {
		t.Fatalf("CLI version exit code = %d", code)
	}
	if !strings.Contains(stdout.String(), "arcpub ") {
		t.Fatalf("version output = %q", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestCLIDependenciesUseApplicationFactory(t *testing.T) {
	t.Parallel()

	rt := New(Options{})
	deps := rt.CLIDependencies()
	if deps.AppFactory == nil {
		t.Fatal("AppFactory is nil")
	}
	if deps.App != nil {
		t.Fatal("runtime CLI should use an AppFactory, not a preconstructed app")
	}
}

func TestCLIUsesInjectedBuildInfo(t *testing.T) {
	oldVersion, oldCommit, oldDate, oldDirty := buildinfo.Version, buildinfo.Commit, buildinfo.Date, buildinfo.Dirty
	t.Cleanup(func() {
		buildinfo.Version = oldVersion
		buildinfo.Commit = oldCommit
		buildinfo.Date = oldDate
		buildinfo.Dirty = oldDirty
	})

	called := false
	buildInfo := func() buildinfo.Info {
		called = true
		buildinfo.Version = "v9.9.9-runtime"
		buildinfo.Commit = "runtime-commit"
		buildinfo.Date = "2026-05-24"
		buildinfo.Dirty = "false"
		return buildinfo.Current()
	}

	rt := NewWithDependencies(Dependencies{BuildInfo: buildInfo}, Options{})
	var stdout, stderr bytes.Buffer
	code := rt.CLI().Run(context.Background(), []string{"version"}, &stdout, &stderr)
	if code != cli.ExitOK {
		t.Fatalf("CLI version exit code = %d, stderr = %q", code, stderr.String())
	}
	if !called {
		t.Fatal("BuildInfo dependency was not called")
	}
	if !strings.Contains(stdout.String(), "v9.9.9-runtime") {
		t.Fatalf("version output = %q", stdout.String())
	}
}
