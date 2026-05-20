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
	"net/http"
	"strings"
	"testing"

	remoteport "arcoris.dev/arcoris-publisher/internal/ports/remote"
)

func TestCreatePullRequest(t *testing.T) {
	provider := newTestProvider(func(r *http.Request) testResponse {
		if r.Method != http.MethodPost || r.URL.Path != "/repos/arcoris/repo/pulls" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		body := readRequestBody(t, r)
		if !strings.Contains(body, `"head":"feature"`) || !strings.Contains(body, `"base":"main"`) || !strings.Contains(body, `"draft":true`) {
			t.Fatalf("unexpected body %s", body)
		}
		return jsonResponse(201, `{"number":7,"html_url":"https://example/pull/7"}`)
	})

	pr, err := provider.CreatePullRequest(context.Background(), remoteport.CreatePullRequestRequest{
		Repository: remoteport.RepositoryRef{Owner: "arcoris", Name: "repo"},
		Title:      "Title",
		Body:       "Body",
		HeadBranch: "feature",
		BaseBranch: "main",
		Draft:      true,
	})
	if err != nil || pr.Number != 7 || pr.URL == "" {
		t.Fatalf("CreatePullRequest() = %#v, %v", pr, err)
	}
}

func TestCreatePullRequestMapsValidationFailure(t *testing.T) {
	provider := newTestProvider(func(r *http.Request) testResponse {
		return jsonResponse(422, `{"message":"Validation Failed"}`)
	})

	_, err := provider.CreatePullRequest(context.Background(), remoteport.CreatePullRequestRequest{
		Repository: remoteport.RepositoryRef{Owner: "arcoris", Name: "repo"},
		Title:      "Title",
		HeadBranch: "feature",
		BaseBranch: "main",
	})
	assertPortCode(t, err, remoteport.CodePullRequestFailed)
}
