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

package porttest

import (
	"context"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/ports/git"
)

func TestGitStatusReturnsDetachedEntries(t *testing.T) {
	fake := NewGit()
	fake.Statuses["/repo"] = git.Status{
		Clean:   false,
		Entries: []git.StatusEntry{{Path: "go.mod", Code: " M"}},
	}

	status, err := fake.Status(context.Background(), "/repo")
	if err != nil {
		t.Fatalf("Status() error = %v", err)
	}
	status.Entries[0].Path = "mutated"

	again, err := fake.Status(context.Background(), "/repo")
	if err != nil {
		t.Fatalf("Status() second error = %v", err)
	}
	if again.Entries[0].Path != "go.mod" {
		t.Fatalf("status entries were attached: %#v", again.Entries)
	}
}

func TestGitRecordsCallsInOrder(t *testing.T) {
	fake := NewGit()
	_ = fake.AddAll(context.Background(), "/repo")
	_, _ = fake.Commit(context.Background(), "/repo", "message", git.CommitOptions{})
	_ = fake.Push(context.Background(), "/repo", "origin", "main:main", git.PushOptions{})

	got := []string{fake.Calls[0].Op, fake.Calls[1].Op, fake.Calls[2].Op}
	want := []string{"add", "commit", "push"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("call[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
