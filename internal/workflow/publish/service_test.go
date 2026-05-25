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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	modulemanifest "arcoris.dev/arcoris-publisher/internal/manifest/module"
	"arcoris.dev/arcoris-publisher/internal/manifest/resolved"
	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/ports/git"
	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
	"arcoris.dev/arcoris-publisher/internal/testutil/publishertest"
	"arcoris.dev/arcoris-publisher/internal/versioning"
	"arcoris.dev/arcoris-publisher/internal/workflow/modulefile"
	"arcoris.dev/arcoris-publisher/internal/workflow/source"
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

func TestBranchRefUsesTargetBranch(t *testing.T) {
	got := branchRef(manifest.BranchName("release/v1"))
	want := "refs/heads/release/v1"

	if got != want {
		t.Fatalf("branchRef() = %q, want %q", got, want)
	}
}

func TestPublishPushesBranchBeforeTag(t *testing.T) {
	req, fakeGit, worktree := publishRequest(t, nil)
	fakeGit.Statuses[worktree] = dirtyStatus()

	result, err := New(Dependencies{Git: fakeGit}, publishOptions(t, Options{})).Publish(context.Background(), req)

	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	if !result.Published() {
		t.Fatal("Published() = false")
	}
	assertCallOrder(t, fakeGit.Calls, "add", "commit", "push", "push", "tag", "push-tag")
}

func TestPublishDoesNotPushTagWhenBranchPushFails(t *testing.T) {
	req, fakeGit, worktree := publishRequest(t, nil)
	fakeGit.Statuses[worktree] = dirtyStatus()
	fakeGit.PushError = errors.New("push failed")

	_, err := New(Dependencies{Git: fakeGit}, publishOptions(t, Options{})).Publish(context.Background(), req)

	if err == nil {
		t.Fatal("Publish() error = nil")
	}
	assertCallAbsent(t, fakeGit.Calls, "tag")
	assertCallAbsent(t, fakeGit.Calls, "push-tag")
}

func TestPublishDryRunDoesNotMutateGit(t *testing.T) {
	req, fakeGit, worktree := publishRequest(t, nil)
	fakeGit.Statuses[worktree] = dirtyStatus()

	result, err := New(Dependencies{Git: fakeGit}, publishOptions(t, Options{DryRun: true})).Publish(context.Background(), req)

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

	result, err := New(Dependencies{Git: fakeGit}, publishOptions(t, Options{})).Publish(context.Background(), req)

	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	if !result.Modules()[0].Skipped() {
		t.Fatal("module was not skipped")
	}
	assertCallAbsent(t, fakeGit.Calls, "add")
}

func TestPublishRejectsMissingSourceModuleBeforeGitMutation(t *testing.T) {
	req, fakeGit, worktree := publishRequest(t, nil)
	req.Source = source.Snapshot{}
	fakeGit.Statuses[worktree] = dirtyStatus()

	_, err := New(Dependencies{Git: fakeGit}, publishOptions(t, Options{})).Publish(context.Background(), req)

	got, ok := err.(*Error)
	if !ok {
		t.Fatalf("error type = %T", err)
	}
	if got.Code != CodeMissingSourceSnapshot {
		t.Fatalf("Code = %q", got.Code)
	}
	assertCallAbsent(t, fakeGit.Calls, "add")
	assertCallAbsent(t, fakeGit.Calls, "commit")
	assertCallAbsent(t, fakeGit.Calls, "push")
	assertCallAbsent(t, fakeGit.Calls, "push-tag")
}

func TestPublishFallsBackToModulefileChangeWhenStatusUnavailable(t *testing.T) {
	req, fakeGit, worktree := publishRequest(t, nil)
	req.ModuleFile = changedModuleFileResult(t)
	fakeGit.StatusErrors[worktree] = errors.New("status unavailable")

	result, err := New(
		Dependencies{Git: fakeGit},
		publishOptions(t, Options{AllowStatusFallback: true}),
	).Publish(context.Background(), req)

	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	if !result.Published() {
		t.Fatal("Published() = false")
	}
}

func TestPublishStatusErrorFailsByDefault(t *testing.T) {
	req, fakeGit, worktree := publishRequest(t, nil)
	fakeGit.StatusErrors[worktree] = errors.New("status unavailable")

	_, err := New(Dependencies{Git: fakeGit}, publishOptions(t, Options{})).Publish(context.Background(), req)

	got, ok := err.(*Error)
	if !ok {
		t.Fatalf("error type = %T", err)
	}
	if got.Code != CodePreflightFailed {
		t.Fatalf("Code = %q", got.Code)
	}
	assertCallAbsent(t, fakeGit.Calls, "add")
}

func TestPublishRejectsMissingCommitIdentityBeforeTransaction(t *testing.T) {
	req, fakeGit, worktree := publishRequest(t, nil)
	fakeGit.Statuses[worktree] = dirtyStatus()
	delete(fakeGit.ConfigValues, worktree+"\x00user.name")
	delete(fakeGit.ConfigValues, worktree+"\x00user.email")
	opts := publishOptions(t, Options{})

	_, err := New(Dependencies{Git: fakeGit}, opts).Publish(context.Background(), req)

	got, ok := err.(*Error)
	if !ok {
		t.Fatalf("error type = %T", err)
	}
	if got.Code != CodePreflightFailed {
		t.Fatalf("Code = %q", got.Code)
	}
	assertCallAbsent(t, fakeGit.Calls, "add")
	assertCallAbsent(t, fakeGit.Calls, "commit")
	assertCallAbsent(t, fakeGit.Calls, "push")
	assertCallAbsent(t, fakeGit.Calls, "push-tag")
	if _, err := os.Stat(filepath.Join(opts.StateDir, "transactions")); !os.IsNotExist(err) {
		t.Fatalf("transactions dir exists after identity failure: %v", err)
	}
	if _, err := os.Stat(filepath.Join(opts.StateDir, "publish.lock")); !os.IsNotExist(err) {
		t.Fatalf("publish lock exists after identity failure: %v", err)
	}
}

func TestPublishUsesCleanGitStatusOverStageResults(t *testing.T) {
	req, fakeGit, worktree := publishRequest(t, nil)
	req.ModuleFile = changedModuleFileResult(t)
	fakeGit.Statuses[worktree] = git.Status{Clean: true}

	result, err := New(Dependencies{Git: fakeGit}, publishOptions(t, Options{})).Publish(context.Background(), req)

	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	if !result.Modules()[0].Skipped() {
		t.Fatal("clean Git status did not skip module")
	}
}

func TestPublishRejectsExistingLocalTagBeforeMutation(t *testing.T) {
	req, fakeGit, worktree := publishRequest(t, nil)
	fakeGit.Statuses[worktree] = dirtyStatus()
	fakeGit.Tags["v0.3.0"] = true

	_, err := New(Dependencies{Git: fakeGit}, publishOptions(t, Options{})).Publish(context.Background(), req)

	got, ok := err.(*Error)
	if !ok {
		t.Fatalf("error type = %T", err)
	}
	if got.Code != CodePreflightFailed {
		t.Fatalf("Code = %q", got.Code)
	}
	assertCallAbsent(t, fakeGit.Calls, "add")
	assertCallAbsent(t, fakeGit.Calls, "push")
}

func TestPublishRejectsExistingRemoteTagBeforeMutation(t *testing.T) {
	req, fakeGit, worktree := publishRequest(t, nil)
	fakeGit.Statuses[worktree] = dirtyStatus()
	fakeGit.RemoteRefs[porttest.RemoteRefKeyForRepo(worktree, "origin", "refs/tags/v0.3.0")] = true

	_, err := New(Dependencies{Git: fakeGit}, publishOptions(t, Options{})).Publish(context.Background(), req)

	got, ok := err.(*Error)
	if !ok {
		t.Fatalf("error type = %T", err)
	}
	if got.Code != CodePreflightFailed {
		t.Fatalf("Code = %q", got.Code)
	}
	assertCallAbsent(t, fakeGit.Calls, "add")
	assertCallAbsent(t, fakeGit.Calls, "push")
}

func TestPublishPreflightsAllModulesBeforeMutation(t *testing.T) {
	req, fakeGit, foundationWorktree := publishRequestForModules(
		t,
		publishertest.Module{Name: "foundation"},
		publishertest.Module{Name: "control"},
	)
	fakeGit.Statuses[foundationWorktree] = dirtyStatus()
	controlWorktree, ok := targetWorktree(req, "control")
	if !ok {
		t.Fatal("control worktree missing")
	}
	fakeGit.Statuses[controlWorktree] = dirtyStatus()
	fakeGit.RemoteRefs[porttest.RemoteRefKeyForRepo(controlWorktree, "origin", "refs/tags/v0.3.0")] = true

	_, err := New(Dependencies{Git: fakeGit}, publishOptions(t, Options{})).Publish(context.Background(), req)

	got, ok := err.(*Error)
	if !ok {
		t.Fatalf("error type = %T", err)
	}
	if got.Code != CodePreflightFailed {
		t.Fatalf("Code = %q", got.Code)
	}
	assertCallAbsent(t, fakeGit.Calls, "add")
	assertCallAbsent(t, fakeGit.Calls, "commit")
	assertCallAbsent(t, fakeGit.Calls, "push")
}

func TestPublishRejectsMultiBranchModule(t *testing.T) {
	req, fakeGit, worktree := publishRequest(t, nil)
	req.Plan = multiBranchPlan(t)
	fakeGit.Statuses[worktree] = dirtyStatus()

	_, err := New(Dependencies{Git: fakeGit}, publishOptions(t, Options{})).Publish(context.Background(), req)

	got, ok := err.(*Error)
	if !ok {
		t.Fatalf("error type = %T", err)
	}
	if got.Code != CodePreflightFailed {
		t.Fatalf("Code = %q", got.Code)
	}
	assertCallAbsent(t, fakeGit.Calls, "add")
}

func TestPublishCommitTrailersIncludeSourceProjectionHash(t *testing.T) {
	req, fakeGit, worktree := publishRequest(t, nil)
	fakeGit.Statuses[worktree] = dirtyStatus()

	_, err := New(Dependencies{Git: fakeGit}, publishOptions(t, Options{})).Publish(context.Background(), req)

	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	message := findCall(fakeGit.Calls, "commit").Ref
	for _, required := range []string{
		"Arcoris-Source-Hash: sha256:",
		"Arcoris-Projection-Hash: sha256:",
	} {
		if !strings.Contains(message, required) {
			t.Fatalf("commit message missing %q:\n%s", required, message)
		}
	}
}

func TestPublishMapsForceWithLease(t *testing.T) {
	pushPolicy := string(manifest.PushPolicyForceWithLease)
	req, fakeGit, worktree := publishRequest(t, &manifest.PublishSpec{
		PushPolicy: &pushPolicy,
	})
	fakeGit.Statuses[worktree] = dirtyStatus()
	fakeGit.RemoteRefHashes[porttest.RemoteRefKeyForRepo(worktree, "origin", "refs/heads/main")] = "base"

	_, err := New(Dependencies{Git: fakeGit}, publishOptions(t, Options{})).Publish(context.Background(), req)

	if err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	push := findPushContaining(fakeGit.Calls, "refs/heads/main")
	if !push.ForceWithLease {
		t.Fatal("branch push did not use force-with-lease")
	}
	if push.ForceWithLeaseRef != "refs/heads/main" || push.ForceWithLeaseExpect == "" {
		t.Fatalf("branch push exact lease = ref %q expect %q", push.ForceWithLeaseRef, push.ForceWithLeaseExpect)
	}
}

func publishRequest(
	t *testing.T,
	publishSpec *manifest.PublishSpec,
) (Request, *porttest.Git, string) {
	t.Helper()

	if publishSpec != nil {
		return publishRequestWithPolicy(t, publishSpec)
	}
	return publishRequestForModules(t, publishertest.Module{Name: "foundation"})
}

func publishOptions(t *testing.T, opts Options) Options {
	t.Helper()
	if opts.StateDir == "" {
		opts.StateDir = t.TempDir()
	}
	if opts.TransactionIDFunc == nil {
		opts.TransactionIDFunc = func(TransactionInput) TransactionID { return TransactionID("tx-test") }
	}
	return opts
}

func publishRequestForModules(
	t *testing.T,
	modules ...publishertest.Module,
) (Request, *porttest.Git, string) {
	t.Helper()

	opts := publishertest.PlanOptions{}
	p, err := publishertest.Plan(opts, modules...)
	if err != nil {
		t.Fatalf("publishertest.Plan() error = %v", err)
	}

	fakeFS := porttest.NewFileSystem()
	for _, mod := range modules {
		fakeFS.AddFile(
			"/repo/staging/src/arcoris.dev/"+mod.Name+"/go.mod",
			[]byte("module arcoris.dev/"+mod.Name+"\n"),
		)
		fakeFS.AddFile(
			"/repo/staging/src/arcoris.dev/"+mod.Name+"/contracts/doc.go",
			[]byte("package contracts\n"),
		)
	}
	fakeFS.AddDir("/target")
	fakeGit := porttest.NewGit()
	snapshot, err := source.New(
		source.Dependencies{FS: fakeFS, Git: fakeGit},
		source.Options{},
	).Inspect(context.Background(), source.Request{
		Plan:          p,
		RepositoryDir: "/repo",
		StagingDir:    "/repo/staging",
	})
	if err != nil {
		t.Fatalf("source.Inspect() error = %v", err)
	}
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
	for _, ws := range targets.Workspaces() {
		setPublishCommitIdentity(fakeGit, ws.WorktreeDir())
	}

	ws, ok := targets.WorkspaceByModule("foundation")
	if !ok {
		t.Fatal("workspace for foundation not found")
	}

	return Request{Plan: p, Source: snapshot, Targets: targets}, fakeGit, ws.WorktreeDir()
}

func setPublishCommitIdentity(fakeGit *porttest.Git, worktree string) {
	fakeGit.ConfigValues[worktree+"\x00user.name"] = "ARCORIS Test"
	fakeGit.ConfigValues[worktree+"\x00user.email"] = "arcoris-test@example.invalid"
}

func publishRequestWithPolicy(
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

	req, fakeGit, worktree := publishRequestForModules(t, publishertest.Module{Name: "foundation"})
	req.Plan = p
	return req, fakeGit, worktree
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

func multiBranchPlan(t *testing.T) plan.Plan {
	t.Helper()

	stagingManifest, err := staging.New(staging.Spec{
		APIVersion: string(manifest.APIVersionV1Alpha1),
		Kind:       string(manifest.KindStagingManifest),
		Metadata:   manifest.MetadataSpec{Name: "arcoris"},
		Source: manifest.SourceSpec{
			Repository:    "arcoris/arcoris",
			DefaultBranch: "main",
		},
		Defaults: staging.DefaultsSpec{
			Branches: []manifest.BranchMappingSpec{
				{Source: "main", Target: "main"},
				{Source: "release", Target: "release"},
			},
		},
		Modules: []staging.ModuleSpec{
			{
				Name:       "foundation",
				SourceDir:  "src/arcoris.dev/foundation",
				Repository: "arcoris/foundation",
			},
		},
	})
	if err != nil {
		t.Fatalf("staging.New() error = %v", err)
	}

	moduleManifest, err := modulemanifest.New(modulemanifest.Spec{
		APIVersion: string(manifest.APIVersionV1Alpha1),
		Kind:       string(manifest.KindModuleManifest),
		Metadata:   manifest.MetadataSpec{Name: "foundation"},
		Module:     manifest.ModuleIdentitySpec{Path: "arcoris.dev/foundation"},
		Publish: modulemanifest.PublishSpec{
			Entries: publishertest.DefaultEntries(),
		},
	})
	if err != nil {
		t.Fatalf("modulemanifest.New() error = %v", err)
	}

	set, err := resolved.Resolve(resolved.ResolveInput{
		Staging: stagingManifest,
		Modules: []modulemanifest.Manifest{
			moduleManifest,
		},
	})
	if err != nil {
		t.Fatalf("resolved.Resolve() error = %v", err)
	}

	p, err := plan.FromPublicationSet(set, versioning.Must("v0.3.0"))
	if err != nil {
		t.Fatalf("plan.FromPublicationSet() error = %v", err)
	}
	return p
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

func findPushContaining(calls []porttest.GitCall, ref string) porttest.GitCall {
	for _, call := range calls {
		if call.Op == "push" && strings.Contains(call.Ref, ref) {
			return call
		}
	}
	return porttest.GitCall{}
}
