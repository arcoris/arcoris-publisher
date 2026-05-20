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

// CreateTag creates an annotated or lightweight tag.
func (c *Client) CreateTag(ctx context.Context, repoDir string, tag gitport.TagName, target gitport.CommitHash, opts gitport.TagOptions) error {
	result, err := c.runner.Run(ctx, c.command(repoDir, createTagArgs(tag, target, opts), nil, true, true))
	if err != nil {
		return wrapGitCommandError("git tag creation failed", result, err)
	}
	return nil
}

// createTagArgs renders tag creation flags in a stable order.
func createTagArgs(tag gitport.TagName, target gitport.CommitHash, opts gitport.TagOptions) []string {
	args := []string{"tag"}
	if opts.Force {
		args = append(args, "-f")
	}
	if opts.Annotated {
		args = append(args, "-a")
	}
	args = append(args, tag.String())
	if target != "" {
		args = append(args, target.String())
	}
	if opts.Annotated && opts.Message != "" {
		args = append(args, "-m", opts.Message)
	}
	return args
}
