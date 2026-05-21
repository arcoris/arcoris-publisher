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

package resolved

import (
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
)

// resolveBranches applies module overrides, then defaults, then source default branch.
func (r *resolver) resolveBranches(path string, sm staging.Module) []manifest.BranchMapping {
	tracePath := path + ".branches"

	if override, ok := sm.BranchesOverride(); ok {
		r.trace.AddStagingModule(
			tracePath,
			branchCount(override),
		)
		return override
	}

	if r.input.Staging.Defaults().BranchesSet() {
		branches := r.input.Staging.Defaults().Branches()
		r.trace.AddStagingDefault(
			tracePath,
			branchCount(branches),
			"defaults.branches",
		)
		return branches
	}

	return r.defaultBranchMapping(path)
}

// defaultBranchMapping maps the staging default branch to itself.
func (r *resolver) defaultBranchMapping(path string) []manifest.BranchMapping {
	branch := r.input.Staging.Source().DefaultBranch()
	mapping, _ := manifest.NewBranchMapping(
		manifest.BranchMappingSpec{
			Source: branch.String(),
			Target: branch.String(),
		},
	)
	r.trace.AddBuiltInDefault(
		path+".branches",
		fmt.Sprintf("%s->%s", branch, branch),
		"source.defaultBranch",
	)
	return []manifest.BranchMapping{mapping}
}

// branchCount formats branch mapping counts for trace output.
func branchCount(branches []manifest.BranchMapping) string {
	return fmt.Sprintf("%d mapping(s)", len(branches))
}
