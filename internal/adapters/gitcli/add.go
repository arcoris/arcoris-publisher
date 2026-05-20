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

// AddAll stages every visible change in the repository.
func (c *Client) AddAll(ctx context.Context, repoDir string) error {
	result, err := c.runner.Run(ctx, c.command(repoDir, []string{"add", "-A"}, nil, true, true))
	if err != nil {
		return wrapGitCommandError("git add failed", result, err)
	}
	return nil
}
