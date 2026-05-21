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

func TestFetchBuildsCommand(t *testing.T) {
	runner := &fakeRunner{}
	client := New(runner, Options{})

	err := client.Fetch(context.Background(), "/repo", "", gitport.FetchOptions{
		Prune:    true,
		Tags:     gitport.FetchTagsNone,
		RefSpecs: []gitport.RefSpec{"refs/heads/main:refs/remotes/origin/main"},
	})
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	assertStringSlice(t, runner.specs[0].Args, []string{
		"fetch",
		"--prune",
		"--no-tags",
		"origin",
		"refs/heads/main:refs/remotes/origin/main",
	})
}

func TestFetchArgsSupportsAllTags(t *testing.T) {
	args := fetchArgs("upstream", gitport.FetchOptions{Tags: gitport.FetchTagsAll})
	assertStringSlice(t, args, []string{"fetch", "--tags", "upstream"})
}
