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

package gitcli

import (
	"context"
	"errors"
	"testing"
	"time"

	gitport "arcoris.dev/arcoris-publisher/internal/ports/git"
	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

func TestCommitRequiresStagedChangesAndDoesNotAddAll(t *testing.T) {
	runner := &fakeRunner{results: []processport.Result{{ExitCode: 0}}}
	client := New(runner, Options{})

	_, err := client.Commit(context.Background(), "/repo", "msg", gitport.CommitOptions{})
	assertPortCode(t, err, gitport.CodeNoChanges)
	if len(runner.specs) != 1 {
		t.Fatalf("Commit() should only check staged diff, ran %#v", runner.specs)
	}
	assertStringSlice(t, runner.specs[0].Args, []string{"diff", "--cached", "--quiet"})
}

func TestCommitBuildsEnvironmentAndReturnsHead(t *testing.T) {
	when := time.Unix(123, 0).UTC()
	runner := &fakeRunner{results: []processport.Result{
		{ExitCode: 1},
		{},
		{Stdout: []byte("abc\n")},
	}}
	client := New(runner, Options{})

	commit, err := client.Commit(context.Background(), "/repo", "msg", gitport.CommitOptions{
		AuthorName:     "Author",
		AuthorEmail:    "author@example.com",
		CommitterName:  "Committer",
		CommitterEmail: "committer@example.com",
		AuthorDate:     when,
		CommitterDate:  when,
	})
	if err != nil || commit != "abc" {
		t.Fatalf("Commit() = %q, %v", commit, err)
	}
	assertStringSlice(t, runner.specs[1].Args, []string{"commit", "-m", "msg"})
	assertStringSlice(t, runner.specs[1].Env, []string{
		"GIT_AUTHOR_NAME=Author",
		"GIT_AUTHOR_EMAIL=author@example.com",
		"GIT_COMMITTER_NAME=Committer",
		"GIT_COMMITTER_EMAIL=committer@example.com",
		"GIT_AUTHOR_DATE=123 +0000",
		"GIT_COMMITTER_DATE=123 +0000",
	})
}

func TestCommitWrapsDiffError(t *testing.T) {
	runner := &fakeRunner{
		results: []processport.Result{{Stderr: []byte("fatal")}},
		errs:    []error{errors.New("failed")},
	}
	client := New(runner, Options{})

	_, err := client.Commit(context.Background(), "/repo", "msg", gitport.CommitOptions{})
	assertPortCode(t, err, gitport.CodeCommandFailed)
}

func TestCommitArgsSupportsAllowEmpty(t *testing.T) {
	args := commitArgs("msg", gitport.CommitOptions{AllowEmpty: true})
	assertStringSlice(t, args, []string{"commit", "-m", "msg", "--allow-empty"})
}
