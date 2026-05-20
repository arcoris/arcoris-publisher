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
	Checkout(ctx context.Context, repoDir string, ref string, opts CheckoutOptions) error
	CreateBranch(ctx context.Context, repoDir string, branch BranchName, startPoint string, opts CreateBranchOptions) error
	ResetHard(ctx context.Context, repoDir string, ref string) error
	Clean(ctx context.Context, repoDir string, opts CleanOptions) error
	AddAll(ctx context.Context, repoDir string) error
	Commit(ctx context.Context, repoDir string, message string, opts CommitOptions) (CommitHash, error)
}
