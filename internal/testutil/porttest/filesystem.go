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

// Package porttest contains small in-memory port implementations for tests.
package porttest

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"arcoris.dev/arcoris-publisher/internal/ports/filesystem"
)

// FileSystem is a deterministic in-memory filesystem port for workflow tests.
type FileSystem struct {
	// Dirs stores known directory paths.
	Dirs map[string]bool

	// Files stores regular file contents by path.
	Files map[string][]byte
}

// NewFileSystem returns an empty in-memory filesystem.
func NewFileSystem() *FileSystem {
	return &FileSystem{
		Dirs:  map[string]bool{},
		Files: map[string][]byte{},
	}
}

// AddDir registers path and its parents as directories.
func (fs *FileSystem) AddDir(path string) {
	fs.addParents(filepath.Clean(path))
}

// AddFile registers path as a regular file with detached contents.
func (fs *FileSystem) AddFile(path string, data []byte) {
	path = filepath.Clean(path)
	fs.addParents(filepath.Dir(path))
	fs.Files[path] = append([]byte(nil), data...)
}

// Exists reports whether path is registered.
func (fs *FileSystem) Exists(_ context.Context, path string) (bool, error) {
	path = filepath.Clean(path)
	return fs.Dirs[path] || fs.Files[path] != nil, nil
}

// IsDir reports whether path is registered as a directory.
func (fs *FileSystem) IsDir(_ context.Context, path string) (bool, error) {
	return fs.Dirs[filepath.Clean(path)], nil
}

// ReadFile returns detached file contents.
func (fs *FileSystem) ReadFile(_ context.Context, path string) ([]byte, error) {
	path = filepath.Clean(path)
	data, ok := fs.Files[path]
	if !ok {
		return nil, fmt.Errorf("file %s not found", path)
	}

	return append([]byte(nil), data...), nil
}

// WriteFile stores detached file contents.
func (fs *FileSystem) WriteFile(
	_ context.Context,
	path string,
	data []byte,
	opts filesystem.WriteFileOptions,
) error {
	path = filepath.Clean(path)
	if fs.Files[path] != nil && !opts.Overwrite {
		return fmt.Errorf("file %s already exists", path)
	}
	if opts.CreateDirs {
		fs.addParents(filepath.Dir(path))
	}

	fs.Files[path] = append([]byte(nil), data...)
	return nil
}

// MkdirAll registers path and its parents as directories.
func (fs *FileSystem) MkdirAll(_ context.Context, path string, _ filesystem.MkdirOptions) error {
	fs.addParents(filepath.Clean(path))
	return nil
}

// RemoveAll removes path and descendants.
func (fs *FileSystem) RemoveAll(_ context.Context, path string, _ filesystem.RemoveOptions) error {
	fs.removeTree(filepath.Clean(path))
	return nil
}

// CleanDir removes directory contents while preserving the directory itself.
func (fs *FileSystem) CleanDir(_ context.Context, dir string, opts filesystem.CleanDirOptions) error {
	dir = filepath.Clean(dir)
	if !fs.Dirs[dir] {
		if opts.AllowMissing {
			return nil
		}
		return fmt.Errorf("directory %s not found", dir)
	}

	for _, child := range fs.children(dir) {
		if shouldPreserve(dir, child, opts.Preserve) {
			continue
		}
		fs.removeTree(child)
	}

	return nil
}

// CopyTree copies files under src into dst.
func (fs *FileSystem) CopyTree(
	_ context.Context,
	src string,
	dst string,
	_ filesystem.CopyTreeOptions,
) (filesystem.CopyTreeResult, error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)
	if !fs.Dirs[src] {
		return filesystem.CopyTreeResult{}, fmt.Errorf("directory %s not found", src)
	}

	fs.addParents(dst)
	result := filesystem.CopyTreeResult{DirectoriesCopied: 1}
	for _, path := range fs.sortedFiles() {
		if !strings.HasPrefix(path, src+string(filepath.Separator)) {
			continue
		}

		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		fs.AddFile(target, fs.Files[path])
		result.FilesCopied++
		result.BytesCopied += int64(len(fs.Files[path]))
	}

	return result, nil
}

// TreeHash returns a stable synthetic hash for present files under root.
func (fs *FileSystem) TreeHash(
	_ context.Context,
	root string,
	_ filesystem.TreeHashOptions,
) (filesystem.TreeHash, error) {
	root = filepath.Clean(root)
	if !fs.Dirs[root] {
		return "", fmt.Errorf("directory %s not found", root)
	}

	count := 0
	for _, path := range fs.sortedFiles() {
		if strings.HasPrefix(path, root+string(filepath.Separator)) {
			count++
		}
	}

	return filesystem.TreeHash(fmt.Sprintf("sha256:test:%s:%d", root, count)), nil
}

func (fs *FileSystem) addParents(path string) {
	for path != "." && path != "" {
		fs.Dirs[path] = true
		parent := filepath.Dir(path)
		if parent == path {
			break
		}
		path = parent
	}
}

func (fs *FileSystem) children(dir string) []string {
	seen := map[string]struct{}{}
	prefix := dir + string(filepath.Separator)
	for path := range fs.Dirs {
		child, ok := directChild(prefix, path)
		if ok {
			seen[child] = struct{}{}
		}
	}
	for path := range fs.Files {
		child, ok := directChild(prefix, path)
		if ok {
			seen[child] = struct{}{}
		}
	}

	children := make([]string, 0, len(seen))
	for child := range seen {
		children = append(children, child)
	}
	sort.Strings(children)
	return children
}

func directChild(prefix, path string) (string, bool) {
	if !strings.HasPrefix(path, prefix) {
		return "", false
	}
	rest := strings.TrimPrefix(path, prefix)
	if rest == "" {
		return "", false
	}
	first := strings.SplitN(rest, string(filepath.Separator), 2)[0]
	return filepath.Join(strings.TrimSuffix(prefix, string(filepath.Separator)), first), true
}

func shouldPreserve(root, path string, preserve []string) bool {
	for _, rel := range preserve {
		if filepath.Join(root, filepath.FromSlash(rel)) == path {
			return true
		}
	}
	return false
}

func (fs *FileSystem) removeTree(path string) {
	delete(fs.Dirs, path)
	delete(fs.Files, path)
	prefix := path + string(filepath.Separator)
	for child := range fs.Dirs {
		if strings.HasPrefix(child, prefix) {
			delete(fs.Dirs, child)
		}
	}
	for child := range fs.Files {
		if strings.HasPrefix(child, prefix) {
			delete(fs.Files, child)
		}
	}
}

func (fs *FileSystem) sortedFiles() []string {
	files := make([]string, 0, len(fs.Files))
	for path := range fs.Files {
		files = append(files, path)
	}
	sort.Strings(files)
	return files
}
