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

// copySymlink handles one symbolic link according to the requested policy.
//
// Following symlinks is deliberately unsupported in this local adapter for now:
// safe following requires root-escape checks after resolving every link in a
// chain, and rejecting it is safer than quietly doing an unsafe partial follow.
func (fs *FileSystem) copySymlink(src, dst string, policy fsport.SymlinkPolicy, result *fsport.CopyTreeResult) error {
	switch policy {
	case fsport.SymlinkReject:
		return pathError(fsport.CodeSymlinkRejected, "symbolic link rejected", nil, src)
	case fsport.SymlinkPreserve:
		target, err := os.Readlink(src)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return err
		}
		_ = os.Remove(dst)
		if err := os.Symlink(target, dst); err != nil {
			return err
		}
		result.FilesCopied++
		return nil
	case fsport.SymlinkFollow:
		return pathError(fsport.CodeSymlinkRejected, "following symbolic links is not supported by this adapter", nil, src)
	default:
		return pathError(fsport.CodeSymlinkRejected, "unknown symbolic link policy", nil, src)
	}
}
