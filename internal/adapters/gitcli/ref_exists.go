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

// RefExists checks local ref existence with rev-parse --verify.
func (c *Client) RefExists(ctx context.Context, repoDir string, ref string) (bool, error) {
	spec := c.command(repoDir, []string{"rev-parse", "--verify", "--quiet", ref}, nil, true, true)
	spec.AllowedExitCodes = []int{0, 1}
	result, err := c.runner.Run(ctx, spec)
	if err != nil {
		return false, wrapGitCommandError("git ref lookup failed", result, err)
	}
	return result.ExitCode == 0, nil
}
