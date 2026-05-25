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

// Push runs git push for one refspec.
func (c *Client) Push(ctx context.Context, repoDir string, remote string, refspec gitport.RefSpec, opts gitport.PushOptions) error {
	result, err := c.runner.Run(ctx, c.command(repoDir, pushArgs(remote, refspec, opts), opts.SensitiveValues, true, true))
	if err != nil {
		return wrapGitCommandError("git push failed", result, err)
	}
	return nil
}

// DeleteRemoteRef deletes one remote ref with git push <remote> :<ref>.
func (c *Client) DeleteRemoteRef(
	ctx context.Context,
	repoDir string,
	remote string,
	ref string,
	opts gitport.PushOptions,
) error {
	result, err := c.runner.Run(ctx, c.command(repoDir, deleteRemoteRefArgs(remote, ref, opts), opts.SensitiveValues, true, true))
	if err != nil {
		return wrapGitCommandError("git remote ref delete failed", result, err)
	}
	return nil
}

// pushArgs renders force and atomic options before the remote/refspec pair.
func pushArgs(remote string, refspec gitport.RefSpec, opts gitport.PushOptions) []string {
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
	return append(args, defaultRemote(remote), refspec.String())
}

func deleteRemoteRefArgs(remote string, ref string, opts gitport.PushOptions) []string {
	return pushArgs(remote, gitport.RefSpec(":"+ref), opts)
}
