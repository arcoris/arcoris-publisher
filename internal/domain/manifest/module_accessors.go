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

// Name returns the module's manifest-local identifier.
func (m Module) Name() ModuleName { return m.name }

// ModulePath returns the public Go module path.
func (m Module) ModulePath() ModulePath { return m.modulePath }

// SourceDir returns the repository-relative staged module directory.
func (m Module) SourceDir() SourceDir { return m.sourceDir }

// Repository returns the target repository reference.
func (m Module) Repository() RepositoryRef { return m.repository }

// Branches returns detached source-to-target branch mappings.
func (m Module) Branches() []BranchMapping {
	return append([]BranchMapping(nil), m.branches...)
}

// Dependencies returns detached direct dependency declarations.
func (m Module) Dependencies() []Dependency {
	return append([]Dependency(nil), m.dependencies...)
}

// Visibility returns the module publication visibility.
func (m Module) Visibility() Visibility { return m.visibility }

// Publishable reports whether the module should be considered for external publication.
func (m Module) Publishable() bool { return m.visibility == VisibilityPublic }
