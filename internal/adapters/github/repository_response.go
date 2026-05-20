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

// repositoryResponse mirrors the subset of GitHub repository JSON used here.
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

// permissionsResponse mirrors GitHub's repository permissions object.
type permissionsResponse struct {
	Pull  bool `json:"pull"`
	Push  bool `json:"push"`
	Admin bool `json:"admin"`
}

// toPort converts provider JSON into the stable remote.Repository contract.
func (r repositoryResponse) toPort(ref remoteport.RepositoryRef) remoteport.Repository {
	return remoteport.Repository{
		Ref:           ref,
		CloneURL:      r.CloneURL,
		SSHURL:        r.SSHURL,
		WebURL:        r.HTMLURL,
		DefaultBranch: r.DefaultBranch,
		Private:       r.Private,
		Archived:      r.Archived,
		Disabled:      r.Disabled,
		Permissions: remoteport.RepositoryPermissions{
			CanRead:  r.Permissions.Pull,
			CanWrite: r.Permissions.Push,
			CanAdmin: r.Permissions.Admin,
		},
	}
}
