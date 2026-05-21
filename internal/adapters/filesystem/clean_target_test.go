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
	"errors"
	"path/filepath"
	"testing"

	fsport "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
)

func TestValidateCleanTargetAllowsConfiguredMissingDirectory(t *testing.T) {
	ok, err := validateCleanTarget(filepath.Join(t.TempDir(), "missing"), fsport.CleanDirOptions{AllowMissing: true})

	if ok || err != nil {
		t.Fatalf("validateCleanTarget() = %v, %v; want false, nil", ok, err)
	}
}

func TestValidateCleanTargetRejectsFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "file.txt")
	writeFile(t, file, "content")

	ok, err := validateCleanTarget(file, fsport.CleanDirOptions{})
	if ok {
		t.Fatalf("validateCleanTarget() should reject files")
	}
	assertPortCode(t, err, fsport.CodePermissionDenied)
}

func TestValidateRequiredGitDir(t *testing.T) {
	dir := t.TempDir()
	ok, err := validateRequiredGitDir(dir)
	if ok {
		t.Fatalf("validateRequiredGitDir() should reject missing .git")
	}
	assertPortCode(t, err, fsport.CodePathNotFound)
}

func TestHandleMissingCleanTargetMapsGenericStatFailure(t *testing.T) {
	ok, err := handleMissingCleanTarget("/repo", fsport.CleanDirOptions{}, errors.New("stat failed"))
	if ok {
		t.Fatalf("handleMissingCleanTarget() should not validate generic stat failure")
	}
	assertPortCode(t, err, fsport.CodePermissionDenied)
}
