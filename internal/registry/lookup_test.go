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

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestLookupByNamePathRepositoryAndSourceDir(t *testing.T) {
	registry := testRegistry(t)

	byName, ok := registry.ModuleByName("control")
	if !ok || byName.Name() != "control" {
		t.Fatalf("ModuleByName returned %s, %v", byName.Name(), ok)
	}

	byPath, ok := registry.ModuleByPath("arcoris.dev/control")
	if !ok || byPath.Name() != "control" {
		t.Fatalf("ModuleByPath returned %s, %v", byPath.Name(), ok)
	}

	byRepository, ok := registry.ModuleByRepository("arcoris/control")
	if !ok || byRepository.Name() != "control" {
		t.Fatalf("ModuleByRepository returned %s, %v", byRepository.Name(), ok)
	}

	bySourceDir, ok := registry.ModuleBySourceDir("src/arcoris.dev/control")
	if !ok || bySourceDir.Name() != "control" {
		t.Fatalf("ModuleBySourceDir returned %s, %v", bySourceDir.Name(), ok)
	}
}

func TestContainsPredicates(t *testing.T) {
	registry := testRegistry(t)

	if !registry.ContainsName("foundation") {
		t.Fatal("expected foundation name")
	}
	if !registry.ContainsPath("arcoris.dev/foundation") {
		t.Fatal("expected foundation path")
	}
	if !registry.ContainsRepository("arcoris/foundation") {
		t.Fatal("expected foundation repository")
	}
	if !registry.ContainsSourceDir("src/arcoris.dev/foundation") {
		t.Fatal("expected foundation sourceDir")
	}
}

func TestLookupReportsMissingValues(t *testing.T) {
	registry := testRegistry(t)

	if _, ok := registry.ModuleByName(manifest.ModuleName("missing")); ok {
		t.Fatal("unexpected module by name")
	}
	if _, ok := registry.ModuleByPath(manifest.ModulePath("missing")); ok {
		t.Fatal("unexpected module by path")
	}
	if _, ok := registry.ModuleByRepository(manifest.RepositoryRef("missing/repo")); ok {
		t.Fatal("unexpected module by repository")
	}
	if _, ok := registry.ModuleBySourceDir(manifest.SourceDir("missing")); ok {
		t.Fatal("unexpected module by sourceDir")
	}
}
