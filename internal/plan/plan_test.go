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
	"arcoris.dev/arcoris-publisher/internal/versioning"
)

func TestBuildCreatesPublicationOrderedPlan(t *testing.T) {
	p := mustPlan(t, "v0.3.0",
		testModule{name: "foundation"},
		testModule{name: "control", dependencies: []string{"foundation"}},
		testModule{name: "scheduler", dependencies: []string{"control"}},
	)

	if p.Empty() {
		t.Fatal("plan is empty")
	}
	if p.Len() != 3 {
		t.Fatalf("Len() = %d, want 3", p.Len())
	}
	assertModuleNames(t, p.ModuleNames(), "foundation", "control", "scheduler")
	if p.PublishPolicy().Mode() != manifest.PublishModeExplicitProjection {
		t.Fatalf("publish mode = %q", p.PublishPolicy().Mode())
	}
}

func TestBuildCopiesModuleExecutionData(t *testing.T) {
	p := mustPlan(t, "v0.3.0",
		testModule{name: "foundation"},
		testModule{
			name:         "control",
			dependencies: []string{"foundation"},
			branches: []manifest.BranchMappingSpec{{
				Source: "main",
				Target: "stable",
			}},
		},
	)
	mod, ok := p.ModuleByName("control")
	if !ok {
		t.Fatal("control not found")
	}
	if mod.ModulePath() != "arcoris.dev/control" {
		t.Fatalf("ModulePath() = %q", mod.ModulePath())
	}
	if mod.ModuleType() != manifest.ModuleTypeGo {
		t.Fatalf("ModuleType() = %q", mod.ModuleType())
	}
	if mod.SourceDir() != "src/arcoris.dev/control" {
		t.Fatalf("SourceDir() = %q", mod.SourceDir())
	}
	if mod.ModuleRoot() != "." {
		t.Fatalf("ModuleRoot() = %q", mod.ModuleRoot())
	}
	if mod.GoMod() != "go.mod" {
		t.Fatalf("GoMod() = %q", mod.GoMod())
	}
	if mod.Repository() != "arcoris/control" {
		t.Fatalf("Repository() = %q", mod.Repository())
	}
	if mod.Version() != versioning.Version("v0.3.0") {
		t.Fatalf("Version() = %q", mod.Version())
	}
	assertBranches(t, mod.Branches(), "main", "stable")
	assertPublishEntries(t, mod.PublishEntries())
	if mod.Verification().Go().WorkspaceMode() != manifest.GoWorkspaceModeOff {
		t.Fatalf("workspace mode = %q", mod.Verification().Go().WorkspaceMode())
	}
}

func TestBuildCarriesDirectDependencyRequirements(t *testing.T) {
	p := mustPlan(t, "v0.3.0",
		testModule{name: "foundation"},
		testModule{name: "control", dependencies: []string{"foundation"}},
	)
	control, ok := p.ModuleByName("control")
	if !ok {
		t.Fatal("control not found")
	}
	reqs := control.Requirements()
	if len(reqs) != 1 {
		t.Fatalf("len(requirements) = %d, want 1", len(reqs))
	}
	assertRequirement(t, reqs[0], "foundation", "arcoris.dev/foundation", "v0.3.0")
}
