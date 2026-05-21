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

	"arcoris.dev/arcoris-publisher/internal/graph"
	"arcoris.dev/arcoris-publisher/internal/manifest"
	modulemanifest "arcoris.dev/arcoris-publisher/internal/manifest/module"
	"arcoris.dev/arcoris-publisher/internal/manifest/resolved"
	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
	"arcoris.dev/arcoris-publisher/internal/registry"
	"arcoris.dev/arcoris-publisher/internal/versioning"
)

type testModule struct {
	// name is the resolved module name used by staging and module manifests.
	name string

	// modulePath overrides the default arcoris.dev/<name> module path.
	modulePath string

	// visibility overrides the default public visibility.
	visibility string

	// dependencies are direct internal module names declared by the module
	// manifest.
	dependencies []string

	// entries overrides the default explicit go.mod publication entry.
	entries []manifest.PublishEntrySpec

	// branches overrides default branch mappings.
	branches []manifest.BranchMappingSpec
}

// stringPtr keeps fixture specs readable when optional manifest fields are set.
func stringPtr(v string) *string { return &v }

// defaultEntries returns the minimal explicit publication content required by
// module manifests.
func defaultEntries() []manifest.PublishEntrySpec {
	return []manifest.PublishEntrySpec{{
		Type: string(manifest.PublishEntryFile),
		From: "go.mod",
		To:   "go.mod",
	}}
}

// mustPublicationSet builds realistic resolved publication sets for plan tests
// without mocking manifest internals.
func mustPublicationSet(
	t *testing.T,
	modules ...testModule,
) resolved.PublicationSet {
	t.Helper()
	stagingModules := make([]staging.ModuleSpec, 0, len(modules))
	moduleManifests := make([]modulemanifest.Manifest, 0, len(modules))
	for _, mod := range modules {
		var visibility *string
		if mod.visibility != "" {
			visibility = stringPtr(mod.visibility)
		}

		modulePath := mod.modulePath
		if modulePath == "" {
			modulePath = "arcoris.dev/" + mod.name
		}

		stagingModules = append(stagingModules, staging.ModuleSpec{
			Name:       mod.name,
			SourceDir:  "src/arcoris.dev/" + mod.name,
			Repository: "arcoris/" + mod.name,
			Visibility: visibility,
			Branches:   mod.branches,
		})
		entries := mod.entries
		if entries == nil {
			entries = defaultEntries()
		}
		moduleSpec := modulemanifest.Spec{
			APIVersion: string(manifest.APIVersionV1Alpha1),
			Kind:       string(manifest.KindModuleManifest),
			Metadata:   manifest.MetadataSpec{Name: mod.name},
			Module: manifest.ModuleIdentitySpec{
				Path: modulePath,
			},
			Dependencies: modulemanifest.DependenciesSpec{Internal: mod.dependencies},
			Publish:      modulemanifest.PublishSpec{Entries: entries},
		}
		moduleManifest, err := modulemanifest.New(moduleSpec)
		if err != nil {
			t.Fatalf("module.New(%s) error = %v", mod.name, err)
		}
		moduleManifests = append(moduleManifests, moduleManifest)
	}
	stagingManifest, err := staging.New(staging.Spec{
		APIVersion: string(manifest.APIVersionV1Alpha1),
		Kind:       string(manifest.KindStagingManifest),
		Metadata:   manifest.MetadataSpec{Name: "arcoris"},
		Source: manifest.SourceSpec{
			Repository:    "arcoris/arcoris",
			DefaultBranch: "main",
		},
		Modules: stagingModules,
	})
	if err != nil {
		t.Fatalf("staging.New() error = %v", err)
	}
	set, err := resolved.Resolve(resolved.ResolveInput{
		Staging: stagingManifest,
		Modules: moduleManifests,
	})
	if err != nil {
		t.Fatalf("resolved.Resolve() error = %v", err)
	}
	return set
}

// mustRequest builds a fully resolved request through the real manifest,
// registry, graph, and versioning constructors.
func mustRequest(t *testing.T, version string, modules ...testModule) Request {
	t.Helper()
	set := mustPublicationSet(t, modules...)
	reg, err := registry.New(set)
	if err != nil {
		t.Fatalf("registry.New() error = %v", err)
	}
	g, err := graph.New(reg)
	if err != nil {
		t.Fatalf("graph.New() error = %v", err)
	}
	parsed, err := versioning.Parse(version)
	if err != nil {
		t.Fatalf("versioning.Parse(%q) error = %v", version, err)
	}
	assignments, err := versioning.Assign(versioning.Request{
		Set:      set,
		Registry: reg,
		Graph:    g,
		Version:  parsed,
	})
	if err != nil {
		t.Fatalf("versioning.Assign() error = %v", err)
	}
	return Request{Set: set, Registry: reg, Graph: g, Assignments: assignments}
}

// mustPlan builds a complete executable plan for tests.
func mustPlan(t *testing.T, version string, modules ...testModule) Plan {
	t.Helper()
	p, err := Build(mustRequest(t, version, modules...))
	if err != nil {
		t.Fatalf("Build() error = %v", err)
	}
	return p
}

// assertModuleNames checks deterministic module-name order.
func assertModuleNames(t *testing.T, got []manifest.ModuleName, want ...string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("len = %d, want %d; got %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != manifest.ModuleName(want[i]) {
			t.Fatalf("got[%d] = %q, want %q; all = %v", i, got[i], want[i], got)
		}
	}
}

// assertBranches checks the first branch mapping in a single-branch fixture.
func assertBranches(
	t *testing.T,
	got []BranchPlan,
	source string,
	target string,
) {
	t.Helper()

	if len(got) != 1 {
		t.Fatalf("len(branches) = %d, want 1; got %#v", len(got), got)
	}
	if got[0].Source() != manifest.BranchName(source) {
		t.Fatalf("branch source = %q, want %q", got[0].Source(), source)
	}
	if got[0].Target() != manifest.BranchName(target) {
		t.Fatalf("branch target = %q, want %q", got[0].Target(), target)
	}
}

// assertPublishEntries checks the default explicit go.mod projection.
func assertPublishEntries(t *testing.T, got []manifest.PublishEntry) {
	t.Helper()

	if len(got) != 1 {
		t.Fatalf("len(entries) = %d, want 1; got %#v", len(got), got)
	}
	if got[0].Kind() != manifest.PublishEntryFile {
		t.Fatalf("entry kind = %q, want file", got[0].Kind())
	}
	if got[0].From() != "go.mod" {
		t.Fatalf("entry from = %q, want go.mod", got[0].From())
	}
}

// assertRequirement checks all public fields on a DependencyRequirement value.
func assertRequirement(
	t *testing.T,
	got DependencyRequirement,
	module string,
	modulePath string,
	version string,
) {
	t.Helper()

	if got.Module() != manifest.ModuleName(module) {
		t.Fatalf("Requirement.Module() = %q, want %q", got.Module(), module)
	}
	if got.ModulePath() != manifest.ModulePath(modulePath) {
		t.Fatalf("Requirement.ModulePath() = %q, want %q", got.ModulePath(), modulePath)
	}
	if got.Version() != versioning.Version(version) {
		t.Fatalf("Requirement.Version() = %q, want %q", got.Version(), version)
	}
}
