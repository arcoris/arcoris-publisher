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
	"os"
	"path/filepath"
	"testing"

	fsport "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
)

func TestCopyTreeOperationCreateRoot(t *testing.T) {
	dst := filepath.Join(t.TempDir(), "dst")
	op := copyTreeOperation{dst: dst}

	if err := op.createRoot(); err != nil {
		t.Fatalf("createRoot() error = %v", err)
	}
	assertExists(t, dst)
}

func TestCopyTreeOperationVisitRootCreatesDestination(t *testing.T) {
	src := t.TempDir()
	dst := filepath.Join(t.TempDir(), "dst")
	entry := readDirEntry(t, filepath.Dir(src), filepath.Base(src))
	result := fsport.CopyTreeResult{}
	op := copyTreeOperation{ctx: context.Background(), src: src, dst: dst, result: &result}

	if err := op.visit(src, entry, nil); err != nil {
		t.Fatalf("visit(root) error = %v", err)
	}
	assertExists(t, dst)
}

func TestCopyTreeOperationIncludeAndSkipEntry(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "skip.tmp"), "tmp")
	must(t, os.Mkdir(filepath.Join(dir, "subdir"), 0o755))
	entries, err := os.ReadDir(dir)
	must(t, err)
	result := fsport.CopyTreeResult{}
	op := copyTreeOperation{opts: fsport.CopyTreeOptions{Exclude: []string{"*.tmp"}}, result: &result}

	for _, entry := range entries {
		if entry.Name() == "skip.tmp" {
			if op.includeEntry(entry.Name(), entry) {
				t.Fatalf("includeEntry() should reject excluded file")
			}
			if err := op.skipEntry(entry); err != nil {
				t.Fatalf("skipEntry(file) error = %v", err)
			}
		}
		if entry.Name() == "subdir" {
			if err := op.skipEntry(entry); !errors.Is(err, filepath.SkipDir) {
				t.Fatalf("skipEntry(dir) = %v, want SkipDir", err)
			}
		}
	}
	if result.FilesSkipped != 1 {
		t.Fatalf("FilesSkipped = %d, want 1", result.FilesSkipped)
	}
}

func TestCopyTreeOperationVisitReturnsWalkAndContextErrors(t *testing.T) {
	walkErr := errors.New("walk failed")
	op := copyTreeOperation{ctx: context.Background()}
	if err := op.visit("", nil, walkErr); !errors.Is(err, walkErr) {
		t.Fatalf("visit(walkErr) = %v, want walk error", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	op = copyTreeOperation{ctx: ctx}
	if err := op.visit("", nil, nil); !errors.Is(err, context.Canceled) {
		t.Fatalf("visit(cancelled) = %v, want context canceled", err)
	}
}
