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

func TestMkdirAllAndRemoveAll(t *testing.T) {
	fs := New()
	root := t.TempDir()
	dir := filepath.Join(root, "dir")

	if err := fs.MkdirAll(context.Background(), dir, fsport.MkdirOptions{}); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := fs.RemoveAll(context.Background(), dir, fsport.RemoveOptions{SafetyRoot: root}); err != nil {
		t.Fatalf("RemoveAll() error = %v", err)
	}
	assertMissing(t, dir)
}

func TestRemoveAllAllowMissingAndUnsafe(t *testing.T) {
	root := t.TempDir()
	fs := New()

	err := fs.RemoveAll(context.Background(), filepath.Join(root, "missing"), fsport.RemoveOptions{SafetyRoot: root, AllowMissing: true})
	if err != nil {
		t.Fatalf("RemoveAll() missing allowed error = %v", err)
	}
	err = fs.RemoveAll(context.Background(), filepath.Dir(root), fsport.RemoveOptions{SafetyRoot: root})
	assertPortCode(t, err, fsport.CodeUnsafeRemove)
}
