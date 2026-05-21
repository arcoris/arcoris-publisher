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

package plan

import (
	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/versioning"
)

// ModulePlan is the executable publication intent for one public module.
//
// ModulePlan deliberately contains only resolved, effective values. Workflow
// packages should not return to raw manifests to discover branch mappings,
// publish entries, verification settings, dependency versions, or source and
// target routing.
type ModulePlan struct {
	// name is the resolved manifest module name.
	name manifest.ModuleName

	// modulePath is the public Go module path.
	modulePath manifest.ModulePath

	// moduleType describes the module layout expected by workflow stages.
	moduleType manifest.ModuleType

	// sourceDir is the source directory relative to the staging root.
	sourceDir manifest.SourceDir

	// moduleRoot is the module root relative to arcpub.module.yaml.
	moduleRoot manifest.RelativePath

	// goMod is the go.mod path relative to ModuleRoot.
	goMod manifest.RelativePath

	// repository is the target repository for publication.
	repository manifest.RepositoryRef

	// visibility is retained so validation and workflow code can assert that
	// executable plans contain only public modules.
	visibility manifest.Visibility

	// version is the assigned publication version.
	version versioning.Version

	// branches are source-to-target branch mappings for this module.
	branches []BranchPlan

	// entries are explicit publication content projections.
	entries []manifest.PublishEntry

	// requirements are direct internal dependency versions for go.mod updates.
	requirements []DependencyRequirement

	// verification is the effective verification policy for this module.
	verification manifest.VerificationPolicy
}

// Name returns the staged module identity.
func (m ModulePlan) Name() manifest.ModuleName { return m.name }

// ModulePath returns the public Go module path.
func (m ModulePlan) ModulePath() manifest.ModulePath { return m.modulePath }

// ModuleType returns the publication module type.
func (m ModulePlan) ModuleType() manifest.ModuleType { return m.moduleType }

// SourceDir returns the module source directory relative to the staging root.
func (m ModulePlan) SourceDir() manifest.SourceDir { return m.sourceDir }

// ModuleRoot returns the module root relative to arcpub.module.yaml.
func (m ModulePlan) ModuleRoot() manifest.RelativePath { return m.moduleRoot }

// GoMod returns the go.mod path relative to ModuleRoot.
func (m ModulePlan) GoMod() manifest.RelativePath { return m.goMod }

// Repository returns the target repository.
func (m ModulePlan) Repository() manifest.RepositoryRef { return m.repository }

// Visibility returns the resolved module visibility.
func (m ModulePlan) Visibility() manifest.Visibility { return m.visibility }

// Version returns the assigned publication version.
func (m ModulePlan) Version() versioning.Version { return m.version }

// Branches returns detached branch publication mappings.
func (m ModulePlan) Branches() []BranchPlan { return cloneBranchPlans(m.branches) }

// PublishEntries returns detached explicit publish entries.
func (m ModulePlan) PublishEntries() []manifest.PublishEntry {
	return manifest.ClonePublishEntries(m.entries)
}

// Requirements returns detached direct internal dependency requirements.
func (m ModulePlan) Requirements() []DependencyRequirement {
	return cloneRequirements(m.requirements)
}

// Verification returns the effective module verification policy.
func (m ModulePlan) Verification() manifest.VerificationPolicy { return m.verification }
