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

// RepositoryRef identifies a Git hosting repository in owner/name form.
type RepositoryRef string

// ParseRepositoryRef validates a repository reference in owner/name form.
func ParseRepositoryRef(value string) (RepositoryRef, error) {
	if strings.TrimSpace(value) == "" {
		return "", fmt.Errorf("repository is required")
	}
	if strings.TrimSpace(value) != value || strings.ContainsAny(value, " \t\n\r") {
		return "", fmt.Errorf("repository must not contain whitespace")
	}
	parts := strings.Split(value, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", fmt.Errorf("repository must be in owner/name form")
	}
	for _, part := range parts {
		if strings.Contains(part, "..") || strings.HasPrefix(part, ".") || strings.HasSuffix(part, ".") {
			return "", fmt.Errorf("repository contains rejected path-like segment")
		}
	}
	return RepositoryRef(value), nil
}

// String returns the repository reference in owner/name form.
func (r RepositoryRef) String() string { return string(r) }
