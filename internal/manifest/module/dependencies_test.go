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

package module_test

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	modulemanifest "arcoris.dev/arcoris-publisher/internal/manifest/module"
)

func TestNewDependenciesAcceptsUniqueInternalModules(t *testing.T) {
	deps, err := modulemanifest.NewDependencies(modulemanifest.DependenciesSpec{Internal: []string{"foundation", "runtime"}})
	if err != nil {
		t.Fatalf("NewDependencies returned error: %v", err)
	}
	got := deps.Internal()
	if len(got) != 2 || got[0] != "foundation" || got[1] != "runtime" {
		t.Fatalf("unexpected dependencies: %#v", got)
	}
}

func TestNewDependenciesRejectsInvalidAndDuplicateNames(t *testing.T) {
	for _, spec := range []modulemanifest.DependenciesSpec{
		{Internal: []string{"Control"}},
		{Internal: []string{"foundation", "foundation"}},
	} {
		if _, err := modulemanifest.NewDependencies(spec); err == nil {
			t.Fatalf("NewDependencies(%#v) returned nil error", spec)
		}
	}
}

func TestDependenciesInternalReturnsDetachedSlice(t *testing.T) {
	deps, err := modulemanifest.NewDependencies(modulemanifest.DependenciesSpec{Internal: []string{"foundation"}})
	if err != nil {
		t.Fatalf("NewDependencies returned error: %v", err)
	}
	got := deps.Internal()
	got[0] = manifest.ModuleName("mutated")
	if deps.Internal()[0] == "mutated" {
		t.Fatalf("Internal accessor leaked internal slice")
	}
}
