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

func TestValidateCopyTreeRejectsUnsafeDestinationFirst(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src")
	must(t, os.MkdirAll(src, 0o755))

	err := validateCopyTree(src, filepath.Join(filepath.Dir(root), "dst"), fsport.CopyTreeOptions{SafetyRoot: root})
	assertPortCode(t, err, fsport.CodePathOutsideRoot)
}

func TestValidateCopyTreeRejectsNonDirectorySource(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "file.txt")
	writeFile(t, src, "content")

	err := validateCopyTree(src, filepath.Join(root, "dst"), fsport.CopyTreeOptions{SafetyRoot: root})
	assertPortCode(t, err, fsport.CodeCopyFailed)
}

func TestCopySourceStatErrorDistinguishesMissingPath(t *testing.T) {
	err := copySourceStatError(filepath.Join(t.TempDir(), "missing"), os.ErrNotExist)
	assertPortCode(t, err, fsport.CodePathNotFound)
}

func TestCopySourceStatErrorMapsGenericStatFailure(t *testing.T) {
	err := copySourceStatError("/repo", os.ErrPermission)
	assertPortCode(t, err, fsport.CodePermissionDenied)
}
