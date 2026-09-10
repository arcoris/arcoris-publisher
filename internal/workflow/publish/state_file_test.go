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

package publish

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestReadBoundedStateFileReadsRegularFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state")
	want := []byte("transaction-state")
	if err := os.WriteFile(path, want, 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := readBoundedStateFile(path, int64(len(want)))
	if err != nil {
		t.Fatalf("readBoundedStateFile() error = %v", err)
	}
	if string(got) != string(want) {
		t.Fatalf("readBoundedStateFile() = %q, want %q", got, want)
	}
}

func TestReadBoundedStateFileRejectsOversizedFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state")
	if err := os.WriteFile(path, []byte("oversized"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	_, err := readBoundedStateFile(path, 4)
	if !errors.Is(err, errStateFileTooLarge) {
		t.Fatalf("readBoundedStateFile() error = %v, want errStateFileTooLarge", err)
	}
}

func TestReadBoundedStateFileRejectsDirectory(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state")
	if err := os.Mkdir(path, 0o700); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}

	_, err := readBoundedStateFile(path, 1024)
	if !errors.Is(err, errStateFileNotRegular) {
		t.Fatalf("readBoundedStateFile() error = %v, want errStateFileNotRegular", err)
	}
}

func TestReadBoundedStateFileRejectsSymlink(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "target")
	if err := os.WriteFile(target, []byte("state"), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	path := filepath.Join(dir, "state")
	if err := os.Symlink(target, path); err != nil {
		t.Skipf("Symlink() unavailable: %v", err)
	}

	_, err := readBoundedStateFile(path, 1024)
	if !errors.Is(err, errStateFileNotRegular) {
		t.Fatalf("readBoundedStateFile() error = %v, want errStateFileNotRegular", err)
	}
}
