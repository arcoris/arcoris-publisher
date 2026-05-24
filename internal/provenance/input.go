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

package provenance

import (
	"arcoris.dev/arcoris-publisher/internal/buildinfo"
	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/workflow/source"
)

// Input contains the resolved runtime values used to render provenance.
//
// Paths in Plan and Source are filtered before rendering so committed
// provenance describes repository-relative intent instead of local checkout or
// worktree locations.
type Input struct {
	// Plan supplies global source and publication policy values.
	Plan plan.Plan

	// Module supplies the public module being constructed or published.
	Module plan.ModulePlan

	// Source supplies repository-wide source Git state.
	Source source.Snapshot

	// SourceModule supplies inspected explicit entries and source hashes for
	// Module.
	SourceModule source.ModuleSnapshot

	// Build supplies publisher build metadata. The zero value renders as the
	// same safe defaults as buildinfo.Current.
	Build buildinfo.Info
}

func (i Input) targetBranches() []string {
	branches := i.Module.Branches()
	out := make([]string, 0, len(branches))
	for _, branch := range branches {
		out = append(out, branch.Target().String())
	}
	return out
}
