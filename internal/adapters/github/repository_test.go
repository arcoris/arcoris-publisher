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
	"testing"

	remoteport "arcoris.dev/arcoris-publisher/internal/ports/remote"
)

func TestRepositoryMapsResponse(t *testing.T) {
	provider := newTestProvider(func(r *http.Request) testResponse {
		if r.URL.Path != "/repos/arcoris/repo" {
			t.Fatalf("path = %q", r.URL.Path)
		}
		return jsonResponse(200, `{
			"clone_url":"https://example/repo.git",
			"ssh_url":"git@example:repo.git",
			"html_url":"https://example/repo",
			"default_branch":"main",
			"private":true,
			"archived":false,
			"disabled":false,
			"permissions":{"pull":true,"push":true,"admin":false}
		}`)
	})

	repo, err := provider.Repository(context.Background(), remoteport.RepositoryRef{Owner: "arcoris", Name: "repo"})
	if err != nil {
		t.Fatalf("Repository() error = %v", err)
	}
	if repo.CloneURL == "" || repo.SSHURL == "" || repo.WebURL == "" || repo.DefaultBranch != "main" || !repo.Private || !repo.Permissions.CanWrite {
		t.Fatalf("unexpected repository %#v", repo)
	}
}

func TestRepositoryMapsNotFound(t *testing.T) {
	provider := newTestProvider(func(r *http.Request) testResponse {
		return jsonResponse(404, `{"message":"Not Found"}`)
	})

	_, err := provider.Repository(context.Background(), remoteport.RepositoryRef{Owner: "missing", Name: "repo"})
	assertPortCode(t, err, remoteport.CodeRepositoryNotFound)
}

func TestRepoPathEscapesOwnerAndName(t *testing.T) {
	got := repoPath(remoteport.RepositoryRef{Owner: "owner space", Name: "repo/name"})
	if got != "/repos/owner%20space/repo%2Fname" {
		t.Fatalf("repoPath() = %q", got)
	}
}
