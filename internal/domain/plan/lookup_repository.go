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

import "arcoris.dev/arcoris-publisher/internal/domain/manifest"

// ModulePlanByRepository returns the plan for repository and whether it was found.
func (p Plan) ModulePlanByRepository(repository manifest.RepositoryRef) (ModulePlan, bool) {
	index, ok := p.byRepo[repository]
	if !ok {
		return ModulePlan{}, false
	}
	return p.modules[index], true
}

// ContainsRepository reports whether repository is planned for publication.
func (p Plan) ContainsRepository(repository manifest.RepositoryRef) bool {
	_, ok := p.byRepo[repository]
	return ok
}
