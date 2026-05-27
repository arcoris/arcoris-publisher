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
	"arcoris.dev/arcoris-publisher/internal/workflow/publish"
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

func TestCheckReportsOperationLock(t *testing.T) {
	deps, req, opts := preflightFixture(t)
	if err := os.MkdirAll(opts.StateDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	content := "schemaVersion=1\noperation=publish\ntoken=token-one\npid=1\nstartedAt=2026-01-01T00:00:00Z\n"
	if err := os.WriteFile(filepath.Join(opts.StateDir, "operation.lock"), []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	result, err := New(deps, opts).Check(context.Background(), req)
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	assertGlobalCheckCode(t, result, "publish-lock", StatusFailed, "operation_lock_exists")
	if _, err := os.Stat(filepath.Join(opts.StateDir, "operation.lock")); err != nil {
		t.Fatalf("operation lock missing: %v", err)
	}
}

func TestCheckOperationLockDiagnostics(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantCode string
	}{
		{
			name:     "present",
			content:  "schemaVersion=1\noperation=publish\ntoken=token-one\npid=1\nstartedAt=2026-01-01T00:00:00Z\n",
			wantCode: "operation_lock_exists",
		},
		{
			name:     "corrupt",
			content:  "schemaVersion=1\noperation=publish\npid=1\nstartedAt=2026-01-01T00:00:00Z\n",
			wantCode: "operation_lock_corrupt",
		},
		{
			name:     "read failed",
			wantCode: "operation_lock_read_failed",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			deps, req, opts := preflightFixture(t)
			if err := os.MkdirAll(opts.StateDir, 0o700); err != nil {
				t.Fatalf("MkdirAll() error = %v", err)
			}
			path := filepath.Join(opts.StateDir, "operation.lock")
			if tt.content == "" {
				if err := os.MkdirAll(path, 0o700); err != nil {
					t.Fatalf("MkdirAll(operation.lock) error = %v", err)
				}
			} else if err := os.WriteFile(path, []byte(tt.content), 0o600); err != nil {
				t.Fatalf("WriteFile() error = %v", err)
			}

			result, err := New(deps, opts).Check(context.Background(), req)
			if err != nil {
				t.Fatalf("Check() error = %v", err)
			}
			assertGlobalCheckCode(t, result, "publish-lock", StatusFailed, tt.wantCode)
			if _, err := os.Stat(path); err != nil {
				t.Fatalf("operation lock missing after preflight: %v", err)
			}
		})
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
	assertGlobalCheckCode(t, result, "pending-transactions", StatusFailed, "pending_transaction")
}

func TestCheckFailsRecoveryTransaction(t *testing.T) {
	for _, tt := range []struct {
		name   string
		status string
	}{
		{name: "failed", status: "failed"},
		{name: "rollback failed", status: "rollback_failed"},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			deps, req, opts := preflightFixture(t)
			writePreflightJournalWithID(t, opts.StateDir, "tx-recovery", tt.status)

			result, err := New(deps, opts).Check(context.Background(), req)
			if err != nil {
				t.Fatalf("Check() error = %v", err)
			}
			assertGlobalCheckCode(t, result, "pending-transactions", StatusFailed, "transaction_recovery_required")
		})
	}
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
	assertGlobalCheckCode(t, result, "pending-transactions", StatusFailed, "transaction_journal_corrupt")
}

func TestCheckFailsJournalReadFailures(t *testing.T) {
	t.Run("directory read failure", func(t *testing.T) {
		deps, req, opts := preflightFixture(t)
		if err := os.WriteFile(filepath.Join(opts.StateDir, "transactions"), []byte("not a directory"), 0o600); err != nil {
			t.Fatalf("WriteFile() error = %v", err)
		}

		result, err := New(deps, opts).Check(context.Background(), req)
		if err != nil {
			t.Fatalf("Check() error = %v", err)
		}
		assertGlobalCheckCode(t, result, "pending-transactions", StatusFailed, "transaction_journal_directory_read_failed")
	})

	t.Run("file read failure", func(t *testing.T) {
		deps, req, opts := preflightFixture(t)
		if err := os.MkdirAll(filepath.Join(opts.StateDir, "transactions", "tx-bad.json"), 0o700); err != nil {
			t.Fatalf("MkdirAll() error = %v", err)
		}

		result, err := New(deps, opts).Check(context.Background(), req)
		if err != nil {
			t.Fatalf("Check() error = %v", err)
		}
		assertGlobalCheckCode(t, result, "pending-transactions", StatusFailed, "transaction_journal_file_read_failed")
	})
}

func TestCheckFailsPublishLock(t *testing.T) {
	deps, req, opts := preflightFixture(t)
	if err := os.MkdirAll(opts.StateDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(opts.StateDir, "publish.lock"), []byte("transaction=tx-other\npid=1\nstartedAt=2026-01-01T00:00:00Z\ncommand=publish\n"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	result, err := New(deps, opts).Check(context.Background(), req)
	if err != nil {
		t.Fatalf("Check() error = %v", err)
	}
	assertGlobalCheck(t, result, "publish-lock", StatusFailed)
}

func TestCheckPublishLockDiagnostics(t *testing.T) {
	tests := []struct {
		name       string
		setup      func(t *testing.T, stateDir string)
		wantStatus Status
		wantCode   string
	}{
		{
			name:       "no lock",
			wantStatus: StatusPassed,
		},
		{
			name: "active journal",
			setup: func(t *testing.T, stateDir string) {
				writePreflightJournalWithID(t, stateDir, "tx-active", "pending")
				writePreflightLock(t, stateDir, "tx-active")
			},
			wantStatus: StatusFailed,
			wantCode:   "publish_lock_exists",
		},
		{
			name: "committed journal",
			setup: func(t *testing.T, stateDir string) {
				writePreflightJournalWithID(t, stateDir, "tx-committed", "committed")
				writePreflightLock(t, stateDir, "tx-committed")
			},
			wantStatus: StatusFailed,
			wantCode:   "stale_publish_lock_terminal_transaction",
		},
		{
			name: "rolled back journal",
			setup: func(t *testing.T, stateDir string) {
				writePreflightJournalWithID(t, stateDir, "tx-rolled-back", "rolled_back")
				writePreflightLock(t, stateDir, "tx-rolled-back")
			},
			wantStatus: StatusFailed,
			wantCode:   "stale_publish_lock_terminal_transaction",
		},
		{
			name: "failed journal",
			setup: func(t *testing.T, stateDir string) {
				writePreflightJournalWithID(t, stateDir, "tx-failed", "failed")
				writePreflightLock(t, stateDir, "tx-failed")
			},
			wantStatus: StatusFailed,
			wantCode:   "publish_lock_recovery_required",
		},
		{
			name: "rollback failed journal",
			setup: func(t *testing.T, stateDir string) {
				writePreflightJournalWithID(t, stateDir, "tx-rollback-failed", "rollback_failed")
				writePreflightLock(t, stateDir, "tx-rollback-failed")
			},
			wantStatus: StatusFailed,
			wantCode:   "publish_lock_recovery_required",
		},
		{
			name: "missing journal",
			setup: func(t *testing.T, stateDir string) {
				writePreflightLock(t, stateDir, "tx-missing")
			},
			wantStatus: StatusFailed,
			wantCode:   "stale_publish_lock_journal_missing",
		},
		{
			name: "corrupt lock",
			setup: func(t *testing.T, stateDir string) {
				writePreflightRawLock(t, stateDir, "pid=1\n")
			},
			wantStatus: StatusFailed,
			wantCode:   "publish_lock_corrupt",
		},
		{
			name: "corrupt journal",
			setup: func(t *testing.T, stateDir string) {
				writePreflightLock(t, stateDir, "tx-corrupt")
				txDir := filepath.Join(stateDir, "transactions")
				if err := os.MkdirAll(txDir, 0o700); err != nil {
					t.Fatalf("MkdirAll() error = %v", err)
				}
				if err := os.WriteFile(filepath.Join(txDir, "tx-corrupt.json"), []byte("{"), 0o600); err != nil {
					t.Fatalf("WriteFile() error = %v", err)
				}
			},
			wantStatus: StatusFailed,
			wantCode:   "publish_lock_journal_corrupt",
		},
		{
			name: "lock read failed",
			setup: func(t *testing.T, stateDir string) {
				if err := os.MkdirAll(filepath.Join(stateDir, "publish.lock"), 0o700); err != nil {
					t.Fatalf("MkdirAll() error = %v", err)
				}
			},
			wantStatus: StatusFailed,
			wantCode:   "publish_lock_read_failed",
		},
		{
			name: "lock journal read failed",
			setup: func(t *testing.T, stateDir string) {
				writePreflightLock(t, stateDir, "tx-unreadable")
				if err := os.MkdirAll(filepath.Join(stateDir, "transactions", "tx-unreadable.json"), 0o700); err != nil {
					t.Fatalf("MkdirAll() error = %v", err)
				}
			},
			wantStatus: StatusFailed,
			wantCode:   "publish_lock_journal_read_failed",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			deps, req, opts := preflightFixture(t)
			if tt.setup != nil {
				tt.setup(t, opts.StateDir)
			}

			result, err := New(deps, opts).Check(context.Background(), req)
			if err != nil {
				t.Fatalf("Check() error = %v", err)
			}
			assertGlobalCheckCode(t, result, "publish-lock", tt.wantStatus, tt.wantCode)
		})
	}
}

func TestLockBlockerCheckFailsSafeForUnknownBlocker(t *testing.T) {
	check := lockBlockerCheck(publish.TransactionStateBlocker{Kind: publish.TransactionStateBlockerKind("future_blocker")})
	if check.Status() != StatusFailed || check.Code() != "lock_lookup_failed" {
		t.Fatalf("check = %#v", check)
	}
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
	writePreflightJournalWithFragment(t, stateDir, "tx-pending", statusFragment)
}

func writePreflightJournalWithID(t *testing.T, stateDir string, id string, status string) {
	t.Helper()
	writePreflightJournalWithFragment(t, stateDir, id, `"status":"`+status+`"`)
}

func writePreflightJournalWithFragment(t *testing.T, stateDir string, id string, statusFragment string) {
	t.Helper()
	txDir := filepath.Join(stateDir, "transactions")
	if err := os.MkdirAll(txDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	data := `{"schemaVersion":1,"id":"` + id + `",` + statusFragment + `,"version":"v0.1.0","startedAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z","modules":[]}`
	if err := os.WriteFile(filepath.Join(txDir, id+".json"), []byte(data), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}

func writePreflightLock(t *testing.T, stateDir string, id string) {
	t.Helper()
	writePreflightRawLock(t, stateDir, "transaction="+id+"\npid=1\nstartedAt=2026-01-01T00:00:00Z\ncommand=publish\n")
}

func writePreflightRawLock(t *testing.T, stateDir string, data string) {
	t.Helper()
	if err := os.MkdirAll(stateDir, 0o700); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(stateDir, "publish.lock"), []byte(data), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}

func assertGlobalCheck(t *testing.T, result Result, name string, status Status) {
	t.Helper()
	assertGlobalCheckCode(t, result, name, status, "")
}

func assertGlobalCheckCode(t *testing.T, result Result, name string, status Status, code string) {
	t.Helper()
	for _, check := range result.Checks() {
		if check.Name() == name && check.Status() == status && (code == "" || check.Code() == code) {
			return
		}
	}
	t.Fatalf("global check %s=%s code %s not found: %#v", name, status, code, result.Checks())
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
