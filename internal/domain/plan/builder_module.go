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
	"errors"
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

// modulePlan builds one publishable module plan from registry and assignment data.
func (b builder) modulePlan(index int, name manifest.ModuleName) (ModulePlan, []Issue) {
	module, ok := b.registryValue.ModuleByName(name)
	if !ok {
		return ModulePlan{}, []Issue{missingModuleIssue(name)}
	}
	version, ok := b.assignments.VersionOfModule(name)
	if !ok {
		return ModulePlan{}, []Issue{missingVersionIssue(name)}
	}
	branches, err := branchPlans(module.Branches())
	if err != nil {
		return ModulePlan{}, []Issue{invalidBranchIssue(name, err)}
	}
	dependencies, err := dependencyPlans(b.registryValue, b.assignments, module)
	if err != nil {
		return ModulePlan{}, dependencyIssues(name, err)
	}
	return ModulePlan{
		module:       module,
		version:      version,
		action:       ActionPublish,
		branches:     branches,
		dependencies: dependencies,
		requirements: requirementMap(dependencies),
		orderIndex:   index,
	}, nil
}

// dependencyIssues adapts dependency planning errors into module-level issues.
func dependencyIssues(name manifest.ModuleName, err error) []Issue {
	var validationErr *ValidationError
	if errors.As(err, &validationErr) {
		return validationErr.Issues
	}
	return []Issue{{
		Code:    IssueInvalidDependency,
		Module:  name,
		Message: fmt.Sprintf("module %q has invalid dependency plan: %v", name, err),
	}}
}
