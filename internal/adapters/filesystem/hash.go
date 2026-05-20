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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	fsport "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
)

type hashEntry struct {
	rel     string
	mode    os.FileMode
	kind    string
	content string
}

// TreeHash walks root and computes a deterministic content digest.
//
// The walker gathers normalized entries first and sorts them before hashing so
// filesystem iteration order cannot affect the digest. File content is hashed
// separately, then folded into the tree hash with relative path and selected
// metadata.
func (fs *FileSystem) TreeHash(ctx context.Context, root string, opts fsport.TreeHashOptions) (fsport.TreeHash, error) {
	policy := symlinkMode(opts.SymlinkPolicy)
	entries := []hashEntry{}
	err := filepath.WalkDir(root, func(pathValue string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		rel, err := filepath.Rel(root, pathValue)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		relSlash := slash(rel)
		if !shouldInclude(relSlash, opts.Include, opts.Exclude) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		mode := info.Mode()
		h := hashEntry{rel: relSlash, mode: mode.Perm()}
		switch {
		case mode&os.ModeSymlink != 0:
			if policy == fsport.SymlinkReject {
				return pathError(fsport.CodeSymlinkRejected, "symbolic link rejected", nil, pathValue)
			}
			if !opts.IncludeSymlinks {
				return nil
			}
			target, err := os.Readlink(pathValue)
			if err != nil {
				return err
			}
			h.kind = "symlink"
			h.content = target
		case entry.IsDir():
			h.kind = "dir"
		case mode.IsRegular():
			digest, err := fileDigest(pathValue)
			if err != nil {
				return err
			}
			h.kind = "file"
			h.content = digest
		default:
			return nil
		}
		entries = append(entries, h)
		return nil
	})
	if err != nil {
		if isPortError(err) {
			return "", err
		}
		if isNotExist(err) {
			return "", pathError(fsport.CodePathNotFound, "tree hash root not found", err, root)
		}
		return "", fsError(fsport.CodeTreeHashFailed, "tree hash failed", err, porterrDetails("root", root))
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].rel < entries[j].rel })
	sum := sha256.New()
	for _, entry := range entries {
		if opts.IncludeFileMode {
			fmt.Fprintf(sum, "%s %s %o %s\n", entry.kind, entry.rel, entry.mode, entry.content)
		} else {
			fmt.Fprintf(sum, "%s %s %s\n", entry.kind, entry.rel, entry.content)
		}
	}
	return fsport.TreeHash("sha256:" + hex.EncodeToString(sum.Sum(nil))), nil
}

// fileDigest returns the sha256 digest of one regular file.
func fileDigest(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	sum := sha256.New()
	if _, err := io.Copy(sum, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(sum.Sum(nil)), nil
}
