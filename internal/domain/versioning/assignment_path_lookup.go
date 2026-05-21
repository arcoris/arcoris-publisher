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

package versioning

import "arcoris.dev/arcoris-publisher/internal/domain/manifest"

// VersionOfPath returns the version assigned to modulePath.
func (a Assignments) VersionOfPath(modulePath manifest.ModulePath) (Version, bool) {
	index, ok := a.byPath[modulePath]
	if !ok {
		return "", false
	}
	return a.items[index].Version(), true
}

// ContainsPath reports whether modulePath has an assigned version.
func (a Assignments) ContainsPath(modulePath manifest.ModulePath) bool {
	_, ok := a.byPath[modulePath]
	return ok
}
