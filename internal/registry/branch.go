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

import "arcoris.dev/arcoris-publisher/internal/manifest"

// BranchMapping returns the effective mapping for module on source branch.
func (r Registry) BranchMapping(
	module manifest.ModuleName,
	source manifest.BranchName,
) (manifest.BranchMapping, bool) {
	publicationModule, ok := r.ModuleByName(module)
	if !ok {
		return manifest.BranchMapping{}, false
	}

	for _, mapping := range publicationModule.Branches() {
		if mapping.Source() == source {
			return mapping, true
		}
	}

	return manifest.BranchMapping{}, false
}

// SourceBranches returns deduplicated source branches in declaration order.
func (r Registry) SourceBranches() []manifest.BranchName {
	return r.branchNames(func(mapping manifest.BranchMapping) manifest.BranchName {
		return mapping.Source()
	})
}

// TargetBranches returns deduplicated target branches in declaration order.
func (r Registry) TargetBranches() []manifest.BranchName {
	return r.branchNames(func(mapping manifest.BranchMapping) manifest.BranchName {
		return mapping.Target()
	})
}

func (r Registry) branchNames(
	selectName func(manifest.BranchMapping) manifest.BranchName,
) []manifest.BranchName {
	seen := map[manifest.BranchName]struct{}{}
	out := make([]manifest.BranchName, 0)

	for _, module := range r.modules {
		for _, mapping := range module.Branches() {
			name := selectName(mapping)
			if _, ok := seen[name]; ok {
				continue
			}

			seen[name] = struct{}{}
			out = append(out, name)
		}
	}

	return out
}
