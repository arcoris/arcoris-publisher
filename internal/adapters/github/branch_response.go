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

package github

import remoteport "arcoris.dev/arcoris-publisher/internal/ports/remote"

// protectionResponse mirrors the GitHub branch protection fields used by plans.
type protectionResponse struct {
	RequiredPullRequestReviews *struct{} `json:"required_pull_request_reviews"`
	RequiredStatusChecks       *struct{} `json:"required_status_checks"`
	AllowForcePushes           *struct {
		Enabled bool `json:"enabled"`
	} `json:"allow_force_pushes"`
	AllowDeletions *struct {
		Enabled bool `json:"enabled"`
	} `json:"allow_deletions"`
}

// toPort converts GitHub protection JSON into remote.BranchProtection.
func (r protectionResponse) toPort() remoteport.BranchProtection {
	protection := remoteport.BranchProtection{
		Protected:            true,
		RequiresPullRequest:  r.RequiredPullRequestReviews != nil,
		RequiresStatusChecks: r.RequiredStatusChecks != nil,
	}
	if r.AllowForcePushes != nil {
		protection.AllowsForcePushes = r.AllowForcePushes.Enabled
	}
	if r.AllowDeletions != nil {
		protection.AllowsDeletions = r.AllowDeletions.Enabled
	}
	return protection
}
