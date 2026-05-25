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

import "context"

// ConfigGet reads an effective Git configuration value for repoDir.
func (c *Client) ConfigGet(ctx context.Context, repoDir string, key string) (string, bool, error) {
	spec := c.command(repoDir, []string{"config", "--get", key}, nil, true, true)
	spec.AllowedExitCodes = []int{0, 1}
	result, err := c.runner.Run(ctx, spec)
	if err != nil {
		return "", false, wrapGitCommandError("git config get failed", result, err)
	}
	if result.ExitCode == 1 {
		return "", false, nil
	}
	value := trimOutput(result.Stdout)
	return value, value != "", nil
}
