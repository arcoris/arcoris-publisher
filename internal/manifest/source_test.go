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

package manifest_test

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestNewSourceAppliesDefaultsAndRoundTripsSpec(t *testing.T) {
	source, err := manifest.NewSource(manifest.SourceSpec{
		Repository:    "arcoris/arcoris",
		DefaultBranch: "main",
	})
	if err != nil {
		t.Fatalf("NewSource returned error: %v", err)
	}
	if source.Repository() != "arcoris/arcoris" || source.DefaultBranch() != "main" {
		t.Fatalf("unexpected source values")
	}
	if source.StagingRoot().String() != "." ||
		source.ModuleRoot().String() != "." ||
		source.DirtyPolicy() != manifest.DirtyPolicyFail {
		t.Fatalf("unexpected source defaults")
	}
	spec := source.Spec()
	if spec.Repository != "arcoris/arcoris" || spec.StagingRoot == nil || *spec.StagingRoot != "." {
		t.Fatalf("unexpected source spec: %#v", spec)
	}
}

func TestNewSourceAcceptsExplicitRootsAndDirtyPolicy(t *testing.T) {
	stagingRoot := "src"
	moduleRoot := "modules"
	dirtyPolicy := string(manifest.DirtyPolicyWarn)
	source, err := manifest.NewSource(manifest.SourceSpec{
		Repository:    "arcoris/arcoris",
		DefaultBranch: "main",
		StagingRoot:   &stagingRoot,
		ModuleRoot:    &moduleRoot,
		DirtyPolicy:   &dirtyPolicy,
	})
	if err != nil {
		t.Fatalf("NewSource returned error: %v", err)
	}
	if source.StagingRoot().String() != stagingRoot ||
		source.ModuleRoot().String() != moduleRoot ||
		source.DirtyPolicy() != manifest.DirtyPolicyWarn {
		t.Fatalf("explicit source values were not applied")
	}
}

func TestNewSourceRejectsInvalidFields(t *testing.T) {
	badDirty := "explode"
	for _, spec := range []manifest.SourceSpec{
		{Repository: "arcoris", DefaultBranch: "main"},
		{Repository: "arcoris/arcoris", DefaultBranch: "bad branch"},
		{Repository: "arcoris/arcoris", DefaultBranch: "main", StagingRoot: stringPtr("../src")},
		{Repository: "arcoris/arcoris", DefaultBranch: "main", ModuleRoot: stringPtr("/modules")},
		{Repository: "arcoris/arcoris", DefaultBranch: "main", DirtyPolicy: &badDirty},
	} {
		if _, err := manifest.NewSource(spec); err == nil {
			t.Fatalf("NewSource(%#v) returned nil error", spec)
		}
	}
}

func TestNewSourceCollectsInvalidFields(t *testing.T) {
	badDirty := "explode"
	_, err := manifest.NewSource(manifest.SourceSpec{
		Repository:    "arcoris",
		DefaultBranch: "bad branch",
		StagingRoot:   stringPtr("../src"),
		ModuleRoot:    stringPtr("/modules"),
		DirtyPolicy:   &badDirty,
	})

	requireValidationIssuePaths(
		t,
		err,
		"repository",
		"defaultBranch",
		"stagingRoot",
		"moduleRoot",
		"dirtyPolicy",
	)
}
