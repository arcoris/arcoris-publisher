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

// validateModuleDependencies checks dependency existence and visibility rules.
func validateModuleDependencies(collector *manifest.IssueCollector, modules []PublicationModule) {
	byName := indexPublicationModules(modules)

	for i, mod := range modules {
		modulePath := fmt.Sprintf("modules[%d]", i)

		for j, depName := range mod.Dependencies() {
			dep, found := byName[depName]
			depPath := fmt.Sprintf("%s.dependencies[%d]", modulePath, j)

			validateModuleDependency(collector, depPath, mod, depName, dep, found)
		}
	}
}

// indexPublicationModules builds the dependency lookup keyed by module name.
func indexPublicationModules(
	modules []PublicationModule,
) map[manifest.ModuleName]PublicationModule {
	byName := make(map[manifest.ModuleName]PublicationModule, len(modules))
	for _, mod := range modules {
		byName[mod.Name()] = mod
	}
	return byName
}

// validateModuleDependency checks one declared dependency edge.
func validateModuleDependency(
	collector *manifest.IssueCollector,
	depPath string,
	mod PublicationModule,
	depName manifest.ModuleName,
	dep PublicationModule,
	found bool,
) {
	if !found {
		collector.Add(manifest.IssueUnknownModule, depPath, "unknown dependency %q", depName)
		return
	}

	if dep.Name() == mod.Name() {
		collector.Add(manifest.IssueInvalidDependency, depPath, "module cannot depend on itself")
	}

	if mod.Visibility() == manifest.VisibilityPublic &&
		dep.Visibility() != manifest.VisibilityPublic {
		collector.Add(
			manifest.IssueInvalidDependency,
			depPath,
			"public module %q depends on non-public module %q",
			mod.Name(),
			dep.Name(),
		)
	}
}
