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

func TestPublicationSetAccessorsReturnResolvedTopLevelValues(t *testing.T) {
	set := resolvePublicationSet(t, stagingManifest(t, baseStagingSpec()), standardModules(t)...)
	if set.Metadata().Name() != "arcoris" || set.Source().Repository() != "arcoris/arcoris" || set.Publish().Mode() == "" {
		t.Fatalf("unexpected top-level publication set values")
	}
}

func TestPublicationSetModulesReturnsDetachedSlice(t *testing.T) {
	set := resolvePublicationSet(t, stagingManifest(t, baseStagingSpec()), standardModules(t)...)
	mods := set.Modules()
	mods[0] = mods[1]
	if set.Modules()[0].Name() != "foundation" {
		t.Fatalf("Modules accessor leaked internal slice")
	}
}
