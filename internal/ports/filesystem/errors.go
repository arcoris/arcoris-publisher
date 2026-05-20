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

import "arcoris.dev/arcoris-publisher/internal/ports/porterr"

const (
	// CodePathNotFound identifies a missing source or target path.
	CodePathNotFound porterr.Code = "fs_path_not_found"
	// CodePathOutsideRoot identifies an attempted access outside SafetyRoot.
	CodePathOutsideRoot porterr.Code = "fs_path_outside_root"
	// CodeUnsafeRemove identifies a removal rejected by safety checks.
	CodeUnsafeRemove porterr.Code = "fs_unsafe_remove"
	// CodeSymlinkRejected identifies a symlink blocked by policy.
	CodeSymlinkRejected porterr.Code = "fs_symlink_rejected"
	// CodeCopyFailed identifies a tree copy failure.
	CodeCopyFailed porterr.Code = "fs_copy_failed"
	// CodeTreeHashFailed identifies a deterministic tree hashing failure.
	CodeTreeHashFailed porterr.Code = "fs_tree_hash_failed"
	// CodePermissionDenied identifies an operating-system permission denial.
	CodePermissionDenied porterr.Code = "fs_permission_denied"
)
