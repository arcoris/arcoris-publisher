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

// ProvenanceSpec is the raw provenance policy declaration.
type ProvenanceSpec struct {
	CommitTrailers *bool   `json:"commitTrailers,omitempty" yaml:"commitTrailers,omitempty"`
	File           *string `json:"file,omitempty" yaml:"file,omitempty"`
}

// ProvenancePolicy is the validated provenance policy.
type ProvenancePolicy struct {
	commitTrailers bool
	file           RelativePath
	fileEnabled    bool
}

// NewProvenancePolicy validates spec and applies built-in provenance defaults.
func NewProvenancePolicy(spec ProvenanceSpec) (ProvenancePolicy, error) {
	commitTrailers := boolValue(spec.CommitTrailers, true)
	if spec.File == nil || *spec.File == "" {
		return ProvenancePolicy{commitTrailers: commitTrailers}, nil
	}

	var collector IssueCollector

	file, err := ParseRelativePath("provenance.file", *spec.File, false)
	collector.AddError("file", err)

	if err := collector.Err(); err != nil {
		return ProvenancePolicy{}, err
	}

	return ProvenancePolicy{commitTrailers: commitTrailers, file: file, fileEnabled: true}, nil
}

// CommitTrailers reports whether provenance trailers must be added to commits.
func (p ProvenancePolicy) CommitTrailers() bool { return p.commitTrailers }

// FileEnabled reports whether a provenance file should be generated.
func (p ProvenancePolicy) FileEnabled() bool { return p.fileEnabled }

// File returns the provenance file path when FileEnabled is true.
func (p ProvenancePolicy) File() RelativePath { return p.file }
