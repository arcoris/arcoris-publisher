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

import "arcoris.dev/arcoris-publisher/internal/domain/manifest"

// newEmptyRegistry allocates every index map used by a valid Registry.
func newEmptyRegistry(modules []manifest.Module) Registry {
	return Registry{
		modules:      cloneModules(modules),
		byName:       make(map[manifest.ModuleName]int, len(modules)),
		byModulePath: make(map[manifest.ModulePath]int, len(modules)),
		bySourceDir:  make(map[manifest.SourceDir]int, len(modules)),
		byRepository: make(map[manifest.RepositoryRef]int, len(modules)),
		byVisibility: make(map[manifest.Visibility][]int, 3),
		byBranch:     make(map[manifest.ModuleName]map[manifest.BranchName]manifest.BranchMapping, len(modules)),
	}
}
