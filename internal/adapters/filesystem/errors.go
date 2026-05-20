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
	"errors"

	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
)

// fsError creates the structured filesystem-port error used by this adapter.
func fsError(code porterr.Code, message string, cause error, details porterr.Details) error {
	return porterr.New(porterr.KindFilesystem, code, message, cause).WithDetails(details)
}

// pathError is the common one-path error shape used by local filesystem calls.
func pathError(code porterr.Code, message string, cause error, path string) error {
	return fsError(code, message, cause, porterr.Details{"path": path})
}

// isPortError detects errors that already carry a stable port code.
//
// Tree walkers wrap many low-level callback errors. This helper lets callers
// preserve a more precise error, such as fs_symlink_rejected, instead of
// replacing it with a generic copy/hash/cleanup failure.
func isPortError(err error) bool {
	var perr *porterr.Error
	return errors.As(err, &perr)
}
