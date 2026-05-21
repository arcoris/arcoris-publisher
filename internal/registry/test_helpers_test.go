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

package registry

import (
	"reflect"
	"testing"
	"unsafe"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	modulemanifest "arcoris.dev/arcoris-publisher/internal/manifest/module"
	"arcoris.dev/arcoris-publisher/internal/manifest/resolved"
	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
)

func testRegistry(t *testing.T) Registry {
	t.Helper()
	return Must(testPublicationSet(t))
}

func testPublicationSet(t *testing.T) resolved.PublicationSet {
	t.Helper()

	set, err := resolved.Resolve(resolved.ResolveInput{
		Staging: stagingManifest(t, stagingSpec()),
		Modules: []modulemanifest.Manifest{
			moduleManifest(t, "foundation", "arcoris.dev/foundation", nil),
			moduleManifest(t, "control", "arcoris.dev/control", []string{"foundation"}),
			moduleManifest(t, "internal-tool", "arcoris.dev/internal-tool", nil),
			moduleManifest(t, "disabled-tool", "arcoris.dev/disabled-tool", nil),
		},
	})
	if err != nil {
		t.Fatalf("resolved.Resolve returned error: %v", err)
	}

	return set
}

func stagingSpec() staging.Spec {
	return staging.Spec{
		APIVersion: string(manifest.APIVersionV1Alpha1),
		Kind:       string(manifest.KindStagingManifest),
		Metadata:   manifest.MetadataSpec{Name: "arcoris"},
		Source: manifest.SourceSpec{
			Repository:    "arcoris/arcoris",
			DefaultBranch: "main",
		},
		Modules: []staging.ModuleSpec{
			{
				Name:       "foundation",
				SourceDir:  "src/arcoris.dev/foundation",
				Repository: "arcoris/foundation",
			},
			{
				Name:       "control",
				SourceDir:  "src/arcoris.dev/control",
				Repository: "arcoris/control",
				Branches: []manifest.BranchMappingSpec{
					{Source: "release", Target: "stable"},
				},
			},
			{
				Name:       "internal-tool",
				SourceDir:  "src/arcoris.dev/internal-tool",
				Repository: "arcoris/internal-tool",
				Visibility: stringPtr("internal"),
			},
			{
				Name:       "disabled-tool",
				SourceDir:  "src/arcoris.dev/disabled-tool",
				Repository: "arcoris/disabled-tool",
				Visibility: stringPtr("disabled"),
			},
		},
	}
}

func stagingManifest(t *testing.T, spec staging.Spec) staging.Manifest {
	t.Helper()

	manifest, err := staging.New(spec)
	if err != nil {
		t.Fatalf("staging.New returned error: %v", err)
	}

	return manifest
}

func moduleManifest(
	t *testing.T,
	name string,
	path string,
	deps []string,
) modulemanifest.Manifest {
	t.Helper()

	manifest, err := modulemanifest.New(modulemanifest.Spec{
		APIVersion: string(manifest.APIVersionV1Alpha1),
		Kind:       string(manifest.KindModuleManifest),
		Metadata: manifest.MetadataSpec{
			Name: name,
		},
		Module: manifest.ModuleIdentitySpec{
			Path: path,
		},
		Dependencies: modulemanifest.DependenciesSpec{
			Internal: deps,
		},
		Publish: modulemanifest.PublishSpec{
			Entries: []manifest.PublishEntrySpec{
				{Type: "file", From: "go.mod", To: "go.mod"},
			},
		},
	})
	if err != nil {
		t.Fatalf("modulemanifest.New returned error: %v", err)
	}

	return manifest
}

func duplicatePublicationSet(
	t *testing.T,
) resolved.PublicationSet {
	t.Helper()

	set := testPublicationSet(t)
	modules := set.Modules()
	modules = append(modules, modules[0])
	return setPublicationModules(set, modules)
}

func setPublicationModules(
	set resolved.PublicationSet,
	modules []resolved.PublicationModule,
) resolved.PublicationSet {
	value := reflect.ValueOf(&set).Elem().FieldByName("modules")
	reflect.NewAt(
		value.Type(),
		unsafe.Pointer(value.UnsafeAddr()),
	).Elem().Set(reflect.ValueOf(modules))

	return set
}

func stringPtr(value string) *string {
	return &value
}
