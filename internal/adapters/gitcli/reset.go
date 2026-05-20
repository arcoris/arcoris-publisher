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

// ResetHard runs git reset --hard for the selected ref.
func (c *Client) ResetHard(ctx context.Context, repoDir string, ref string) error {
	result, err := c.runner.Run(ctx, c.command(repoDir, resetHardArgs(ref), nil, true, true))
	if err != nil {
		return wrapGitCommandError("git reset --hard failed", result, err)
	}
	return nil
}

// resetHardArgs defaults an empty ref to HEAD, matching git's common usage.
func resetHardArgs(ref string) []string {
	if ref == "" {
		ref = "HEAD"
	}
	return []string{"reset", "--hard", ref}
}
