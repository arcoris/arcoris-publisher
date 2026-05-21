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

package staging_test

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
)

func TestNewDefaultsAppliesOptionalBranchesAndVerification(t *testing.T) {
	localPolicy := string(manifest.LocalReplacePolicyWarn)
	defaults, err := staging.NewDefaults(staging.DefaultsSpec{
		Branches:     []manifest.BranchMappingSpec{{Source: "main", Target: "release"}},
		Verification: manifest.VerificationSpec{LocalReplacePolicy: &localPolicy},
	})
	if err != nil {
		t.Fatalf("NewDefaults returned error: %v", err)
	}
	if !defaults.BranchesSet() || len(defaults.Branches()) != 1 || defaults.Branches()[0].Target() != "release" {
		t.Fatalf("unexpected branches default")
	}
	merged := manifest.MergeVerification(manifest.BuiltInVerificationPolicy(), defaults.Verification())
	if merged.LocalReplacePolicy() != manifest.LocalReplacePolicyWarn {
		t.Fatalf("verification override was not applied")
	}
}

func TestNewDefaultsDistinguishesAbsentAndEmptyBranches(t *testing.T) {
	defaults, err := staging.NewDefaults(staging.DefaultsSpec{})
	if err != nil {
		t.Fatalf("NewDefaults returned error: %v", err)
	}
	if defaults.BranchesSet() || len(defaults.Branches()) != 0 {
		t.Fatalf("absent branches should not be marked set")
	}
	defaults, err = staging.NewDefaults(staging.DefaultsSpec{Branches: []manifest.BranchMappingSpec{}})
	if err != nil {
		t.Fatalf("NewDefaults returned error: %v", err)
	}
	if !defaults.BranchesSet() || len(defaults.Branches()) != 0 {
		t.Fatalf("explicit empty branches should be marked set")
	}
}

func TestDefaultsBranchesReturnsDetachedSlice(t *testing.T) {
	defaults, err := staging.NewDefaults(staging.DefaultsSpec{Branches: []manifest.BranchMappingSpec{{Source: "main", Target: "release"}}})
	if err != nil {
		t.Fatalf("NewDefaults returned error: %v", err)
	}
	branches := defaults.Branches()
	branches[0], _ = manifest.NewBranchMapping(manifest.BranchMappingSpec{Source: "dev", Target: "dev"})
	if defaults.Branches()[0].Source() != "main" {
		t.Fatalf("Branches accessor leaked internal slice")
	}
}

func TestNewDefaultsRejectsInvalidNestedDefaults(t *testing.T) {
	for _, spec := range []staging.DefaultsSpec{
		{Branches: []manifest.BranchMappingSpec{{Source: "bad branch", Target: "release"}}},
		{ModuleManifest: staging.ModuleManifestDefaultsSpec{Path: stringPtr("../module.yaml")}},
		{Verification: manifest.VerificationSpec{LocalReplacePolicy: stringPtr("ignore")}},
	} {
		if _, err := staging.NewDefaults(spec); err == nil {
			t.Fatalf("NewDefaults(%#v) returned nil error", spec)
		}
	}
}
