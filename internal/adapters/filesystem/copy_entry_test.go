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
	"os"
	"path/filepath"
	"testing"

	fsport "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
)

func TestCopyIncludedEntryCopiesRegularFile(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src")
	dst := filepath.Join(root, "dst")
	writeFile(t, filepath.Join(src, "file.txt"), "content")
	entry := readDirEntry(t, src, "file.txt")
	result := fsport.CopyTreeResult{}
	op := copyTreeOperation{fs: New(), dst: dst, result: &result}

	if err := op.copyIncludedEntry(filepath.Join(src, "file.txt"), "file.txt", entry); err != nil {
		t.Fatalf("copyIncludedEntry() error = %v", err)
	}
	if result.FilesCopied != 1 || result.BytesCopied != int64(len("content")) {
		t.Fatalf("copyIncludedEntry() result = %#v", result)
	}
	assertExists(t, filepath.Join(dst, "file.txt"))
}

func TestCopyDirectoryRecordsMetric(t *testing.T) {
	target := filepath.Join(t.TempDir(), "dst")
	result := fsport.CopyTreeResult{}
	op := copyTreeOperation{result: &result}

	if err := op.copyDirectory(target, 0o700); err != nil {
		t.Fatalf("copyDirectory() error = %v", err)
	}
	if result.DirectoriesCopied != 1 {
		t.Fatalf("DirectoriesCopied = %d, want 1", result.DirectoriesCopied)
	}
}

func TestCopyRegularFilePreservesTimesWhenRequested(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src.txt")
	dst := filepath.Join(root, "nested", "dst.txt")
	writeFile(t, src, "content")
	info, err := os.Stat(src)
	must(t, err)
	result := fsport.CopyTreeResult{}
	op := copyTreeOperation{opts: fsport.CopyTreeOptions{PreserveTimes: true}, result: &result}

	if err := op.copyRegularFile(src, dst, info.Mode(), info); err != nil {
		t.Fatalf("copyRegularFile() error = %v", err)
	}
	if result.FilesCopied != 1 || result.BytesCopied != int64(len("content")) {
		t.Fatalf("copyRegularFile() result = %#v", result)
	}
	assertExists(t, dst)
}

func readDirEntry(t *testing.T, dir string, name string) os.DirEntry {
	t.Helper()
	entries, err := os.ReadDir(dir)
	must(t, err)
	for _, entry := range entries {
		if entry.Name() == name {
			return entry
		}
	}
	t.Fatalf("entry %s not found in %s", name, dir)
	return nil
}
