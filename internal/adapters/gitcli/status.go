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

// Status reads porcelain-v1 status and converts it into port status entries.
func (c *Client) Status(ctx context.Context, repoDir string) (gitport.Status, error) {
	result, err := c.runner.Run(ctx, c.command(repoDir, []string{"status", "--porcelain=v1", "-z"}, nil, true, true))
	if err != nil {
		return gitport.Status{}, wrapGitCommandError("git status failed", result, err)
	}
	entries := parseStatus(result.Stdout)
	return gitport.Status{Clean: len(entries) == 0, Entries: entries}, nil
}
