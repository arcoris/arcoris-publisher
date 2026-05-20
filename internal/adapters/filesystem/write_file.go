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
	if err := prepareWriteParent(path, opts); err != nil {
		return err
	}
	if err := checkOverwriteAllowed(path, opts); err != nil {
		return err
	}
	if err := os.WriteFile(path, data, writeFilePerm(opts)); err != nil {
		return pathError(fsport.CodePermissionDenied, "file write failed", err, path)
	}
	return nil
}

// prepareWriteParent creates parent directories when the caller requested it.
func prepareWriteParent(path string, opts fsport.WriteFileOptions) error {
	if !opts.CreateDirs {
		return nil
	}
	parent := filepath.Dir(path)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return pathError(fsport.CodePermissionDenied, "parent directory creation failed", err, parent)
	}
	return nil
}

// checkOverwriteAllowed enforces WriteFileOptions.Overwrite.
func checkOverwriteAllowed(path string, opts fsport.WriteFileOptions) error {
	if opts.Overwrite {
		return nil
	}
	_, err := os.Lstat(path)
	if err == nil {
		return pathError(fsport.CodePermissionDenied, "file already exists", nil, path)
	}
	if !isNotExist(err) {
		return pathError(fsport.CodePermissionDenied, "filesystem stat failed", err, path)
	}
	return nil
}

// writeFilePerm applies the adapter default for newly written files.
func writeFilePerm(opts fsport.WriteFileOptions) os.FileMode {
	if opts.Perm != 0 {
		return opts.Perm
	}
	return 0o644
}
