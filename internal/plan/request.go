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
	"arcoris.dev/arcoris-publisher/internal/graph"
	"arcoris.dev/arcoris-publisher/internal/manifest/resolved"
	"arcoris.dev/arcoris-publisher/internal/registry"
	"arcoris.dev/arcoris-publisher/internal/versioning"
)

// Request contains all deterministic inputs required to build a publication
// plan.
//
// The request intentionally carries already-resolved and already-indexed
// inputs. Planning does not decode manifests, rebuild defaults, assign versions,
// or inspect the filesystem. Those responsibilities belong to earlier or later
// layers.
type Request struct {
	// Set is the resolved publication model that supplies global context.
	Set resolved.PublicationSet

	// Registry is the lookup index built from Set.
	Registry registry.Registry

	// Graph supplies deterministic publication order and dependency topology.
	Graph graph.Graph

	// Assignments supplies module versions and dependency requirements.
	Assignments versioning.Assignments
}
