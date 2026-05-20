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
	"testing"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

func mustGraph(t *testing.T, moduleSpecs []manifest.ModuleSpec) Graph {
	t.Helper()
	manifestValue := manifest.Must(manifest.Spec{
		Version: string(manifest.VersionV1),
		Source:  sourceSpec(),
		Policy:  manifest.PolicySpec{},
		Modules: moduleSpecs,
	})
	graph, err := FromManifest(manifestValue)
	if err != nil {
		t.Fatalf("FromManifest() error = %v", err)
	}
	return graph
}

func sourceSpec() manifest.SourceSpec {
	return manifest.SourceSpec{Repository: "arcoris/arcoris", DefaultBranch: "main"}
}

func moduleSpec(module string, deps ...string) manifest.ModuleSpec {
	return moduleSpecWithVisibility(module, "", deps...)
}

func moduleSpecWithVisibility(module string, visibility string, deps ...string) manifest.ModuleSpec {
	return manifest.ModuleSpec{
		Name:         module,
		ModulePath:   "arcoris.dev/" + module,
		SourceDir:    "staging/src/arcoris.dev/" + module,
		Repository:   "arcoris/" + module,
		Branches:     []manifest.BranchMappingSpec{{Source: "main", Target: "main"}},
		Dependencies: deps,
		Visibility:   visibility,
	}
}

func name(value string) manifest.ModuleName {
	name, err := manifest.ParseModuleName(value)
	if err != nil {
		panic(err)
	}
	return name
}

func assertNames(t *testing.T, got []manifest.ModuleName, want ...string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("names len = %d, want %d: got %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != name(want[i]) {
			t.Fatalf("names[%d] = %q, want %q; full = %v", i, got[i], want[i], got)
		}
	}
}
