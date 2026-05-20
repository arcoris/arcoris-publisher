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
	"testing"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

func TestNewBuildsRegistryInDeclarationOrder(t *testing.T) {
	registry := mustRegistry(t, []manifest.ModuleSpec{
		moduleSpec("foundation"),
		moduleSpec("control", withDependency("foundation")),
		moduleSpec("scheduler", withDependency("control")),
	})

	if registry.Len() != 3 {
		t.Fatalf("Len() = %d, want 3", registry.Len())
	}
	if registry.Empty() {
		t.Fatalf("Empty() = true, want false")
	}
	assertNames(t, registry.ModuleNames(), "foundation", "control", "scheduler")
	assertModules(t, registry.Modules(), "foundation", "control", "scheduler")
}

func TestNewAllowsEmptyRegistry(t *testing.T) {
	registry, err := New(nil)
	if err != nil {
		t.Fatalf("New(nil) error = %v", err)
	}
	if !registry.Empty() || registry.Len() != 0 {
		t.Fatalf("empty registry Empty/Len = %v/%d", registry.Empty(), registry.Len())
	}
	if len(registry.Modules()) != 0 || len(registry.ModuleNames()) != 0 {
		t.Fatalf("empty registry returned modules")
	}
}

func TestFromManifestBuildsRegistry(t *testing.T) {
	manifestValue := manifest.Must(manifest.Spec{
		Version: "v1",
		Source:  manifest.SourceSpec{Repository: "arcoris/arcoris", DefaultBranch: "main"},
		Policy:  manifest.PolicySpec{},
		Modules: []manifest.ModuleSpec{
			moduleSpec("foundation"),
			moduleSpec("control", withDependency("foundation")),
		},
	})

	registry, err := FromManifest(manifestValue)
	if err != nil {
		t.Fatalf("FromManifest() error = %v", err)
	}
	assertNames(t, registry.ModuleNames(), "foundation", "control")
}

func TestMustPanicsForInvalidModules(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("Must() did not panic")
		}
	}()
	module := mustModule(t, moduleSpec("foundation"))
	Must([]manifest.Module{module, module})
}

func TestMustReturnsRegistryForValidModules(t *testing.T) {
	module := mustModule(t, moduleSpec("foundation"))

	registry := Must([]manifest.Module{module})

	assertNames(t, registry.ModuleNames(), "foundation")
}
