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

func TestCheckoutAndBranchCommands(t *testing.T) {
	runner := &fakeRunner{}
	client := New(runner, Options{})

	if err := client.Checkout(context.Background(), "/repo", "main", gitport.CheckoutOptions{Force: true, Create: true}); err != nil {
		t.Fatalf("Checkout() error = %v", err)
	}
	if err := client.CreateBranch(context.Background(), "/repo", gitport.BranchName("next"), "HEAD", gitport.CreateBranchOptions{Force: true}); err != nil {
		t.Fatalf("CreateBranch() error = %v", err)
	}
	assertStringSlice(t, runner.specs[0].Args, []string{"checkout", "--force", "-b", "main"})
	assertStringSlice(t, runner.specs[1].Args, []string{"branch", "-f", "next", "HEAD"})
}

func TestCleanFlags(t *testing.T) {
	tests := []struct {
		name string
		opts gitport.CleanOptions
		want []string
	}{
		{name: "noop", opts: gitport.CleanOptions{}, want: nil},
		{name: "untracked", opts: gitport.CleanOptions{RemoveUntracked: true}, want: []string{"clean", "-f"}},
		{name: "ignored only", opts: gitport.CleanOptions{RemoveIgnored: true, Directories: true}, want: []string{"clean", "-fdX"}},
		{name: "all forced", opts: gitport.CleanOptions{RemoveUntracked: true, RemoveIgnored: true, Directories: true, Force: true}, want: []string{"clean", "-ffdx"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := &fakeRunner{}
			client := New(runner, Options{})
			if err := client.Clean(context.Background(), "/repo", tt.opts); err != nil {
				t.Fatalf("Clean() error = %v", err)
			}
			if tt.want == nil {
				if len(runner.specs) != 0 {
					t.Fatalf("Clean() should not run command: %#v", runner.specs)
				}
				return
			}
			assertStringSlice(t, runner.specs[0].Args, tt.want)
		})
	}
}

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
