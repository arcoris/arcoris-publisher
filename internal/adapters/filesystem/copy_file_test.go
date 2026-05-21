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
)

func TestCopyFileCopiesBytesAndAppliesMode(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src.txt")
	dst := filepath.Join(root, "dst.txt")
	writeFile(t, src, "content")

	written, err := copyFile(src, dst, 0o600, true)
	if err != nil {
		t.Fatalf("copyFile() error = %v", err)
	}
	if written != int64(len("content")) {
		t.Fatalf("copyFile() bytes = %d, want %d", written, len("content"))
	}
	data, err := os.ReadFile(dst)
	if err != nil || string(data) != "content" {
		t.Fatalf("copied data = %q, %v", data, err)
	}
}

func TestCopyFilePerm(t *testing.T) {
	if got := copyFilePerm(0o700, false); got != 0o644 {
		t.Fatalf("copyFilePerm(default) = %v, want 0644", got)
	}
	if got := copyFilePerm(0o700, true); got != 0o700 {
		t.Fatalf("copyFilePerm(preserve) = %v, want 0700", got)
	}
}

func TestCopyFileReturnsOpenError(t *testing.T) {
	_, err := copyFile(filepath.Join(t.TempDir(), "missing.txt"), filepath.Join(t.TempDir(), "dst.txt"), 0o644, false)
	if err == nil {
		t.Fatalf("copyFile() should return source open error")
	}
}
