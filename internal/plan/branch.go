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

import "arcoris.dev/arcoris-publisher/internal/manifest"

// BranchPlan describes one source-to-target branch publication mapping for a
// module.
type BranchPlan struct {
	// source is the branch read from the source repository.
	source manifest.BranchName

	// target is the branch written in the target repository.
	target manifest.BranchName
}

// newBranchPlan converts a resolved manifest mapping into a plan value.
func newBranchPlan(mapping manifest.BranchMapping) BranchPlan {
	return BranchPlan{source: mapping.Source(), target: mapping.Target()}
}

// Source returns the source repository branch to read from.
func (b BranchPlan) Source() manifest.BranchName { return b.source }

// Target returns the target repository branch to publish to.
func (b BranchPlan) Target() manifest.BranchName { return b.target }

// cloneBranchPlans detaches branch slices before returning them to callers.
func cloneBranchPlans(in []BranchPlan) []BranchPlan {
	out := make([]BranchPlan, len(in))
	copy(out, in)
	return out
}
