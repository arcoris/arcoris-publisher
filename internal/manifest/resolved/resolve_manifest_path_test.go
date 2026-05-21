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

import "testing"

func TestResolveUsesModuleManifestPathOverride(t *testing.T) {
	spec := baseStagingSpec()
	spec.Modules[1].Manifest = stringPtr("publisher.yaml")

	set := resolvePublicationSet(
		t,
		stagingManifest(t, spec),
		standardModules(t)...,
	)

	if set.Modules()[1].ManifestPath().String() != "publisher.yaml" {
		t.Fatalf("module manifest path override was not applied")
	}
}
