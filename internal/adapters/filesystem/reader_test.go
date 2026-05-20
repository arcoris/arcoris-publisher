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

func TestReaderExistsIsDirAndReadFile(t *testing.T) {
	fs := New()
	dir := t.TempDir()
	file := filepath.Join(dir, "file.txt")
	writeFile(t, file, "content")

	exists, err := fs.Exists(context.Background(), file)
	if err != nil || !exists {
		t.Fatalf("Exists() = %v, %v; want true, nil", exists, err)
	}
	isDir, err := fs.IsDir(context.Background(), dir)
	if err != nil || !isDir {
		t.Fatalf("IsDir() = %v, %v; want true, nil", isDir, err)
	}
	data, err := fs.ReadFile(context.Background(), file)
	if err != nil || string(data) != "content" {
		t.Fatalf("ReadFile() = %q, %v; want content, nil", data, err)
	}
}

func TestReadFileMissing(t *testing.T) {
	_, err := New().ReadFile(context.Background(), filepath.Join(t.TempDir(), "missing"))
	assertPortCode(t, err, fsport.CodePathNotFound)
}
