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
	"strings"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

// Cycle describes one directed dependency cycle.
//
// Nodes returns a path where the first and last names are equal.
type Cycle struct {
	nodes []manifest.ModuleName
}

// NewCycle creates a cycle with detached node storage.
func NewCycle(nodes []manifest.ModuleName) Cycle {
	return Cycle{nodes: append([]manifest.ModuleName(nil), nodes...)}
}

// Nodes returns a detached cycle path.
func (c Cycle) Nodes() []manifest.ModuleName {
	return append([]manifest.ModuleName(nil), c.nodes...)
}

// String returns a human-readable cycle path.
func (c Cycle) String() string {
	parts := make([]string, 0, len(c.nodes))
	for _, node := range c.nodes {
		parts = append(parts, string(node))
	}
	return strings.Join(parts, " -> ")
}
