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
	"testing"

	gitport "arcoris.dev/arcoris-publisher/internal/ports/git"
	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

func TestHeadAndCurrentBranch(t *testing.T) {
	runner := &fakeRunner{results: []processport.Result{
		{Stdout: []byte("abc\n")},
		{Stdout: []byte("main\n")},
	}}
	client := New(runner, Options{})

	head, err := client.Head(context.Background(), "/repo")
	if err != nil || head != "abc" {
		t.Fatalf("Head() = %q, %v", head, err)
	}
	branch, err := client.CurrentBranch(context.Background(), "/repo")
	if err != nil || branch != "main" {
		t.Fatalf("CurrentBranch() = %q, %v", branch, err)
	}
}

func TestCurrentBranchDetachedHead(t *testing.T) {
	client := New(&fakeRunner{results: []processport.Result{{Stdout: []byte("\n")}}}, Options{})
	_, err := client.CurrentBranch(context.Background(), "/repo")
	assertPortCode(t, err, gitport.CodeRefNotFound)
}

func TestRefExistsAndRemoteRefExists(t *testing.T) {
	runner := &fakeRunner{results: []processport.Result{
		{ExitCode: 0},
		{ExitCode: 1},
		{ExitCode: 0},
		{ExitCode: 2},
	}}
	client := New(runner, Options{})

	exists, err := client.RefExists(context.Background(), "/repo", "HEAD")
	if err != nil || !exists {
		t.Fatalf("RefExists() = %v, %v", exists, err)
	}
	exists, err = client.RefExists(context.Background(), "/repo", "missing")
	if err != nil || exists {
		t.Fatalf("RefExists() missing = %v, %v", exists, err)
	}
	exists, err = client.RemoteRefExists(context.Background(), "/repo", "", "main")
	if err != nil || !exists {
		t.Fatalf("RemoteRefExists() = %v, %v", exists, err)
	}
	exists, err = client.RemoteRefExists(context.Background(), "/repo", "", "missing")
	if err != nil || exists {
		t.Fatalf("RemoteRefExists() missing = %v, %v", exists, err)
	}
	assertStringSlice(t, runner.specs[2].Args, []string{"ls-remote", "--exit-code", "origin", "main"})
}

func TestCommitMessageDefaultsToHead(t *testing.T) {
	runner := &fakeRunner{results: []processport.Result{{Stdout: []byte("subject\n\nbody\n")}}}
	client := New(runner, Options{})

	message, err := client.CommitMessage(context.Background(), "/repo", "")
	if err != nil || message != "subject\n\nbody\n" {
		t.Fatalf("CommitMessage() = %q, %v", message, err)
	}
	assertStringSlice(t, runner.specs[0].Args, []string{"log", "-1", "--format=%B", "HEAD"})
}
