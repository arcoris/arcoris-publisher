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
)

func TestResolvedModuleCarriesModuleManifestIdentityAndPublishData(t *testing.T) {
	set := resolvePublicationSet(
		t,
		stagingManifest(t, baseStagingSpec()),
		standardModules(t)...,
	)
	control := set.Modules()[1]

	if control.ModuleType() != manifest.ModuleTypeGo ||
		control.ModuleRoot().String() != "." ||
		control.GoMod().String() != "go.mod" {
		t.Fatalf("unexpected module identity data")
	}

	if len(control.PublishEntries()) != 2 {
		t.Fatalf("unexpected publish entry count")
	}

	if control.Repository() != "arcoris/control" ||
		control.SourceDir().String() != "src/arcoris.dev/control" {
		t.Fatalf("unexpected staging routing data")
	}
}
