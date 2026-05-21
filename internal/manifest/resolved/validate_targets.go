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
)

// validateUniqueModuleTargets checks cross-module uniqueness constraints.
func validateUniqueModuleTargets(collector *manifest.IssueCollector, modules []PublicationModule) {
	modulePaths := make(map[manifest.ModulePath]int, len(modules))
	repositories := make(map[manifest.RepositoryRef]int, len(modules))

	for i, mod := range modules {
		path := fmt.Sprintf("modules[%d]", i)

		collectDuplicateModulePath(collector, modulePaths, i, path, mod)
		collectDuplicatePublicRepository(collector, repositories, i, path, mod)
	}
}

// collectDuplicateModulePath records duplicate module.path values.
func collectDuplicateModulePath(
	collector *manifest.IssueCollector,
	seen map[manifest.ModulePath]int,
	index int,
	path string,
	mod PublicationModule,
) {
	if prev, ok := seen[mod.ModulePath()]; ok {
		collector.Add(
			manifest.IssueDuplicateValue,
			path+".module.path",
			"duplicate module path %q previously declared at modules[%d]",
			mod.ModulePath(),
			prev,
		)
		return
	}

	seen[mod.ModulePath()] = index
}

// collectDuplicatePublicRepository records duplicate public target repositories.
func collectDuplicatePublicRepository(
	collector *manifest.IssueCollector,
	seen map[manifest.RepositoryRef]int,
	index int,
	path string,
	mod PublicationModule,
) {
	if mod.Visibility() != manifest.VisibilityPublic {
		return
	}

	if prev, ok := seen[mod.Repository()]; ok {
		collector.Add(
			manifest.IssueDuplicateValue,
			path+".repository",
			"duplicate public target repository %q previously declared at modules[%d]",
			mod.Repository(),
			prev,
		)
		return
	}

	seen[mod.Repository()] = index
}
