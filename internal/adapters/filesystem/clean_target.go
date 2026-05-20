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

	fsport "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
)

// validateCleanTarget verifies that CleanDir may safely operate on dir.
//
// The boolean return is false only for the allowed "missing directory" case,
// where CleanDir should become a no-op rather than an error.
func validateCleanTarget(dir string, opts fsport.CleanDirOptions) (bool, error) {
	if err := ensureSafeRemove(dir, opts.SafetyRoot); err != nil {
		return false, pathError(fsport.CodeUnsafeRemove, "unsafe clean target", err, dir)
	}
	info, err := os.Lstat(dir)
	if err != nil {
		return handleMissingCleanTarget(dir, opts, err)
	}
	if !info.IsDir() {
		return false, pathError(fsport.CodePermissionDenied, "clean target is not a directory", nil, dir)
	}
	if opts.RequireGitDir {
		return validateRequiredGitDir(dir)
	}
	return true, nil
}

// handleMissingCleanTarget maps Lstat failures into the CleanDir contract.
func handleMissingCleanTarget(dir string, opts fsport.CleanDirOptions, err error) (bool, error) {
	if isNotExist(err) && opts.AllowMissing {
		return false, nil
	}
	if isNotExist(err) {
		return false, pathError(fsport.CodePathNotFound, "directory not found", err, dir)
	}
	return false, pathError(fsport.CodePermissionDenied, "directory stat failed", err, dir)
}

// validateRequiredGitDir enforces the opt-in guard for repository cleanup.
func validateRequiredGitDir(dir string) (bool, error) {
	gitDir := filepath.Join(dir, ".git")
	if _, err := os.Lstat(gitDir); err != nil {
		return false, pathError(fsport.CodePathNotFound, "required .git directory is missing", err, gitDir)
	}
	return true, nil
}
