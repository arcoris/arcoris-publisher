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

	fsport "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
)

// validateCopyTree verifies source and destination constraints before walking.
//
// Destination safety is checked before source inspection so an unsafe target is
// rejected even if the source also happens to be invalid.
func validateCopyTree(src string, dst string, opts fsport.CopyTreeOptions) error {
	if err := ensureInside(dst, opts.SafetyRoot); err != nil {
		return pathError(fsport.CodePathOutsideRoot, "copy destination is outside safety root", err, dst)
	}
	sourceInfo, err := os.Lstat(src)
	if err != nil {
		return copySourceStatError(src, err)
	}
	if !sourceInfo.IsDir() {
		return pathError(fsport.CodeCopyFailed, "copy source is not a directory", nil, src)
	}
	return nil
}

// copySourceStatError maps source Lstat failures into precise port errors.
func copySourceStatError(src string, err error) error {
	if isNotExist(err) {
		return pathError(fsport.CodePathNotFound, "copy source not found", err, src)
	}
	return pathError(fsport.CodePermissionDenied, "copy source stat failed", err, src)
}
