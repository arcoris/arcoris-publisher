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

// ModulePlan returns the plan for module and whether it was found.
func (p Plan) ModulePlan(module manifest.ModuleName) (ModulePlan, bool) {
	index, ok := p.byModule[module]
	if !ok {
		return ModulePlan{}, false
	}
	return p.modules[index], true
}

// ContainsModule reports whether module is planned for publication.
func (p Plan) ContainsModule(module manifest.ModuleName) bool {
	_, ok := p.byModule[module]
	return ok
}
