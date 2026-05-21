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
	"os"
	"path/filepath"
	"strings"

	fsport "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
)

// absClean normalizes a path before safety comparisons.
//
// Empty paths are rejected because filepath.Abs("") would otherwise resolve to
// the current working directory, which is too surprising for destructive
// filesystem operations.
func absClean(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", errors.New("empty path")
	}
	return filepath.Abs(filepath.Clean(path))
}

// ensureInside verifies that path is equal to or below root.
//
// An empty root means the caller opted into the adapter's default safety policy.
// The check uses filepath.Rel instead of string prefixes so sibling paths like
// /tmp/root-other are not accepted for /tmp/root.
func ensureInside(path, root string) error {
	if root == "" {
		return nil
	}
	absPath, err := absClean(path)
	if err != nil {
		return err
	}
	absRoot, err := absClean(root)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(absRoot, absPath)
	if err != nil {
		return err
	}

	if isPathInsideRoot(rel) {
		return nil
	}

	return errors.New("path is outside safety root")
}

// ensureSafeRemove applies extra checks for destructive removal.
//
// It rejects the filesystem root even when no SafetyRoot is supplied. Callers
// should still provide SafetyRoot whenever possible for stronger protection.
func ensureSafeRemove(path, root string) error {
	abs, err := absClean(path)
	if err != nil {
		return err
	}
	if abs == string(filepath.Separator) {
		return errors.New("refusing to remove filesystem root")
	}
	if err := ensureInside(abs, root); err != nil {
		return err
	}
	return nil
}

// isNotExist centralizes os.ErrNotExist matching for wrapped filesystem errors.
func isNotExist(err error) bool { return errors.Is(err, os.ErrNotExist) }

// symlinkMode chooses the adapter default for unspecified symlink handling.
//
// The ports package makes the zero SymlinkPolicy invalid. This concrete adapter
// defaults to rejection because following or preserving symlinks can be unsafe
// unless the caller has explicitly asked for that behavior.
func symlinkMode(policy fsport.SymlinkPolicy) fsport.SymlinkPolicy {
	if policy == "" {
		return fsport.SymlinkReject
	}
	return policy
}

func isPathInsideRoot(rel string) bool {
	if rel == "." {
		return true
	}
	return !strings.HasPrefix(rel, "..") && !filepath.IsAbs(rel)
}
