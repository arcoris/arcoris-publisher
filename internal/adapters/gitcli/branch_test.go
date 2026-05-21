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

func TestCreateBranchBuildsCommand(t *testing.T) {
	runner := &fakeRunner{}
	client := New(runner, Options{})

	if err := client.CreateBranch(context.Background(), "/repo", gitport.BranchName("next"), "HEAD", gitport.CreateBranchOptions{Force: true}); err != nil {
		t.Fatalf("CreateBranch() error = %v", err)
	}
	assertStringSlice(t, runner.specs[0].Args, []string{"branch", "-f", "next", "HEAD"})
}

func TestCreateBranchArgsOmitEmptyStartPoint(t *testing.T) {
	assertStringSlice(t, createBranchArgs(gitport.BranchName("next"), "", gitport.CreateBranchOptions{}), []string{"branch", "next"})
}
