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

// GoVerificationOverride is a validated partial Go verification policy.
type GoVerificationOverride struct {
	workspaceMode *GoWorkspaceMode
	list          *bool
	test          *bool
	tidy          *bool
	patterns      []string
	patternsSet   bool
}

// NewGoVerificationOverride validates a partial Go verification policy.
func NewGoVerificationOverride(spec GoVerifySpec) (GoVerificationOverride, error) {
	var collector IssueCollector
	var override GoVerificationOverride

	if spec.WorkspaceMode != nil {
		mode, err := ParseGoWorkspaceMode(*spec.WorkspaceMode)
		collector.AddError("workspaceMode", err)
		if err == nil {
			override.workspaceMode = &mode
		}
	}

	override.list = spec.List
	override.test = spec.Test
	override.tidy = spec.Tidy

	if spec.Patterns != nil {
		patterns := make([]string, 0, len(spec.Patterns))
		for i, pattern := range spec.Patterns {
			if pattern == "" {
				collector.Add(IssueInvalidValue, fmt.Sprintf("patterns[%d]", i), "must not be empty")
				continue
			}
			patterns = append(patterns, pattern)
		}
		override.patterns = patterns
		override.patternsSet = true
	}

	if err := collector.Err(); err != nil {
		return GoVerificationOverride{}, err
	}

	return override, nil
}
