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

package config

import (
	"fmt"
	"strings"
)

// splitYAMLKeyValue separates "key: value" lines while rejecting the limited
// YAML subset's unsupported key syntax.
func splitYAMLKeyValue(line yamlLine) (string, string, bool, error) {
	idx := strings.Index(line.text, ":")
	if idx < 0 {
		return "", "", false, expectedYAMLKeyValueError(line)
	}

	key := strings.TrimSpace(line.text[:idx])
	if key == "" {
		return "", "", false, emptyYAMLKeyError(line)
	}
	if strings.ContainsAny(key, " \t[]{}") {
		return "", "", false, unsupportedYAMLKeyError(line, key)
	}

	raw := strings.TrimSpace(line.text[idx+1:])
	return key, raw, raw != "", nil
}

// validateListItemLine reports whether the current line belongs to the active
// list and validates that it is aligned with the list indentation.
func validateListItemLine(line yamlLine, indent int) (bool, error) {
	if line.indent < indent {
		return false, nil
	}
	if line.indent > indent {
		return false, unexpectedIndentationError(line)
	}
	if !isYAMLListItem(line) {
		return false, nil
	}

	return true, nil
}

// isYAMLListItem recognizes the two dash forms supported by the parser.
func isYAMLListItem(line yamlLine) bool {
	return strings.HasPrefix(line.text, "- ") || line.text == "-"
}

// looksLikeInlineMapItem recognizes compact list map entries such as
// "- name: foundation".
func looksLikeInlineMapItem(item string) bool {
	idx := strings.Index(item, ":")
	if idx <= 0 {
		return false
	}

	key := strings.TrimSpace(item[:idx])
	if key == "" {
		return false
	}

	return !strings.ContainsAny(key, " \t[]{}")
}

// expectedYAMLKeyValueError reports a non-empty map line that has no colon.
func expectedYAMLKeyValueError(line yamlLine) error {
	return fmt.Errorf("line %d: expected key-value pair", line.number)
}

// emptyYAMLKeyError reports a colon with no key before it.
func emptyYAMLKeyError(line yamlLine) error {
	return fmt.Errorf("line %d: key must not be empty", line.number)
}

// unsupportedYAMLKeyError reports keys outside the intentionally small manifest
// YAML subset.
func unsupportedYAMLKeyError(line yamlLine, key string) error {
	return fmt.Errorf("line %d: unsupported key %q", line.number, key)
}
