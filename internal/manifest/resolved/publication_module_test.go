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

func TestPublicationModuleCollectionAccessorsAreDetached(t *testing.T) {
	set := resolvePublicationSet(
		t,
		stagingManifest(t, baseStagingSpec()),
		standardModules(t)...,
	)
	control := set.Modules()[1]

	branches := control.Branches()
	branches[0], _ = manifest.NewBranchMapping(
		manifest.BranchMappingSpec{
			Source: "dev",
			Target: "dev",
		},
	)
	if control.Branches()[0].Source() == "dev" {
		t.Fatalf("Branches accessor leaked internal slice")
	}

	deps := control.Dependencies()
	deps[0] = "mutated"
	if control.Dependencies()[0] == "mutated" {
		t.Fatalf("Dependencies accessor leaked internal slice")
	}

	entries := control.PublishEntries()
	entries[0], _ = manifest.NewPublishEntry(
		manifest.PublishEntrySpec{
			Type: "file",
			From: "README.md",
			To:   "README.md",
		},
	)
	if control.PublishEntries()[0].From() == "README.md" {
		t.Fatalf("PublishEntries accessor leaked internal slice")
	}
}
