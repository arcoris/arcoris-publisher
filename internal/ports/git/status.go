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

package git

// Status describes a Git working tree status.
//
// Clean is the adapter's summarized view, while Entries contains path-level
// diagnostics. A defensive caller can use IsDirty, which treats either a false
// Clean flag or non-empty Entries as dirty.
type Status struct {
	// Clean reports whether the working tree has no known changes.
	Clean bool
	// Entries contains changed, staged, unstaged, or untracked paths.
	Entries []StatusEntry
}

// StatusEntry describes one changed, staged, unstaged, or untracked path.
type StatusEntry struct {
	// Path is the repository-relative path reported by Git.
	Path string
	// Code is the porcelain status code associated with Path.
	Code string
}

// HasEntries reports whether the status contains path-level changes.
func (s Status) HasEntries() bool {
	return len(s.Entries) > 0
}

// IsDirty reports whether the status represents a non-clean working tree.
//
// The method is intentionally conservative: if an adapter marks Clean as false
// without listing entries, callers still treat the tree as dirty.
func (s Status) IsDirty() bool {
	return !s.Clean || s.HasEntries()
}
