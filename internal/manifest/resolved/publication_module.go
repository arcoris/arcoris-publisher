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

package resolved

import "arcoris.dev/arcoris-publisher/internal/manifest"

// PublicationModule is the effective resolved model for one staged module.
type PublicationModule struct {
	name         manifest.ModuleName
	sourceDir    manifest.SourceDir
	manifestPath manifest.RelativePath
	repository   manifest.RepositoryRef
	visibility   manifest.Visibility
	branches     []manifest.BranchMapping

	moduleType manifest.ModuleType
	modulePath manifest.ModulePath
	moduleRoot manifest.RelativePath
	goMod      manifest.RelativePath

	dependencies []manifest.ModuleName
	entries      []manifest.PublishEntry
	verification manifest.VerificationPolicy
}

// Name returns the publication module name.
func (m PublicationModule) Name() manifest.ModuleName { return m.name }

// SourceDir returns the staged source directory.
func (m PublicationModule) SourceDir() manifest.SourceDir { return m.sourceDir }

// ManifestPath returns the module manifest path relative to SourceDir.
func (m PublicationModule) ManifestPath() manifest.RelativePath { return m.manifestPath }

// Repository returns the target repository.
func (m PublicationModule) Repository() manifest.RepositoryRef { return m.repository }

// Visibility returns the effective module visibility.
func (m PublicationModule) Visibility() manifest.Visibility { return m.visibility }

// Branches returns detached effective branch mappings.
func (m PublicationModule) Branches() []manifest.BranchMapping {
	return manifest.CloneBranchMappings(m.branches)
}

// ModuleType returns the module type.
func (m PublicationModule) ModuleType() manifest.ModuleType { return m.moduleType }

// ModulePath returns the published module path.
func (m PublicationModule) ModulePath() manifest.ModulePath { return m.modulePath }

// ModuleRoot returns the module root relative to arcpub.module.yaml.
func (m PublicationModule) ModuleRoot() manifest.RelativePath { return m.moduleRoot }

// GoMod returns the go.mod path relative to ModuleRoot.
func (m PublicationModule) GoMod() manifest.RelativePath { return m.goMod }

// Dependencies returns detached internal module dependencies.
func (m PublicationModule) Dependencies() []manifest.ModuleName {
	return manifest.CloneModuleNames(m.dependencies)
}

// PublishEntries returns detached explicit publication entries.
func (m PublicationModule) PublishEntries() []manifest.PublishEntry {
	return manifest.ClonePublishEntries(m.entries)
}

// Verification returns the effective verification policy.
func (m PublicationModule) Verification() manifest.VerificationPolicy { return m.verification }
