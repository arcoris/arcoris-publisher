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

// RemoveAll removes a path after applying destructive-operation safety checks.
func (fs *FileSystem) RemoveAll(ctx context.Context, path string, opts fsport.RemoveOptions) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := ensureSafeRemove(path, opts.SafetyRoot); err != nil {
		return pathError(fsport.CodeUnsafeRemove, "unsafe remove target", err, path)
	}
	if err := validateRemoveTarget(path, opts); err != nil {
		return err
	}
	if err := os.RemoveAll(path); err != nil {
		return pathError(fsport.CodePermissionDenied, "remove failed", err, path)
	}
	return nil
}

// validateRemoveTarget maps missing paths according to RemoveOptions.AllowMissing.
func validateRemoveTarget(path string, opts fsport.RemoveOptions) error {
	if _, err := os.Lstat(path); err != nil {
		if isNotExist(err) && opts.AllowMissing {
			return nil
		}
		if isNotExist(err) {
			return pathError(fsport.CodePathNotFound, "path not found", err, path)
		}
		return pathError(fsport.CodePermissionDenied, "filesystem stat failed", err, path)
	}
	return nil
}
