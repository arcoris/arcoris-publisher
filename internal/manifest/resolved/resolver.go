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

// resolver owns one resolution pass and its trace.
type resolver struct {
	input ResolveInput
	trace ResolutionTrace
}

// resolve binds staging modules to module manifests and validates the result.
func (r *resolver) resolve() (ResolveResult, error) {
	moduleByName, err := r.indexModuleManifests()
	if err != nil {
		return ResolveResult{}, err
	}
	stagingModules := r.input.Staging.Modules()
	resolvedModules := make([]PublicationModule, 0, len(stagingModules))
	var collector manifest.IssueCollector
	used := make(map[manifest.ModuleName]struct{}, len(stagingModules))
	for i, stagingModule := range stagingModules {
		path := fmt.Sprintf("modules[%d]", i)
		moduleManifest, ok := moduleByName[stagingModule.Name()]
		if !ok {
			collector.Add(manifest.IssueMissingModuleManifest, path, "missing module manifest for %q", stagingModule.Name())
			continue
		}
		used[stagingModule.Name()] = struct{}{}
		publicationModule := r.resolveModule(path, stagingModule, moduleManifest)
		resolvedModules = append(resolvedModules, publicationModule)
	}
	for name := range moduleByName {
		if _, ok := used[name]; !ok {
			collector.Add(manifest.IssueUnknownModule, "moduleManifests", "module manifest %q is not referenced by staging manifest", name)
		}
	}
	if err := collector.Err(); err != nil {
		return ResolveResult{}, err
	}
	set := PublicationSet{metadata: r.input.Staging.Metadata(), source: r.input.Staging.Source(), publish: r.input.Staging.Publish(), modules: resolvedModules}
	if err := validatePublicationSet(set); err != nil {
		return ResolveResult{}, err
	}
	return ResolveResult{Set: set, Trace: r.trace}, nil
}
