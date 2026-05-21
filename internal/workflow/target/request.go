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

package target

import "arcoris.dev/arcoris-publisher/internal/plan"

// Request describes target worktrees to prepare for a plan.
type Request struct {
	// Plan is the executable publication plan whose modules need worktrees.
	Plan plan.Plan

	// RootDir is the local directory that contains all target worktrees.
	RootDir string
}
