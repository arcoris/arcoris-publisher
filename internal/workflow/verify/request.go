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

package verify

import (
	"time"

	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/workflow/construct"
	"arcoris.dev/arcoris-publisher/internal/workflow/modulefile"
	"arcoris.dev/arcoris-publisher/internal/workflow/target"
)

// Request describes constructed targets to verify before publication.
type Request struct {
	// Plan supplies effective module identity and verification policy.
	Plan plan.Plan

	// Targets contains prepared target repository worktrees.
	Targets target.WorkspaceSet

	// Construct contains prior construction results.
	Construct construct.Result

	// ModuleFile contains prior module-file rewrite results.
	ModuleFile modulefile.Result
}

// Options configures verification command behavior.
type Options struct {
	// GoBinary overrides the executable used for Go toolchain checks.
	GoBinary string

	// Timeout limits Go toolchain checks when greater than zero.
	Timeout time.Duration

	// RequireClean fails verification when checks leave target worktrees dirty.
	RequireClean bool
}

// DefaultOptions returns conservative verification defaults.
func DefaultOptions() Options {
	return Options{RequireClean: true}
}
