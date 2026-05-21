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

func TestWriteFileCreatesParentsAndRespectsOverwrite(t *testing.T) {
	fs := New()
	file := filepath.Join(t.TempDir(), "nested", "file.txt")

	err := fs.WriteFile(context.Background(), file, []byte("one"), fsport.WriteFileOptions{CreateDirs: true})
	if err != nil {
		t.Fatalf("WriteFile() create error = %v", err)
	}
	err = fs.WriteFile(context.Background(), file, []byte("two"), fsport.WriteFileOptions{})
	assertPortCode(t, err, fsport.CodePermissionDenied)

	err = fs.WriteFile(context.Background(), file, []byte("two"), fsport.WriteFileOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("WriteFile() overwrite error = %v", err)
	}
	data, err := os.ReadFile(file)
	if err != nil || string(data) != "two" {
		t.Fatalf("ReadFile() = %q, %v; want two, nil", data, err)
	}
}

func TestWriteFileHelpers(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "nested", "file.txt")

	if err := prepareWriteParent(file, fsport.WriteFileOptions{CreateDirs: true}); err != nil {
		t.Fatalf("prepareWriteParent() error = %v", err)
	}
	if got := writeFilePerm(fsport.WriteFileOptions{}); got != 0o644 {
		t.Fatalf("writeFilePerm() = %v, want 0644", got)
	}
	if got := writeFilePerm(fsport.WriteFileOptions{Perm: 0o600}); got != 0o600 {
		t.Fatalf("writeFilePerm(custom) = %v, want 0600", got)
	}
}

func TestWriteFileContextCancelledAndMissingParent(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := New().WriteFile(ctx, filepath.Join(t.TempDir(), "file.txt"), []byte("data"), fsport.WriteFileOptions{}); err == nil {
		t.Fatalf("WriteFile() should return context cancellation")
	}

	err := New().WriteFile(context.Background(), filepath.Join(t.TempDir(), "missing", "file.txt"), []byte("data"), fsport.WriteFileOptions{})
	assertPortCode(t, err, fsport.CodePermissionDenied)
}

func TestCheckOverwriteAllowedAllowsMissingFile(t *testing.T) {
	if err := checkOverwriteAllowed(filepath.Join(t.TempDir(), "missing.txt"), fsport.WriteFileOptions{}); err != nil {
		t.Fatalf("checkOverwriteAllowed(missing) error = %v", err)
	}
}

func TestPrepareWriteParentMapsCreateDirectoryFailure(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "file")
	writeFile(t, file, "content")

	err := prepareWriteParent(filepath.Join(file, "child.txt"), fsport.WriteFileOptions{CreateDirs: true})
	assertPortCode(t, err, fsport.CodePermissionDenied)
}
