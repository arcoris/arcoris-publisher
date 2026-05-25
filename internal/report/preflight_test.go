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

package report

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
	"arcoris.dev/arcoris-publisher/internal/testutil/publishertest"
	"arcoris.dev/arcoris-publisher/internal/workflow/preflight"
	"arcoris.dev/arcoris-publisher/internal/workflow/target"
)

func TestPreflightReportHidesLocalPathsByDefault(t *testing.T) {
	result := preflightResultFixture()

	report := BuildPreflightReport(result, Options{})

	if report.Modules[0].WorktreeDir != "" {
		t.Fatalf("WorktreeDir = %q, want hidden", report.Modules[0].WorktreeDir)
	}
	if report.Modules[0].Checks[0].Path != "" {
		t.Fatalf("check path = %q, want hidden", report.Modules[0].Checks[0].Path)
	}
}

func TestPreflightReportIncludesLocalPathsWhenRequested(t *testing.T) {
	result := preflightResultFixture()

	report := BuildPreflightReport(result, Options{IncludeLocalPaths: true})

	if report.Modules[0].WorktreeDir != "/target/arcoris__foundation" {
		t.Fatalf("WorktreeDir = %q", report.Modules[0].WorktreeDir)
	}
}

func TestPreflightTextRendersStatus(t *testing.T) {
	var buf bytes.Buffer
	result := preflightResultFixture()

	if err := New(Options{Format: FormatText}).Preflight(&buf, result); err != nil {
		t.Fatalf("Preflight() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Preflight") || !strings.Contains(output, "Status: passed") {
		t.Fatalf("preflight text output = %q", output)
	}
}

func preflightResultFixture() preflight.Result {
	p, err := publishertest.Plan(
		publishertest.PlanOptions{Version: "v0.1.0"},
		publishertest.Module{Name: "foundation"},
	)
	if err != nil {
		panic(err)
	}

	fs := porttest.NewFileSystem()
	fs.AddDir("/repo")
	fs.AddDir("/repo/staging")
	fs.AddFile("/repo/staging/src/arcoris.dev/foundation/go.mod", []byte("module arcoris.dev/foundation\n"))
	fs.AddFile("/repo/staging/src/arcoris.dev/foundation/contracts/doc.go", []byte("package contracts\n"))
	fs.AddDir("/target")
	worktree := target.RepositoryWorktree("/target", "arcoris/foundation")
	fs.AddDir(worktree)

	fakeGit := porttest.NewGit()
	fakeGit.Refs[worktree+"\x00refs/heads/main"] = true

	result, err := preflight.New(
		preflight.Dependencies{FS: fs, Git: fakeGit},
		preflight.Options{},
	).Check(context.Background(), preflight.Request{
		Plan:                p,
		SourceRepositoryDir: "/repo",
		StagingDir:          "/repo/staging",
		TargetRootDir:       "/target",
	})
	if err != nil {
		panic(err)
	}
	return result
}
