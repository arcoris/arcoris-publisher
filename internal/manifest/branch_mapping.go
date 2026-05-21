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

// BranchMappingSpec is the raw source-to-target branch mapping declaration.
type BranchMappingSpec struct {
	Source string `json:"source" yaml:"source"`
	Target string `json:"target" yaml:"target"`
}

// BranchMapping maps one source repository branch to one target repository branch.
//
// Staging manifests can publish a module from a source branch into a different
// target branch. The validated value keeps both sides as BranchName so callers
// never have to re-check Git-ref safety before using the mapping.
type BranchMapping struct {
	source BranchName
	target BranchName
}

// NewBranchMapping validates spec and returns a branch mapping.
func NewBranchMapping(spec BranchMappingSpec) (BranchMapping, error) {
	source, err := ParseBranchName(spec.Source)
	if err != nil {
		return BranchMapping{}, fmt.Errorf("source: %w", err)
	}
	target, err := ParseBranchName(spec.Target)
	if err != nil {
		return BranchMapping{}, fmt.Errorf("target: %w", err)
	}
	return BranchMapping{source: source, target: target}, nil
}

// Source returns the source repository branch.
func (m BranchMapping) Source() BranchName { return m.source }

// Target returns the target repository branch.
func (m BranchMapping) Target() BranchName { return m.target }

// Spec returns a serializable branch mapping representation.
func (m BranchMapping) Spec() BranchMappingSpec {
	return BranchMappingSpec{Source: string(m.source), Target: string(m.target)}
}
