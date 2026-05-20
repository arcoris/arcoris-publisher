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
	"io"
	"os"
	"path/filepath"

	fsport "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
)

func (fs *FileSystem) CopyTree(ctx context.Context, src string, dst string, opts fsport.CopyTreeOptions) (fsport.CopyTreeResult, error) {
	var result fsport.CopyTreeResult
	if err := ctx.Err(); err != nil {
		return result, err
	}
	if err := ensureInside(dst, opts.SafetyRoot); err != nil {
		return result, pathError(fsport.CodePathOutsideRoot, "copy destination is outside safety root", err, dst)
	}
	sourceInfo, err := os.Lstat(src)
	if err != nil {
		if isNotExist(err) {
			return result, pathError(fsport.CodePathNotFound, "copy source not found", err, src)
		}
		return result, pathError(fsport.CodePermissionDenied, "copy source stat failed", err, src)
	}
	if !sourceInfo.IsDir() {
		return result, pathError(fsport.CodeCopyFailed, "copy source is not a directory", nil, src)
	}
	policy := symlinkMode(opts.SymlinkPolicy)
	err = filepath.WalkDir(src, func(pathValue string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		rel, err := filepath.Rel(src, pathValue)
		if err != nil {
			return err
		}
		if rel == "." {
			return os.MkdirAll(dst, 0o755)
		}
		relSlash := slash(rel)
		if !shouldInclude(relSlash, opts.Include, opts.Exclude) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			result.FilesSkipped++
			return nil
		}
		target := filepath.Join(dst, rel)
		info, err := entry.Info()
		if err != nil {
			return err
		}
		mode := info.Mode()
		if mode&os.ModeSymlink != 0 {
			return fs.copySymlink(pathValue, target, policy, &result)
		}
		if entry.IsDir() {
			perm := os.FileMode(0o755)
			if opts.PreserveMode {
				perm = mode.Perm()
			}
			if err := os.MkdirAll(target, perm); err != nil {
				return err
			}
			result.DirectoriesCopied++
			return nil
		}
		if !mode.IsRegular() {
			result.FilesSkipped++
			return nil
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		bytes, err := copyFile(pathValue, target, mode, opts.PreserveMode)
		if err != nil {
			return err
		}
		result.FilesCopied++
		result.BytesCopied += bytes
		if opts.PreserveTimes {
			_ = os.Chtimes(target, info.ModTime(), info.ModTime())
		}
		return nil
	})
	if err != nil {
		if isPortError(err) {
			return result, err
		}
		return result, fsError(fsport.CodeCopyFailed, "tree copy failed", err, porterrDetails("src", src, "dst", dst))
	}
	return result, nil
}

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
	perm := os.FileMode(0o644)
	if preserveMode {
		perm = mode.Perm()
	}
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

// porterrDetails builds small non-secret detail maps for structured errors.
func porterrDetails(kv ...string) map[string]string {
	m := map[string]string{}
	for i := 0; i+1 < len(kv); i += 2 {
		m[kv[i]] = kv[i+1]
	}
	return m
}
