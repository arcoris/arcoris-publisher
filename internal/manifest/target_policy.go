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

// TargetSpec is the raw top-level target-worktree policy declaration.
type TargetSpec struct {
	RemoteTemplate *string `json:"remoteTemplate,omitempty" yaml:"remoteTemplate,omitempty"`
}

// TargetPolicy is the validated target-worktree preparation policy.
type TargetPolicy struct {
	remoteTemplate RemoteTemplate
	hasTemplate    bool
}

// NewTargetPolicy validates spec and applies empty defaults.
func NewTargetPolicy(spec TargetSpec) (TargetPolicy, error) {
	if spec.RemoteTemplate == nil {
		return TargetPolicy{}, nil
	}

	var collector IssueCollector
	template, err := ParseRemoteTemplate(*spec.RemoteTemplate)
	collector.AddError("remoteTemplate", err)
	if err := collector.Err(); err != nil {
		return TargetPolicy{}, err
	}
	return TargetPolicy{remoteTemplate: template, hasTemplate: true}, nil
}

// RemoteTemplate returns the configured clone/fetch URL template.
func (p TargetPolicy) RemoteTemplate() (RemoteTemplate, bool) {
	return p.remoteTemplate, p.hasTemplate
}
