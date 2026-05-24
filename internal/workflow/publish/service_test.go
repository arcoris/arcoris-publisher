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

package publish

import (
	"context"
	"errors"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/ports/git"
	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
	"arcoris.dev/arcoris-publisher/internal/testutil/publishertest"
	"arcoris.dev/arcoris-publisher/internal/workflow/modulefile"
	"arcoris.dev/arcoris-publisher/internal/workflow/target"
)

func TestPublishRejectsInvalidRequest(t *testing.T) {
	_, err := New(Dependencies{}, Options{}).Publish(context.Background(), Request{})

	got, ok := err.(*Error)
	if !ok {
		t.Fatalf("error type = %T", err)
	}
	if got.Code != CodeInvalidRequest {
		t.Fatalf("Code = %q", got.Code)
	}
}

func TestBranchRefspecUsesTargetBranch(t *testing.T) {
	got := branchRefspec(manifest.BranchName("release/v1"))
	want := git.RefSpec("refs/heads/release/v1:refs/heads/release/v1")

	if got != want {
		t.Fatalf("branchRefspec() = %q, want %q", got, want)
	}
}

func TestPublishPushesBranchBeforeTag(t *testing.T) {
	req, fakeGit, worktree := publishRequest(t, nil)
	fakeGit.Statuses[worktree] = dirtyStatus()

	result, err := New(Dependencies{Git: fakeGit}, Options{}).Publish(context.Background(), req)

	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	if !result.Published() {
		t.Fatal("Published() = false")
	}
	assertCallOrder(t, fakeGit.Calls, "add", "commit", "push", "tag", "push-tag")
}

func TestPublishDoesNotPushTagWhenBranchPushFails(t *testing.T) {
	req, fakeGit, worktree := publishRequest(t, nil)
	fakeGit.Statuses[worktree] = dirtyStatus()
	fakeGit.PushError = errors.New("push failed")

	_, err := New(Dependencies{Git: fakeGit}, Options{}).Publish(context.Background(), req)

	if err == nil {
		t.Fatal("Publish() error = nil")
	}
	assertCallAbsent(t, fakeGit.Calls, "tag")
	assertCallAbsent(t, fakeGit.Calls, "push-tag")
}

func TestPublishDryRunDoesNotMutateGit(t *testing.T) {
	req, fakeGit, worktree := publishRequest(t, nil)
	fakeGit.Statuses[worktree] = dirtyStatus()

	result, err := New(Dependencies{Git: fakeGit}, Options{DryRun: true}).Publish(context.Background(), req)

	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	if result.Modules()[0].Skipped() {
		t.Fatal("dry run result was marked skipped")
	}
	assertCallAbsent(t, fakeGit.Calls, "add")
	assertCallAbsent(t, fakeGit.Calls, "push")
	assertCallAbsent(t, fakeGit.Calls, "push-tag")
}

func TestPublishSkipsCleanWorktree(t *testing.T) {
	req, fakeGit, worktree := publishRequest(t, nil)
	fakeGit.Statuses[worktree] = git.Status{Clean: true}

	result, err := New(Dependencies{Git: fakeGit}, Options{}).Publish(context.Background(), req)

	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	if !result.Modules()[0].Skipped() {
		t.Fatal("module was not skipped")
	}
	assertCallAbsent(t, fakeGit.Calls, "add")
}

func TestPublishFallsBackToModulefileChangeWhenStatusUnavailable(t *testing.T) {
	req, fakeGit, worktree := publishRequest(t, nil)
	req.ModuleFile = changedModuleFileResult(t)
	fakeGit.StatusErrors[worktree] = errors.New("status unavailable")

	result, err := New(Dependencies{Git: fakeGit}, Options{}).Publish(context.Background(), req)

	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	if !result.Published() {
		t.Fatal("Published() = false")
	}
}

func TestPublishUsesCleanGitStatusOverStageResults(t *testing.T) {
	req, fakeGit, worktree := publishRequest(t, nil)
	req.ModuleFile = changedModuleFileResult(t)
	fakeGit.Statuses[worktree] = git.Status{Clean: true}

	result, err := New(Dependencies{Git: fakeGit}, Options{}).Publish(context.Background(), req)

	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	if !result.Modules()[0].Skipped() {
		t.Fatal("clean Git status did not skip module")
	}
}

func TestPublishMapsForceWithLease(t *testing.T) {
	pushPolicy := string(manifest.PushPolicyForceWithLease)
	req, fakeGit, worktree := publishRequest(t, &manifest.PublishSpec{
		PushPolicy: &pushPolicy,
	})
	fakeGit.Statuses[worktree] = dirtyStatus()

	_, err := New(Dependencies{Git: fakeGit}, Options{}).Publish(context.Background(), req)

	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	push := findCall(fakeGit.Calls, "push")
	if !push.ForceWithLease {
		t.Fatal("branch push did not use force-with-lease")
	}
}

func publishRequest(
	t *testing.T,
	publishSpec *manifest.PublishSpec,
) (Request, *porttest.Git, string) {
	t.Helper()

	opts := publishertest.PlanOptions{}
	if publishSpec != nil {
		opts.Publish = *publishSpec
	}
	p, err := publishertest.Plan(opts, publishertest.Module{Name: "foundation"})
	if err != nil {
		t.Fatalf("publishertest.Plan() error = %v", err)
	}

	fakeFS := porttest.NewFileSystem()
	fakeFS.AddDir("/target")
	fakeGit := porttest.NewGit()
	targets, err := target.New(
		target.Dependencies{FS: fakeFS, Git: fakeGit},
		target.Options{CreateMissing: true},
	).Prepare(context.Background(), target.Request{
		Plan:    p,
		RootDir: "/target",
	})
	if err != nil {
		t.Fatalf("target.Prepare() error = %v", err)
	}

	ws, ok := targets.WorkspaceByModule("foundation")
	if !ok {
		t.Fatal("workspace for foundation not found")
	}

	return Request{Plan: p, Targets: targets}, fakeGit, ws.WorktreeDir()
}

func changedModuleFileResult(t *testing.T) modulefile.Result {
	t.Helper()

	p, err := publishertest.Plan(
		publishertest.PlanOptions{},
		publishertest.Module{Name: "foundation"},
	)
	if err != nil {
		t.Fatalf("publishertest.Plan() error = %v", err)
	}
	fakeFS := porttest.NewFileSystem()
	fakeGit := porttest.NewGit()
	targets, err := target.New(
		target.Dependencies{FS: fakeFS, Git: fakeGit},
		target.Options{CreateMissing: true},
	).Prepare(context.Background(), target.Request{
		Plan:    p,
		RootDir: "/target",
	})
	if err != nil {
		t.Fatalf("target.Prepare() error = %v", err)
	}

	ws, _ := targets.WorkspaceByModule("foundation")
	fakeFS.AddFile(ws.WorktreeDir()+"/go.mod", []byte("module wrong.example/foundation\n"))
	result, err := modulefile.New(
		modulefile.Dependencies{FS: fakeFS},
		modulefile.Options{},
	).Rewrite(context.Background(), modulefile.Request{
		Plan:    p,
		Targets: targets,
	})
	if err != nil {
		t.Fatalf("modulefile.Rewrite() error = %v", err)
	}

	return result
}

func dirtyStatus() git.Status {
	return git.Status{
		Clean:   false,
		Entries: []git.StatusEntry{{Path: "go.mod", Code: " M"}},
	}
}

func assertCallOrder(t *testing.T, calls []porttest.GitCall, want ...string) {
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
			t.Fatalf("call[%d] = %q, want %q; all calls = %v", i, got[i], want[i], got)
		}
	}
}

func assertCallAbsent(t *testing.T, calls []porttest.GitCall, op string) {
	t.Helper()
	if call := findCall(calls, op); call.Op != "" {
		t.Fatalf("unexpected %s call: %#v", op, call)
	}
}

func findCall(calls []porttest.GitCall, op string) porttest.GitCall {
	for _, call := range calls {
		if call.Op == op {
			return call
		}
	}
	return porttest.GitCall{}
}
