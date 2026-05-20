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
	Head(ctx context.Context, repoDir string) (CommitHash, error)
	CurrentBranch(ctx context.Context, repoDir string) (BranchName, error)
	Status(ctx context.Context, repoDir string) (Status, error)
	RefExists(ctx context.Context, repoDir string, ref string) (bool, error)
	RemoteRefExists(ctx context.Context, repoDir string, remote string, ref string) (bool, error)
	CommitMessage(ctx context.Context, repoDir string, ref string) (string, error)
}
