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

func TestRepositoryResponseToPort(t *testing.T) {
	ref := remoteport.RepositoryRef{Owner: "arcoris", Name: "repo"}
	repo := repositoryResponse{
		CloneURL:      "https://example/repo.git",
		SSHURL:        "git@example:repo.git",
		HTMLURL:       "https://example/repo",
		DefaultBranch: "main",
		Private:       true,
		Archived:      true,
		Permissions:   permissionsResponse{Pull: true, Push: true, Admin: false},
	}.toPort(ref)

	if repo.Ref != ref ||
		repo.CloneURL == "" ||
		repo.SSHURL == "" ||
		repo.WebURL == "" ||
		repo.DefaultBranch != "main" ||
		!repo.Private ||
		!repo.Archived ||
		!repo.Permissions.CanWrite {
		t.Fatalf("toPort() = %#v", repo)
	}
}
