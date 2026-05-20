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
)

// Clone runs git clone with transport-related options.
func (c *Client) Clone(ctx context.Context, remoteURL string, dir string, opts gitport.CloneOptions) error {
	args := []string{"clone"}
	if opts.NoTags {
		args = append(args, "--no-tags")
	}
	if opts.Depth > 0 {
		args = append(args, "--depth", strconv.Itoa(opts.Depth))
	}
	if opts.Bare {
		args = append(args, "--bare")
	}
	if opts.Mirror {
		args = append(args, "--mirror")
	}
	args = append(args, remoteURL, dir)
	result, err := c.runner.Run(ctx, c.command("", args, opts.SensitiveValues, true, true))
	if err != nil {
		return wrapGitCommandError("git clone failed", result, err)
	}
	return nil
}

// Fetch runs git fetch for one remote.
func (c *Client) Fetch(ctx context.Context, repoDir string, remote string, opts gitport.FetchOptions) error {
	if remote == "" {
		remote = "origin"
	}
	args := []string{"fetch"}
	if opts.Prune {
		args = append(args, "--prune")
	}
	switch opts.Tags {
	case gitport.FetchTagsAll:
		args = append(args, "--tags")
	case gitport.FetchTagsNone:
		args = append(args, "--no-tags")
	}
	args = append(args, remote)
	args = append(args, stringsOf(opts.RefSpecs)...)
	result, err := c.runner.Run(ctx, c.command(repoDir, args, opts.SensitiveValues, true, true))
	if err != nil {
		return wrapGitCommandError("git fetch failed", result, err)
	}
	return nil
}

// Push runs git push for one refspec.
func (c *Client) Push(ctx context.Context, repoDir string, remote string, refspec gitport.RefSpec, opts gitport.PushOptions) error {
	if remote == "" {
		remote = "origin"
	}
	args := []string{"push"}
	if opts.Force {
		args = append(args, "--force")
	}
	if opts.ForceWithLease {
		args = append(args, "--force-with-lease")
	}
	if opts.Atomic {
		args = append(args, "--atomic")
	}
	args = append(args, remote, refspec.String())
	result, err := c.runner.Run(ctx, c.command(repoDir, args, opts.SensitiveValues, true, true))
	if err != nil {
		return wrapGitCommandError("git push failed", result, err)
	}
	return nil
}
