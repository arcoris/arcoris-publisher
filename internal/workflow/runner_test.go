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

package workflow

import (
	"context"
	"errors"
	"strings"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/ports/git"
	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
	"arcoris.dev/arcoris-publisher/internal/testutil/publishertest"
	"arcoris.dev/arcoris-publisher/internal/workflow/construct"
	"arcoris.dev/arcoris-publisher/internal/workflow/modulefile"
	"arcoris.dev/arcoris-publisher/internal/workflow/publish"
	"arcoris.dev/arcoris-publisher/internal/workflow/source"
	"arcoris.dev/arcoris-publisher/internal/workflow/target"
	"arcoris.dev/arcoris-publisher/internal/workflow/verify"
)

func TestRunnerVerifyOnlyRun(t *testing.T) {
	runner, req, fakeGit := workflowFixture(t)

	result, err := runner.Run(context.Background(), req)

	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Source().Repository().Head() == "" {
		t.Fatal("source snapshot is empty")
	}
	if result.Verify().Failed() {
		t.Fatal("verification failed")
	}
	if result.Publish().Published() {
		t.Fatal("publish stage ran for verify-only request")
	}
	assertNoWorkflowCall(t, fakeGit.Calls, "push")
}

func TestRunnerPublishRun(t *testing.T) {
	runner, req, fakeGit := workflowFixture(t)
	req.Publish = true
	fakeGit.Statuses["/target/arcoris__foundation"] = dirtyStatus()

	result, err := runner.Run(context.Background(), req)

	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !result.Publish().Published() {
		t.Fatal("publish stage did not publish")
	}
	assertWorkflowCallOrder(t, fakeGit.Calls, "add", "commit", "push", "tag", "push-tag")
}

func TestRunnerVerificationFailurePreventsPublish(t *testing.T) {
	runner, req, fakeGit := workflowFixture(t)
	req.Publish = true
	fakeGit.Statuses["/target/arcoris__foundation"] = dirtyStatus()
	runner.deps.Verify.Go = porttest.GoToolchain{ModTidyError: errors.New("tidy failed")}

	result, err := runner.Run(context.Background(), req)

	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if !result.Verify().Failed() {
		t.Fatal("verification did not fail")
	}
	if result.Publish().Published() {
		t.Fatal("publish stage ran after verification failure")
	}
	assertNoWorkflowCall(t, fakeGit.Calls, "push")
}

func TestRunnerDryRunDoesNotPublish(t *testing.T) {
	runner, req, fakeGit := workflowFixture(t)
	req.Publish = true
	runner.opts.DryRun = true
	fakeGit.Statuses["/target/arcoris__foundation"] = dirtyStatus()

	result, err := runner.Run(context.Background(), req)

	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Publish().Published() {
		t.Fatal("dry-run published")
	}
	assertNoWorkflowCall(t, fakeGit.Calls, "push")
}

func TestRunnerWrapsSourceErrors(t *testing.T) {
	runner, req, _ := workflowFixture(t)
	runner.deps.Source.Git = nil

	_, err := runner.Run(context.Background(), req)

	if err == nil || !strings.Contains(err.Error(), "source:") {
		t.Fatalf("Run() error = %v", err)
	}
}

func workflowFixture(t *testing.T) (Runner, Request, *porttest.Git) {
	t.Helper()

	p, err := publishertest.Plan(
		publishertest.PlanOptions{},
		publishertest.Module{Name: "foundation"},
	)
	if err != nil {
		t.Fatalf("publishertest.Plan() error = %v", err)
	}

	fs := workflowFS()
	fakeGit := porttest.NewGit()
	deps := Dependencies{
		Source:     source.Dependencies{FS: fs, Git: fakeGit},
		Target:     target.Dependencies{FS: fs, Git: fakeGit},
		Construct:  construct.Dependencies{FS: fs},
		ModuleFile: modulefile.Dependencies{FS: fs},
		Verify:     verify.Dependencies{FS: fs, Go: porttest.GoToolchain{}},
		Publish:    publish.Dependencies{Git: fakeGit},
	}
	opts := Options{
		Target: target.Options{
			CreateMissing: true,
			RequireClean:  false,
		},
		Construct: construct.Options{PreserveGitDir: true},
		Verify:    verify.Options{RequireClean: false},
	}

	return New(deps, opts), Request{
		Plan:                p,
		SourceRepositoryDir: "/repo",
		StagingDir:          "/repo/staging",
		TargetRootDir:       "/target",
	}, fakeGit
}

func workflowFS() *porttest.FileSystem {
	fs := porttest.NewFileSystem()
	fs.AddFile(
		"/repo/staging/src/arcoris.dev/foundation/go.mod",
		[]byte("module arcoris.dev/foundation\n"),
	)
	fs.AddFile(
		"/repo/staging/src/arcoris.dev/foundation/contracts/doc.go",
		[]byte("package contracts\n"),
	)
	return fs
}

func dirtyStatus() git.Status {
	return git.Status{
		Clean:   false,
		Entries: []git.StatusEntry{{Path: "go.mod", Code: " M"}},
	}
}

func assertWorkflowCallOrder(t *testing.T, calls []porttest.GitCall, want ...string) {
	t.Helper()

	var got []string
	for _, call := range calls {
		switch call.Op {
		case "add", "commit", "push", "tag", "push-tag":
			got = append(got, call.Op)
		}
	}
	if len(got) != len(want) {
		t.Fatalf("calls = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("call[%d] = %q, want %q; calls = %v", i, got[i], want[i], got)
		}
	}
}

func assertNoWorkflowCall(t *testing.T, calls []porttest.GitCall, op string) {
	t.Helper()
	for _, call := range calls {
		if call.Op == op {
			t.Fatalf("unexpected %s call: %#v", op, call)
		}
	}
}
