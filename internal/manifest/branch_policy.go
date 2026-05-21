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

import "strings"

// hasRejectedBranchShape detects branch names that are unsafe for Git refs.
//
// The policy is intentionally a small manifest-level guard rather than a full
// reimplementation of git-check-ref-format. It rejects the dangerous shapes
// that could be confused with revision syntax, path traversal, globbing, or
// lock files before any workflow reaches Git.
func hasRejectedBranchShape(value string) bool {
	if strings.ContainsAny(value, "~^:?*[\\") {
		return true
	}
	if strings.Contains(value, "..") ||
		strings.Contains(value, "//") ||
		strings.Contains(value, "@{") {
		return true
	}
	if strings.HasPrefix(value, "/") ||
		strings.HasSuffix(value, "/") ||
		strings.HasSuffix(value, ".") {
		return true
	}
	return value == "." || strings.HasSuffix(value, ".lock")
}
