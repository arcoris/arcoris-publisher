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

package preflight

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/ports/git"
	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
	"arcoris.dev/arcoris-publisher/internal/testutil/publishertest"
	"arcoris.dev/arcoris-publisher/internal/workflow/target"
)

func TestCheckPassesWithoutMutatingGit(t *testing.T) {
	deps, req, opts := preflightFixture(t)

	result, err := New(deps, opts).Check(context.Background(), req)
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if result.Status() != StatusPassed {
		t.Fatalf("status = %q, want passed: %#v", result.Status(), result)
	}
	if len(deps.Git.(*porttest.Git).Calls) != 0 {
		t.Fatalf("preflight recorded mutating Git calls: %#v", deps.Git.(*porttest.Git).Calls)
	}
}

func TestCheckFailsDirtyTarget(t *testing.T) {
	deps, req, opts := preflightFixture(t)
	worktree := target.RepositoryWorktree("/target", "arcoris/foundation")
	deps.Git.(*porttest.Git).Statuses[worktree] = git.Status{
		Clean:   false,
		Entries: []git.StatusEntry{{Path: "go.mod", Code: " M"}},
	}

	result, err := New(deps, opts).Check(context.Background(), req)
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	if result.Status() != StatusFailed {
		t.Fatalf("status = %q, want failed", result.Status())
	}
	assertModuleCheck(t, result, "target-status", StatusFailed)
}

func TestCheckFailsPendingTransaction(t *testing.T) {
	deps, req, opts := preflightFixture(t)
	writePreflightJournal(t, opts.StateDir, `"status":"pending"`)

	result, err := New(deps, opts).Check(context.Background(), req)
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	assertGlobalCheck(t, result, "pending-transactions", StatusFailed)
}

func TestCheckFailsCorruptedJournal(t *testing.T) {
	deps, req, opts := preflightFixture(t)
	txDir := filepath.Join(opts.StateDir, "transactions")
	if err := os.MkdirAll(txDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(txDir, "tx-bad.json"), []byte("{"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	result, err := New(deps, opts).Check(context.Background(), req)
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	assertGlobalCheck(t, result, "pending-transactions", StatusFailed)
}

func TestCheckFailsPublishLock(t *testing.T) {
	deps, req, opts := preflightFixture(t)
	if err := os.MkdirAll(opts.StateDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(opts.StateDir, "publish.lock"), []byte("transaction=tx-other\npid=1\nstartedAt=now\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	result, err := New(deps, opts).Check(context.Background(), req)
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	assertGlobalCheck(t, result, "publish-lock", StatusFailed)
}

func preflightFixture(t *testing.T) (Dependencies, Request, Options) {
	t.Helper()
	p, err := publishertest.Plan(
		publishertest.PlanOptions{Version: "v0.1.0"},
		publishertest.Module{Name: "foundation"},
	)
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
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
	setFakeCommitIdentity(fakeGit, worktree)

	return Dependencies{FS: fs, Git: fakeGit}, Request{
		Plan:                p,
		SourceRepositoryDir: "/repo",
		StagingDir:          "/repo/staging",
		TargetRootDir:       "/target",
	}, Options{StateDir: t.TempDir()}
}

func TestCheckFailsMissingCommitIdentity(t *testing.T) {
	deps, req, opts := preflightFixture(t)
	fakeGit := deps.Git.(*porttest.Git)
	worktree := target.RepositoryWorktree("/target", "arcoris/foundation")
	delete(fakeGit.ConfigValues, worktree+"\x00user.name")
	delete(fakeGit.ConfigValues, worktree+"\x00user.email")

	result, err := New(deps, opts).Check(context.Background(), req)
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	assertModuleCheckCode(t, result, "commit-identity", StatusFailed, target.CommitIdentityCodeMissingBoth)
}

func TestCheckFailsMissingCommitEmail(t *testing.T) {
	deps, req, opts := preflightFixture(t)
	fakeGit := deps.Git.(*porttest.Git)
	worktree := target.RepositoryWorktree("/target", "arcoris/foundation")
	delete(fakeGit.ConfigValues, worktree+"\x00user.email")

	result, err := New(deps, opts).Check(context.Background(), req)
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	assertModuleCheckCode(t, result, "commit-identity", StatusFailed, target.CommitIdentityCodeMissingEmail)
}

func TestCheckFailsMissingCommitName(t *testing.T) {
	deps, req, opts := preflightFixture(t)
	fakeGit := deps.Git.(*porttest.Git)
	worktree := target.RepositoryWorktree("/target", "arcoris/foundation")
	delete(fakeGit.ConfigValues, worktree+"\x00user.name")

	result, err := New(deps, opts).Check(context.Background(), req)
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	assertModuleCheckCode(t, result, "commit-identity", StatusFailed, target.CommitIdentityCodeMissingName)
}

func TestCheckFailsCommitIdentityReadError(t *testing.T) {
	deps, req, opts := preflightFixture(t)
	fakeGit := deps.Git.(*porttest.Git)
	worktree := target.RepositoryWorktree("/target", "arcoris/foundation")
	fakeGit.ConfigErrors[worktree+"\x00user.name"] = os.ErrPermission

	result, err := New(deps, opts).Check(context.Background(), req)
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	assertModuleCheckCode(t, result, "commit-identity", StatusFailed, target.CommitIdentityCodeReadFailed)
}

func setFakeCommitIdentity(fakeGit *porttest.Git, worktree string) {
	fakeGit.ConfigValues[worktree+"\x00user.name"] = "ARCORIS Test"
	fakeGit.ConfigValues[worktree+"\x00user.email"] = "arcoris-test@example.invalid"
}

func writePreflightJournal(t *testing.T, stateDir string, statusFragment string) {
	t.Helper()
	txDir := filepath.Join(stateDir, "transactions")
	if err := os.MkdirAll(txDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	data := `{"schemaVersion":1,"id":"tx-pending",` + statusFragment + `,"version":"v0.1.0","startedAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z","modules":[]}`
	if err := os.WriteFile(filepath.Join(txDir, "tx-pending.json"), []byte(data), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}

func assertGlobalCheck(t *testing.T, result Result, name string, status Status) {
	t.Helper()
	for _, check := range result.Checks() {
		if check.Name() == name && check.Status() == status {
			return
		}
	}
	t.Fatalf("global check %s=%s not found: %#v", name, status, result.Checks())
}

func assertModuleCheck(t *testing.T, result Result, name string, status Status) {
	t.Helper()
	assertModuleCheckCode(t, result, name, status, "")
}

func assertModuleCheckCode(t *testing.T, result Result, name string, status Status, code string) {
	t.Helper()
	for _, mod := range result.Modules() {
		for _, check := range mod.Checks() {
			if check.Name() == name && check.Status() == status && (code == "" || check.Code() == code) {
				return
			}
		}
	}
	t.Fatalf("module check %s=%s code %s not found: %#v", name, status, code, result.Modules())
}
