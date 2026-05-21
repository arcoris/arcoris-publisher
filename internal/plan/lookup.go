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

package plan

import "arcoris.dev/arcoris-publisher/internal/manifest"

// ModuleByName returns the planned module with name.
func (p Plan) ModuleByName(name manifest.ModuleName) (ModulePlan, bool) {
	idx, ok := p.byName[name]
	if !ok {
		return ModulePlan{}, false
	}
	return p.modules[idx], true
}

// ModuleByPath returns the planned module with module path.
func (p Plan) ModuleByPath(path manifest.ModulePath) (ModulePlan, bool) {
	idx, ok := p.byPath[path]
	if !ok {
		return ModulePlan{}, false
	}
	return p.modules[idx], true
}

// ModuleByRepository returns the planned module targeting repository.
func (p Plan) ModuleByRepository(repository manifest.RepositoryRef) (ModulePlan, bool) {
	idx, ok := p.byRepository[repository]
	if !ok {
		return ModulePlan{}, false
	}
	return p.modules[idx], true
}

// ContainsName reports whether name is planned.
func (p Plan) ContainsName(name manifest.ModuleName) bool {
	_, ok := p.byName[name]
	return ok
}

// ContainsPath reports whether module path is planned.
func (p Plan) ContainsPath(path manifest.ModulePath) bool {
	_, ok := p.byPath[path]
	return ok
}

// ContainsRepository reports whether repository is targeted by a planned module.
func (p Plan) ContainsRepository(repository manifest.RepositoryRef) bool {
	_, ok := p.byRepository[repository]
	return ok
}
