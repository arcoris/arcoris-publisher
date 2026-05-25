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

	gitport "arcoris.dev/arcoris-publisher/internal/ports/git"
)

// RemoteRefExists checks a remote ref with ls-remote.
func (c *Client) RemoteRefExists(ctx context.Context, repoDir string, remote string, ref string) (bool, error) {
	spec := c.command(repoDir, []string{"ls-remote", "--exit-code", defaultRemote(remote), ref}, nil, true, true)
	spec.AllowedExitCodes = []int{0, 2}
	result, err := c.runner.Run(ctx, spec)
	if err != nil {
		return false, wrapGitCommandError("git remote ref lookup failed", result, err)
	}
	return result.ExitCode == 0, nil
}

// RemoteRefHash returns the object hash for one remote ref.
func (c *Client) RemoteRefHash(
	ctx context.Context,
	repoDir string,
	remote string,
	ref string,
) (gitport.CommitHash, bool, error) {
	spec := c.command(repoDir, []string{"ls-remote", "--exit-code", defaultRemote(remote), ref}, nil, true, true)
	spec.AllowedExitCodes = []int{0, 2}
	result, err := c.runner.Run(ctx, spec)
	if err != nil {
		return "", false, wrapGitCommandError("git remote ref lookup failed", result, err)
	}
	if result.ExitCode == 2 {
		return "", false, nil
	}
	hash, ok := parseRemoteRefHash(string(result.Stdout), ref)
	if !ok {
		return "", false, gitError(
			gitport.CodeRemoteRefNotFound,
			"git remote ref lookup returned no matching ref",
			nil,
			nil,
		)
	}
	return hash, true, nil
}

func parseRemoteRefHash(output string, ref string) (gitport.CommitHash, bool) {
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 || fields[1] != ref {
			continue
		}
		if fields[0] == "" {
			return "", false
		}
		return gitport.CommitHash(fields[0]), true
	}
	return "", false
}

// defaultRemote applies Git's common origin default at the adapter boundary.
func defaultRemote(remote string) string {
	if remote == "" {
		return "origin"
	}
	return remote
}
