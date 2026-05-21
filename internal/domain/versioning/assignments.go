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

// Assignments is an immutable-by-convention set of publish version assignments.
//
// Items preserve registry declaration order while lookup maps provide constant
// time reads by module name and module path. The maps are rebuilt by Validate,
// which keeps manually assembled test values and constructed values consistent.
type Assignments struct {
	policy manifest.VersionPolicy
	items  []ModuleVersion

	byModule map[manifest.ModuleName]int
	byPath   map[manifest.ModulePath]int
}
