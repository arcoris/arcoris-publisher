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

package staging

import (
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

// DefaultsSpec is the raw top-level defaults declaration from arcpub.yaml.
type DefaultsSpec struct {
	Branches       []manifest.BranchMappingSpec `json:"branches,omitempty" yaml:"branches,omitempty"`
	ModuleManifest ModuleManifestDefaultsSpec   `json:"moduleManifest,omitempty" yaml:"moduleManifest,omitempty"`
	Verification   manifest.VerificationSpec    `json:"verification,omitempty" yaml:"verification,omitempty"`
}

// Defaults is the validated top-level defaults declaration.
type Defaults struct {
	branches       []manifest.BranchMapping
	branchesSet    bool
	moduleManifest ModuleManifestDefaults
	verification   manifest.VerificationOverride
}

// NewDefaults validates spec and returns Defaults.
func NewDefaults(spec DefaultsSpec) (Defaults, error) {
	var collector manifest.IssueCollector
	branches := make([]manifest.BranchMapping, 0, len(spec.Branches))
	for i, branchSpec := range spec.Branches {
		branch, err := manifest.NewBranchMapping(branchSpec)
		if err != nil {
			collector.AddError(fmt.Sprintf("defaults.branches[%d]", i), err)
			continue
		}
		branches = append(branches, branch)
	}
	moduleDefaults, err := NewModuleManifestDefaults(spec.ModuleManifest)
	collector.AddError("defaults.moduleManifest", err)
	verification, err := manifest.NewVerificationOverride(spec.Verification)
	collector.AddError("defaults.verification", err)
	if err := collector.Err(); err != nil {
		return Defaults{}, err
	}
	return Defaults{branches: branches, branchesSet: spec.Branches != nil, moduleManifest: moduleDefaults, verification: verification}, nil
}

// Branches returns detached default branch mappings.
func (d Defaults) Branches() []manifest.BranchMapping {
	return manifest.CloneBranchMappings(d.branches)
}

// BranchesSet reports whether defaults.branches was explicitly declared.
func (d Defaults) BranchesSet() bool { return d.branchesSet }

// ModuleManifest returns the default module manifest location policy.
func (d Defaults) ModuleManifest() ModuleManifestDefaults { return d.moduleManifest }

// Verification returns the top-level verification defaults as an override.
func (d Defaults) Verification() manifest.VerificationOverride { return d.verification }
