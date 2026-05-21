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

func TestCreateAnnotatedTagBuildsCommand(t *testing.T) {
	runner := &fakeRunner{}
	client := New(runner, Options{})

	err := client.CreateTag(
		context.Background(),
		"/repo",
		gitport.TagName("v1.0.0"),
		gitport.CommitHash("abc"),
		gitport.TagOptions{Annotated: true, Message: "release", Force: true},
	)
	if err != nil {
		t.Fatalf("CreateTag() error = %v", err)
	}
	assertStringSlice(t, runner.specs[0].Args, []string{
		"tag",
		"-f",
		"-a",
		"v1.0.0",
		"abc",
		"-m",
		"release",
	})
}

func TestCreateTagArgsOmitsEmptyOptionalParts(t *testing.T) {
	args := createTagArgs(gitport.TagName("v1.0.0"), "", gitport.TagOptions{})
	assertStringSlice(t, args, []string{"tag", "v1.0.0"})
}
