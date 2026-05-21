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
	"testing"

	remoteport "arcoris.dev/arcoris-publisher/internal/ports/remote"
)

func TestNewCreatePullRequestBody(t *testing.T) {
	body := newCreatePullRequestBody(remoteport.CreatePullRequestRequest{
		Title:      "Title",
		Body:       "Body",
		HeadBranch: "feature",
		BaseBranch: "main",
		Draft:      true,
	})

	if body.Title != "Title" || body.Body != "Body" || body.Head != "feature" || body.Base != "main" || !body.Draft {
		t.Fatalf("newCreatePullRequestBody() = %#v", body)
	}
}

func TestPullRequestResponseToPort(t *testing.T) {
	pr := pullRequestResponse{Number: 7, HTMLURL: "https://example/pull/7"}.toPort()
	if pr.Number != 7 || pr.URL == "" {
		t.Fatalf("toPort() = %#v", pr)
	}
}
