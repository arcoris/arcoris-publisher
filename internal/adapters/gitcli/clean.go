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

// Clean runs git clean for the selected untracked/ignored path categories.
//
// Git distinguishes -X (ignored only) from -x (ignored and ordinary untracked);
// the port options preserve that distinction.
func (c *Client) Clean(ctx context.Context, repoDir string, opts gitport.CleanOptions) error {
	args, ok := cleanArgs(opts)
	if !ok {
		return nil
	}
	result, err := c.runner.Run(ctx, c.command(repoDir, args, nil, true, true))
	if err != nil {
		return wrapGitCommandError("git clean failed", result, err)
	}
	return nil
}

// cleanArgs returns false when options request no cleanup at all.
func cleanArgs(opts gitport.CleanOptions) ([]string, bool) {
	if !opts.RemoveUntracked && !opts.RemoveIgnored {
		return nil, false
	}
	return []string{"clean", cleanFlagString(opts)}, true
}

// cleanFlagString builds git clean's compact flag group.
func cleanFlagString(opts gitport.CleanOptions) string {
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
	return flags
}
