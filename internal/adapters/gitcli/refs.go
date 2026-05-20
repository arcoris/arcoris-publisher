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

// Head returns the object name for HEAD.
func (c *Client) Head(ctx context.Context, repoDir string) (gitport.CommitHash, error) {
	result, err := c.runner.Run(ctx, c.command(repoDir, []string{"rev-parse", "HEAD"}, nil, true, true))
	if err != nil {
		return "", wrapGitCommandError("git head lookup failed", result, err)
	}
	return gitport.CommitHash(trimOutput(result.Stdout)), nil
}

// CurrentBranch returns the current branch name and rejects detached HEAD.
func (c *Client) CurrentBranch(ctx context.Context, repoDir string) (gitport.BranchName, error) {
	result, err := c.runner.Run(ctx, c.command(repoDir, []string{"branch", "--show-current"}, nil, true, true))
	if err != nil {
		return "", wrapGitCommandError("git current branch lookup failed", result, err)
	}
	branch := trimOutput(result.Stdout)
	if branch == "" {
		return "", gitError(gitport.CodeRefNotFound, "git current branch lookup failed because HEAD is detached", nil, porterr.Details{"repo": repoDir})
	}
	return gitport.BranchName(branch), nil
}

// RefExists checks local ref existence with rev-parse --verify.
func (c *Client) RefExists(ctx context.Context, repoDir string, ref string) (bool, error) {
	spec := c.command(repoDir, []string{"rev-parse", "--verify", "--quiet", ref}, nil, true, true)
	spec.AllowedExitCodes = []int{0, 1}
	result, err := c.runner.Run(ctx, spec)
	if err != nil {
		return false, wrapGitCommandError("git ref lookup failed", result, err)
	}
	return result.ExitCode == 0, nil
}

// RemoteRefExists checks a remote ref with ls-remote.
func (c *Client) RemoteRefExists(ctx context.Context, repoDir string, remote string, ref string) (bool, error) {
	if remote == "" {
		remote = "origin"
	}
	spec := c.command(repoDir, []string{"ls-remote", "--exit-code", remote, ref}, nil, true, true)
	spec.AllowedExitCodes = []int{0, 2}
	result, err := c.runner.Run(ctx, spec)
	if err != nil {
		return false, wrapGitCommandError("git remote ref lookup failed", result, err)
	}
	return result.ExitCode == 0, nil
}

// CommitMessage returns the raw body for a commit selected by ref.
func (c *Client) CommitMessage(ctx context.Context, repoDir string, ref string) (string, error) {
	if ref == "" {
		ref = "HEAD"
	}
	result, err := c.runner.Run(ctx, c.command(repoDir, []string{"log", "-1", "--format=%B", ref}, nil, true, true))
	if err != nil {
		return "", gitError(gitport.CodeRefNotFound, "git commit message lookup failed", err, porterr.Details{"repo": repoDir, "ref": ref})
	}
	return string(result.Stdout), nil
}
