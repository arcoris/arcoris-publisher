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

// Checkout runs git checkout with the caller's explicit mode flags.
func (c *Client) Checkout(ctx context.Context, repoDir string, ref string, opts gitport.CheckoutOptions) error {
	result, err := c.runner.Run(ctx, c.command(repoDir, checkoutArgs(ref, opts), nil, true, true))
	if err != nil {
		return wrapGitCommandError("git checkout failed", result, err)
	}
	return nil
}

// checkoutArgs converts typed checkout options into stable git CLI flags.
func checkoutArgs(ref string, opts gitport.CheckoutOptions) []string {
	args := []string{"checkout"}
	if opts.Force {
		args = append(args, "--force")
	}
	if opts.Detach {
		args = append(args, "--detach")
	}
	if opts.Orphan {
		args = append(args, "--orphan")
	}
	if opts.Create {
		args = append(args, "-b")
	}
	return append(args, ref)
}
