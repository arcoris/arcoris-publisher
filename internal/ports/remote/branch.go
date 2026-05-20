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

package remote

// BranchProtection describes provider-level branch protection metadata.
type BranchProtection struct {
	// Protected reports whether the provider has protection enabled.
	Protected bool
	// RequiresPullRequest reports whether changes must arrive through a PR.
	RequiresPullRequest bool
	// RequiresStatusChecks reports whether checks must pass before merge.
	RequiresStatusChecks bool
	// AllowsForcePushes reports whether forced updates are allowed.
	AllowsForcePushes bool
	// AllowsDeletions reports whether branch deletion is allowed.
	AllowsDeletions bool
}

// BlocksDirectPush reports whether the protection settings imply that a direct
// branch push is likely to be rejected.
func (p BranchProtection) BlocksDirectPush() bool {
	return p.Protected && p.RequiresPullRequest
}
