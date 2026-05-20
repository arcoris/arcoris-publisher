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

func TestCheckAccessAllowsRequestedPermission(t *testing.T) {
	provider := newTestProvider(func(r *http.Request) testResponse {
		return jsonResponse(200, `{"default_branch":"main","permissions":{"pull":true,"push":true}}`)
	})

	err := provider.CheckAccess(context.Background(), remoteport.RepositoryRef{Owner: "arcoris", Name: "repo"}, remoteport.AccessWrite)
	if err != nil {
		t.Fatalf("CheckAccess() error = %v", err)
	}
}

func TestCheckAccessRejectsInsufficientPermission(t *testing.T) {
	provider := newTestProvider(func(r *http.Request) testResponse {
		return jsonResponse(200, `{"default_branch":"main","permissions":{"pull":true}}`)
	})

	err := provider.CheckAccess(context.Background(), remoteport.RepositoryRef{Owner: "arcoris", Name: "repo"}, remoteport.AccessWrite)
	assertPortCode(t, err, remoteport.CodeAccessDenied)
}

func TestDefaultBranchUsesRepositoryMetadata(t *testing.T) {
	provider := newTestProvider(func(r *http.Request) testResponse {
		return jsonResponse(200, `{"default_branch":"trunk","permissions":{"pull":true}}`)
	})

	branch, err := provider.DefaultBranch(context.Background(), remoteport.RepositoryRef{Owner: "arcoris", Name: "repo"})
	if err != nil || branch != "trunk" {
		t.Fatalf("DefaultBranch() = %q, %v; want trunk, nil", branch, err)
	}
}
