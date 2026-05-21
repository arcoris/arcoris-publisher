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

// ModuleSpec is one raw top-level module routing declaration.
type ModuleSpec struct {
	Name       string                       `json:"name" yaml:"name"`
	SourceDir  string                       `json:"sourceDir" yaml:"sourceDir"`
	Manifest   *string                      `json:"manifest,omitempty" yaml:"manifest,omitempty"`
	Repository string                       `json:"repository" yaml:"repository"`
	Visibility *string                      `json:"visibility,omitempty" yaml:"visibility,omitempty"`
	Branches   []manifest.BranchMappingSpec `json:"branches,omitempty" yaml:"branches,omitempty"`
}

// Module is a validated top-level module routing declaration.
type Module struct {
	name          manifest.ModuleName
	sourceDir     manifest.SourceDir
	manifestPath  manifest.RelativePath
	manifestSet   bool
	repository    manifest.RepositoryRef
	visibility    manifest.Visibility
	visibilitySet bool
	branches      []manifest.BranchMapping
	branchesSet   bool
}

// NewModule validates one top-level module routing declaration.
func NewModule(spec ModuleSpec) (Module, error) {
	var collector manifest.IssueCollector
	name, err := manifest.ParseModuleName(spec.Name)
	collector.AddError("name", err)
	sourceDir, err := manifest.ParseSourceDir(spec.SourceDir)
	collector.AddError("sourceDir", err)
	repository, err := manifest.ParseRepositoryRef(spec.Repository)
	collector.AddError("repository", err)
	var manifestPath manifest.RelativePath
	manifestSet := false
	if spec.Manifest != nil {
		manifestSet = true
		manifestPath, err = manifest.ParseRelativePath("manifest", *spec.Manifest, false)
		collector.AddError("manifest", err)
	}
	var visibility manifest.Visibility
	visibilitySet := false
	if spec.Visibility != nil {
		visibilitySet = true
		visibility, err = manifest.ParseVisibility(*spec.Visibility)
		collector.AddError("visibility", err)
	}
	branches := make([]manifest.BranchMapping, 0, len(spec.Branches))
	for i, branchSpec := range spec.Branches {
		branch, err := manifest.NewBranchMapping(branchSpec)
		if err != nil {
			collector.AddError(fmt.Sprintf("branches[%d]", i), err)
			continue
		}
		branches = append(branches, branch)
	}
	if err := collector.Err(); err != nil {
		return Module{}, err
	}
	return Module{name: name, sourceDir: sourceDir, manifestPath: manifestPath, manifestSet: manifestSet, repository: repository, visibility: visibility, visibilitySet: visibilitySet, branches: branches, branchesSet: spec.Branches != nil}, nil
}

// Name returns the top-level module name.
func (m Module) Name() manifest.ModuleName { return m.name }

// SourceDir returns the staged source directory of the module.
func (m Module) SourceDir() manifest.SourceDir { return m.sourceDir }

// Repository returns the target repository for the module.
func (m Module) Repository() manifest.RepositoryRef { return m.repository }

// ManifestPathOverride returns an explicitly declared module manifest path.
func (m Module) ManifestPathOverride() (manifest.RelativePath, bool) {
	return m.manifestPath, m.manifestSet
}

// VisibilityOverride returns an explicitly declared module visibility.
func (m Module) VisibilityOverride() (manifest.Visibility, bool) {
	return m.visibility, m.visibilitySet
}

// BranchesOverride returns explicitly declared module branch mappings.
func (m Module) BranchesOverride() ([]manifest.BranchMapping, bool) {
	return manifest.CloneBranchMappings(m.branches), m.branchesSet
}
