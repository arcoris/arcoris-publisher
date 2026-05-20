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

	remoteport "arcoris.dev/arcoris-publisher/internal/ports/remote"
)

type createPullRequestBody struct {
	Title string `json:"title"`
	Body  string `json:"body,omitempty"`
	Head  string `json:"head"`
	Base  string `json:"base"`
	Draft bool   `json:"draft,omitempty"`
}

type pullRequestResponse struct {
	Number  int    `json:"number"`
	HTMLURL string `json:"html_url"`
}

// CreatePullRequest creates a GitHub pull request and returns its public handle.
func (p *Provider) CreatePullRequest(ctx context.Context, req remoteport.CreatePullRequestRequest) (remoteport.PullRequest, error) {
	body := createPullRequestBody{Title: req.Title, Body: req.Body, Head: req.HeadBranch, Base: req.BaseBranch, Draft: req.Draft}
	var out pullRequestResponse
	if err := p.do(ctx, "POST", repoPath(req.Repository)+"/pulls", body, &out); err != nil {
		return remoteport.PullRequest{}, wrapRemoteOperationError(remoteport.CodePullRequestFailed, "github pull request creation failed", err, nil)
	}
	return remoteport.PullRequest{Number: out.Number, URL: out.HTMLURL}, nil
}
