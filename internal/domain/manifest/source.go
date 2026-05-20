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

// SourceSpec is the raw authoritative source repository declaration.
type SourceSpec struct {
	Repository    string `json:"repository" yaml:"repository"`
	DefaultBranch string `json:"default_branch" yaml:"default_branch"`
}

// Source is the validated authoritative source repository declaration.
type Source struct {
	repository    RepositoryRef
	defaultBranch BranchName
}

// NewSource validates spec and returns Source.
func NewSource(spec SourceSpec) (Source, error) {
	repository, err := ParseRepositoryRef(spec.Repository)
	if err != nil {
		return Source{}, fmt.Errorf("repository: %w", err)
	}
	branch, err := ParseBranchName(spec.DefaultBranch)
	if err != nil {
		return Source{}, fmt.Errorf("default_branch: %w", err)
	}
	return Source{repository: repository, defaultBranch: branch}, nil
}

// Repository returns the authoritative source repository.
func (s Source) Repository() RepositoryRef { return s.repository }

// DefaultBranch returns the default source branch.
func (s Source) DefaultBranch() BranchName { return s.defaultBranch }

// Spec returns a serializable source declaration.
func (s Source) Spec() SourceSpec {
	return SourceSpec{Repository: string(s.repository), DefaultBranch: string(s.defaultBranch)}
}
