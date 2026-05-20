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

import (
	"context"
	"net/url"

	remoteport "arcoris.dev/arcoris-publisher/internal/ports/remote"
)

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

// BranchProtection maps GitHub branch protection metadata into the remote port.
//
// GitHub returns 404 both for missing resources and for unprotected branches on
// this endpoint. For publisher planning, treating that as "not protected" is the
// least surprising behavior; actual write failures are still surfaced by Git or
// later provider calls.
func (p *Provider) BranchProtection(ctx context.Context, ref remoteport.RepositoryRef, branch string) (remoteport.BranchProtection, error) {
	var out protectionResponse
	err := p.do(ctx, "GET", repoPath(ref)+"/branches/"+url.PathEscape(branch)+"/protection", nil, &out)
	if err != nil {
		if isRemoteCode(err, remoteport.CodeRepositoryNotFound) {
			return remoteport.BranchProtection{}, nil
		}
		return remoteport.BranchProtection{}, err
	}
	protection := remoteport.BranchProtection{Protected: true, RequiresPullRequest: out.RequiredPullRequestReviews != nil, RequiresStatusChecks: out.RequiredStatusChecks != nil}
	if out.AllowForcePushes != nil {
		protection.AllowsForcePushes = out.AllowForcePushes.Enabled
	}
	if out.AllowDeletions != nil {
		protection.AllowsDeletions = out.AllowDeletions.Enabled
	}
	return protection, nil
}
