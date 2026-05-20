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

// PushTag pushes one tag ref to a remote.
func (c *Client) PushTag(ctx context.Context, repoDir string, remote string, tag gitport.TagName, opts gitport.PushOptions) error {
	refspec := gitport.RefSpec(tagRef(tag) + ":" + tagRef(tag))
	return c.Push(ctx, repoDir, remote, refspec, opts)
}

// tagRef returns the fully qualified Git reference for one tag name.
func tagRef(tag gitport.TagName) string {
	return "refs/tags/" + tag.String()
}
