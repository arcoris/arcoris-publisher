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

import (
	"arcoris.dev/arcoris-publisher/internal/graph"
	"arcoris.dev/arcoris-publisher/internal/manifest/resolved"
	"arcoris.dev/arcoris-publisher/internal/registry"
)

// Request contains all deterministic inputs required to assign module versions.
//
// Set carries the effective publication policy, Registry provides module lookup,
// Graph provides dependency and publication order, and Version is the release or
// snapshot version supplied by the caller. Runtime concerns such as selected
// modules, dry-run flags, Git refs, and filesystem paths intentionally do not
// belong to versioning.
type Request struct {
	// Set is the resolved publication model that owns the effective version
	// policy and module declarations.
	Set resolved.PublicationSet

	// Registry is the lookup index built from Set. Assign validates that it
	// still contains every publishable module from Set before using it.
	Registry registry.Registry

	// Graph is the dependency graph built from Registry. It supplies the
	// deterministic publication order and direct dependency lists.
	Graph graph.Graph

	// Version is the release or snapshot version assigned to every publishable
	// module in the request.
	Version Version
}
