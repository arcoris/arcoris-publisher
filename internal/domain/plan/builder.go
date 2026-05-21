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

// builder coordinates the pure domain projections needed to create a plan.
type builder struct {
	manifestValue manifest.Manifest
	registryValue registry.Registry
	graphValue    graph.Graph
	assignments   versioning.Assignments
}

// newBuilder captures constructor inputs for step-by-step plan assembly.
func newBuilder(
	manifestValue manifest.Manifest,
	registryValue registry.Registry,
	graphValue graph.Graph,
	assignments versioning.Assignments,
) builder {
	return builder{
		manifestValue: manifestValue,
		registryValue: registryValue,
		graphValue:    graphValue,
		assignments:   assignments,
	}
}

// build validates inputs, creates module plans, indexes them, and validates the result.
func (b builder) build() (Plan, error) {
	if err := b.validateInputs(); err != nil {
		return Plan{}, err
	}
	order, err := b.publishOrder()
	if err != nil {
		return Plan{}, err
	}
	modules, err := b.modulePlans(order)
	if err != nil {
		return Plan{}, err
	}
	plan := b.newPlan(modules)
	indexPlan(&plan)
	if err := plan.Validate(); err != nil {
		return Plan{}, err
	}
	return plan, nil
}
