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

func TestCopyTreeCopiesIncludedFilesAndSkipsGitDir(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src")
	dst := filepath.Join(root, "dst")
	writeFile(t, filepath.Join(src, "go.mod"), "module example.com/m\n")
	writeFile(t, filepath.Join(src, "pkg", "x.go"), "package pkg\n")
	writeFile(t, filepath.Join(src, "pkg", "x.tmp"), "tmp\n")
	writeFile(t, filepath.Join(src, ".git", "config"), "git\n")

	result, err := New().CopyTree(context.Background(), src, dst, fsport.CopyTreeOptions{
		SafetyRoot: root,
		Exclude:    []string{"**/*.tmp"},
	})
	if err != nil {
		t.Fatalf("CopyTree() error = %v", err)
	}
	if result.FilesCopied != 2 || result.FilesSkipped != 1 {
		t.Fatalf("CopyTree() result = %#v", result)
	}
	assertExists(t, filepath.Join(dst, "go.mod"))
	assertExists(t, filepath.Join(dst, "pkg", "x.go"))
	assertMissing(t, filepath.Join(dst, "pkg", "x.tmp"))
	assertMissing(t, filepath.Join(dst, ".git"))
}

func TestCopyTreeRejectsDestinationOutsideSafetyRoot(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src")
	must(t, os.MkdirAll(src, 0o755))

	_, err := New().CopyTree(context.Background(), src, filepath.Join(filepath.Dir(root), "dst"), fsport.CopyTreeOptions{SafetyRoot: root})
	assertPortCode(t, err, fsport.CodePathOutsideRoot)
}

func TestCopyTreeMissingSource(t *testing.T) {
	root := t.TempDir()
	_, err := New().CopyTree(
		context.Background(),
		filepath.Join(root, "missing"),
		filepath.Join(root, "dst"),
		fsport.CopyTreeOptions{SafetyRoot: root},
	)
	assertPortCode(t, err, fsport.CodePathNotFound)
}

func TestCopyTreeContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := New().CopyTree(ctx, t.TempDir(), filepath.Join(t.TempDir(), "dst"), fsport.CopyTreeOptions{})
	if err == nil {
		t.Fatalf("CopyTree() should return context cancellation")
	}
}

func TestCopyTreeRejectsFileSource(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "file.txt")
	writeFile(t, src, "content")

	_, err := New().CopyTree(context.Background(), src, filepath.Join(root, "dst"), fsport.CopyTreeOptions{SafetyRoot: root})
	assertPortCode(t, err, fsport.CodeCopyFailed)
}
