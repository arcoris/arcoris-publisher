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

package manifest

import (
	"fmt"
	"path"
	"strings"
)

// RelativePath is a normalized slash-separated relative path from a manifest.
type RelativePath string

// ParseRelativePath validates and normalizes a slash-separated relative path.
func ParseRelativePath(field string, value string, allowDot bool) (RelativePath, error) {
	if strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("%s is required", field)
	}
	if strings.Contains(value, "\\") {
		return "", fmt.Errorf("%s must use '/' as path separator", field)
	}
	if strings.TrimSpace(value) != value {
		return "", fmt.Errorf("%s must not have surrounding whitespace", field)
	}
	if strings.HasPrefix(value, "/") || looksWindowsAbsolute(value) {
		return "", fmt.Errorf("%s must be relative", field)
	}
	if hasParentTraversalSegment(value) {
		return "", fmt.Errorf("%s must not contain path traversal", field)
	}
	cleaned := path.Clean(value)
	if cleaned == "." && !allowDot {
		return "", fmt.Errorf("%s must not be '.'", field)
	}
	return RelativePath(cleaned), nil
}

// String returns the normalized path string.
func (p RelativePath) String() string { return string(p) }
