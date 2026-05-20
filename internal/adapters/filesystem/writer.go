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

// WriteFile creates or replaces one file.
//
// The method maps "already exists" into permission_denied because the port does
// not currently expose a dedicated overwrite-rejected code. Callers can opt into
// replacement with WriteFileOptions.Overwrite.
func (fs *FileSystem) WriteFile(ctx context.Context, path string, data []byte, opts fsport.WriteFileOptions) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if opts.CreateDirs {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return pathError(fsport.CodePermissionDenied, "parent directory creation failed", err, filepath.Dir(path))
		}
	}
	if !opts.Overwrite {
		_, err := os.Lstat(path)
		if err == nil {
			return pathError(fsport.CodePermissionDenied, "file already exists", nil, path)
		}
		if !isNotExist(err) {
			return pathError(fsport.CodePermissionDenied, "filesystem stat failed", err, path)
		}
	}
	perm := opts.Perm
	if perm == 0 {
		perm = 0o644
	}
	if err := os.WriteFile(path, data, perm); err != nil {
		return pathError(fsport.CodePermissionDenied, "file write failed", err, path)
	}
	return nil
}

// MkdirAll creates a directory tree using the adapter's default permission when
// the caller leaves MkdirOptions.Perm at zero.
func (fs *FileSystem) MkdirAll(ctx context.Context, path string, opts fsport.MkdirOptions) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	perm := opts.Perm
	if perm == 0 {
		perm = 0o755
	}
	if err := os.MkdirAll(path, perm); err != nil {
		return pathError(fsport.CodePermissionDenied, "directory creation failed", err, path)
	}
	return nil
}

// RemoveAll removes a path after applying destructive-operation safety checks.
func (fs *FileSystem) RemoveAll(ctx context.Context, path string, opts fsport.RemoveOptions) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := ensureSafeRemove(path, opts.SafetyRoot); err != nil {
		return pathError(fsport.CodeUnsafeRemove, "unsafe remove target", err, path)
	}
	if _, err := os.Lstat(path); err != nil {
		if isNotExist(err) && opts.AllowMissing {
			return nil
		}
		if isNotExist(err) {
			return pathError(fsport.CodePathNotFound, "path not found", err, path)
		}
		return pathError(fsport.CodePermissionDenied, "filesystem stat failed", err, path)
	}
	if err := os.RemoveAll(path); err != nil {
		return pathError(fsport.CodePermissionDenied, "remove failed", err, path)
	}
	return nil
}
