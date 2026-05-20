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
	"path/filepath"

	fsport "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
)

// CopyTree copies a directory tree from src to dst with filtering and symlink policy.
//
// The method now delegates validation, walking, entry handling, and file-copy
// mechanics to focused helpers. That keeps the public operation readable while
// preserving the detailed behavior required by the filesystem port.
func (fs *FileSystem) CopyTree(ctx context.Context, src string, dst string, opts fsport.CopyTreeOptions) (fsport.CopyTreeResult, error) {
	var result fsport.CopyTreeResult
	if err := ctx.Err(); err != nil {
		return result, err
	}
	if err := validateCopyTree(src, dst, opts); err != nil {
		return result, err
	}
	op := copyTreeOperation{fs: fs, ctx: ctx, src: src, dst: dst, opts: opts, policy: symlinkMode(opts.SymlinkPolicy), result: &result}
	err := filepath.WalkDir(src, op.visit)
	if err != nil {
		if isPortError(err) {
			return result, err
		}
		return result, fsError(fsport.CodeCopyFailed, "tree copy failed", err, porterrDetails("src", src, "dst", dst))
	}
	return result, nil
}
