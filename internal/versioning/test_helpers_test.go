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

package versioning

import (
	"errors"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/graph"
	"arcoris.dev/arcoris-publisher/internal/manifest"
	modulemanifest "arcoris.dev/arcoris-publisher/internal/manifest/module"
	"arcoris.dev/arcoris-publisher/internal/manifest/resolved"
	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
	"arcoris.dev/arcoris-publisher/internal/registry"
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
}

// stringPtr keeps fixture specs readable when optional manifest fields are set.
func stringPtr(v string) *string { return &v }

// mustRequest builds a fully resolved request through the real manifest,
// registry, and graph constructors.
func mustRequest(t *testing.T, version string, policy string, modules ...testModule) Request {
	t.Helper()
	set := mustPublicationSet(t, policy, modules...)
	reg, err := registry.New(set)
	if err != nil {
		t.Fatalf("registry.New() error = %v", err)
	}
	g, err := graph.New(reg)
	if err != nil {
		t.Fatalf("graph.New() error = %v", err)
	}
	parsed, err := Parse(version)
	if err != nil {
		t.Fatalf("Parse(%q) error = %v", version, err)
	}
	return Request{Set: set, Registry: reg, Graph: g, Version: parsed}
}

// mustPublicationSet builds realistic resolved publication sets for versioning
// tests without mocking manifest internals.
func mustPublicationSet(
	t *testing.T,
	policy string,
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
		})
		moduleSpec := modulemanifest.Spec{
			APIVersion: string(manifest.APIVersionV1Alpha1),
			Kind:       string(manifest.KindModuleManifest),
			Metadata:   manifest.MetadataSpec{Name: mod.name},
			Module: manifest.ModuleIdentitySpec{
				Path: modulePath,
			},
			Dependencies: modulemanifest.DependenciesSpec{Internal: mod.dependencies},
			Publish: modulemanifest.PublishSpec{Entries: []manifest.PublishEntrySpec{{
				Type: string(manifest.PublishEntryFile),
				From: "go.mod",
				To:   "go.mod",
			}}},
		}
		moduleManifest, err := modulemanifest.New(moduleSpec)
		if err != nil {
			t.Fatalf("module.New(%s) error = %v", mod.name, err)
		}
		moduleManifests = append(moduleManifests, moduleManifest)
	}
	var versionPolicy *string
	if policy != "" {
		versionPolicy = stringPtr(policy)
	}
	stagingManifest, err := staging.New(staging.Spec{
		APIVersion: string(manifest.APIVersionV1Alpha1),
		Kind:       string(manifest.KindStagingManifest),
		Metadata:   manifest.MetadataSpec{Name: "arcoris"},
		Source: manifest.SourceSpec{
			Repository:    "arcoris/arcoris",
			DefaultBranch: "main",
		},
		Publish: manifest.PublishSpec{VersionPolicy: versionPolicy},
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

// assertModuleVersion checks all public fields on a ModuleVersion value.
func assertModuleVersion(
	t *testing.T,
	got ModuleVersion,
	module string,
	modulePath string,
	version string,
) {
	t.Helper()

	if got.Module() != manifest.ModuleName(module) {
		t.Fatalf("ModuleVersion.Module() = %q, want %q", got.Module(), module)
	}
	if got.ModulePath() != manifest.ModulePath(modulePath) {
		t.Fatalf("ModuleVersion.ModulePath() = %q, want %q", got.ModulePath(), modulePath)
	}
	if got.Version() != Version(version) {
		t.Fatalf("ModuleVersion.Version() = %q, want %q", got.Version(), version)
	}
}

// assertRequirement checks all public fields on a Requirement value.
func assertRequirement(
	t *testing.T,
	got Requirement,
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
	if got.Version() != Version(version) {
		t.Fatalf("Requirement.Version() = %q, want %q", got.Version(), version)
	}
}

// assertValidationHas verifies that err is a ValidationError containing code.
func assertValidationHas(t *testing.T, err error, code IssueCode) {
	t.Helper()

	var validation *ValidationError
	if !errors.As(err, &validation) || !validation.Has(code) {
		t.Fatalf("error = %v, want validation issue %q", err, code)
	}
}
