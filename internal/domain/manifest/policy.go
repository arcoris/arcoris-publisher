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

// PolicySpec is the raw global publication policy declaration.
type PolicySpec struct {
	VersionPolicy string `json:"version_policy" yaml:"version_policy"`
	PushPolicy    string `json:"push_policy" yaml:"push_policy"`
}

// Policy is the validated global publication policy.
type Policy struct {
	versionPolicy VersionPolicy
	pushPolicy    PushPolicy
}

// NewPolicy validates spec and returns Policy.
func NewPolicy(spec PolicySpec) (Policy, error) {
	versionPolicy, err := ParseVersionPolicy(defaultString(spec.VersionPolicy, string(VersionPolicyReleaseTrain)))
	if err != nil {
		return Policy{}, fmt.Errorf("version_policy: %w", err)
	}
	pushPolicy, err := ParsePushPolicy(defaultString(spec.PushPolicy, string(PushPolicyFastForwardOnly)))
	if err != nil {
		return Policy{}, fmt.Errorf("push_policy: %w", err)
	}
	return Policy{versionPolicy: versionPolicy, pushPolicy: pushPolicy}, nil
}

// VersionPolicy returns the global module version assignment policy.
func (p Policy) VersionPolicy() VersionPolicy { return p.versionPolicy }

// PushPolicy returns the global Git push policy.
func (p Policy) PushPolicy() PushPolicy { return p.pushPolicy }

// Spec returns a serializable policy declaration.
func (p Policy) Spec() PolicySpec {
	return PolicySpec{VersionPolicy: string(p.versionPolicy), PushPolicy: string(p.pushPolicy)}
}
