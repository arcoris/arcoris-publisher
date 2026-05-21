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

package graph

import (
	"reflect"
	"testing"
	"unsafe"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	modulemanifest "arcoris.dev/arcoris-publisher/internal/manifest/module"
	"arcoris.dev/arcoris-publisher/internal/manifest/resolved"
	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
	"arcoris.dev/arcoris-publisher/internal/registry"
)

type testModule struct {
	name         string
	visibility   string
	dependencies []string
}

func mustGraph(t *testing.T, modules ...testModule) Graph {
	t.Helper()
	set := mustPublicationSet(t, modules...)
	reg, err := registry.New(set)
	if err != nil {
		t.Fatalf("registry.New() error = %v", err)
	}
	g, err := New(reg)
	if err != nil {
		t.Fatalf("graph.New() error = %v", err)
	}
	return g
}

func mustPublicationSet(t *testing.T, modules ...testModule) resolved.PublicationSet {
	t.Helper()
	stagingModules := make([]staging.ModuleSpec, 0, len(modules))
	moduleManifests := make([]modulemanifest.Manifest, 0, len(modules))
	for _, mod := range modules {
		visibility := mod.visibility
		var visibilityPtr *string
		if visibility != "" {
			visibilityPtr = &visibility
		}
		stagingModules = append(stagingModules, staging.ModuleSpec{
			Name:       mod.name,
			SourceDir:  "src/arcoris.dev/" + mod.name,
			Repository: "arcoris/" + mod.name,
			Visibility: visibilityPtr,
		})
		moduleSpec := modulemanifest.Spec{
			APIVersion: string(manifest.APIVersionV1Alpha1),
			Kind:       string(manifest.KindModuleManifest),
			Metadata:   manifest.MetadataSpec{Name: mod.name},
			Module: manifest.ModuleIdentitySpec{
				Path: "arcoris.dev/" + mod.name,
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

func mustRegistry(t *testing.T, modules ...testModule) registry.Registry {
	t.Helper()

	reg, err := registry.New(mustPublicationSet(t, modules...))
	if err != nil {
		t.Fatalf("registry.New() error = %v", err)
	}

	return reg
}

func registryWithDuplicateModule(t *testing.T) registry.Registry {
	t.Helper()

	reg := mustRegistry(t,
		testModule{name: "foundation"},
		testModule{name: "control"},
	)
	modules := reg.Modules()
	modules = append(modules, modules[0])
	setRegistryModules(&reg, modules)

	return reg
}

func publicationSetWithDuplicateModule(
	t *testing.T,
) resolved.PublicationSet {
	t.Helper()

	set := mustPublicationSet(t,
		testModule{name: "foundation"},
		testModule{name: "control"},
	)
	modules := set.Modules()
	modules = append(modules, modules[0])
	setPublicationModules(&set, modules)

	return set
}

func setPublicationModules(
	set *resolved.PublicationSet,
	modules []resolved.PublicationModule,
) {
	value := reflect.ValueOf(set).Elem().FieldByName("modules")
	reflect.NewAt(
		value.Type(),
		unsafe.Pointer(value.UnsafeAddr()),
	).Elem().Set(reflect.ValueOf(modules))
}

func setRegistryModules(
	reg *registry.Registry,
	modules []resolved.PublicationModule,
) {
	value := reflect.ValueOf(reg).Elem().FieldByName("modules")
	reflect.NewAt(
		value.Type(),
		unsafe.Pointer(value.UnsafeAddr()),
	).Elem().Set(reflect.ValueOf(modules))
}

func names(values ...string) []manifest.ModuleName {
	out := make([]manifest.ModuleName, 0, len(values))
	for _, value := range values {
		out = append(out, manifest.ModuleName(value))
	}
	return out
}

func assertNames(t *testing.T, got []manifest.ModuleName, want ...string) {
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
