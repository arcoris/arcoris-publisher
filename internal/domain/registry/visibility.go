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

// ModulesByVisibility returns modules with visibility in manifest declaration order.
//
// Visibility indexes store declaration positions rather than module copies.
// This method resolves those positions into a fresh slice so callers can safely
// mutate the returned collection while the registry remains stable.
func (r Registry) ModulesByVisibility(visibility manifest.Visibility) []manifest.Module {
	indexes := r.byVisibility[visibility]
	modules := make([]manifest.Module, 0, len(indexes))
	for _, index := range indexes {
		modules = append(modules, r.modules[index])
	}
	return modules
}
