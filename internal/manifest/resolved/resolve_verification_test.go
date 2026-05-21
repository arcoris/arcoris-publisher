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

func TestResolveMergesVerificationInPrecedenceOrder(t *testing.T) {
	spec := baseStagingSpec()
	localPolicy := string(manifest.LocalReplacePolicyWarn)
	spec.Defaults.Verification = manifest.VerificationSpec{
		LocalReplacePolicy: &localPolicy,
	}

	controlSpec := moduleManifestSpec(
		"control",
		"arcoris.dev/control",
		[]string{"foundation"},
	)
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
	if verification.LocalReplacePolicy() != manifest.LocalReplacePolicyWarn ||
		verification.Go().List() {
		t.Fatalf("verification precedence was not applied")
	}
}
