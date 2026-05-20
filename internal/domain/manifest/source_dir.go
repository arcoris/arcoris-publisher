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

// SourceDir is a repository-relative staged module directory.
type SourceDir string

// ParseSourceDir validates a repository-relative source directory.
func ParseSourceDir(value string) (SourceDir, error) {
	if err := validateNonEmptyToken("source directory", value); err != nil {
		return "", err
	}
	if strings.Contains(value, "\\") {
		return "", fmt.Errorf("source directory must use slash separators")
	}
	if path.IsAbs(value) {
		return "", fmt.Errorf("source directory must be relative")
	}
	clean := path.Clean(value)
	if cleanEscapesRoot(clean) {
		return "", fmt.Errorf("source directory must not escape repository root")
	}
	return SourceDir(clean), nil
}

// cleanEscapesRoot reports whether a cleaned relative path escapes the root.
func cleanEscapesRoot(clean string) bool {
	return clean == "." || clean == ".." || strings.HasPrefix(clean, "../") || strings.Contains(clean, "/../")
}
