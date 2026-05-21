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

package filesystem

import (
	"path"
	"strings"
)

// slash converts a relative filesystem path to stable slash form.
//
// Include, exclude, preserve, and hash paths all use slash form so tests and
// publisher behavior remain stable across operating systems.
func slash(pathValue string) string {
	return strings.TrimPrefix(strings.ReplaceAll(pathValue, "\\", "/"), "./")
}

// shouldInclude applies publisher-style include/exclude filtering to rel.
//
// Include patterns narrow the candidate set first, then exclude patterns remove
// entries from that set. The .git directory is always skipped because publisher
// tree operations should copy content, not repository metadata.
func shouldInclude(rel string, include, exclude []string) bool {
	rel = slash(rel)
	if isGitMetadataPath(rel) {
		return false
	}

	if len(include) > 0 {
		if !matchesAny(include, rel) {
			return false
		}
	}

	if matchesAny(exclude, rel) {
		return false
	}

	return true
}

// matchPattern implements the small glob dialect used by this adapter.
//
// It supports exact relative paths, filepath-style globs, a trailing "/**" for
// subtrees, and a leading "**/" for basename matches anywhere in the tree.
func matchPattern(pattern, name string) bool {
	pattern = slash(pattern)
	name = slash(name)
	if pattern == "" {
		return false
	}
	if pattern == name {
		return true
	}

	if isSubtreePatternMatch(pattern, name) {
		return true
	}

	if strings.HasPrefix(pattern, "**/") {
		tail := strings.TrimPrefix(pattern, "**/")
		if ok, _ := path.Match(tail, path.Base(name)); ok {
			return true
		}
		if ok, _ := path.Match(tail, name); ok {
			return true
		}
	}
	ok, _ := path.Match(pattern, name)
	if ok {
		return true
	}
	ok, _ = path.Match(pattern, path.Base(name))
	return ok
}

func isGitMetadataPath(rel string) bool {
	return rel == ".git" || strings.HasPrefix(rel, ".git/")
}

func matchesAny(patterns []string, name string) bool {
	for _, pattern := range patterns {
		if matchPattern(pattern, name) {
			return true
		}
	}
	return false
}

func isSubtreePatternMatch(pattern string, name string) bool {
	if !strings.HasSuffix(pattern, "/**") {
		return false
	}
	prefix := strings.TrimSuffix(pattern, "/**") + "/"
	return strings.HasPrefix(name, prefix)
}
