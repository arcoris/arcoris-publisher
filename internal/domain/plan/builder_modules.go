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

// modulePlans builds all module plans in dependency-first publish order.
func (b builder) modulePlans(order []manifest.ModuleName) ([]ModulePlan, error) {
	plans := make([]ModulePlan, 0, len(order))
	var issues []Issue
	for index, name := range order {
		modulePlan, moduleIssues := b.modulePlan(index, name)
		if len(moduleIssues) > 0 {
			issues = append(issues, moduleIssues...)
			continue
		}
		plans = append(plans, modulePlan)
	}
	if len(issues) > 0 {
		return nil, &ValidationError{Issues: issues}
	}
	return plans, nil
}
