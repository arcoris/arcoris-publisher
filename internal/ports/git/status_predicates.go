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
