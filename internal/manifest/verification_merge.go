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

// MergeVerification applies override over base and returns a complete policy.
func MergeVerification(base VerificationPolicy, override VerificationOverride) VerificationPolicy {
	out := base
	if override.localReplacePolicy != nil {
		out.localReplacePolicy = *override.localReplacePolicy
	}
	out.goPolicy = MergeGoVerification(base.goPolicy, override.goPolicy)
	return out
}

// MergeGoVerification applies override over base and returns a complete Go policy.
func MergeGoVerification(
	base GoVerificationPolicy,
	override GoVerificationOverride,
) GoVerificationPolicy {
	out := base
	if override.workspaceMode != nil {
		out.workspaceMode = *override.workspaceMode
	}
	if override.list != nil {
		out.list = *override.list
	}
	if override.test != nil {
		out.test = *override.test
	}
	if override.tidy != nil {
		out.tidy = *override.tidy
	}
	if override.patternsSet {
		out.patterns = cloneStrings(override.patterns)
	}
	return out
}
