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
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
	"arcoris.dev/arcoris-publisher/internal/domain/registry"
	"arcoris.dev/arcoris-publisher/internal/domain/versioning"
)

// dependencyPlans resolves direct module dependencies to publish-time versions.
func dependencyPlans(registryValue registry.Registry, assignments versioning.Assignments, module manifest.Module) ([]DependencyPlan, error) {
	versions, err := assignments.DependencyVersions(registryValue, module)
	if err != nil {
		return nil, &ValidationError{Issues: []Issue{{
			Code:    IssueInvalidDependency,
			Module:  module.Name(),
			Message: fmt.Sprintf("module %q dependency versions are invalid: %v", module.Name(), err),
		}}}
	}
	plans := make([]DependencyPlan, 0, len(versions))
	for _, dependency := range versions {
		plans = append(plans, DependencyPlan{module: dependency.Module(), modulePath: dependency.ModulePath(), version: dependency.Version()})
	}
	return plans, nil
}
