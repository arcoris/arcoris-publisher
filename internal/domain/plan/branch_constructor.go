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

// NewBranchPlan constructs a branch plan from a validated manifest mapping.
func NewBranchPlan(mapping manifest.BranchMapping) (BranchPlan, error) {
	return newBranchPlan(mapping.Source(), mapping.Target())
}

// newBranchPlan validates branch names already extracted from a mapping.
func newBranchPlan(source manifest.BranchName, target manifest.BranchName) (BranchPlan, error) {
	candidate := BranchPlan{source: source, target: target}
	if candidate.source == "" {
		return BranchPlan{}, validationErrorf(IssueInvalidBranch, "", "source branch is required")
	}
	if candidate.target == "" {
		return BranchPlan{}, validationErrorf(IssueInvalidBranch, "", "target branch is required")
	}
	return candidate, nil
}
