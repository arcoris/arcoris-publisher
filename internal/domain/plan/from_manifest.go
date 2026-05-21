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

// FromManifest builds registry and dependency graph values from manifest before
// constructing a publish plan.
func FromManifest(manifestValue manifest.Manifest, assignments versioning.Assignments) (Plan, error) {
	registryValue, err := registry.FromManifest(manifestValue)
	if err != nil {
		return Plan{}, validationErrorf(IssueInvalidRegistry, "", "invalid registry: %s", err)
	}
	graphValue, err := graph.FromManifest(manifestValue)
	if err != nil {
		return Plan{}, validationErrorf(IssueInvalidGraph, "", "invalid graph: %s", err)
	}
	return New(manifestValue, registryValue, graphValue, assignments)
}
