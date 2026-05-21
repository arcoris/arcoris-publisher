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

func TestNewBuildsValidatedStagingManifest(t *testing.T) {
	m, err := staging.New(validSpec())
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	if m.APIVersion() != manifest.APIVersionV1Alpha1 || m.Kind() != manifest.KindStagingManifest {
		t.Fatalf("unexpected api metadata")
	}
	if m.Metadata().Name() != "arcoris" || m.Source().Repository() != "arcoris/arcoris" {
		t.Fatalf("unexpected metadata or source")
	}
	if m.Publish().Mode() != manifest.PublishModeExplicitProjection {
		t.Fatalf("unexpected publish mode")
	}
	if m.Defaults().ModuleManifest().Path().String() != "arcpub.module.yaml" {
		t.Fatalf("unexpected module manifest default")
	}
	if len(m.Modules()) != 2 {
		t.Fatalf("expected two modules")
	}
}

func TestManifestModulesReturnsDetachedSlice(t *testing.T) {
	m, err := staging.New(validSpec())
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	mods := m.Modules()
	mods[0] = mods[1]
	if m.Modules()[0].Name() != "foundation" {
		t.Fatalf("Modules accessor leaked internal slice")
	}
}

func TestNewRejectsMissingModulesAndInvalidSections(t *testing.T) {
	for _, spec := range []staging.Spec{
		func() staging.Spec {
			spec := validSpec()
			spec.Modules = nil
			return spec
		}(),
		func() staging.Spec {
			spec := validSpec()
			spec.APIVersion = "v1"
			spec.Source.Repository = "arcoris"
			return spec
		}(),
		func() staging.Spec {
			spec := validSpec()
			spec.Modules[0].Name = "Foundation"
			return spec
		}(),
	} {
		if _, err := staging.New(spec); err == nil {
			t.Fatalf("staging.New(%#v) returned nil error", spec)
		}
	}
}
