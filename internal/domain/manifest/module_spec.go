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

// Spec returns a detached serializable representation of the module.
func (m Module) Spec() ModuleSpec {
	return ModuleSpec{
		Name:         string(m.name),
		ModulePath:   string(m.modulePath),
		SourceDir:    string(m.sourceDir),
		Repository:   string(m.repository),
		Branches:     branchSpecs(m.branches),
		Dependencies: dependencySpecs(m.dependencies),
		Visibility:   string(m.visibility),
	}
}

// branchSpecs converts validated branch mappings back into DTOs.
func branchSpecs(branches []BranchMapping) []BranchMappingSpec {
	out := make([]BranchMappingSpec, 0, len(branches))
	for _, branch := range branches {
		out = append(out, branch.Spec())
	}
	return out
}

// dependencySpecs converts validated dependencies back into raw module names.
func dependencySpecs(dependencies []Dependency) []string {
	out := make([]string, 0, len(dependencies))
	for _, dependency := range dependencies {
		out = append(out, string(dependency.Module()))
	}
	return out
}
