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

// planValidator checks module plans, assignment coverage, and lookup indexes.
type planValidator struct {
	plan Plan

	expected planIndexes
	issues   []Issue
}

// newPlanValidator prepares deterministic validation for one plan.
func newPlanValidator(plan Plan) planValidator {
	return planValidator{
		plan:     plan,
		expected: newPlanIndexes(len(plan.modules)),
	}
}

// validate returns nil when every plan invariant holds.
func (v *planValidator) validate() error {
	v.validateModules()
	v.validateAssignments()
	v.validateIndexes()
	if len(v.issues) > 0 {
		return &ValidationError{Issues: v.issues}
	}
	return nil
}

// addIssue records one validation issue in deterministic discovery order.
func (v *planValidator) addIssue(issue Issue) {
	v.issues = append(v.issues, issue)
}

// addInvalidIndexIssue records a lookup index mismatch.
func (v *planValidator) addInvalidIndexIssue(index string, message string) {
	v.addIssue(Issue{Code: IssueInvalidIndex, Index: index, Message: message})
}

// addExpectedModule indexes a module plan and reports duplicate natural keys.
func (v *planValidator) addExpectedModule(modulePlan ModulePlan, index int) {
	if previous, exists := v.expected.byModule[modulePlan.Name()]; exists {
		v.addIssue(duplicateModuleIssue(modulePlan.Name(), previous))
	} else {
		v.expected.byModule[modulePlan.Name()] = index
	}
	if previous, exists := v.expected.byPath[modulePlan.ModulePath()]; exists {
		v.addIssue(duplicateModulePathIssue(modulePlan.Name(), modulePlan.ModulePath(), previous))
	} else {
		v.expected.byPath[modulePlan.ModulePath()] = index
	}
	if previous, exists := v.expected.byRepo[modulePlan.Repository()]; exists {
		v.addIssue(duplicateRepositoryIssue(modulePlan.Name(), modulePlan.Repository(), previous))
	} else {
		v.expected.byRepo[modulePlan.Repository()] = index
	}
}

// hasExpectedModule reports whether a module name belongs to the plan modules.
func (v planValidator) hasExpectedModule(module manifest.ModuleName) bool {
	_, ok := v.expected.byModule[module]
	return ok
}
