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

func TestCopyTreePreservesSymlink(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src")
	dst := filepath.Join(root, "dst")
	must(t, os.MkdirAll(src, 0o755))
	if err := os.Symlink("target", filepath.Join(src, "link")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}

	result, err := New().CopyTree(context.Background(), src, dst, fsport.CopyTreeOptions{
		SafetyRoot:    root,
		SymlinkPolicy: fsport.SymlinkPreserve,
	})
	if err != nil {
		t.Fatalf("CopyTree() error = %v", err)
	}
	if result.FilesCopied != 1 {
		t.Fatalf("FilesCopied = %d, want 1", result.FilesCopied)
	}
	target, err := os.Readlink(filepath.Join(dst, "link"))
	if err != nil || target != "target" {
		t.Fatalf("Readlink() = %q, %v; want target, nil", target, err)
	}
}

func TestCopyTreeRejectsSymlinkWithOriginalCode(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src")
	dst := filepath.Join(root, "dst")
	must(t, os.MkdirAll(src, 0o755))
	if err := os.Symlink("target", filepath.Join(src, "link")); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}

	_, err := New().CopyTree(context.Background(), src, dst, fsport.CopyTreeOptions{
		SafetyRoot:    root,
		SymlinkPolicy: fsport.SymlinkReject,
	})
	assertPortCode(t, err, fsport.CodeSymlinkRejected)
}

func TestCopySymlinkRejectsFollowAndUnknownPolicy(t *testing.T) {
	root := t.TempDir()
	link := filepath.Join(root, "link")
	if err := os.Symlink("target", link); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	result := fsport.CopyTreeResult{}

	assertPortCode(
		t,
		New().copySymlink(link, filepath.Join(root, "dst"), fsport.SymlinkFollow, &result),
		fsport.CodeSymlinkRejected,
	)
	assertPortCode(
		t,
		New().copySymlink(link, filepath.Join(root, "dst"), fsport.SymlinkPolicy("unknown"), &result),
		fsport.CodeSymlinkRejected,
	)
}

func TestCopySymlinkReturnsReadlinkError(t *testing.T) {
	result := fsport.CopyTreeResult{}
	err := New().copySymlink(filepath.Join(t.TempDir(), "missing"), filepath.Join(t.TempDir(), "dst"), fsport.SymlinkPreserve, &result)
	if err == nil {
		t.Fatalf("copySymlink() should return readlink error")
	}
}
