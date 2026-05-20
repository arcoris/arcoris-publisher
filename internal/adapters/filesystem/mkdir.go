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

// MkdirAll creates a directory tree using the adapter's default permission when
// the caller leaves MkdirOptions.Perm at zero.
func (fs *FileSystem) MkdirAll(ctx context.Context, path string, opts fsport.MkdirOptions) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := os.MkdirAll(path, mkdirPerm(opts)); err != nil {
		return pathError(fsport.CodePermissionDenied, "directory creation failed", err, path)
	}
	return nil
}

// mkdirPerm applies the adapter default for created directories.
func mkdirPerm(opts fsport.MkdirOptions) os.FileMode {
	if opts.Perm != 0 {
		return opts.Perm
	}
	return 0o755
}
