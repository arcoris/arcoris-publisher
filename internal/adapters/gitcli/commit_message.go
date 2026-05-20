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
	"arcoris.dev/arcoris-publisher/internal/ports/porterr"
)

// CommitMessage returns the raw body for a commit selected by ref.
func (c *Client) CommitMessage(ctx context.Context, repoDir string, ref string) (string, error) {
	ref = defaultRef(ref)
	result, err := c.runner.Run(ctx, c.command(repoDir, []string{"log", "-1", "--format=%B", ref}, nil, true, true))
	if err != nil {
		return "", gitError(gitport.CodeRefNotFound, "git commit message lookup failed", err, porterr.Details{"repo": repoDir, "ref": ref})
	}
	return string(result.Stdout), nil
}

// defaultRef applies HEAD when an operation accepts an optional ref.
func defaultRef(ref string) string {
	if ref == "" {
		return "HEAD"
	}
	return ref
}
