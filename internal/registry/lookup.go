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

// ModuleByName returns the module with name.
func (r Registry) ModuleByName(
	name manifest.ModuleName,
) (resolved.PublicationModule, bool) {
	index, ok := r.byName[name]
	if !ok {
		return resolved.PublicationModule{}, false
	}

	return r.lookup(index)
}

// ModuleByPath returns the module with module path.
func (r Registry) ModuleByPath(
	path manifest.ModulePath,
) (resolved.PublicationModule, bool) {
	index, ok := r.byPath[path]
	if !ok {
		return resolved.PublicationModule{}, false
	}

	return r.lookup(index)
}

// ModuleByRepository returns the module targeting repo.
func (r Registry) ModuleByRepository(
	repo manifest.RepositoryRef,
) (resolved.PublicationModule, bool) {
	index, ok := r.byRepository[repo]
	if !ok {
		return resolved.PublicationModule{}, false
	}

	return r.lookup(index)
}

// ModuleBySourceDir returns the module staged at dir.
func (r Registry) ModuleBySourceDir(
	dir manifest.SourceDir,
) (resolved.PublicationModule, bool) {
	index, ok := r.bySourceDir[dir]
	if !ok {
		return resolved.PublicationModule{}, false
	}

	return r.lookup(index)
}

// ContainsName reports whether name exists in the registry.
func (r Registry) ContainsName(name manifest.ModuleName) bool {
	_, ok := r.byName[name]
	return ok
}

// ContainsPath reports whether path exists in the registry.
func (r Registry) ContainsPath(path manifest.ModulePath) bool {
	_, ok := r.byPath[path]
	return ok
}

// ContainsRepository reports whether repo exists in the registry.
func (r Registry) ContainsRepository(repo manifest.RepositoryRef) bool {
	_, ok := r.byRepository[repo]
	return ok
}

// ContainsSourceDir reports whether dir exists in the registry.
func (r Registry) ContainsSourceDir(dir manifest.SourceDir) bool {
	_, ok := r.bySourceDir[dir]
	return ok
}

func (r Registry) lookup(index int) (resolved.PublicationModule, bool) {
	if index < 0 || index >= len(r.modules) {
		return resolved.PublicationModule{}, false
	}

	return r.modules[index], true
}
