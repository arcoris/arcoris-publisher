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

package staging_test

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
)

func TestNewModuleTracksOptionalOverrides(t *testing.T) {
	mod, err := staging.NewModule(staging.ModuleSpec{
		Name:       "control",
		SourceDir:  "src/arcoris.dev/control",
		Manifest:   stringPtr("publisher.yaml"),
		Repository: "arcoris/control",
		Visibility: stringPtr("internal"),
		Branches:   []manifest.BranchMappingSpec{{Source: "main", Target: "release"}},
	})
	if err != nil {
		t.Fatalf("NewModule returned error: %v", err)
	}
	if mod.Name() != "control" || mod.SourceDir().String() != "src/arcoris.dev/control" || mod.Repository() != "arcoris/control" {
		t.Fatalf("unexpected module core fields")
	}
	if got, ok := mod.ManifestPathOverride(); !ok || got.String() != "publisher.yaml" {
		t.Fatalf("manifest override not tracked")
	}
	if got, ok := mod.VisibilityOverride(); !ok || got != manifest.VisibilityInternal {
		t.Fatalf("visibility override not tracked")
	}
	if branches, ok := mod.BranchesOverride(); !ok || len(branches) != 1 || branches[0].Target() != "release" {
		t.Fatalf("branch override not tracked")
	}
}

func TestNewModuleDistinguishesAbsentOverrides(t *testing.T) {
	mod, err := staging.NewModule(staging.ModuleSpec{Name: "control", SourceDir: "src/arcoris.dev/control", Repository: "arcoris/control"})
	if err != nil {
		t.Fatalf("NewModule returned error: %v", err)
	}
	if _, ok := mod.ManifestPathOverride(); ok {
		t.Fatalf("manifest override unexpectedly set")
	}
	if _, ok := mod.VisibilityOverride(); ok {
		t.Fatalf("visibility override unexpectedly set")
	}
	if _, ok := mod.BranchesOverride(); ok {
		t.Fatalf("branches override unexpectedly set")
	}
}

func TestModuleBranchesOverrideReturnsDetachedSlice(t *testing.T) {
	mod, err := staging.NewModule(staging.ModuleSpec{
		Name:       "control",
		SourceDir:  "src/arcoris.dev/control",
		Repository: "arcoris/control",
		Branches:   []manifest.BranchMappingSpec{{Source: "main", Target: "release"}},
	})
	if err != nil {
		t.Fatalf("NewModule returned error: %v", err)
	}
	branches, ok := mod.BranchesOverride()
	if !ok {
		t.Fatalf("expected branch override")
	}
	branches[0], _ = manifest.NewBranchMapping(manifest.BranchMappingSpec{Source: "dev", Target: "dev"})
	got, _ := mod.BranchesOverride()
	if got[0].Source() != "main" {
		t.Fatalf("BranchesOverride leaked internal slice")
	}
}

func TestNewModuleRejectsInvalidFields(t *testing.T) {
	for _, spec := range []staging.ModuleSpec{
		{Name: "Control", SourceDir: "src/arcoris.dev/control", Repository: "arcoris/control"},
		{Name: "control", SourceDir: ".", Repository: "arcoris/control"},
		{Name: "control", SourceDir: "src/arcoris.dev/control", Repository: "arcoris"},
		{Name: "control", SourceDir: "src/arcoris.dev/control", Repository: "arcoris/control", Manifest: stringPtr("../publisher.yaml")},
		{Name: "control", SourceDir: "src/arcoris.dev/control", Repository: "arcoris/control", Visibility: stringPtr("private")},
		{Name: "control", SourceDir: "src/arcoris.dev/control", Repository: "arcoris/control", Branches: []manifest.BranchMappingSpec{{Source: "bad branch", Target: "main"}}},
	} {
		if _, err := staging.NewModule(spec); err == nil {
			t.Fatalf("NewModule(%#v) returned nil error", spec)
		}
	}
}
