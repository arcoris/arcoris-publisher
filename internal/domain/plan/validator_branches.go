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

// validateModuleBranches checks that a module has valid branch plans.
func (v *planValidator) validateModuleBranches(module manifest.ModuleName, modulePlan ModulePlan) {
	if len(modulePlan.branches) == 0 {
		v.addIssue(Issue{Code: IssueInvalidBranch, Module: module, Message: fmt.Sprintf("module %q has no branch plans", module)})
		return
	}
	for index, branch := range modulePlan.branches {
		if branch.Source() == "" {
			v.addIssue(Issue{Code: IssueInvalidBranch, Module: module, Message: fmt.Sprintf("module %q branch[%d] has empty source", module, index)})
		}
		if branch.Target() == "" {
			v.addIssue(Issue{Code: IssueInvalidBranch, Module: module, Message: fmt.Sprintf("module %q branch[%d] has empty target", module, index)})
		}
	}
}
