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
)

// copyIncludedEntry dispatches a filtered-in source entry by filesystem kind.
func (op copyTreeOperation) copyIncludedEntry(pathValue string, rel string, entry os.DirEntry) error {
	target := filepath.Join(op.dst, rel)
	info, err := entry.Info()
	if err != nil {
		return err
	}
	mode := info.Mode()
	if mode&os.ModeSymlink != 0 {
		return op.fs.copySymlink(pathValue, target, op.policy, op.result)
	}
	if entry.IsDir() {
		return op.copyDirectory(target, mode)
	}
	if mode.IsRegular() {
		return op.copyRegularFile(pathValue, target, mode, info)
	}
	op.result.FilesSkipped++
	return nil
}

// copyDirectory creates one destination directory and records it in metrics.
func (op copyTreeOperation) copyDirectory(target string, mode os.FileMode) error {
	perm := os.FileMode(0o755)
	if op.opts.PreserveMode {
		perm = mode.Perm()
	}
	if err := os.MkdirAll(target, perm); err != nil {
		return err
	}
	op.result.DirectoriesCopied++
	return nil
}

// copyRegularFile copies one regular file, then applies optional metadata.
func (op copyTreeOperation) copyRegularFile(pathValue string, target string, mode os.FileMode, info os.FileInfo) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	bytes, err := copyFile(pathValue, target, mode, op.opts.PreserveMode)
	if err != nil {
		return err
	}
	op.result.FilesCopied++
	op.result.BytesCopied += bytes
	if op.opts.PreserveTimes {
		_ = os.Chtimes(target, info.ModTime(), info.ModTime())
	}
	return nil
}
