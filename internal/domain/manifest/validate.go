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

// Validate verifies aggregate-level manifest invariants.
//
// Constructors validate individual value objects. This aggregate pass checks
// relationships between modules: uniqueness across collections and dependency
// references that require knowledge of the full manifest.
func (m Manifest) Validate() error {
	validator := newManifestValidator(m)
	validator.validate()
	if len(validator.issues) > 0 {
		return &ValidationError{Issues: validator.issues}
	}
	return nil
}

// manifestValidator owns mutable validation indexes for one aggregate pass.
type manifestValidator struct {
	manifest     Manifest
	issues       []Issue
	moduleNames  map[ModuleName]int
	modulePaths  map[ModulePath]int
	sourceDirs   map[SourceDir]int
	repositories map[RepositoryRef]int
}

// newManifestValidator prepares uniqueness indexes sized for the manifest.
func newManifestValidator(manifest Manifest) manifestValidator {
	size := len(manifest.modules)
	return manifestValidator{
		manifest:     manifest,
		moduleNames:  make(map[ModuleName]int, size),
		modulePaths:  make(map[ModulePath]int, size),
		sourceDirs:   make(map[SourceDir]int, size),
		repositories: make(map[RepositoryRef]int, size),
	}
}

// validate runs every aggregate-level validation phase in stable order.
func (v *manifestValidator) validate() {
	v.validateRequiredFields()
	v.indexModules()
	v.validateDependencies()
}

// validateRequiredFields checks manifest-wide fields that cannot be empty.
func (v *manifestValidator) validateRequiredFields() {
	if v.manifest.version == "" {
		v.issues = append(v.issues, Issue{Code: IssueUnsupportedVersion, Path: "version", Message: "version is required"})
	}
	if len(v.manifest.modules) == 0 {
		v.issues = append(v.issues, Issue{Code: IssueInvalidModule, Path: "modules", Message: "at least one module is required"})
	}
}

// indexModules records uniqueness constraints for each module declaration.
func (v *manifestValidator) indexModules() {
	for i, module := range v.manifest.modules {
		path := issuePath("modules", i)
		v.recordModuleName(path, module.name, i)
		v.recordModulePath(path, module.modulePath, i)
		v.recordSourceDir(path, module.sourceDir, i)
		v.recordRepository(path, module.repository, i)
	}
}

// validateDependencies checks every dependency against the final module-name index.
func (v *manifestValidator) validateDependencies() {
	for i, module := range v.manifest.modules {
		modulePath := issuePath("modules", i)
		for j, dependency := range module.dependencies {
			v.validateDependency(module, dependency, modulePath, j)
		}
	}
}

// validateDependency checks one dependency edge inside one module declaration.
func (v *manifestValidator) validateDependency(module Module, dependency Dependency, modulePath string, index int) {
	dependencyPath := issuePath(modulePath+".dependencies", index)
	if dependency.Module() == module.name {
		v.issues = append(v.issues, Issue{Code: IssueInvalidDependency, Path: dependencyPath, Message: "module cannot depend on itself"})
		return
	}
	if _, exists := v.moduleNames[dependency.Module()]; !exists {
		v.issues = append(v.issues, Issue{Code: IssueUnknownDependency, Path: dependencyPath, Message: "unknown dependency " + quoteModuleName(dependency.Module())})
	}
}
