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

// TagPolicySpec is the raw tag policy declaration.
type TagPolicySpec struct {
	Enabled *bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Mode    *string `json:"mode,omitempty" yaml:"mode,omitempty"`
}

// TagMode describes the target repository tag naming mode.
type TagMode string

const (
	// TagModeSemver uses plain Go module semantic-version tags such as v0.3.0.
	TagModeSemver TagMode = "semver"
)

// TagPolicy is the validated tag policy.
type TagPolicy struct {
	enabled bool
	mode    TagMode
}

// NewTagPolicy validates spec and applies built-in tag defaults.
func NewTagPolicy(spec TagPolicySpec) (TagPolicy, error) {
	mode, err := ParseTagMode(stringValue(spec.Mode, string(TagModeSemver)))
	if err != nil {
		return TagPolicy{}, err
	}
	return TagPolicy{enabled: boolValue(spec.Enabled, true), mode: mode}, nil
}

// ParseTagMode validates a tag mode.
func ParseTagMode(value string) (TagMode, error) {
	switch TagMode(value) {
	case TagModeSemver:
		return TagMode(value), nil
	default:
		return "", fmt.Errorf("unsupported tag mode %q", value)
	}
}

// Enabled reports whether tags should be created for release publication.
func (p TagPolicy) Enabled() bool { return p.enabled }

// Mode returns the tag naming mode.
func (p TagPolicy) Mode() TagMode { return p.mode }
