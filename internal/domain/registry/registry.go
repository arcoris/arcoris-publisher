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

// Registry is an immutable-by-convention module lookup index.
//
// Registry preserves manifest declaration order while providing constant-time
// lookup by the identifiers used by later planning and workflow packages. The
// registry owns detached copies of all slices and maps constructed from input
// modules; accessors return detached slices so callers cannot mutate registry
// state after construction.
type Registry struct {
	modules []manifest.Module

	byName       map[manifest.ModuleName]int
	byModulePath map[manifest.ModulePath]int
	bySourceDir  map[manifest.SourceDir]int
	byRepository map[manifest.RepositoryRef]int
	byVisibility map[manifest.Visibility][]int
	byBranch     map[manifest.ModuleName]map[manifest.BranchName]manifest.BranchMapping
}
