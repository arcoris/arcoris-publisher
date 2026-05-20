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

// Fetch runs git fetch for one remote.
func (c *Client) Fetch(ctx context.Context, repoDir string, remote string, opts gitport.FetchOptions) error {
	result, err := c.runner.Run(ctx, c.command(repoDir, fetchArgs(remote, opts), opts.SensitiveValues, true, true))
	if err != nil {
		return wrapGitCommandError("git fetch failed", result, err)
	}
	return nil
}

// fetchArgs renders fetch flags, remote name, and optional refspecs.
func fetchArgs(remote string, opts gitport.FetchOptions) []string {
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
	args = append(args, defaultRemote(remote))
	return append(args, stringsOf(opts.RefSpecs)...)
}
