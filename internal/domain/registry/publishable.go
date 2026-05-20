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

import "arcoris.dev/arcoris-publisher/internal/domain/manifest"

// PublishableModules returns public modules in manifest declaration order.
//
// Public modules are the only modules eligible for external publication. The
// name "publishable" mirrors Module.Publishable and keeps publication planning
// code away from raw visibility constants.
func (r Registry) PublishableModules() []manifest.Module {
	return r.ModulesByVisibility(manifest.VisibilityPublic)
}

// PublishableNames returns public module names in manifest declaration order.
func (r Registry) PublishableNames() []manifest.ModuleName {
	modules := r.PublishableModules()
	names := make([]manifest.ModuleName, 0, len(modules))
	for _, module := range modules {
		names = append(names, module.Name())
	}
	return names
}
