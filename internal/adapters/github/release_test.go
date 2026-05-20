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

func TestCreateRelease(t *testing.T) {
	provider := newTestProvider(func(r *http.Request) testResponse {
		if r.Method != http.MethodPost || r.URL.Path != "/repos/arcoris/repo/releases" {
			t.Fatalf("request = %s %s", r.Method, r.URL.Path)
		}
		body := readRequestBody(t, r)
		if !strings.Contains(body, `"tag_name":"v1.0.0"`) || !strings.Contains(body, `"prerelease":true`) {
			t.Fatalf("unexpected body %s", body)
		}
		return jsonResponse(201, `{"id":9,"html_url":"https://example/releases/9"}`)
	})

	release, err := provider.CreateRelease(context.Background(), remoteport.CreateReleaseRequest{
		Repository: remoteport.RepositoryRef{Owner: "arcoris", Name: "repo"},
		TagName:    "v1.0.0",
		Name:       "Release",
		Body:       "Notes",
		Prerelease: true,
	})
	if err != nil || release.ID != "9" || release.URL == "" {
		t.Fatalf("CreateRelease() = %#v, %v", release, err)
	}
}

func TestCreateReleasePreservesRepositoryNotFound(t *testing.T) {
	provider := newTestProvider(func(r *http.Request) testResponse {
		return jsonResponse(404, `{"message":"Not Found"}`)
	})

	_, err := provider.CreateRelease(context.Background(), remoteport.CreateReleaseRequest{
		Repository: remoteport.RepositoryRef{Owner: "arcoris", Name: "repo"},
		TagName:    "v1.0.0",
	})
	assertPortCode(t, err, remoteport.CodeRepositoryNotFound)
}
