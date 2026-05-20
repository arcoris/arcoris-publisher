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

// RepositoryReader describes read-only Git repository operations.
type RepositoryReader interface {
	// Head returns the commit currently checked out in repoDir.
	//
	// Detached HEAD states should still return the checked-out commit hash.
	Head(ctx context.Context, repoDir string) (CommitHash, error)
	// CurrentBranch returns the selected local branch name.
	//
	// Detached HEAD states should return a structured error rather than guessing
	// a branch from refs that happen to point at the same commit.
	CurrentBranch(ctx context.Context, repoDir string) (BranchName, error)
	// Status returns a normalized working tree status for repoDir.
	//
	// Adapters should preserve enough entry detail for diagnostics while keeping
	// Status.Clean and Status.Entries internally consistent.
	Status(ctx context.Context, repoDir string) (Status, error)
	// RefExists reports whether ref resolves in the local repository.
	//
	// Missing refs should return (false, nil); invalid repositories or command
	// failures should return an error.
	RefExists(ctx context.Context, repoDir string, ref string) (bool, error)
	// RemoteRefExists reports whether ref exists on the named remote.
	//
	// Implementations may query cached remote refs or contact the remote, but
	// they should document that behavior because it affects freshness.
	RemoteRefExists(ctx context.Context, repoDir string, remote string, ref string) (bool, error)
	// CommitMessage returns the commit message for ref.
	//
	// The returned string should preserve the message body as Git reports it,
	// without adding adapter-specific prefixes or diagnostics.
	CommitMessage(ctx context.Context, repoDir string, ref string) (string, error)
}
