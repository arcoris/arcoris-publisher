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

// ModuleSpec is the raw declaration of a publishable or known module.
type ModuleSpec struct {
	// Name is the stable manifest-local module identifier.
	Name string `json:"name" yaml:"name"`
	// ModulePath is the public Go module path written into go.mod.
	ModulePath string `json:"module_path" yaml:"module_path"`
	// SourceDir is the repository-relative staged module directory.
	SourceDir string `json:"source_dir" yaml:"source_dir"`
	// Repository is the target repository in owner/name form.
	Repository string `json:"repository" yaml:"repository"`
	// Branches maps source branches to target branches.
	Branches []BranchMappingSpec `json:"branches" yaml:"branches"`
	// Dependencies declares allowed direct dependencies on other manifest modules.
	Dependencies []string `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	// Visibility controls whether the module is externally publishable.
	Visibility string `json:"visibility,omitempty" yaml:"visibility,omitempty"`
}

// Module is a validated manifest module declaration.
type Module struct {
	name         ModuleName
	modulePath   ModulePath
	sourceDir    SourceDir
	repository   RepositoryRef
	branches     []BranchMapping
	dependencies []Dependency
	visibility   Visibility
}

// NewModule validates spec and returns a Module.
func NewModule(spec ModuleSpec) (Module, error) {
	identity, err := parseModuleIdentity(spec)
	if err != nil {
		return Module{}, err
	}
	branches, err := parseModuleBranches(spec.Branches)
	if err != nil {
		return Module{}, err
	}
	dependencies, err := parseModuleDependencies(spec.Dependencies, identity.name)
	if err != nil {
		return Module{}, err
	}
	visibility, err := parseModuleVisibility(spec.Visibility)
	if err != nil {
		return Module{}, err
	}

	module := Module{
		name:         identity.name,
		modulePath:   identity.modulePath,
		sourceDir:    identity.sourceDir,
		repository:   identity.repository,
		branches:     branches,
		dependencies: dependencies,
		visibility:   visibility,
	}
	if err := module.validateLocalUniqueness(); err != nil {
		return Module{}, err
	}
	return module, nil
}
