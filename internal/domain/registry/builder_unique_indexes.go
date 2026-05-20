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

// indexModuleName records the manifest-local module name lookup.
func (b *builder) indexModuleName(index int, name manifest.ModuleName) {
	if previous, exists := b.registry.byName[name]; exists {
		b.addIssue(duplicateModuleNameIssue(name, b.modules[previous].Name()))
		return
	}
	b.registry.byName[name] = index
}

// indexModulePath records the public Go module path lookup.
func (b *builder) indexModulePath(index int, name manifest.ModuleName, modulePath manifest.ModulePath) {
	if previous, exists := b.registry.byModulePath[modulePath]; exists {
		b.addIssue(duplicateModulePathIssue(name, b.modules[previous].Name(), modulePath))
		return
	}
	b.registry.byModulePath[modulePath] = index
}

// indexSourceDir records the staged source directory lookup.
func (b *builder) indexSourceDir(index int, name manifest.ModuleName, sourceDir manifest.SourceDir) {
	if previous, exists := b.registry.bySourceDir[sourceDir]; exists {
		b.addIssue(duplicateSourceDirIssue(name, b.modules[previous].Name(), sourceDir))
		return
	}
	b.registry.bySourceDir[sourceDir] = index
}

// indexRepository records the target repository lookup.
func (b *builder) indexRepository(index int, name manifest.ModuleName, repository manifest.RepositoryRef) {
	if previous, exists := b.registry.byRepository[repository]; exists {
		b.addIssue(duplicateRepositoryIssue(name, b.modules[previous].Name(), repository))
		return
	}
	b.registry.byRepository[repository] = index
}
