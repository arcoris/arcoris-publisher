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

package plan

import "arcoris.dev/arcoris-publisher/internal/domain/manifest"

// planIndexes stores expected lookup maps derived from module plans.
type planIndexes struct {
	byModule map[manifest.ModuleName]int
	byPath   map[manifest.ModulePath]int
	byRepo   map[manifest.RepositoryRef]int
}

// newPlanIndexes allocates lookup maps sized for the module plan count.
func newPlanIndexes(size int) planIndexes {
	return planIndexes{
		byModule: make(map[manifest.ModuleName]int, size),
		byPath:   make(map[manifest.ModulePath]int, size),
		byRepo:   make(map[manifest.RepositoryRef]int, size),
	}
}
