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
	"unicode"
)

// validateNonEmptyToken rejects empty and all-whitespace manifest scalar values.
func validateNonEmptyToken(name string, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", name)
	}
	return nil
}

// hasASCIIWhitespace reports whether value contains any Unicode whitespace.
func hasASCIIWhitespace(value string) bool {
	for _, r := range value {
		if unicode.IsSpace(r) {
			return true
		}
	}
	return false
}

// hasPathTraversal reports whether a slash path escapes its root.
func hasPathTraversal(value string) bool {
	clean := path.Clean(value)
	return clean == ".." || strings.HasPrefix(clean, "../") || strings.Contains(clean, "/../")
}

// containsAny reports whether value contains any rune in chars.
func containsAny(value string, chars string) bool {
	return strings.ContainsAny(value, chars)
}

// defaultString returns fallback when value is empty.
func defaultString(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
