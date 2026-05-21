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

package versioning

import (
	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
	"arcoris.dev/arcoris-publisher/internal/domain/registry"
)

// RequirementMap returns dependency module paths mapped to assigned versions for
// module's declared direct dependencies.
func (a Assignments) RequirementMap(registryValue registry.Registry, module manifest.Module) (map[manifest.ModulePath]Version, error) {
	dependencies, err := a.DependencyVersions(registryValue, module)
	if err != nil {
		return nil, err
	}
	requirements := make(map[manifest.ModulePath]Version, len(dependencies))
	for _, dependency := range dependencies {
		requirements[dependency.ModulePath()] = dependency.Version()
	}
	return requirements, nil
}
