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
	"io"
	"os"
)

// copyFile copies one regular file and returns the number of bytes written.
//
// os.OpenFile applies process umask to new files, so Chmod is called after the
// copy to make PreserveMode deterministic when the platform supports it.
func copyFile(src, dst string, mode os.FileMode, preserveMode bool) (int64, error) {
	in, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer in.Close()
	perm := copyFilePerm(mode, preserveMode)
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return 0, err
	}
	defer out.Close()
	written, err := io.Copy(out, in)
	if err != nil {
		return written, err
	}
	if err := out.Chmod(perm); err != nil {
		return written, err
	}
	return written, nil
}

// copyFilePerm chooses the destination mode for a copied regular file.
func copyFilePerm(mode os.FileMode, preserveMode bool) os.FileMode {
	if preserveMode {
		return mode.Perm()
	}
	return 0o644
}
