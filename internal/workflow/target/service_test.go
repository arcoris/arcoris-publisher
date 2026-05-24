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
	"testing"

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
