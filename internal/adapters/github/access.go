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

	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
	remoteport "arcoris.dev/arcoris-publisher/internal/ports/remote"
)

// CheckAccess loads repository permissions and verifies the requested level.
func (p *Provider) CheckAccess(ctx context.Context, ref remoteport.RepositoryRef, access remoteport.AccessLevel) error {
	repo, err := p.Repository(ctx, ref)
	if err != nil {
		return err
	}
	if !repo.Permissions.Allows(access) {
		return remoteError(remoteport.CodeAccessDenied, "repository access level is insufficient", nil, porterr.Details{"repository": ref.FullName(), "access": access.String()})
	}
	return nil
}

// DefaultBranch returns the default branch reported by GitHub repository metadata.
func (p *Provider) DefaultBranch(ctx context.Context, ref remoteport.RepositoryRef) (string, error) {
	repo, err := p.Repository(ctx, ref)
	if err != nil {
		return "", err
	}
	return repo.DefaultBranch, nil
}
