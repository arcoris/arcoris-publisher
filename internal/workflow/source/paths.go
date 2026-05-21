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

package source

import (
	"fmt"
	"path/filepath"
	"strings"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

// cleanAbs rejects empty input and returns a cleaned absolute path.
func cleanAbs(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("path is required")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.Clean(abs), nil
}

// ensureInside rejects paths that escape root after filepath.Rel normalization.
func ensureInside(root string, path string) error {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return err
	}
	if rel == "." {
		return nil
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("%q escapes root %q", path, root)
	}
	return nil
}

// joinRelative joins manifest slash-separated relative paths to a platform path.
func joinRelative(root string, rel fmt.Stringer) string {
	return filepath.Clean(filepath.Join(root, filepath.FromSlash(rel.String())))
}

// resolveModuleDir resolves a module SourceDir under the staging root.
func resolveModuleDir(stagingDir string, sourceDir manifest.SourceDir) string {
	return filepath.Clean(filepath.Join(stagingDir, filepath.FromSlash(sourceDir.String())))
}

// resolveModuleRootDir resolves a module root under the module source dir.
func resolveModuleRootDir(moduleDir string, moduleRoot manifest.RelativePath) string {
	return joinRelative(moduleDir, moduleRoot)
}

// resolveEntrySource resolves one explicit publish entry source path.
func resolveEntrySource(moduleRootDir string, entry manifest.PublishEntry) string {
	return joinRelative(moduleRootDir, entry.From())
}
