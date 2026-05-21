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

// VerificationOverride is a validated partial verification policy.
type VerificationOverride struct {
	localReplacePolicy *LocalReplacePolicy
	goPolicy           GoVerificationOverride
}

// NewVerificationOverride validates a partial verification policy declaration.
func NewVerificationOverride(spec VerificationSpec) (VerificationOverride, error) {
	var collector IssueCollector
	var override VerificationOverride

	if spec.LocalReplacePolicy != nil {
		policy, err := ParseLocalReplacePolicy(*spec.LocalReplacePolicy)
		collector.AddError("localReplacePolicy", err)
		if err == nil {
			override.localReplacePolicy = &policy
		}
	}

	goOverride, err := NewGoVerificationOverride(spec.Go)
	collector.AddError("go", err)
	override.goPolicy = goOverride

	if err := collector.Err(); err != nil {
		return VerificationOverride{}, err
	}

	return override, nil
}
