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
	"testing"

	gitport "arcoris.dev/arcoris-publisher/internal/ports/git"
)

func TestCommandDetachesSlicesAndAppliesOptions(t *testing.T) {
	runner := &fakeRunner{}
	client := New(runner, Options{GitBinary: "custom-git", Env: []string{"A=1"}})
	args := []string{"status"}
	sensitive := []string{"token"}

	spec := client.command("/repo", args, sensitive, true, true)
	args[0] = "mutated"
	sensitive[0] = "mutated"

	if spec.Name != "custom-git" || spec.Dir != "/repo" {
		t.Fatalf("unexpected spec identity: %#v", spec)
	}
	assertStringSlice(t, spec.Args, []string{"status"})
	assertStringSlice(t, spec.Env, []string{"A=1"})
	assertStringSlice(t, spec.SensitiveValues, []string{"token"})
	if !spec.CaptureStdout || !spec.CaptureStderr {
		t.Fatalf("expected captured stdout/stderr")
	}
}

func TestStringsOf(t *testing.T) {
	got := stringsOf([]gitport.RefSpec{"a:b", "c:d"})
	assertStringSlice(t, got, []string{"a:b", "c:d"})
}
