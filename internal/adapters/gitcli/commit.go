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

// Commit creates a commit from already-staged changes.
//
// The adapter intentionally does not call AddAll here. The port exposes AddAll
// separately so workflow code controls exactly which changes enter a commit.
func (c *Client) Commit(ctx context.Context, repoDir string, message string, opts gitport.CommitOptions) (gitport.CommitHash, error) {
	if err := c.ensureStagedChanges(ctx, repoDir, opts); err != nil {
		return "", err
	}
	if err := c.runCommit(ctx, repoDir, message, opts); err != nil {
		return "", err
	}
	return c.Head(ctx, repoDir)
}

// ensureStagedChanges prevents accidental empty commits unless explicitly allowed.
func (c *Client) ensureStagedChanges(ctx context.Context, repoDir string, opts gitport.CommitOptions) error {
	diff := c.command(repoDir, []string{"diff", "--cached", "--quiet"}, nil, false, true)
	diff.AllowedExitCodes = []int{0, 1}
	diffResult, err := c.runner.Run(ctx, diff)
	if err != nil {
		return wrapGitCommandError("git staged diff failed", diffResult, err)
	}
	if diffResult.ExitCode == 0 && !opts.AllowEmpty {
		return gitError(gitport.CodeNoChanges, "git commit skipped because there are no staged changes", nil, porterr.Details{"repo": repoDir})
	}
	return nil
}

// runCommit executes git commit with message, empty-commit flag, and identity env.
func (c *Client) runCommit(ctx context.Context, repoDir string, message string, opts gitport.CommitOptions) error {
	spec := c.command(repoDir, commitArgs(message, opts), nil, true, true)
	spec.Env = append(spec.Env, commitEnv(opts)...)
	result, err := c.runner.Run(ctx, spec)
	if err != nil {
		return wrapGitCommandError("git commit failed", result, err)
	}
	return nil
}

// commitArgs renders the git commit invocation itself.
func commitArgs(message string, opts gitport.CommitOptions) []string {
	args := []string{"commit", "-m", message}
	if opts.AllowEmpty {
		args = append(args, "--allow-empty")
	}
	return args
}
