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

// BranchMapping returns the branch mapping for module and source branch.
//
// The registry indexes mappings by source branch because publication workflows
// start from a checked-out source branch and need to resolve the corresponding
// target branch before pushing.
func (r Registry) BranchMapping(module manifest.ModuleName, source manifest.BranchName) (manifest.BranchMapping, bool) {
	branches, ok := r.byBranch[module]
	if !ok {
		return manifest.BranchMapping{}, false
	}
	branch, ok := branches[source]
	return branch, ok
}
