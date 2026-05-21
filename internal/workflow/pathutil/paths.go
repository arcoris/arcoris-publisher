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

package pathutil

import (
	"fmt"
	"path/filepath"
	"strings"
)

// CleanAbs rejects blank paths and returns a cleaned absolute path.
func CleanAbs(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		return "", fmt.Errorf("path is required")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.Clean(abs), nil
}

// EnsureInside rejects a resolved path that escapes root.
func EnsureInside(root string, path string) error {
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

// JoinRelative joins a slash-separated relative path under root.
func JoinRelative(root string, rel fmt.Stringer) string {
	if rel.String() == "." {
		return filepath.Clean(root)
	}
	return filepath.Clean(filepath.Join(root, filepath.FromSlash(rel.String())))
}
