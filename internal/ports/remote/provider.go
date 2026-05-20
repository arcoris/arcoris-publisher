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

// Package remote defines the infrastructure port for remote hosting provider
// metadata APIs such as repository lookup, access checks, pull requests, and
// releases.
//
// This package models provider APIs, not Git transport. A Git push belongs to
// the git port; creating a pull request, checking branch protection, or reading
// effective permissions belongs here. Keeping those boundaries separate makes it
// possible to combine any Git adapter with any supported hosting provider.
//
// Implementations should normalize provider-specific responses into these
// request and result types, preserve browser-facing URLs, and translate API
// failures into structured porterr.Error values using remote error codes.
package remote

import "context"

// Provider describes remote hosting provider operations.
//
// This port is intentionally separate from the Git port. Git clone, fetch, and
// push are Git transport operations; repository metadata, branch protection,
// pull requests, and releases belong to this provider API port.
type Provider interface {
	// Repository returns normalized metadata for ref.
	//
	// Missing repositories should use CodeRepositoryNotFound. Permission problems
	// should use CodeAccessDenied or CodeAuthenticationFailed so callers can
	// distinguish nonexistent targets from inaccessible ones.
	Repository(ctx context.Context, ref RepositoryRef) (Repository, error)
	// CheckAccess verifies that the current principal has the requested access.
	//
	// Successful return means the access level is satisfied. Providers that
	// expose only coarse permissions should map them through RepositoryPermissions.
	CheckAccess(ctx context.Context, ref RepositoryRef, access AccessLevel) error
	// DefaultBranch returns the provider's default branch name for ref.
	//
	// Implementations should return the provider value as-is, without adding
	// refs/heads/ prefixes.
	DefaultBranch(ctx context.Context, ref RepositoryRef) (string, error)
	// BranchProtection returns protection metadata for branch.
	//
	// Providers without branch protection support should return a zero
	// BranchProtection and nil only when the branch is effectively unprotected.
	BranchProtection(ctx context.Context, ref RepositoryRef, branch string) (BranchProtection, error)
	// CreatePullRequest creates a pull request from req.HeadBranch to req.BaseBranch.
	//
	// The returned PullRequest should include the provider-visible number and a
	// browser URL whenever the provider exposes one.
	CreatePullRequest(ctx context.Context, req CreatePullRequestRequest) (PullRequest, error)
	// CreateRelease creates a release object for req.TagName.
	//
	// This operation describes provider release metadata only; creating or pushing
	// the Git tag belongs to the git port.
	CreateRelease(ctx context.Context, req CreateReleaseRequest) (Release, error)
}
