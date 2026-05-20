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

import "fmt"

// ModuleName is a stable manifest-local module identifier.
type ModuleName string

// ParseModuleName validates a module identifier.
func ParseModuleName(value string) (ModuleName, error) {
	if err := validateNonEmptyToken("module name", value); err != nil {
		return "", err
	}
	for _, r := range value {
		if isModuleNameRune(r) {
			continue
		}
		return "", fmt.Errorf("module name %q contains invalid character %q", value, r)
	}
	if startsWithPunctuation(value) {
		return "", fmt.Errorf("module name must not start with punctuation")
	}
	return ModuleName(value), nil
}

// isModuleNameRune reports whether r is allowed inside a module name.
func isModuleNameRune(r rune) bool {
	return r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' || r == '_'
}

// startsWithPunctuation reports whether a token starts with a punctuation char forbidden here.
func startsWithPunctuation(value string) bool {
	return value[0] == '-' || value[0] == '_'
}

// quoteModuleName renders a module name for validation messages.
func quoteModuleName(name ModuleName) string {
	return fmt.Sprintf("%q", name)
}
