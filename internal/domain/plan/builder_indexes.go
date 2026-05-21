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

import "arcoris.dev/arcoris-publisher/internal/domain/manifest"

// indexPlan rebuilds lookup maps for module name, module path, and repository.
func indexPlan(plan *Plan) {
	plan.byModule = make(map[manifest.ModuleName]int, len(plan.modules))
	plan.byPath = make(map[manifest.ModulePath]int, len(plan.modules))
	plan.byRepo = make(map[manifest.RepositoryRef]int, len(plan.modules))

	for index, modulePlan := range plan.modules {
		indexModulePlan(plan, index, modulePlan)
	}
}

// indexModulePlan adds one module plan to all lookup indexes.
func indexModulePlan(plan *Plan, index int, modulePlan ModulePlan) {
	if _, exists := plan.byModule[modulePlan.Name()]; !exists {
		plan.byModule[modulePlan.Name()] = index
	}
	if _, exists := plan.byPath[modulePlan.ModulePath()]; !exists {
		plan.byPath[modulePlan.ModulePath()] = index
	}
	if _, exists := plan.byRepo[modulePlan.Repository()]; !exists {
		plan.byRepo[modulePlan.Repository()] = index
	}
}
