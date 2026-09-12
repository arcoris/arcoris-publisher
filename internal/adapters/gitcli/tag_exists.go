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

// TagExists checks whether a local tag ref exists.
func (c *Client) TagExists(ctx context.Context, repoDir string, tag gitport.TagName) (bool, error) {
	if err := validateGitPositional("tag", tag.String()); err != nil {
		return false, err
	}
	return c.RefExists(ctx, repoDir, tagRef(tag))
}

// TagTargetHash resolves both lightweight and annotated tags to their target
// commit. The ^{} suffix asks Git to peel annotated tag objects before the
// caller compares ownership with a transaction-created commit.
func (c *Client) TagTargetHash(ctx context.Context, repoDir string, tag gitport.TagName) (gitport.CommitHash, bool, error) {
	if err := validateGitPositional("tag", tag.String()); err != nil {
		return "", false, err
	}
	ref := tagRef(tag) + "^{}"
	spec := c.command(repoDir, []string{"rev-parse", "--verify", "--quiet", ref}, nil, true, true)
	spec.AllowedExitCodes = []int{0, 1}
	result, err := c.runner.Run(ctx, spec)
	if err != nil {
		return "", false, wrapGitCommandError("git tag target lookup failed", result, err)
	}
	if result.ExitCode == 1 {
		return "", false, nil
	}
	hash := strings.TrimSpace(string(result.Stdout))
	if hash == "" {
		return "", false, gitError(
			gitport.CodeRefNotFound,
			"git tag target lookup returned an empty object id",
			nil,
			nil,
		)
	}
	return gitport.CommitHash(hash), true, nil
}

// DeleteTag removes one local tag if Git accepts the deletion.
func (c *Client) DeleteTag(ctx context.Context, repoDir string, tag gitport.TagName) error {
	if err := validateGitPositional("tag", tag.String()); err != nil {
		return err
	}
	spec := c.command(repoDir, []string{"tag", "-d", tag.String()}, nil, true, true)
	spec.AllowedExitCodes = []int{0, 1}
	result, err := c.runner.Run(ctx, spec)
	if err != nil {
		return wrapGitCommandError("git tag delete failed", result, err)
	}
	return nil
}
