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
	"path/filepath"

	fsport "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
)

// copyTreeOperation carries immutable copy settings plus the mutable result.
//
// WalkDir callbacks otherwise need a long parameter list. Keeping state in this
// small value makes each helper read as a step in one copy operation.
type copyTreeOperation struct {
	fs     *FileSystem
	ctx    context.Context
	src    string
	dst    string
	opts   fsport.CopyTreeOptions
	policy fsport.SymlinkPolicy
	result *fsport.CopyTreeResult
}

// visit handles one WalkDir entry for CopyTree.
func (op copyTreeOperation) visit(pathValue string, entry os.DirEntry, walkErr error) error {
	if walkErr != nil {
		return walkErr
	}
	if err := op.ctx.Err(); err != nil {
		return err
	}
	rel, err := filepath.Rel(op.src, pathValue)
	if err != nil {
		return err
	}
	if rel == "." {
		return op.createRoot()
	}
	if !op.includeEntry(rel, entry) {
		return op.skipEntry(entry)
	}
	return op.copyIncludedEntry(pathValue, rel, entry)
}

// createRoot ensures the destination exists before copying children into it.
func (op copyTreeOperation) createRoot() error {
	return os.MkdirAll(op.dst, 0o755)
}

// includeEntry applies include/exclude matching to one relative path.
func (op copyTreeOperation) includeEntry(rel string, entry os.DirEntry) bool {
	return shouldInclude(slash(rel), op.opts.Include, op.opts.Exclude)
}

// skipEntry updates metrics and prunes filtered directories.
func (op copyTreeOperation) skipEntry(entry os.DirEntry) error {
	if entry.IsDir() {
		return filepath.SkipDir
	}
	op.result.FilesSkipped++
	return nil
}
