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

import "fmt"

// validateModules checks every module plan and builds expected lookup indexes.
func (v *planValidator) validateModules() {
	for index, modulePlan := range v.plan.modules {
		v.validateModule(index, modulePlan)
		v.addExpectedModule(modulePlan, index)
	}
}

// validateModule checks one module plan's local invariants.
func (v *planValidator) validateModule(index int, modulePlan ModulePlan) {
	name := modulePlan.Name()
	if name == "" {
		v.addIssue(Issue{Code: IssueMissingModule, Message: fmt.Sprintf("modules[%d] has empty module name", index)})
	}
	if modulePlan.Version() == "" {
		v.addIssue(Issue{Code: IssueMissingVersion, Module: name, Message: fmt.Sprintf("module %q has no version", name)})
	}
	if modulePlan.Action() != ActionPublish {
		v.addIssue(Issue{Code: IssueInvalidAction, Module: name, Message: fmt.Sprintf("module %q has unsupported action %q", name, modulePlan.Action())})
	}
	v.validateModuleBranches(name, modulePlan)
	v.validateModuleDependencies(name, modulePlan)
}
