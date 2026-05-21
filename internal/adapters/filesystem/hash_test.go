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

func TestTreeHashIsDeterministicAndSensitiveToContent(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "a.txt"), "one")
	writeFile(t, filepath.Join(dir, "b.txt"), "two")

	hash1, err := New().TreeHash(context.Background(), dir, fsport.TreeHashOptions{IncludeFileMode: true})
	if err != nil {
		t.Fatalf("TreeHash() error = %v", err)
	}
	hash2, err := New().TreeHash(context.Background(), dir, fsport.TreeHashOptions{IncludeFileMode: true})
	if err != nil {
		t.Fatalf("TreeHash() second error = %v", err)
	}
	if hash1 != hash2 {
		t.Fatalf("TreeHash() not deterministic: %q != %q", hash1, hash2)
	}

	writeFile(t, filepath.Join(dir, "a.txt"), "changed")
	hash3, err := New().TreeHash(context.Background(), dir, fsport.TreeHashOptions{IncludeFileMode: true})
	if err != nil {
		t.Fatalf("TreeHash() after change error = %v", err)
	}
	if hash1 == hash3 {
		t.Fatalf("TreeHash() should change when content changes")
	}
}

func TestTreeHashFiltersAndSymlinkPolicy(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "keep.go"), "package main\n")
	writeFile(t, filepath.Join(dir, "skip.txt"), "skip\n")
	if err := os.Symlink("keep.go", filepath.Join(dir, "link")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}

	_, err := New().TreeHash(context.Background(), dir, fsport.TreeHashOptions{
		Include:       []string{"*.go"},
		SymlinkPolicy: fsport.SymlinkReject,
	})
	if err != nil {
		t.Fatalf("TreeHash() should not visit excluded symlink: %v", err)
	}

	_, err = New().TreeHash(context.Background(), dir, fsport.TreeHashOptions{
		IncludeSymlinks: true,
		SymlinkPolicy:   fsport.SymlinkReject,
	})
	assertPortCode(t, err, fsport.CodeSymlinkRejected)
}

func TestTreeHashMissingRoot(t *testing.T) {
	_, err := New().TreeHash(context.Background(), filepath.Join(t.TempDir(), "missing"), fsport.TreeHashOptions{})
	assertPortCode(t, err, fsport.CodePathNotFound)
}

func TestTreeHashContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := New().TreeHash(ctx, t.TempDir(), fsport.TreeHashOptions{})
	if err == nil {
		t.Fatalf("TreeHash() should return context cancellation")
	}
}
