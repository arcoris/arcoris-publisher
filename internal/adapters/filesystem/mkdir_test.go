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
	"path/filepath"
	"testing"

	fsport "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
)

func TestMkdirAllCreatesDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "dir")

	if err := New().MkdirAll(context.Background(), dir, fsport.MkdirOptions{}); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	assertExists(t, dir)
}

func TestMkdirPerm(t *testing.T) {
	if got := mkdirPerm(fsport.MkdirOptions{}); got != 0o755 {
		t.Fatalf("mkdirPerm() = %v, want 0755", got)
	}
	if got := mkdirPerm(fsport.MkdirOptions{Perm: 0o700}); got != 0o700 {
		t.Fatalf("mkdirPerm(custom) = %v, want 0700", got)
	}
}

func TestMkdirAllContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if err := New().MkdirAll(ctx, filepath.Join(t.TempDir(), "dir"), fsport.MkdirOptions{}); err == nil {
		t.Fatalf("MkdirAll() should return context cancellation")
	}
}

func TestMkdirAllMapsFilesystemFailure(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "file")
	writeFile(t, file, "content")

	err := New().MkdirAll(context.Background(), filepath.Join(file, "child"), fsport.MkdirOptions{})
	assertPortCode(t, err, fsport.CodePermissionDenied)
}
