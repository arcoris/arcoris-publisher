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
)

func TestPushBuildsCommand(t *testing.T) {
	runner := &fakeRunner{}
	client := New(runner, Options{})

	err := client.Push(
		context.Background(),
		"/repo",
		"",
		gitport.RefSpec("main:main"),
		gitport.PushOptions{ForceWithLease: true, Atomic: true},
	)
	if err != nil {
		t.Fatalf("Push() error = %v", err)
	}
	assertStringSlice(t, runner.specs[0].Args, []string{
		"push",
		"--force-with-lease",
		"--atomic",
		"origin",
		"main:main",
	})
}

func TestPushArgsSupportsForce(t *testing.T) {
	args := pushArgs("upstream", gitport.RefSpec("main"), gitport.PushOptions{Force: true})
	assertStringSlice(t, args, []string{"push", "--force", "upstream", "main"})
}

func TestPushArgsSupportsExactForceWithLease(t *testing.T) {
	args := pushArgs("origin", gitport.RefSpec("new:refs/heads/main"), gitport.PushOptions{
		ForceWithLeaseRef:    "refs/heads/main",
		ForceWithLeaseExpect: "old",
	})
	assertStringSlice(t, args, []string{"push", "--force-with-lease=refs/heads/main:old", "origin", "new:refs/heads/main"})
}

func TestDeleteRemoteRefBuildsDeleteRefspec(t *testing.T) {
	runner := &fakeRunner{}
	client := New(runner, Options{})

	err := client.DeleteRemoteRef(context.Background(), "/repo", "origin", "refs/heads/main", gitport.PushOptions{})
	if err != nil {
		t.Fatalf("DeleteRemoteRef() error = %v", err)
	}
	assertStringSlice(t, runner.specs[0].Args, []string{"push", "origin", ":refs/heads/main"})
}
