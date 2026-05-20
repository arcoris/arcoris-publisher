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

	fsport "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
)

// CleanDir removes every child of dir except explicitly preserved descendants.
//
// The method is intentionally arranged as validate-then-walk. Validation keeps
// destructive-operation safety checks in one place, while the walker focuses on
// deciding which entries are deleted, skipped, or preserved.
func (fs *FileSystem) CleanDir(ctx context.Context, dir string, opts fsport.CleanDirOptions) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	validated, err := validateCleanTarget(dir, opts)
	if err != nil || !validated {
		return err
	}
	preserve := normalizePreserve(opts.Preserve)
	err = filepath.WalkDir(dir, func(pathValue string, entry os.DirEntry, walkErr error) error {
		return cleanEntry(ctx, dir, pathValue, entry, walkErr, preserve)
	})
	if err != nil {
		if isPortError(err) {
			return err
		}
		return pathError(fsport.CodePermissionDenied, "directory cleanup failed", err, dir)
	}
	return nil
}

// cleanEntry applies the preserve policy to one WalkDir entry.
//
// A directory that is fully removed returns filepath.SkipDir so WalkDir does not
// continue into entries that no longer exist. A directory that contains a
// preserved descendant is kept and traversal continues deeper.
func cleanEntry(ctx context.Context, root string, pathValue string, entry os.DirEntry, walkErr error, preserve []string) error {
	if walkErr != nil {
		return walkErr
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if pathValue == root {
		return nil
	}
	rel, err := filepath.Rel(root, pathValue)
	if err != nil {
		return err
	}
	relSlash := slash(rel)
	if shouldPreserve(relSlash, preserve) {
		if entry.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}
	if entry.IsDir() && hasPreservedChild(relSlash, preserve) {
		return nil
	}
	if err := os.RemoveAll(pathValue); err != nil {
		return pathError(fsport.CodePermissionDenied, "directory cleanup failed", err, pathValue)
	}
	if entry.IsDir() {
		return filepath.SkipDir
	}
	return nil
}
