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
	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
)

// CurrentBranch returns the current branch name and rejects detached HEAD.
func (c *Client) CurrentBranch(ctx context.Context, repoDir string) (gitport.BranchName, error) {
	result, err := c.runner.Run(ctx, c.command(repoDir, []string{"branch", "--show-current"}, nil, true, true))
	if err != nil {
		return "", wrapGitCommandError("git current branch lookup failed", result, err)
	}
	branch := trimOutput(result.Stdout)
	if branch == "" {
		return "", gitError(
			gitport.CodeRefNotFound,
			"git current branch lookup failed because HEAD is detached",
			nil,
			porterr.Details{"repo": repoDir},
		)
	}
	return gitport.BranchName(branch), nil
}
