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

// branchPlans converts manifest branch mappings into publication branch plans.
func branchPlans(branches []manifest.BranchMapping) ([]BranchPlan, error) {
	plans := make([]BranchPlan, 0, len(branches))
	for _, branch := range branches {
		branchPlan, err := NewBranchPlan(branch)
		if err != nil {
			return nil, err
		}
		plans = append(plans, branchPlan)
	}
	return plans, nil
}
