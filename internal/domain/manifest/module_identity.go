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

// moduleIdentity contains scalar fields that uniquely identify a module target.
type moduleIdentity struct {
	name       ModuleName
	modulePath ModulePath
	sourceDir  SourceDir
	repository RepositoryRef
}

// parseModuleIdentity validates the scalar identity fields of a module spec.
func parseModuleIdentity(spec ModuleSpec) (moduleIdentity, error) {
	name, err := ParseModuleName(spec.Name)
	if err != nil {
		return moduleIdentity{}, fmt.Errorf("name: %w", err)
	}
	modulePath, err := ParseModulePath(spec.ModulePath)
	if err != nil {
		return moduleIdentity{}, fmt.Errorf("module_path: %w", err)
	}
	sourceDir, err := ParseSourceDir(spec.SourceDir)
	if err != nil {
		return moduleIdentity{}, fmt.Errorf("source_dir: %w", err)
	}
	repository, err := ParseRepositoryRef(spec.Repository)
	if err != nil {
		return moduleIdentity{}, fmt.Errorf("repository: %w", err)
	}
	return moduleIdentity{name: name, modulePath: modulePath, sourceDir: sourceDir, repository: repository}, nil
}
