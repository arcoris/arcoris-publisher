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
	"path/filepath"
	"testing"
)

func TestFileDigestHashesRegularFileContent(t *testing.T) {
	file := filepath.Join(t.TempDir(), "file.txt")
	writeFile(t, file, "content")

	got, err := fileDigest(file)
	if err != nil {
		t.Fatalf("fileDigest() error = %v", err)
	}
	want := "ed7002b439e9ac845f22357d822bac1444730fbdb6016d3ec9432297b9ec9f73"
	if got != want {
		t.Fatalf("fileDigest() = %s, want %s", got, want)
	}
}

func TestFileDigestReturnsOpenError(t *testing.T) {
	if _, err := fileDigest(filepath.Join(t.TempDir(), "missing.txt")); err == nil {
		t.Fatalf("fileDigest() should return source open error")
	}
}
