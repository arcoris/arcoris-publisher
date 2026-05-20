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
	"net/url"

	remoteport "arcoris.dev/arcoris-publisher/internal/ports/remote"
)

type repositoryResponse struct {
	Name          string              `json:"name"`
	FullName      string              `json:"full_name"`
	CloneURL      string              `json:"clone_url"`
	SSHURL        string              `json:"ssh_url"`
	HTMLURL       string              `json:"html_url"`
	DefaultBranch string              `json:"default_branch"`
	Private       bool                `json:"private"`
	Archived      bool                `json:"archived"`
	Disabled      bool                `json:"disabled"`
	Permissions   permissionsResponse `json:"permissions"`
}

type permissionsResponse struct {
	Pull  bool `json:"pull"`
	Push  bool `json:"push"`
	Admin bool `json:"admin"`
}

// Repository loads and normalizes GitHub repository metadata.
func (p *Provider) Repository(ctx context.Context, ref remoteport.RepositoryRef) (remoteport.Repository, error) {
	var out repositoryResponse
	if err := p.do(ctx, "GET", repoPath(ref), nil, &out); err != nil {
		return remoteport.Repository{}, err
	}
	return remoteport.Repository{Ref: ref, CloneURL: out.CloneURL, SSHURL: out.SSHURL, WebURL: out.HTMLURL, DefaultBranch: out.DefaultBranch, Private: out.Private, Archived: out.Archived, Disabled: out.Disabled, Permissions: remoteport.RepositoryPermissions{CanRead: out.Permissions.Pull, CanWrite: out.Permissions.Push, CanAdmin: out.Permissions.Admin}}, nil
}

// repoPath builds the escaped GitHub repository API path.
func repoPath(ref remoteport.RepositoryRef) string {
	return "/repos/" + url.PathEscape(ref.Owner) + "/" + url.PathEscape(ref.Name)
}
