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

	fsport "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
)

// Exists reports whether path exists without following symlinks.
func (fs *FileSystem) Exists(ctx context.Context, path string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	_, err := os.Lstat(path)
	if err == nil {
		return true, nil
	}
	if isNotExist(err) {
		return false, nil
	}
	return false, pathError(fsport.CodePermissionDenied, "filesystem stat failed", err, path)
}

// IsDir reports whether path exists and is a directory.
//
// Symlinks are not followed; a symlink to a directory therefore reports false.
func (fs *FileSystem) IsDir(ctx context.Context, path string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	info, err := os.Lstat(path)
	if err == nil {
		return info.IsDir(), nil
	}
	if isNotExist(err) {
		return false, nil
	}
	return false, pathError(fsport.CodePermissionDenied, "filesystem stat failed", err, path)
}

// ReadFile reads a regular file through os.ReadFile and maps failures to the
// filesystem port's stable error codes.
func (fs *FileSystem) ReadFile(ctx context.Context, path string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err == nil {
		return data, nil
	}
	if isNotExist(err) {
		return nil, pathError(fsport.CodePathNotFound, "file not found", err, path)
	}
	return nil, pathError(fsport.CodePermissionDenied, "file read failed", err, path)
}
