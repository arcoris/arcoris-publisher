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

// createPullRequestBody is the JSON payload accepted by GitHub's pulls API.
type createPullRequestBody struct {
	Title string `json:"title"`
	Body  string `json:"body,omitempty"`
	Head  string `json:"head"`
	Base  string `json:"base"`
	Draft bool   `json:"draft,omitempty"`
}

// pullRequestResponse is the subset of GitHub PR response data returned by the port.
type pullRequestResponse struct {
	Number  int    `json:"number"`
	HTMLURL string `json:"html_url"`
}

// newCreatePullRequestBody converts the port request into GitHub JSON.
func newCreatePullRequestBody(req remoteport.CreatePullRequestRequest) createPullRequestBody {
	return createPullRequestBody{
		Title: req.Title,
		Body:  req.Body,
		Head:  req.HeadBranch,
		Base:  req.BaseBranch,
		Draft: req.Draft,
	}
}

// toPort converts GitHub response JSON into the public pull request handle.
func (r pullRequestResponse) toPort() remoteport.PullRequest {
	return remoteport.PullRequest{Number: r.Number, URL: r.HTMLURL}
}
