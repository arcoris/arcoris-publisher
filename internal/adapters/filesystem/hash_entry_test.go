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

func TestCollectHashEntriesConvertsIncludedFiles(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, "keep.go"), "package main\n")
	writeFile(t, filepath.Join(dir, "skip.txt"), "skip")

	entries, err := collectHashEntries(context.Background(), dir, fsport.TreeHashOptions{Include: []string{"*.go"}})
	if err != nil {
		t.Fatalf("collectHashEntries() error = %v", err)
	}
	if len(entries) != 1 || entries[0].rel != "keep.go" || entries[0].kind != "file" {
		t.Fatalf("collectHashEntries() = %#v", entries)
	}
}

func TestHashCollectorSymlinkPolicy(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "file.txt")
	link := filepath.Join(dir, "link")
	writeFile(t, file, "content")
	if err := os.Symlink("file.txt", link); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	entry := readDirEntry(t, dir, "link")
	collector := hashCollector{
		ctx:    context.Background(),
		root:   dir,
		opts:   fsport.TreeHashOptions{IncludeSymlinks: true},
		policy: fsport.SymlinkPreserve,
	}

	item, include, err := collector.entry(link, entry, nil)
	if err != nil || !include || item.kind != "symlink" || item.content != "file.txt" {
		t.Fatalf("entry(symlink) = %#v, %v, %v", item, include, err)
	}
}

func TestHashCollectorIncludedDirectoryAndRejectedSymlink(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "subdir")
	must(t, os.Mkdir(subdir, 0o755))
	dirEntry := readDirEntry(t, dir, "subdir")
	collector := hashCollector{ctx: context.Background(), root: dir, opts: fsport.TreeHashOptions{}, policy: fsport.SymlinkReject}

	item, include, err := collector.includedEntry(subdir, dirEntry, "subdir")
	if err != nil || !include || item.kind != "dir" {
		t.Fatalf("includedEntry(dir) = %#v, %v, %v", item, include, err)
	}

	link := filepath.Join(dir, "link")
	if err := os.Symlink("subdir", link); err != nil {
		t.Skipf("symlink unavailable: %v", err)
	}
	_, _, err = collector.entry(link, readDirEntry(t, dir, "link"), nil)
	assertPortCode(t, err, fsport.CodeSymlinkRejected)
}

func TestRegularFileHashEntryReturnsDigestError(t *testing.T) {
	_, _, err := regularFileHashEntry(filepath.Join(t.TempDir(), "missing.txt"), hashEntry{})
	if err == nil {
		t.Fatalf("regularFileHashEntry() should return digest error")
	}
}
