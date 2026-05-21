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

package publish

import (
	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/workflow/construct"
	"arcoris.dev/arcoris-publisher/internal/workflow/modulefile"
	"arcoris.dev/arcoris-publisher/internal/workflow/source"
	"arcoris.dev/arcoris-publisher/internal/workflow/target"
	"arcoris.dev/arcoris-publisher/internal/workflow/verify"
)

// Request describes verified targets ready for publication.
type Request struct {
	// Plan supplies publish policy, branch mappings, tags, and versions.
	Plan plan.Plan

	// Source supplies source provenance for commit messages.
	Source source.Snapshot

	// Targets contains prepared target repository worktrees.
	Targets target.WorkspaceSet

	// Construct records target tree changes.
	Construct construct.Result

	// ModuleFile records go.mod rewrite changes.
	ModuleFile modulefile.Result

	// Verify must contain no failed checks before publication.
	Verify verify.Result
}

// Options configures publication Git operations.
type Options struct {
	// RemoteName is the remote to push branches and tags to.
	RemoteName string

	// AllowEmptyCommits permits Git commits when no file content changed.
	AllowEmptyCommits bool

	// DryRun records publish intent without mutating Git repositories.
	DryRun bool
}

// DefaultOptions returns conservative publication defaults.
func DefaultOptions() Options {
	return Options{RemoteName: "origin"}
}
