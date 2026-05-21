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

package resolved_test

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	modulemanifest "arcoris.dev/arcoris-publisher/internal/manifest/module"
	"arcoris.dev/arcoris-publisher/internal/manifest/resolved"
	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
)

func baseStagingSpec() staging.Spec {
	return staging.Spec{
		APIVersion: string(manifest.APIVersionV1Alpha1),
		Kind:       string(manifest.KindStagingManifest),
		Metadata:   manifest.MetadataSpec{Name: "arcoris"},
		Source:     manifest.SourceSpec{Repository: "arcoris/arcoris", DefaultBranch: "main"},
		Modules: []staging.ModuleSpec{
			{Name: "foundation", SourceDir: "src/arcoris.dev/foundation", Repository: "arcoris/foundation"},
			{Name: "control", SourceDir: "src/arcoris.dev/control", Repository: "arcoris/control"},
		},
	}
}

func stagingManifest(t *testing.T, spec staging.Spec) staging.Manifest {
	t.Helper()
	m, err := staging.New(spec)
	if err != nil {
		t.Fatalf("staging.New returned error: %v", err)
	}
	return m
}

func moduleManifestSpec(name string, path string, deps []string) modulemanifest.Spec {
	return modulemanifest.Spec{
		APIVersion:   string(manifest.APIVersionV1Alpha1),
		Kind:         string(manifest.KindModuleManifest),
		Metadata:     manifest.MetadataSpec{Name: name},
		Module:       manifest.ModuleIdentitySpec{Path: path},
		Dependencies: modulemanifest.DependenciesSpec{Internal: deps},
		Publish: modulemanifest.PublishSpec{Entries: []manifest.PublishEntrySpec{
			{Type: "file", From: "go.mod", To: "go.mod"},
			{Type: "directory", From: "contracts", To: "contracts"},
		}},
	}
}

func moduleManifest(t *testing.T, name string, path string, deps []string) modulemanifest.Manifest {
	t.Helper()
	m, err := modulemanifest.New(moduleManifestSpec(name, path, deps))
	if err != nil {
		t.Fatalf("module.New returned error: %v", err)
	}
	return m
}

func resolvePublicationSet(t *testing.T, stg staging.Manifest, modules ...modulemanifest.Manifest) resolved.PublicationSet {
	t.Helper()
	set, err := resolved.Resolve(resolved.ResolveInput{Staging: stg, Modules: modules})
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	return set
}

func standardModules(t *testing.T) []modulemanifest.Manifest {
	t.Helper()
	return []modulemanifest.Manifest{
		moduleManifest(t, "foundation", "arcoris.dev/foundation", nil),
		moduleManifest(t, "control", "arcoris.dev/control", []string{"foundation"}),
	}
}

func stringPtr(value string) *string { return &value }

func boolPtr(value bool) *bool { return &value }
