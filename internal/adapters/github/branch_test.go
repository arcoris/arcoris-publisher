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

func TestBranchProtectionMapsResponse(t *testing.T) {
	provider := newTestProvider(func(r *http.Request) testResponse {
		if r.URL.EscapedPath() != "/repos/arcoris/repo/branches/feature%2Fx/protection" {
			t.Fatalf("escaped path = %q", r.URL.EscapedPath())
		}
		return jsonResponse(200, `{
			"required_pull_request_reviews": {},
			"required_status_checks": {},
			"allow_force_pushes": {"enabled": true},
			"allow_deletions": {"enabled": false}
		}`)
	})

	protection, err := provider.BranchProtection(
		context.Background(),
		remoteport.RepositoryRef{Owner: "arcoris", Name: "repo"},
		"feature/x",
	)
	if err != nil {
		t.Fatalf("BranchProtection() error = %v", err)
	}
	if !protection.Protected ||
		!protection.RequiresPullRequest ||
		!protection.RequiresStatusChecks ||
		!protection.AllowsForcePushes ||
		protection.AllowsDeletions {
		t.Fatalf("unexpected protection %#v", protection)
	}
}

func TestBranchProtectionNotFoundMeansUnprotected(t *testing.T) {
	provider := newTestProvider(func(r *http.Request) testResponse {
		return jsonResponse(404, `{"message":"Branch not protected"}`)
	})

	protection, err := provider.BranchProtection(context.Background(), remoteport.RepositoryRef{Owner: "arcoris", Name: "repo"}, "main")
	if err != nil {
		t.Fatalf("BranchProtection() error = %v", err)
	}
	if protection != (remoteport.BranchProtection{}) {
		t.Fatalf("BranchProtection() = %#v, want zero", protection)
	}
}
