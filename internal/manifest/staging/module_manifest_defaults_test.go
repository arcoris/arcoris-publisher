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

	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
)

func TestNewModuleManifestDefaultsAppliesBuiltInDefaults(t *testing.T) {
	defaults, err := staging.NewModuleManifestDefaults(staging.ModuleManifestDefaultsSpec{})
	if err != nil {
		t.Fatalf("NewModuleManifestDefaults returned error: %v", err)
	}
	if defaults.Path().String() != "arcpub.module.yaml" || !defaults.Required() {
		t.Fatalf("unexpected module manifest defaults")
	}
}

func TestNewModuleManifestDefaultsAcceptsExplicitValues(t *testing.T) {
	defaults, err := staging.NewModuleManifestDefaults(staging.ModuleManifestDefaultsSpec{Path: stringPtr("publisher.yaml"), Required: boolPtr(false)})
	if err != nil {
		t.Fatalf("NewModuleManifestDefaults returned error: %v", err)
	}
	if defaults.Path().String() != "publisher.yaml" || defaults.Required() {
		t.Fatalf("explicit module manifest defaults were not applied")
	}
}

func TestNewModuleManifestDefaultsRejectsUnsafePath(t *testing.T) {
	if _, err := staging.NewModuleManifestDefaults(staging.ModuleManifestDefaultsSpec{Path: stringPtr("../publisher.yaml")}); err == nil {
		t.Fatalf("expected invalid module manifest path")
	}
}
