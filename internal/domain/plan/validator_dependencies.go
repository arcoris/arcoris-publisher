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
)

// validateModuleDependencies checks dependency plans and requirement map parity.
func (v *planValidator) validateModuleDependencies(module manifest.ModuleName, modulePlan ModulePlan) {
	expectedRequirements := make(map[manifest.ModulePath]struct{}, len(modulePlan.dependencies))
	for index, dependency := range modulePlan.dependencies {
		v.validateDependencyPlan(module, index, dependency)
		expectedRequirements[dependency.ModulePath()] = struct{}{}
		if got, ok := modulePlan.requirements[dependency.ModulePath()]; !ok {
			v.addIssue(Issue{Code: IssueInvalidDependency, Module: module, Dependency: dependency.Module(), Path: dependency.ModulePath(), Message: fmt.Sprintf("module %q dependency %q is missing from requirements", module, dependency.Module())})
		} else if got != dependency.Version() {
			v.addIssue(Issue{Code: IssueInvalidDependency, Module: module, Dependency: dependency.Module(), Path: dependency.ModulePath(), Message: fmt.Sprintf("module %q dependency %q requirement version %q does not match dependency version %q", module, dependency.Module(), got, dependency.Version())})
		}
	}
	for modulePath := range modulePlan.requirements {
		if _, ok := expectedRequirements[modulePath]; !ok {
			v.addIssue(Issue{Code: IssueInvalidDependency, Module: module, Path: modulePath, Message: fmt.Sprintf("module %q has unexpected requirement for %q", module, modulePath)})
		}
	}
}

// validateDependencyPlan checks one direct dependency plan.
func (v *planValidator) validateDependencyPlan(module manifest.ModuleName, index int, dependency DependencyPlan) {
	if dependency.Module() == "" {
		v.addIssue(Issue{Code: IssueInvalidDependency, Module: module, Message: fmt.Sprintf("module %q dependency[%d] has empty module", module, index)})
	}
	if dependency.ModulePath() == "" {
		v.addIssue(Issue{Code: IssueInvalidDependency, Module: module, Dependency: dependency.Module(), Message: fmt.Sprintf("module %q dependency[%d] has empty module path", module, index)})
	}
	if dependency.Version() == "" {
		v.addIssue(Issue{Code: IssueInvalidDependency, Module: module, Dependency: dependency.Module(), Path: dependency.ModulePath(), Message: fmt.Sprintf("module %q dependency[%d] has empty version", module, index)})
	}
}
