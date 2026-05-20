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
	"sort"
	"strings"
)

// normalizePreserve canonicalizes preserve paths for stable comparison.
//
// Values are slash-form, relative, non-empty paths. Sorting is not required for
// correctness, but it makes debugging and tests deterministic.
func normalizePreserve(values []string) []string {
	out := make([]string, 0, len(values))
	for _, value := range values {
		value = slash(value)
		value = strings.Trim(value, "/")
		if value != "" && value != "." {
			out = append(out, value)
		}
	}
	sort.Strings(out)
	return out
}

// shouldPreserve reports whether rel is the preserved path or lives below it.
func shouldPreserve(rel string, preserve []string) bool {
	for _, item := range preserve {
		if rel == item || strings.HasPrefix(rel, item+"/") {
			return true
		}
	}
	return false
}

// hasPreservedChild reports whether a directory contains a preserved descendant.
//
// CleanDir uses this to avoid deleting a parent directory that contains one
// nested file or directory selected by Preserve.
func hasPreservedChild(rel string, preserve []string) bool {
	prefix := strings.TrimSuffix(rel, "/") + "/"
	for _, item := range preserve {
		if strings.HasPrefix(item, prefix) {
			return true
		}
	}
	return false
}
