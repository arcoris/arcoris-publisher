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

package manifest

import "fmt"

// parseModuleBranches validates branch mappings and requires at least one mapping.
func parseModuleBranches(specs []BranchMappingSpec) ([]BranchMapping, error) {
	if len(specs) == 0 {
		return nil, fmt.Errorf("branches: at least one branch mapping is required")
	}
	branches := make([]BranchMapping, 0, len(specs))
	for i, spec := range specs {
		branch, err := NewBranchMapping(spec)
		if err != nil {
			return nil, fmt.Errorf("branches[%d]: %w", i, err)
		}
		branches = append(branches, branch)
	}
	return branches, nil
}

// parseModuleDependencies validates direct dependency declarations.
func parseModuleDependencies(specs []string, moduleName ModuleName) ([]Dependency, error) {
	dependencies := make([]Dependency, 0, len(specs))
	for i, spec := range specs {
		dependency, err := NewDependency(spec)
		if err != nil {
			return nil, fmt.Errorf("dependencies[%d]: %w", i, err)
		}
		if dependency.Module() == moduleName {
			return nil, fmt.Errorf("dependencies[%d]: module %q cannot depend on itself", i, moduleName)
		}
		dependencies = append(dependencies, dependency)
	}
	return dependencies, nil
}

// parseModuleVisibility validates module visibility with the default policy.
func parseModuleVisibility(value string) (Visibility, error) {
	visibility, err := ParseVisibility(value)
	if err != nil {
		return "", fmt.Errorf("visibility: %w", err)
	}
	return visibility, nil
}
