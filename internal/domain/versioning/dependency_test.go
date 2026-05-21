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
	"testing"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
	"arcoris.dev/arcoris-publisher/internal/domain/registry"
)

func TestDependencyVersions(t *testing.T) {
	registryValue := testRegistry(t)
	assignments, err := New(registryValue, AssignmentSpec{Release: "v0.3.0"})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	control := mustModule(t, registryValue, "control")
	dependencies, err := assignments.DependencyVersions(registryValue, control)
	if err != nil {
		t.Fatalf("DependencyVersions returned error: %v", err)
	}
	if len(dependencies) != 1 {
		t.Fatalf("unexpected dependency count: %d", len(dependencies))
	}
	if dependencies[0].Module() != testModuleName(t, "foundation") {
		t.Fatalf("unexpected dependency module: %q", dependencies[0].Module())
	}
	if dependencies[0].ModulePath() != "arcoris.dev/foundation" || dependencies[0].Version() != "v0.3.0" {
		t.Fatalf("unexpected dependency requirement: %q %q", dependencies[0].ModulePath(), dependencies[0].Version())
	}
	requirements, err := assignments.RequirementMap(registryValue, control)
	if err != nil {
		t.Fatalf("RequirementMap returned error: %v", err)
	}
	if requirements["arcoris.dev/foundation"] != "v0.3.0" {
		t.Fatalf("unexpected requirement map: %#v", requirements)
	}
}

func TestDependencyVersionsFailForUnassignedDependency(t *testing.T) {
	manifestValue := manifest.Must(manifest.Spec{
		Version: "v1",
		Source:  manifest.SourceSpec{Repository: "arcoris/arcoris", DefaultBranch: "main"},
		Policy:  manifest.PolicySpec{VersionPolicy: "release-train", PushPolicy: "fast-forward-only"},
		Modules: []manifest.ModuleSpec{
			{Name: "public", ModulePath: "arcoris.dev/public", SourceDir: "staging/public", Repository: "arcoris/public", Branches: []manifest.BranchMappingSpec{{Source: "main", Target: "main"}}, Dependencies: []string{"internal"}},
			{Name: "internal", ModulePath: "arcoris.dev/internal", SourceDir: "staging/internal", Repository: "arcoris/internal", Branches: []manifest.BranchMappingSpec{{Source: "main", Target: "main"}}, Visibility: "internal"},
		},
	})
	registryValue := registryFromManifest(t, manifestValue)
	assignments, err := New(registryValue, AssignmentSpec{Release: "v0.1.0"})
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	public := mustModule(t, registryValue, "public")
	if _, err := assignments.DependencyVersions(registryValue, public); err == nil {
		t.Fatal("DependencyVersions succeeded, expected unassigned dependency error")
	}
}

func registryFromManifest(t *testing.T, manifestValue manifest.Manifest) registry.Registry {
	t.Helper()
	registryValue, err := registry.New(manifestValue.Modules())
	if err != nil {
		t.Fatal(err)
	}
	return registryValue
}
