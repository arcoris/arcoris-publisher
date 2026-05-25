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

package target

import (
	"context"
	"errors"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/ports/git"
	"arcoris.dev/arcoris-publisher/internal/testutil/porttest"
	"arcoris.dev/arcoris-publisher/internal/testutil/publishertest"
)

func TestPrepareRejectsInvalidRequest(t *testing.T) {
	_, err := New(Dependencies{}, Options{}).Prepare(context.Background(), Request{})

	validation, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("error type = %T", err)
	}
	if !validation.Has(IssueInvalidRequest) {
		t.Fatalf("validation issues = %v", validation.Issues)
	}
}

func TestPrepareRejectsDirtyWorktreeWhenRequireClean(t *testing.T) {
	p, err := publishertest.Plan(
		publishertest.PlanOptions{},
		publishertest.Module{Name: "foundation"},
	)
	if err != nil {
		t.Fatalf("publishertest.Plan() error = %v", err)
	}

	fs := porttest.NewFileSystem()
	worktree := "/target/arcoris__foundation"
	fs.AddDir(worktree)
	fakeGit := porttest.NewGit()
	fakeGit.Statuses[worktree] = git.Status{
		Clean:   false,
		Entries: []git.StatusEntry{{Path: "old.txt", Code: "??"}},
	}

	_, err = New(
		Dependencies{FS: fs, Git: fakeGit},
		Options{RequireClean: true},
	).Prepare(context.Background(), Request{Plan: p, RootDir: "/target"})

	validation, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("error type = %T", err)
	}
	if !validation.Has(IssueWorktreeDirty) {
		t.Fatalf("validation issues = %v", validation.Issues)
	}
}

func TestPrepareRejectsStatusError(t *testing.T) {
	p, fs, fakeGit, worktree := targetFixture(t)
	fakeGit.StatusErrors[worktree] = errors.New("not a git repository")

	_, err := New(
		Dependencies{FS: fs, Git: fakeGit},
		Options{},
	).Prepare(context.Background(), Request{Plan: p, RootDir: "/target"})

	validation, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("error type = %T", err)
	}
	if !validation.Has(IssueWorktreeStatusFailed) {
		t.Fatalf("validation issues = %v", validation.Issues)
	}
}

func TestPrepareFetchFailurePolicy(t *testing.T) {
	tests := []struct {
		name string
		opts Options
		want IssueCode
	}{
		{
			name: "required",
			opts: Options{Fetch: true, FetchRequired: true},
			want: IssueFetchFailed,
		},
		{
			name: "best effort",
			opts: Options{Fetch: true},
		},
		{
			name: "disabled",
			opts: Options{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, fs, fakeGit, _ := targetFixture(t)
			fakeGit.FetchError = errors.New("offline")

			_, err := New(Dependencies{FS: fs, Git: fakeGit}, tt.opts).
				Prepare(context.Background(), Request{Plan: p, RootDir: "/target"})

			if tt.want == "" {
				if err != nil {
					t.Fatalf("Prepare() error = %v", err)
				}
				return
			}

			validation, ok := err.(*ValidationError)
			if !ok {
				t.Fatalf("error type = %T", err)
			}
			if !validation.Has(tt.want) {
				t.Fatalf("validation issues = %v", validation.Issues)
			}
		})
	}
}

func TestPrepareCleanGitWorktreeSucceeds(t *testing.T) {
	p, fs, fakeGit, _ := targetFixture(t)

	result, err := New(Dependencies{FS: fs, Git: fakeGit}, Options{}).
		Prepare(context.Background(), Request{Plan: p, RootDir: "/target"})

	if err != nil {
		t.Fatalf("Prepare() error = %v", err)
	}
	if result.Len() != 1 {
		t.Fatalf("Len() = %d, want 1", result.Len())
	}
}

func targetFixture(t *testing.T) (plan.Plan, *porttest.FileSystem, *porttest.Git, string) {
	t.Helper()
	p, err := publishertest.Plan(
		publishertest.PlanOptions{},
		publishertest.Module{Name: "foundation"},
	)
	if err != nil {
		t.Fatalf("publishertest.Plan() error = %v", err)
	}

	fs := porttest.NewFileSystem()
	worktree := "/target/arcoris__foundation"
	fs.AddDir(worktree)
	return p, fs, porttest.NewGit(), worktree
}
