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

// hashEntry is the normalized representation folded into the final tree hash.
//
// The content field stores either a regular-file digest, a symlink target, or an
// empty string for directories. The final renderer decides whether mode is part
// of the digest based on TreeHashOptions.
type hashEntry struct {
	rel     string
	mode    os.FileMode
	kind    string
	content string
}

// collectHashEntries walks root and converts included filesystem entries.
func collectHashEntries(ctx context.Context, root string, opts fsport.TreeHashOptions) ([]hashEntry, error) {
	collector := hashCollector{ctx: ctx, root: root, opts: opts, policy: symlinkMode(opts.SymlinkPolicy)}
	entries := []hashEntry{}
	err := filepath.WalkDir(root, func(pathValue string, entry os.DirEntry, walkErr error) error {
		item, include, err := collector.entry(pathValue, entry, walkErr)
		if err != nil || !include {
			return err
		}
		entries = append(entries, item)
		return nil
	})
	return entries, err
}

// hashCollector holds the stable options needed while walking the tree.
type hashCollector struct {
	ctx    context.Context
	root   string
	opts   fsport.TreeHashOptions
	policy fsport.SymlinkPolicy
}

// entry converts one WalkDir record into a hashEntry when it should participate.
func (c hashCollector) entry(pathValue string, entry os.DirEntry, walkErr error) (hashEntry, bool, error) {
	if walkErr != nil {
		return hashEntry{}, false, walkErr
	}
	if err := c.ctx.Err(); err != nil {
		return hashEntry{}, false, err
	}
	rel, err := filepath.Rel(c.root, pathValue)
	if err != nil {
		return hashEntry{}, false, err
	}
	if rel == "." {
		return hashEntry{}, false, nil
	}
	relSlash := slash(rel)
	if !shouldInclude(relSlash, c.opts.Include, c.opts.Exclude) {
		if entry.IsDir() {
			return hashEntry{}, false, filepath.SkipDir
		}
		return hashEntry{}, false, nil
	}
	return c.includedEntry(pathValue, entry, relSlash)
}

// includedEntry converts an already-filtered filesystem entry by kind.
func (c hashCollector) includedEntry(pathValue string, entry os.DirEntry, relSlash string) (hashEntry, bool, error) {
	info, err := entry.Info()
	if err != nil {
		return hashEntry{}, false, err
	}
	mode := info.Mode()
	item := hashEntry{rel: relSlash, mode: mode.Perm()}
	switch {
	case mode&os.ModeSymlink != 0:
		return c.symlinkEntry(pathValue, item)
	case entry.IsDir():
		item.kind = "dir"
		return item, true, nil
	case mode.IsRegular():
		return regularFileHashEntry(pathValue, item)
	default:
		return hashEntry{}, false, nil
	}
}

// symlinkEntry applies symlink policy before adding a link to the hash stream.
func (c hashCollector) symlinkEntry(pathValue string, item hashEntry) (hashEntry, bool, error) {
	if c.policy == fsport.SymlinkReject {
		return hashEntry{}, false, pathError(fsport.CodeSymlinkRejected, "symbolic link rejected", nil, pathValue)
	}
	if !c.opts.IncludeSymlinks {
		return hashEntry{}, false, nil
	}
	target, err := os.Readlink(pathValue)
	if err != nil {
		return hashEntry{}, false, err
	}
	item.kind = "symlink"
	item.content = target
	return item, true, nil
}

// regularFileHashEntry hashes file content and creates the file entry.
func regularFileHashEntry(pathValue string, item hashEntry) (hashEntry, bool, error) {
	digest, err := fileDigest(pathValue)
	if err != nil {
		return hashEntry{}, false, err
	}
	item.kind = "file"
	item.content = digest
	return item, true, nil
}
