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

// ModuleNames returns planned module names in dependency-first publish order.
func (p Plan) ModuleNames() []manifest.ModuleName {
	values := make([]manifest.ModuleName, 0, len(p.modules))
	for _, module := range p.modules {
		values = append(values, module.Name())
	}
	return values
}

// PublishOrder is an alias for ModuleNames.
func (p Plan) PublishOrder() []manifest.ModuleName { return p.ModuleNames() }
