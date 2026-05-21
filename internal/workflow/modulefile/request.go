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

package modulefile

import (
	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/workflow/construct"
	"arcoris.dev/arcoris-publisher/internal/workflow/target"
)

// Request describes constructed targets whose module files should be rewritten.
type Request struct {
	// Plan supplies effective module paths and dependency requirements.
	Plan plan.Plan

	// Targets contains prepared target repository worktrees.
	Targets target.WorkspaceSet

	// Construct contains prior construction results.
	Construct construct.Result
}

// Options configures go.mod rewriting behavior.
type Options struct {
	// RemoveLocalReplaces removes local replace directives for managed internal deps.
	RemoveLocalReplaces bool
}

// DefaultOptions returns conservative module-file rewrite defaults.
func DefaultOptions() Options {
	return Options{RemoveLocalReplaces: true}
}
