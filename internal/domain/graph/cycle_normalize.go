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

import "arcoris.dev/arcoris-publisher/internal/domain/manifest"

// normalizeCycle rotates a cycle so equivalent cycles get the same string form.
//
// DFS can discover the same directed cycle from different starting nodes. By
// rotating to the lexically smallest node, the cycle can be deduplicated by its
// rendered path without losing the closed first==last shape.
func normalizeCycle(nodes []manifest.ModuleName) Cycle {
	if len(nodes) <= 2 {
		return NewCycle(nodes)
	}
	body := append([]manifest.ModuleName(nil), nodes[:len(nodes)-1]...)
	minIndex := minModuleNameIndex(body)
	rotated := make([]manifest.ModuleName, 0, len(nodes))
	rotated = append(rotated, body[minIndex:]...)
	rotated = append(rotated, body[:minIndex]...)
	rotated = append(rotated, rotated[0])
	return NewCycle(rotated)
}

// minModuleNameIndex returns the first lexical minimum in names.
func minModuleNameIndex(names []manifest.ModuleName) int {
	minIndex := 0
	for i := 1; i < len(names); i++ {
		if names[i] < names[minIndex] {
			minIndex = i
		}
	}
	return minIndex
}
