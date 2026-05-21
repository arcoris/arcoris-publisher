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

package registry

import "arcoris.dev/arcoris-publisher/internal/manifest"

// buildIndexes creates stable value-to-module indexes after validation has
// already proven every indexed value unique.
func (r *Registry) buildIndexes() {
	r.byName = make(map[manifest.ModuleName]int, len(r.modules))
	r.byPath = make(map[manifest.ModulePath]int, len(r.modules))
	r.byRepository = make(map[manifest.RepositoryRef]int, len(r.modules))
	r.bySourceDir = make(map[manifest.SourceDir]int, len(r.modules))

	for i, module := range r.modules {
		r.byName[module.Name()] = i
		r.byPath[module.ModulePath()] = i
		r.byRepository[module.Repository()] = i
		r.bySourceDir[module.SourceDir()] = i
	}
}
