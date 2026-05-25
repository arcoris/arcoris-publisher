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

// TagExists checks whether a local tag ref exists.
func (c *Client) TagExists(ctx context.Context, repoDir string, tag gitport.TagName) (bool, error) {
	return c.RefExists(ctx, repoDir, tagRef(tag))
}

// DeleteTag removes one local tag if Git accepts the deletion.
func (c *Client) DeleteTag(ctx context.Context, repoDir string, tag gitport.TagName) error {
	spec := c.command(repoDir, []string{"tag", "-d", tag.String()}, nil, true, true)
	spec.AllowedExitCodes = []int{0, 1}
	result, err := c.runner.Run(ctx, spec)
	if err != nil {
		return wrapGitCommandError("git tag delete failed", result, err)
	}
	return nil
}
