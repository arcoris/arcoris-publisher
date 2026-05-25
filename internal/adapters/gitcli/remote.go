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
	"strings"
)

// RemoteURL reads the configured URL for a named remote.
func (c *Client) RemoteURL(ctx context.Context, repoDir string, remote string) (string, bool, error) {
	spec := c.command(repoDir, []string{"remote", "get-url", defaultRemote(remote)}, nil, true, true)
	spec.AllowedExitCodes = []int{0, 2}
	result, err := c.runner.Run(ctx, spec)
	if err != nil {
		return "", false, wrapGitCommandError("git remote get-url failed", result, err)
	}
	if result.ExitCode == 2 {
		return "", false, nil
	}
	return strings.TrimSpace(string(result.Stdout)), true, nil
}

// AddRemote adds a named remote with url.
func (c *Client) AddRemote(ctx context.Context, repoDir string, remote string, url string) error {
	result, err := c.runner.Run(ctx, c.command(repoDir, []string{"remote", "add", defaultRemote(remote), url}, []string{url}, true, true))
	if err != nil {
		return wrapGitCommandError("git remote add failed", result, err)
	}
	return nil
}
