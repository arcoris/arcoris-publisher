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
	modulemanifest "arcoris.dev/arcoris-publisher/internal/manifest/module"
)

// resolveModules binds every staging module to its matching module manifest.
func (r *resolver) resolveModules(
	moduleByName map[manifest.ModuleName]modulemanifest.Manifest,
) ([]PublicationModule, error) {
	var collector manifest.IssueCollector

	stagingModules := r.input.Staging.Modules()
	resolvedModules := make([]PublicationModule, 0, len(stagingModules))
	used := make(map[manifest.ModuleName]struct{}, len(stagingModules))

	for i, stagingModule := range stagingModules {
		path := fmt.Sprintf("modules[%d]", i)
		moduleManifest, found := moduleByName[stagingModule.Name()]

		if !found {
			collector.Add(
				manifest.IssueMissingModuleManifest,
				path,
				"missing module manifest for %q",
				stagingModule.Name(),
			)
			continue
		}

		used[stagingModule.Name()] = struct{}{}
		resolvedModules = append(resolvedModules, r.resolveModule(path, stagingModule, moduleManifest))
	}

	collectUnusedModuleManifests(&collector, moduleByName, used)

	if err := collector.Err(); err != nil {
		return nil, err
	}

	return resolvedModules, nil
}

// collectUnusedModuleManifests reports module manifests without staging entries.
func collectUnusedModuleManifests(
	collector *manifest.IssueCollector,
	moduleByName map[manifest.ModuleName]modulemanifest.Manifest,
	used map[manifest.ModuleName]struct{},
) {
	for name := range moduleByName {
		if _, ok := used[name]; !ok {
			collector.Add(
				manifest.IssueUnknownModule,
				"moduleManifests",
				"module manifest %q is not referenced by staging manifest",
				name,
			)
		}
	}
}
