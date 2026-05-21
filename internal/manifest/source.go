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

// SourceSpec is the raw top-level source repository declaration.
type SourceSpec struct {
	Repository    string  `json:"repository" yaml:"repository"`
	DefaultBranch string  `json:"defaultBranch" yaml:"defaultBranch"`
	StagingRoot   *string `json:"stagingRoot,omitempty" yaml:"stagingRoot,omitempty"`
	ModuleRoot    *string `json:"moduleRoot,omitempty" yaml:"moduleRoot,omitempty"`
	DirtyPolicy   *string `json:"dirtyPolicy,omitempty" yaml:"dirtyPolicy,omitempty"`
}

// Source is the validated top-level source repository declaration.
type Source struct {
	repository    RepositoryRef
	defaultBranch BranchName
	stagingRoot   RelativePath
	moduleRoot    RelativePath
	dirtyPolicy   DirtyPolicy
}

// NewSource validates spec and applies safe built-in source defaults.
func NewSource(spec SourceSpec) (Source, error) {
	repository, err := ParseRepositoryRef(spec.Repository)
	if err != nil {
		return Source{}, fmt.Errorf("repository: %w", err)
	}
	defaultBranch, err := ParseBranchName(spec.DefaultBranch)
	if err != nil {
		return Source{}, fmt.Errorf("defaultBranch: %w", err)
	}
	stagingRootValue := stringValue(spec.StagingRoot, ".")
	stagingRoot, err := ParseRelativePath("stagingRoot", stagingRootValue, true)
	if err != nil {
		return Source{}, err
	}
	moduleRootValue := stringValue(spec.ModuleRoot, ".")
	moduleRoot, err := ParseRelativePath("moduleRoot", moduleRootValue, true)
	if err != nil {
		return Source{}, err
	}
	dirtyPolicy, err := ParseDirtyPolicy(stringValue(spec.DirtyPolicy, string(DirtyPolicyFail)))
	if err != nil {
		return Source{}, err
	}
	return Source{repository: repository, defaultBranch: defaultBranch, stagingRoot: stagingRoot, moduleRoot: moduleRoot, dirtyPolicy: dirtyPolicy}, nil
}

// Repository returns the authoritative source repository reference.
func (s Source) Repository() RepositoryRef { return s.repository }

// DefaultBranch returns the default source branch.
func (s Source) DefaultBranch() BranchName { return s.defaultBranch }

// StagingRoot returns the staging root relative to arcpub.yaml.
func (s Source) StagingRoot() RelativePath { return s.stagingRoot }

// ModuleRoot returns the expected module root under the staging root.
func (s Source) ModuleRoot() RelativePath { return s.moduleRoot }

// DirtyPolicy returns the source dirty repository policy.
func (s Source) DirtyPolicy() DirtyPolicy { return s.dirtyPolicy }

// Spec returns a serializable source declaration.
func (s Source) Spec() SourceSpec {
	stagingRoot := string(s.stagingRoot)
	moduleRoot := string(s.moduleRoot)
	dirtyPolicy := string(s.dirtyPolicy)
	return SourceSpec{Repository: string(s.repository), DefaultBranch: string(s.defaultBranch), StagingRoot: &stagingRoot, ModuleRoot: &moduleRoot, DirtyPolicy: &dirtyPolicy}
}
