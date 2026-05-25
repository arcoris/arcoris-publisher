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

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
	"arcoris.dev/arcoris-publisher/internal/testutil/publishertest"
	"arcoris.dev/arcoris-publisher/internal/workflow/target"
)

func TestTargetPrepareReportHidesLocalPathsByDefault(t *testing.T) {
	result := targetPrepareFixture(t, "file:///tmp/remotes/{name}.git")

	report := BuildTargetPrepareReport(result, Options{})
	if report.TargetRoot != "" || report.Modules[0].WorktreeDir != "" || report.Modules[0].RemoteURL != "" {
		t.Fatalf("report leaked paths: %#v", report)
	}
}

func TestTargetPrepareReportIncludesLocalPathsWhenRequested(t *testing.T) {
	result := targetPrepareFixture(t, "file:///tmp/remotes/{name}.git")

	report := BuildTargetPrepareReport(result, Options{IncludeLocalPaths: true})
	if report.TargetRoot == "" || report.Modules[0].WorktreeDir == "" || report.Modules[0].RemoteURL == "" {
		t.Fatalf("report hid requested paths: %#v", report)
	}
}

func TestTargetPrepareReportRedactsCredentialRemoteURL(t *testing.T) {
	result := targetPrepareFixture(t, "https://token@example.com/{repository}.git")

	report := BuildTargetPrepareReport(result, Options{IncludeLocalPaths: true})
	if strings.Contains(report.Modules[0].RemoteURL, "token") {
		t.Fatalf("remote URL was not redacted: %#v", report.Modules[0].RemoteURL)
	}
	if !strings.Contains(report.Modules[0].RemoteURL, "redacted@example.com") {
		t.Fatalf("remote URL lost useful host/path context: %#v", report.Modules[0].RemoteURL)
	}
}

func TestTargetPrepareTextRendersStatus(t *testing.T) {
	var buf bytes.Buffer
	result := targetPrepareFixture(t, "git@github.com:{repository}.git")

	if err := New(Options{Format: FormatText}).TargetPrepare(&buf, result); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "Target prepare") || !strings.Contains(buf.String(), "foundation") {
		t.Fatalf("text output = %q", buf.String())
	}
}

func targetPrepareFixture(t *testing.T, rawTemplate string) target.PrepareResult {
	t.Helper()
	p, err := publishertest.Plan(
		publishertest.PlanOptions{},
		publishertest.Module{Name: "foundation"},
	)
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}
	fs := porttest.NewFileSystem()
	git := porttest.NewGit()
	worktree := "/tmp/targets/arcoris__foundation"
	git.RemoteRefs[porttest.RemoteRefKeyForRepo(worktree, "origin", "refs/heads/main")] = true
	template, err := manifest.ParseRemoteTemplate(rawTemplate)
	if err != nil {
		t.Fatal(err)
	}
	result, err := target.New(target.Dependencies{FS: fs, Git: git}, target.DefaultOptions()).PrepareTargets(
		context.Background(),
		target.PrepareRequest{Plan: p, RootDir: "/tmp/targets", RemoteTemplate: template, HasRemoteTemplate: true},
	)
	if err != nil {
		t.Fatal(err)
	}
	if result.Failed() {
		t.Fatalf("fixture target prepare failed: %#v", result)
	}
	return result
}
