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

package app

import (
	"context"
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/graph"
	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/registry"
	"arcoris.dev/arcoris-publisher/internal/versioning"
)

// BuildPlan loads manifests and builds the immutable executable plan.
func (a App) BuildPlan(
	ctx context.Context,
	manifestPath string,
	version versioning.Version,
) (plan.Plan, error) {
	set, err := a.loader.LoadPublicationSet(ctx, manifestPath)
	if err != nil {
		return plan.Plan{}, fmt.Errorf("config: %w", err)
	}

	reg, err := registry.New(set)
	if err != nil {
		return plan.Plan{}, fmt.Errorf("registry: %w", err)
	}

	g, err := graph.New(reg)
	if err != nil {
		return plan.Plan{}, fmt.Errorf("graph: %w", err)
	}

	assignments, err := versioning.Assign(versioning.Request{
		Set:      set,
		Registry: reg,
		Graph:    g,
		Version:  version,
	})
	if err != nil {
		return plan.Plan{}, fmt.Errorf("versioning: %w", err)
	}

	out, err := plan.Build(plan.Request{
		Set:         set,
		Registry:    reg,
		Graph:       g,
		Assignments: assignments,
	})
	if err != nil {
		return plan.Plan{}, fmt.Errorf("plan: %w", err)
	}

	return out, nil
}
