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
	"errors"
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

func TestCleanDirRejectsFileTarget(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "file.txt")
	writeFile(t, file, "content")

	err := New().CleanDir(context.Background(), file, fsport.CleanDirOptions{SafetyRoot: root})
	assertPortCode(t, err, fsport.CodePermissionDenied)
}

func TestCleanDirContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := New().CleanDir(ctx, t.TempDir(), fsport.CleanDirOptions{})
	if err == nil {
		t.Fatalf("CleanDir() should return context cancellation")
	}
}

func TestCleanEntryPreservesDirectoryAndRemovesFile(t *testing.T) {
	root := t.TempDir()
	keep := filepath.Join(root, "keep")
	drop := filepath.Join(root, "drop.txt")
	writeFile(t, filepath.Join(keep, "nested.txt"), "keep")
	writeFile(t, drop, "drop")

	keepEntry := readDirEntry(t, root, "keep")
	err := cleanEntry(context.Background(), root, keep, keepEntry, nil, []string{"keep"})
	if !errors.Is(err, filepath.SkipDir) {
		t.Fatalf("cleanEntry(preserve dir) = %v, want SkipDir", err)
	}
	assertExists(t, keep)

	dropEntry := readDirEntry(t, root, "drop.txt")
	if err := cleanEntry(context.Background(), root, drop, dropEntry, nil, nil); err != nil {
		t.Fatalf("cleanEntry(file) error = %v", err)
	}
	assertMissing(t, drop)
}

func TestCleanEntryReturnsWalkAndContextErrors(t *testing.T) {
	walkErr := errors.New("walk failed")
	if err := cleanEntry(context.Background(), "/root", "/root/file", nil, walkErr, nil); !errors.Is(err, walkErr) {
		t.Fatalf("cleanEntry(walkErr) = %v, want walk error", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := cleanEntry(ctx, "/root", "/root/file", nil, nil, nil); !errors.Is(err, context.Canceled) {
		t.Fatalf("cleanEntry(cancelled) = %v, want context canceled", err)
	}
}
