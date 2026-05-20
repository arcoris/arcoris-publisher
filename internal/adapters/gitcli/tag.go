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

// TagExists checks refs/tags/<tag> in the local repository.
func (c *Client) TagExists(ctx context.Context, repoDir string, tag gitport.TagName) (bool, error) {
	return c.RefExists(ctx, repoDir, "refs/tags/"+tag.String())
}

// CreateTag creates a lightweight or annotated tag.
func (c *Client) CreateTag(ctx context.Context, repoDir string, tag gitport.TagName, target gitport.CommitHash, opts gitport.TagOptions) error {
	args := []string{"tag"}
	if opts.Force {
		args = append(args, "-f")
	}
	if opts.Annotated {
		args = append(args, "-a", tag.String())
		if target != "" {
			args = append(args, target.String())
		}
		args = append(args, "-m", opts.Message)
	} else {
		args = append(args, tag.String())
		if target != "" {
			args = append(args, target.String())
		}
	}
	result, err := c.runner.Run(ctx, c.command(repoDir, args, nil, true, true))
	if err != nil {
		return wrapGitCommandError("git tag creation failed", result, err)
	}
	return nil
}

// PushTag publishes one tag by building an explicit tag refspec.
func (c *Client) PushTag(ctx context.Context, repoDir string, remote string, tag gitport.TagName, opts gitport.PushOptions) error {
	return c.Push(ctx, repoDir, remote, gitport.RefSpec("refs/tags/"+tag.String()+":refs/tags/"+tag.String()), opts)
}
