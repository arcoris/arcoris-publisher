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

	gitport "arcoris.dev/arcoris-publisher/internal/ports/git"
	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
)

func TestPushRejectsOptionLikeRemoteBeforeExecution(t *testing.T) {
	runner := &fakeRunner{}
	client := New(runner, Options{})

	err := client.Push(context.Background(), "/repo", "--upload-pack=helper", "main:main", gitport.PushOptions{})
	assertInvalidGitArgument(t, err)
	if len(runner.specs) != 0 {
		t.Fatalf("runner invoked for invalid remote: %#v", runner.specs)
	}
}

func TestPushRejectsOptionLikeRefspecBeforeExecution(t *testing.T) {
	runner := &fakeRunner{}
	client := New(runner, Options{})

	err := client.Push(context.Background(), "/repo", "origin", "--mirror", gitport.PushOptions{})
	assertInvalidGitArgument(t, err)
	if len(runner.specs) != 0 {
		t.Fatalf("runner invoked for invalid refspec: %#v", runner.specs)
	}
}

func TestRemoteRefLookupRejectsOptionLikeRemoteBeforeExecution(t *testing.T) {
	runner := &fakeRunner{}
	client := New(runner, Options{})

	_, err := client.RemoteRefExists(context.Background(), "/repo", "--server-option=x", "refs/heads/main")
	assertInvalidGitArgument(t, err)
	if len(runner.specs) != 0 {
		t.Fatalf("runner invoked for invalid remote: %#v", runner.specs)
	}
}

func TestAddRemoteRejectsOptionLikeURLBeforeExecution(t *testing.T) {
	runner := &fakeRunner{}
	client := New(runner, Options{})

	err := client.AddRemote(context.Background(), "/repo", "origin", "--upload-pack=helper")
	assertInvalidGitArgument(t, err)
	if len(runner.specs) != 0 {
		t.Fatalf("runner invoked for invalid URL: %#v", runner.specs)
	}
}

func assertInvalidGitArgument(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("error = nil")
	}
	var portError *porterr.Error
	if !errors.As(err, &portError) {
		t.Fatalf("error type = %T, want *porterr.Error", err)
	}
	if portError.Code != gitport.CodeInvalidArgument {
		t.Fatalf("error code = %q, want %q", portError.Code, gitport.CodeInvalidArgument)
	}
}
