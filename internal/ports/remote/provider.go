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
package remote

import "context"

// Provider describes remote hosting provider operations.
//
// This port is intentionally separate from the Git port. Git clone, fetch, and
// push are Git transport operations; repository metadata, branch protection,
// pull requests, and releases belong to this provider API port.
type Provider interface {
	Repository(ctx context.Context, ref RepositoryRef) (Repository, error)
	CheckAccess(ctx context.Context, ref RepositoryRef, access AccessLevel) error
	DefaultBranch(ctx context.Context, ref RepositoryRef) (string, error)
	BranchProtection(ctx context.Context, ref RepositoryRef, branch string) (BranchProtection, error)
	CreatePullRequest(ctx context.Context, req CreatePullRequestRequest) (PullRequest, error)
	CreateRelease(ctx context.Context, req CreateReleaseRequest) (Release, error)
}
