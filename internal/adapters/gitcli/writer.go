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
	"strconv"

	gitport "arcoris.dev/arcoris-publisher/internal/ports/git"
	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
)

// Checkout runs git checkout with the caller's explicit mode flags.
func (c *Client) Checkout(ctx context.Context, repoDir string, ref string, opts gitport.CheckoutOptions) error {
	args := []string{"checkout"}
	if opts.Force {
		args = append(args, "--force")
	}
	if opts.Detach {
		args = append(args, "--detach")
	}
	if opts.Orphan {
		args = append(args, "--orphan")
	}
	if opts.Create {
		args = append(args, "-b")
	}
	args = append(args, ref)
	result, err := c.runner.Run(ctx, c.command(repoDir, args, nil, true, true))
	if err != nil {
		return wrapGitCommandError("git checkout failed", result, err)
	}
	return nil
}

// CreateBranch runs git branch to create or update a local branch.
func (c *Client) CreateBranch(ctx context.Context, repoDir string, branch gitport.BranchName, startPoint string, opts gitport.CreateBranchOptions) error {
	args := []string{"branch"}
	if opts.Force {
		args = append(args, "-f")
	}
	args = append(args, branch.String())
	if startPoint != "" {
		args = append(args, startPoint)
	}
	result, err := c.runner.Run(ctx, c.command(repoDir, args, nil, true, true))
	if err != nil {
		return wrapGitCommandError("git branch creation failed", result, err)
	}
	return nil
}

// ResetHard runs git reset --hard for the selected ref.
func (c *Client) ResetHard(ctx context.Context, repoDir string, ref string) error {
	if ref == "" {
		ref = "HEAD"
	}
	result, err := c.runner.Run(ctx, c.command(repoDir, []string{"reset", "--hard", ref}, nil, true, true))
	if err != nil {
		return wrapGitCommandError("git reset --hard failed", result, err)
	}
	return nil
}

// Clean runs git clean for the selected untracked/ignored path categories.
//
// Git distinguishes -X (ignored only) from -x (ignored and ordinary untracked);
// the port options preserve that distinction.
func (c *Client) Clean(ctx context.Context, repoDir string, opts gitport.CleanOptions) error {
	if !opts.RemoveUntracked && !opts.RemoveIgnored {
		return nil
	}
	flags := "-f"
	if opts.Force {
		flags = "-ff"
	}
	if opts.Directories {
		flags += "d"
	}
	if opts.RemoveIgnored {
		if opts.RemoveUntracked {
			flags += "x"
		} else {
			flags += "X"
		}
	}
	result, err := c.runner.Run(ctx, c.command(repoDir, []string{"clean", flags}, nil, true, true))
	if err != nil {
		return wrapGitCommandError("git clean failed", result, err)
	}
	return nil
}

// AddAll stages every visible change in the repository.
func (c *Client) AddAll(ctx context.Context, repoDir string) error {
	result, err := c.runner.Run(ctx, c.command(repoDir, []string{"add", "-A"}, nil, true, true))
	if err != nil {
		return wrapGitCommandError("git add failed", result, err)
	}
	return nil
}

// Commit creates a commit from already-staged changes.
//
// The adapter intentionally does not call AddAll here. The port exposes AddAll
// separately so workflow code controls exactly which changes enter a commit.
func (c *Client) Commit(ctx context.Context, repoDir string, message string, opts gitport.CommitOptions) (gitport.CommitHash, error) {
	diff := c.command(repoDir, []string{"diff", "--cached", "--quiet"}, nil, false, true)
	diff.AllowedExitCodes = []int{0, 1}
	diffResult, err := c.runner.Run(ctx, diff)
	if err != nil {
		return "", wrapGitCommandError("git staged diff failed", diffResult, err)
	}
	if diffResult.ExitCode == 0 && !opts.AllowEmpty {
		return "", gitError(gitport.CodeNoChanges, "git commit skipped because there are no staged changes", nil, porterr.Details{"repo": repoDir})
	}
	args := []string{"commit", "-m", message}
	if opts.AllowEmpty {
		args = append(args, "--allow-empty")
	}
	env := c.envForCommit(opts)
	spec := c.command(repoDir, args, nil, true, true)
	spec.Env = append(spec.Env, env...)
	result, err := c.runner.Run(ctx, spec)
	if err != nil {
		return "", wrapGitCommandError("git commit failed", result, err)
	}
	return c.Head(ctx, repoDir)
}

// envForCommit converts commit identity options into Git environment variables.
func (c *Client) envForCommit(opts gitport.CommitOptions) []string {
	env := make([]string, 0, 6)
	if opts.AuthorName != "" {
		env = append(env, "GIT_AUTHOR_NAME="+opts.AuthorName)
	}
	if opts.AuthorEmail != "" {
		env = append(env, "GIT_AUTHOR_EMAIL="+opts.AuthorEmail)
	}
	if opts.CommitterName != "" {
		env = append(env, "GIT_COMMITTER_NAME="+opts.CommitterName)
	}
	if opts.CommitterEmail != "" {
		env = append(env, "GIT_COMMITTER_EMAIL="+opts.CommitterEmail)
	}
	if !opts.AuthorDate.IsZero() {
		env = append(env, "GIT_AUTHOR_DATE="+strconv.FormatInt(opts.AuthorDate.Unix(), 10)+" +0000")
	}
	if !opts.CommitterDate.IsZero() {
		env = append(env, "GIT_COMMITTER_DATE="+strconv.FormatInt(opts.CommitterDate.Unix(), 10)+" +0000")
	}
	return env
}
