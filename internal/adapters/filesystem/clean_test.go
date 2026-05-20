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

package filesystem

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	fsport "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
)

func TestCleanDirPreservesNestedPaths(t *testing.T) {
	ctx := context.Background()
	fs := New()
	root := t.TempDir()
	dir := filepath.Join(root, "repo")
	writeFile(t, filepath.Join(dir, ".git", "config"), "git")
	writeFile(t, filepath.Join(dir, "keep", "nested.txt"), "keep")
	writeFile(t, filepath.Join(dir, "keep", "remove.txt"), "remove")
	writeFile(t, filepath.Join(dir, "drop", "file.txt"), "drop")

	err := fs.CleanDir(ctx, dir, fsport.CleanDirOptions{
		Preserve:      []string{".git", "keep/nested.txt"},
		RequireGitDir: true,
		SafetyRoot:    root,
	})
	if err != nil {
		t.Fatalf("CleanDir() error = %v", err)
	}

	assertExists(t, filepath.Join(dir, ".git", "config"))
	assertExists(t, filepath.Join(dir, "keep", "nested.txt"))
	assertMissing(t, filepath.Join(dir, "keep", "remove.txt"))
	assertMissing(t, filepath.Join(dir, "drop"))
}

func TestCleanDirAllowMissing(t *testing.T) {
	err := New().CleanDir(context.Background(), filepath.Join(t.TempDir(), "missing"), fsport.CleanDirOptions{AllowMissing: true})
	if err != nil {
		t.Fatalf("CleanDir() missing with AllowMissing error = %v", err)
	}
}

func TestCleanDirRejectsUnsafeTarget(t *testing.T) {
	root := t.TempDir()
	err := New().CleanDir(context.Background(), filepath.Dir(root), fsport.CleanDirOptions{SafetyRoot: root})
	assertPortCode(t, err, fsport.CodeUnsafeRemove)
}

func TestCleanDirRequiresGitDir(t *testing.T) {
	dir := t.TempDir()
	err := New().CleanDir(context.Background(), dir, fsport.CleanDirOptions{RequireGitDir: true, SafetyRoot: dir})
	assertPortCode(t, err, fsport.CodePathNotFound)
}

func TestNormalizePreserve(t *testing.T) {
	got := normalizePreserve([]string{"./b/c", "", ".", "a"})
	want := []string{"a", "b/c"}
	if len(got) != len(want) {
		t.Fatalf("normalizePreserve() = %#v, want %#v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("normalizePreserve()[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestCleanDirContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := New().CleanDir(ctx, t.TempDir(), fsport.CleanDirOptions{})
	if err == nil {
		t.Fatalf("CleanDir() should return context cancellation")
	}
}

func assertExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Lstat(path); err != nil {
		t.Fatalf("expected %s to exist: %v", path, err)
	}
}

func assertMissing(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Lstat(path); !os.IsNotExist(err) {
		t.Fatalf("expected %s to be missing, stat err = %v", path, err)
	}
}
