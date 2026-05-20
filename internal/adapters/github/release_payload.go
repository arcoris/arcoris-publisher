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
	"strconv"

	remoteport "arcoris.dev/arcoris-publisher/internal/ports/remote"
)

// createReleaseBody is the JSON payload accepted by GitHub's releases API.
type createReleaseBody struct {
	TagName    string `json:"tag_name"`
	Name       string `json:"name,omitempty"`
	Body       string `json:"body,omitempty"`
	Draft      bool   `json:"draft"`
	Prerelease bool   `json:"prerelease"`
}

// releaseResponse is the subset of GitHub release response data returned by the port.
type releaseResponse struct {
	ID      int64  `json:"id"`
	HTMLURL string `json:"html_url"`
}

// newCreateReleaseBody converts the port request into GitHub JSON.
func newCreateReleaseBody(req remoteport.CreateReleaseRequest) createReleaseBody {
	return createReleaseBody{
		TagName:    req.TagName,
		Name:       req.Name,
		Body:       req.Body,
		Draft:      req.Draft,
		Prerelease: req.Prerelease,
	}
}

// toPort converts GitHub response JSON into the public release handle.
func (r releaseResponse) toPort() remoteport.Release {
	return remoteport.Release{ID: strconv.FormatInt(r.ID, 10), URL: r.HTMLURL}
}
