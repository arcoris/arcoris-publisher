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
)

func TestResolveUsesBuiltInBranchMappingWhenNoDefaultsAreSet(t *testing.T) {
	set := resolvePublicationSet(t, stagingManifest(t, baseStagingSpec()), standardModules(t)...)
	branches := set.Modules()[0].Branches()
	if len(branches) != 1 || branches[0].Source() != "main" || branches[0].Target() != "main" {
		t.Fatalf("unexpected built-in branch mapping: %#v", branches)
	}
}

func TestResolveUsesStagingDefaultBranches(t *testing.T) {
	spec := baseStagingSpec()
	spec.Defaults.Branches = []manifest.BranchMappingSpec{{Source: "main", Target: "release"}}
	set := resolvePublicationSet(t, stagingManifest(t, spec), standardModules(t)...)
	if set.Modules()[0].Branches()[0].Target() != "release" {
		t.Fatalf("staging default branch was not applied")
	}
}

func TestResolveUsesModuleLevelRoutingOverrides(t *testing.T) {
	spec := baseStagingSpec()
	spec.Modules[1].Manifest = stringPtr("publisher.yaml")
	spec.Modules[1].Visibility = stringPtr("internal")
	spec.Modules[1].Branches = []manifest.BranchMappingSpec{{Source: "main", Target: "control"}}
	set := resolvePublicationSet(t, stagingManifest(t, spec), standardModules(t)...)
	control := set.Modules()[1]
	if control.ManifestPath().String() != "publisher.yaml" || control.Visibility() != manifest.VisibilityInternal || control.Branches()[0].Target() != "control" {
		t.Fatalf("module-level overrides were not applied")
	}
}

func TestResolveMergesVerificationInPrecedenceOrder(t *testing.T) {
	spec := baseStagingSpec()
	localPolicy := string(manifest.LocalReplacePolicyWarn)
	spec.Defaults.Verification = manifest.VerificationSpec{LocalReplacePolicy: &localPolicy}
	controlSpec := moduleManifestSpec("control", "arcoris.dev/control", []string{"foundation"})
	controlSpec.Verification.Go.List = boolPtr(false)
	control, err := modulemanifest.New(controlSpec)
	if err != nil {
		t.Fatalf("module.New returned error: %v", err)
	}
	set, err := resolved.Resolve(resolved.ResolveInput{
		Staging: stagingManifest(t, spec),
		Modules: []modulemanifest.Manifest{
			moduleManifest(t, "foundation", "arcoris.dev/foundation", nil),
			control,
		},
	})
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	verification := set.Modules()[1].Verification()
	if verification.LocalReplacePolicy() != manifest.LocalReplacePolicyWarn || verification.Go().List() {
		t.Fatalf("verification precedence was not applied")
	}
}
