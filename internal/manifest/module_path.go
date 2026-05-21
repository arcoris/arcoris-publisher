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

// ModulePath identifies a published Go module path.
type ModulePath string

// ParseModulePath validates a Go module path used by a module manifest.
func ParseModulePath(value string) (ModulePath, error) {
	if strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("module path is required")
	}
	if strings.TrimSpace(value) != value || strings.ContainsAny(value, " \t\n\r") {
		return "", fmt.Errorf("module path must not contain whitespace")
	}
	if strings.Contains(value, "\\") {
		return "", fmt.Errorf("module path must use '/' as separator")
	}
	if !strings.Contains(value, "/") {
		return "", fmt.Errorf("module path must contain at least one slash")
	}
	if strings.Contains(value, "..") || strings.Contains(value, "//") {
		return "", fmt.Errorf("module path contains rejected path-like sequence")
	}
	if hasUnsafeModulePathSegment(value) {
		return "", fmt.Errorf("module path contains rejected path-like segment")
	}
	return ModulePath(value), nil
}

// String returns the module path string.
func (p ModulePath) String() string { return string(p) }
