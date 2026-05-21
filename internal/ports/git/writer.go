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

package git

import "context"

// RepositoryWriter describes mutating local Git repository operations.
type RepositoryWriter interface {
	// Checkout switches repoDir to ref using opts.
	//
	// Implementations should make force, detach, create, and orphan behavior
	// explicit through CheckoutOptions instead of inferring it from ref syntax.
	Checkout(ctx context.Context, repoDir string, ref string, opts CheckoutOptions) error
	// CreateBranch creates branch at startPoint.
	//
	// startPoint may be a commit, tag, local branch, or remote-tracking ref as
	// accepted by the adapter's Git implementation.
	CreateBranch(
		ctx context.Context,
		repoDir string,
		branch BranchName,
		startPoint string,
		opts CreateBranchOptions,
	) error
	// ResetHard resets repoDir to ref and discards tracked working tree changes.
	//
	// Because this is destructive, adapters should surface dirty worktree or
	// invalid ref errors with stable codes.
	ResetHard(ctx context.Context, repoDir string, ref string) error
	// Clean removes untracked or ignored files according to opts.
	//
	// Implementations must not remove files unless opts explicitly selects what
	// categories of paths are eligible and opts.Force satisfies Git's safety
	// requirements.
	Clean(ctx context.Context, repoDir string, opts CleanOptions) error
	// AddAll stages all currently visible changes in repoDir.
	AddAll(ctx context.Context, repoDir string) error
	// Commit creates a commit from staged changes and returns its hash.
	//
	// When no changes are staged, adapters should return CodeNoChanges unless
	// opts.AllowEmpty requests an empty commit.
	Commit(ctx context.Context, repoDir string, message string, opts CommitOptions) (CommitHash, error)
}
