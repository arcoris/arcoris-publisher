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
	result, err := c.runner.Run(ctx, c.command("", cloneArgs(remoteURL, dir, opts), opts.SensitiveValues, true, true))
	if err != nil {
		return wrapGitCommandError("git clone failed", result, err)
	}
	return nil
}

// cloneArgs renders clone-specific flags in deterministic order.
func cloneArgs(remoteURL string, dir string, opts gitport.CloneOptions) []string {
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
	return append(args, remoteURL, dir)
}
