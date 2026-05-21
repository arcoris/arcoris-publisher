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

import (
	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/manifest/resolved"
)

// Registry is a read-optimized view over an effective publication set.
type Registry struct {
	modules []resolved.PublicationModule

	byName       map[manifest.ModuleName]int
	byPath       map[manifest.ModulePath]int
	byRepository map[manifest.RepositoryRef]int
	bySourceDir  map[manifest.SourceDir]int
}

// New validates set and builds deterministic lookup indexes.
func New(set resolved.PublicationSet) (Registry, error) {
	modules := set.Modules()
	if err := validateModules(modules); err != nil {
		return Registry{}, err
	}

	registry := Registry{
		modules: modules,
	}
	registry.buildIndexes()

	return registry, nil
}

// Must builds a Registry and panics when validation fails.
func Must(set resolved.PublicationSet) Registry {
	registry, err := New(set)
	if err != nil {
		panic(err)
	}

	return registry
}

// Modules returns all modules in declaration order.
func (r Registry) Modules() []resolved.PublicationModule {
	return cloneModules(r.modules)
}

func cloneModules(
	modules []resolved.PublicationModule,
) []resolved.PublicationModule {
	out := make([]resolved.PublicationModule, len(modules))
	copy(out, modules)
	return out
}
