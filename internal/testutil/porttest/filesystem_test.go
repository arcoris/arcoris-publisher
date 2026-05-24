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

package porttest

import (
	"context"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/ports/filesystem"
)

func TestFileSystemReadWriteDetachedData(t *testing.T) {
	fs := NewFileSystem()
	fs.AddFile("/repo/go.mod", []byte("module example.com/repo\n"))

	data, err := fs.ReadFile(context.Background(), "/repo/go.mod")
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	data[0] = 'X'

	again, err := fs.ReadFile(context.Background(), "/repo/go.mod")
	if err != nil {
		t.Fatalf("ReadFile() second error = %v", err)
	}
	if string(again) != "module example.com/repo\n" {
		t.Fatalf("stored data mutated: %q", again)
	}
}

func TestFileSystemCleanDirPreservesConfiguredChildren(t *testing.T) {
	fs := NewFileSystem()
	fs.AddFile("/repo/.git/config", []byte("git"))
	fs.AddFile("/repo/go.mod", []byte("module example.com/repo\n"))

	err := fs.CleanDir(context.Background(), "/repo", cleanOptions(".git"))
	if err != nil {
		t.Fatalf("CleanDir() error = %v", err)
	}

	if ok, _ := fs.Exists(context.Background(), "/repo/.git/config"); !ok {
		t.Fatal(".git/config was removed")
	}
	if ok, _ := fs.Exists(context.Background(), "/repo/go.mod"); ok {
		t.Fatal("go.mod was preserved unexpectedly")
	}
}

func cleanOptions(preserve ...string) filesystem.CleanDirOptions {
	return filesystem.CleanDirOptions{Preserve: preserve}
}
