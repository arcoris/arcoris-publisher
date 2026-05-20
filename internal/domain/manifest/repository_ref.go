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

// RepositoryRef is a target repository reference in owner/name form.
type RepositoryRef string

// ParseRepositoryRef validates a repository reference.
func ParseRepositoryRef(value string) (RepositoryRef, error) {
	if err := validateNonEmptyToken("repository", value); err != nil {
		return "", err
	}
	if hasASCIIWhitespace(value) {
		return "", fmt.Errorf("repository must not contain whitespace")
	}
	parts := strings.Split(value, "/")
	if !isOwnerNamePair(parts) {
		return "", fmt.Errorf("repository must use owner/name form")
	}
	if repositorySegmentsTraverse(parts) {
		return "", fmt.Errorf("repository segments must not be path traversal markers")
	}
	return RepositoryRef(value), nil
}

// Owner returns the repository owner segment.
func (r RepositoryRef) Owner() string {
	parts := strings.SplitN(string(r), "/", 2)
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

// Name returns the repository name segment.
func (r RepositoryRef) Name() string {
	parts := strings.SplitN(string(r), "/", 2)
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// isOwnerNamePair reports whether parts has non-empty owner and name values.
func isOwnerNamePair(parts []string) bool {
	return len(parts) == 2 && parts[0] != "" && parts[1] != ""
}

// repositorySegmentsTraverse rejects "." and ".." as repository path segments.
func repositorySegmentsTraverse(parts []string) bool {
	for _, part := range parts {
		if part == "." || part == ".." {
			return true
		}
	}
	return false
}
