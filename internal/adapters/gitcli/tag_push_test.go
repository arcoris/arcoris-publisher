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

func TestPushTagBuildsRefSpec(t *testing.T) {
	runner := &fakeRunner{}
	client := New(runner, Options{})

	err := client.PushTag(context.Background(), "/repo", "upstream", gitport.TagName("v1.0.0"), gitport.PushOptions{})
	if err != nil {
		t.Fatalf("PushTag() error = %v", err)
	}
	assertStringSlice(t, runner.specs[0].Args, []string{"push", "upstream", "refs/tags/v1.0.0:refs/tags/v1.0.0"})
}

func TestTagRef(t *testing.T) {
	if got := tagRef(gitport.TagName("v1.0.0")); got != "refs/tags/v1.0.0" {
		t.Fatalf("tagRef() = %q, want refs/tags/v1.0.0", got)
	}
}
