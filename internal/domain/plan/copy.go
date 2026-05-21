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
	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
	"arcoris.dev/arcoris-publisher/internal/domain/versioning"
)

// cloneModulePlans returns a detached module-plan slice.
func cloneModulePlans(modules []ModulePlan) []ModulePlan {
	return append([]ModulePlan(nil), modules...)
}

// cloneSkippedModules returns a detached skipped-module slice.
func cloneSkippedModules(modules []SkippedModule) []SkippedModule {
	return append([]SkippedModule(nil), modules...)
}

// cloneBranchPlans returns a detached branch-plan slice.
func cloneBranchPlans(branches []BranchPlan) []BranchPlan {
	return append([]BranchPlan(nil), branches...)
}

// cloneDependencyPlans returns a detached dependency-plan slice.
func cloneDependencyPlans(dependencies []DependencyPlan) []DependencyPlan {
	return append([]DependencyPlan(nil), dependencies...)
}

// cloneRequirements returns a detached dependency requirement map.
func cloneRequirements(requirements map[manifest.ModulePath]versioning.Version) map[manifest.ModulePath]versioning.Version {
	copy := make(map[manifest.ModulePath]versioning.Version, len(requirements))
	for modulePath, version := range requirements {
		copy[modulePath] = version
	}
	return copy
}
