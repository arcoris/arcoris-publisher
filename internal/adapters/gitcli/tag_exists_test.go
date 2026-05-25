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

func TestTagExistsUsesRefsTags(t *testing.T) {
	runner := &fakeRunner{results: []processport.Result{{ExitCode: 0}}}
	client := New(runner, Options{})

	exists, err := client.TagExists(context.Background(), "/repo", gitport.TagName("v1.0.0"))
	if err != nil || !exists {
		t.Fatalf("TagExists() = %v, %v", exists, err)
	}
	assertStringSlice(t, runner.specs[0].Args, []string{"rev-parse", "--verify", "--quiet", "refs/tags/v1.0.0"})
}

func TestDeleteTagBuildsCommand(t *testing.T) {
	runner := &fakeRunner{results: []processport.Result{{ExitCode: 0}}}
	client := New(runner, Options{})

	err := client.DeleteTag(context.Background(), "/repo", gitport.TagName("v1.0.0"))
	if err != nil {
		t.Fatalf("DeleteTag() error = %v", err)
	}
	assertStringSlice(t, runner.specs[0].Args, []string{"tag", "-d", "v1.0.0"})
}
