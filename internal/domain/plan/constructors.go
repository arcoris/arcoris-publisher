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

// New builds a publish plan from validated domain inputs.
//
// The constructor treats manifest, registry, graph, and assignments as separate
// domain projections and validates that they agree before any workflow layer can
// act on them.
func New(manifestValue manifest.Manifest, registryValue registry.Registry, graphValue graph.Graph, assignments versioning.Assignments) (Plan, error) {
	return newBuilder(manifestValue, registryValue, graphValue, assignments).build()
}

// Must constructs a plan and panics on validation failure.
//
// Must is intended for tests and static wiring. Runtime code should call New and
// return diagnostics to the caller.
func Must(manifestValue manifest.Manifest, registryValue registry.Registry, graphValue graph.Graph, assignments versioning.Assignments) Plan {
	plan, err := New(manifestValue, registryValue, graphValue, assignments)
	if err != nil {
		panic(err)
	}
	return plan
}
