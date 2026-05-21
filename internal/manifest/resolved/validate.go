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

func validatePublicationSet(set PublicationSet) error {
	var collector manifest.IssueCollector
	modulePaths := make(map[manifest.ModulePath]int, len(set.modules))
	repositories := make(map[manifest.RepositoryRef]int, len(set.modules))
	byName := make(map[manifest.ModuleName]PublicationModule, len(set.modules))
	for i, mod := range set.modules {
		path := fmt.Sprintf("modules[%d]", i)
		byName[mod.Name()] = mod
		if prev, ok := modulePaths[mod.ModulePath()]; ok {
			collector.Add(manifest.IssueDuplicateValue, path+".module.path", "duplicate module path %q previously declared at modules[%d]", mod.ModulePath(), prev)
		} else {
			modulePaths[mod.ModulePath()] = i
		}
		if mod.Visibility() == manifest.VisibilityPublic {
			if prev, ok := repositories[mod.Repository()]; ok {
				collector.Add(manifest.IssueDuplicateValue, path+".repository", "duplicate public target repository %q previously declared at modules[%d]", mod.Repository(), prev)
			} else {
				repositories[mod.Repository()] = i
			}
		}
	}
	for i, mod := range set.modules {
		modulePath := fmt.Sprintf("modules[%d]", i)
		for j, depName := range mod.Dependencies() {
			dep, ok := byName[depName]
			depPath := fmt.Sprintf("%s.dependencies[%d]", modulePath, j)
			if !ok {
				collector.Add(manifest.IssueUnknownModule, depPath, "unknown dependency %q", depName)
				continue
			}
			if dep.Name() == mod.Name() {
				collector.Add(manifest.IssueInvalidDependency, depPath, "module cannot depend on itself")
			}
			if mod.Visibility() == manifest.VisibilityPublic && dep.Visibility() != manifest.VisibilityPublic {
				collector.Add(manifest.IssueInvalidDependency, depPath, "public module %q depends on non-public module %q", mod.Name(), dep.Name())
			}
		}
	}
	return collector.Err()
}
