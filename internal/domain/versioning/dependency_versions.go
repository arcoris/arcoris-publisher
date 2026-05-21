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
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
	"arcoris.dev/arcoris-publisher/internal/domain/registry"
)

// DependencyVersions returns version requirements for module's declared direct
// dependencies.
//
// The method fails if a dependency is known to the registry but does not have a
// publish version assignment. That usually means a public module depends on an
// internal or disabled module and the manifest policy must be corrected before a
// publish plan can be built.
func (a Assignments) DependencyVersions(registryValue registry.Registry, module manifest.Module) ([]DependencyVersion, error) {
	dependencies := module.Dependencies()
	versions := make([]DependencyVersion, 0, len(dependencies))
	for i, dependency := range dependencies {
		dependencyVersion, err := a.dependencyVersion(registryValue, dependency, i)
		if err != nil {
			return nil, err
		}
		versions = append(versions, dependencyVersion)
	}
	return versions, nil
}

// dependencyVersion resolves one direct dependency to its required version.
func (a Assignments) dependencyVersion(registryValue registry.Registry, dependency manifest.Dependency, index int) (DependencyVersion, error) {
	path := fmt.Sprintf("dependencies[%d]", index)
	dependencyModule, ok := registryValue.ModuleByName(dependency.Module())
	if !ok {
		return DependencyVersion{}, validationErrorf(IssueUnknownDependency, path, "unknown dependency %q", dependency.Module())
	}
	version, ok := a.VersionOfModule(dependency.Module())
	if !ok {
		return DependencyVersion{}, validationErrorf(IssueUnassignedDependency, path, "dependency %q has no assigned publish version", dependency.Module())
	}
	return DependencyVersion{module: dependency.Module(), modulePath: dependencyModule.ModulePath(), version: version}, nil
}
