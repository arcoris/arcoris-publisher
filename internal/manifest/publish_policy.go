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

// PublishSpec is the raw top-level publication policy declaration.
type PublishSpec struct {
	Mode          *string        `json:"mode,omitempty" yaml:"mode,omitempty"`
	VersionPolicy *string        `json:"versionPolicy,omitempty" yaml:"versionPolicy,omitempty"`
	PushPolicy    *string        `json:"pushPolicy,omitempty" yaml:"pushPolicy,omitempty"`
	Tags          TagPolicySpec  `json:"tags,omitempty" yaml:"tags,omitempty"`
	Provenance    ProvenanceSpec `json:"provenance,omitempty" yaml:"provenance,omitempty"`
}

// PublishPolicy is the validated top-level publication policy.
type PublishPolicy struct {
	mode          PublishMode
	versionPolicy VersionPolicy
	pushPolicy    PushPolicy
	tags          TagPolicy
	provenance    ProvenancePolicy
}

// NewPublishPolicy validates spec and applies safe built-in publication defaults.
func NewPublishPolicy(spec PublishSpec) (PublishPolicy, error) {
	mode, err := ParsePublishMode(stringValue(spec.Mode, string(PublishModeExplicitProjection)))
	if err != nil {
		return PublishPolicy{}, err
	}
	versionPolicy, err := ParseVersionPolicy(stringValue(spec.VersionPolicy, string(VersionPolicyReleaseTrain)))
	if err != nil {
		return PublishPolicy{}, err
	}
	pushPolicy, err := ParsePushPolicy(stringValue(spec.PushPolicy, string(PushPolicyFastForwardOnly)))
	if err != nil {
		return PublishPolicy{}, err
	}
	tags, err := NewTagPolicy(spec.Tags)
	if err != nil {
		return PublishPolicy{}, err
	}
	provenance, err := NewProvenancePolicy(spec.Provenance)
	if err != nil {
		return PublishPolicy{}, err
	}
	return PublishPolicy{mode: mode, versionPolicy: versionPolicy, pushPolicy: pushPolicy, tags: tags, provenance: provenance}, nil
}

// Mode returns the publication construction mode.
func (p PublishPolicy) Mode() PublishMode { return p.mode }

// VersionPolicy returns the module version assignment policy.
func (p PublishPolicy) VersionPolicy() VersionPolicy { return p.versionPolicy }

// PushPolicy returns the Git push policy.
func (p PublishPolicy) PushPolicy() PushPolicy { return p.pushPolicy }

// Tags returns the tag creation policy.
func (p PublishPolicy) Tags() TagPolicy { return p.tags }

// Provenance returns the provenance policy.
func (p PublishPolicy) Provenance() ProvenancePolicy { return p.provenance }
