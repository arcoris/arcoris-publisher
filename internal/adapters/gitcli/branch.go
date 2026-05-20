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

package gitcli

import (
	"context"

	gitport "arcoris.dev/arcoris-publisher/internal/ports/git"
)

// CreateBranch runs git branch to create or update a local branch.
func (c *Client) CreateBranch(ctx context.Context, repoDir string, branch gitport.BranchName, startPoint string, opts gitport.CreateBranchOptions) error {
	result, err := c.runner.Run(ctx, c.command(repoDir, createBranchArgs(branch, startPoint, opts), nil, true, true))
	if err != nil {
		return wrapGitCommandError("git branch creation failed", result, err)
	}
	return nil
}

// createBranchArgs renders the narrow subset of git branch used by the port.
func createBranchArgs(branch gitport.BranchName, startPoint string, opts gitport.CreateBranchOptions) []string {
	args := []string{"branch"}
	if opts.Force {
		args = append(args, "-f")
	}
	args = append(args, branch.String())
	if startPoint != "" {
		args = append(args, startPoint)
	}
	return args
}
