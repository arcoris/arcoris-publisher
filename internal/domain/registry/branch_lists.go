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

package registry

import "arcoris.dev/arcoris-publisher/internal/domain/manifest"

// SourceBranches returns source branches declared for module in module declaration order.
func (r Registry) SourceBranches(module manifest.ModuleName) []manifest.BranchName {
	value, ok := r.ModuleByName(module)
	if !ok {
		return nil
	}
	return branchSources(value.Branches())
}

// TargetBranches returns target branches declared for module in module declaration order.
func (r Registry) TargetBranches(module manifest.ModuleName) []manifest.BranchName {
	value, ok := r.ModuleByName(module)
	if !ok {
		return nil
	}
	return branchTargets(value.Branches())
}

// branchSources extracts source branch names from branch mappings.
func branchSources(branches []manifest.BranchMapping) []manifest.BranchName {
	result := make([]manifest.BranchName, 0, len(branches))
	for _, branch := range branches {
		result = append(result, branch.Source())
	}
	return result
}

// branchTargets extracts target branch names from branch mappings.
func branchTargets(branches []manifest.BranchMapping) []manifest.BranchName {
	result := make([]manifest.BranchName, 0, len(branches))
	for _, branch := range branches {
		result = append(result, branch.Target())
	}
	return result
}
