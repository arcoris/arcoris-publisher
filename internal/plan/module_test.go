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

package plan

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestAccessorsReturnDetachedSlices(t *testing.T) {
	p := mustPlan(t, "v0.3.0",
		testModule{name: "foundation"},
		testModule{name: "control", dependencies: []string{"foundation"}},
	)
	mods := p.Modules()
	mods[0] = ModulePlan{}
	assertModuleNames(t, p.ModuleNames(), "foundation", "control")

	control, _ := p.ModuleByName("control")
	branches := control.Branches()
	branches[0] = BranchPlan{}
	branchesAgain := control.Branches()
	if branchesAgain[0].Source() == "" {
		t.Fatal("Branches() returned mutable internal slice")
	}

	entries := control.PublishEntries()
	entries[0], _ = manifest.NewPublishEntry(manifest.PublishEntrySpec{
		Type: string(manifest.PublishEntryFile),
		From: "README.md",
		To:   "README.md",
	})
	entriesAgain := control.PublishEntries()
	if entriesAgain[0].From() != "go.mod" {
		t.Fatal("PublishEntries() returned mutable internal slice")
	}

	reqs := control.Requirements()
	reqs[0] = DependencyRequirement{}
	reqsAgain := control.Requirements()
	if reqsAgain[0].Module() != "foundation" {
		t.Fatal("Requirements() returned mutable internal slice")
	}
}
