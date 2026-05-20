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
	"strings"
)

// ModulePath is the public Go module path for a module.
type ModulePath string

// ParseModulePath validates a public Go module path.
func ParseModulePath(value string) (ModulePath, error) {
	if err := validateNonEmptyToken("module path", value); err != nil {
		return "", err
	}
	if hasASCIIWhitespace(value) {
		return "", fmt.Errorf("module path must not contain whitespace")
	}
	if strings.Contains(value, "\\") {
		return "", fmt.Errorf("module path must use slash separators")
	}
	if !strings.Contains(value, "/") && !strings.Contains(value, ".") {
		return "", fmt.Errorf("module path should contain a domain or slash-qualified path")
	}
	if hasPathTraversal(value) || strings.HasPrefix(value, "/") {
		return "", fmt.Errorf("module path must not be absolute or contain path traversal")
	}
	return ModulePath(value), nil
}
