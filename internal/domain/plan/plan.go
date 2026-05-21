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

import (
	"arcoris.dev/arcoris-publisher/internal/domain/graph"
	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
	"arcoris.dev/arcoris-publisher/internal/domain/registry"
	"arcoris.dev/arcoris-publisher/internal/domain/versioning"
)

// Plan is an immutable-by-convention publish transaction domain snapshot.
//
// A plan contains only static publication decisions: which public modules should
// be published, in which dependency order, with which versions, repositories,
// branch mappings, and direct dependency requirements. It intentionally does not
// contain source commits, target repository HEADs, filesystem paths, process
// commands, or workflow execution state.
type Plan struct {
	manifest    manifest.Manifest
	registry    registry.Registry
	graph       graph.Graph
	assignments versioning.Assignments

	modules []ModulePlan
	skipped []SkippedModule

	byModule map[manifest.ModuleName]int
	byPath   map[manifest.ModulePath]int
	byRepo   map[manifest.RepositoryRef]int
}
