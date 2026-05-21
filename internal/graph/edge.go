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

import "arcoris.dev/arcoris-publisher/internal/manifest"

// Edge is a directed dependency-to-dependent graph edge.
//
// From is the dependency. To is the dependent. For example, if control depends
// on foundation, the graph contains foundation -> control.
type Edge struct {
	from manifest.ModuleName
	to   manifest.ModuleName
}

// NewEdge returns a dependency-to-dependent edge.
func NewEdge(from manifest.ModuleName, to manifest.ModuleName) Edge {
	return Edge{from: from, to: to}
}

// From returns the dependency module name.
func (e Edge) From() manifest.ModuleName { return e.from }

// To returns the dependent module name.
func (e Edge) To() manifest.ModuleName { return e.to }
