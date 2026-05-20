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

package graph

import (
	"sort"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

// sortedModuleNames returns a detached lexical copy of module names.
//
// Graph construction uses this only for adjacency lists; declaration order is
// still preserved separately in Graph.order for stable topological tie-breaks.
func sortedModuleNames(values []manifest.ModuleName) []manifest.ModuleName {
	copy := append([]manifest.ModuleName(nil), values...)
	sort.Slice(copy, func(i, j int) bool { return copy[i] < copy[j] })
	return copy
}
