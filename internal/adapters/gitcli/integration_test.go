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
	"os"
	osexec "os/exec"
	"path/filepath"
	"testing"

	execadapter "arcoris.dev/arcoris-publisher/internal/adapters/exec"
	gitport "arcoris.dev/arcoris-publisher/internal/ports/git"
)

func TestClientLocalRepositoryLifecycle(t *testing.T) {
	if _, err := osexec.LookPath("git"); err != nil {
		t.Skip("git binary is unavailable")
	}
	ctx := context.Background()
	client := New(execadapter.New(execadapter.Options{}), Options{})
	root := t.TempDir()
	source := filepath.Join(root, "source")
	bare := filepath.Join(root, "remote.git")

	runGit(t, root, "init", source)
	runGit(t, source, "config", "user.email", "test@example.com")
	runGit(t, source, "config", "user.name", "Test User")
	must(t, os.WriteFile(filepath.Join(source, "file.txt"), []byte("one"), 0o644))
	runGit(t, source, "add", ".")
	runGit(t, source, "commit", "-m", "initial")
	runGit(t, root, "init", "--bare", bare)
	runGit(t, source, "remote", "add", "origin", bare)
	runGit(t, source, "push", "origin", "HEAD:main")

	clone := filepath.Join(root, "clone")
	if err := client.Clone(ctx, bare, clone, gitport.CloneOptions{}); err != nil {
		t.Fatalf("Clone() error = %v", err)
	}
	if err := client.Fetch(ctx, clone, "origin", gitport.FetchOptions{Prune: true, Tags: gitport.FetchTagsNone}); err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if err := client.Checkout(ctx, clone, "main", gitport.CheckoutOptions{}); err != nil {
		t.Fatalf("Checkout() error = %v", err)
	}
	head, err := client.Head(ctx, clone)
	if err != nil || head == "" {
		t.Fatalf("Head() = %q, %v", head, err)
	}
	status, err := client.Status(ctx, clone)
	if err != nil || !status.Clean {
		t.Fatalf("Status() = %#v, %v", status, err)
	}
	must(t, os.WriteFile(filepath.Join(clone, "file.txt"), []byte("two"), 0o644))
	if err := client.AddAll(ctx, clone); err != nil {
		t.Fatalf("AddAll() error = %v", err)
	}
	commit, err := client.Commit(ctx, clone, "update", gitport.CommitOptions{
		AuthorName:     "Test User",
		AuthorEmail:    "test@example.com",
		CommitterName:  "Test User",
		CommitterEmail: "test@example.com",
	})
	if err != nil || commit == "" {
		t.Fatalf("Commit() = %q, %v", commit, err)
	}
	if err := client.CreateTag(ctx, clone, gitport.TagName("v0.0.1"), commit, gitport.TagOptions{}); err != nil {
		t.Fatalf("CreateTag() error = %v", err)
	}
	exists, err := client.TagExists(ctx, clone, gitport.TagName("v0.0.1"))
	if err != nil || !exists {
		t.Fatalf("TagExists() = %v, %v", exists, err)
	}
}
