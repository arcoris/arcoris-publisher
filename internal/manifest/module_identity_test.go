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

func TestNewModuleIdentityAppliesDefaults(t *testing.T) {
	identity, err := manifest.NewModuleIdentity(manifest.ModuleIdentitySpec{
		Path: "arcoris.dev/control",
	})
	if err != nil {
		t.Fatalf("NewModuleIdentity returned error: %v", err)
	}
	if identity.Type() != manifest.ModuleTypeGo ||
		identity.Path() != "arcoris.dev/control" ||
		identity.Root().String() != "." ||
		identity.GoMod().String() != "go.mod" {
		t.Fatalf("unexpected identity defaults")
	}
}

func TestNewModuleIdentityAcceptsOverrides(t *testing.T) {
	moduleType := string(manifest.ModuleTypeGo)
	root := "packages/control"
	goMod := "module/go.mod"
	identity, err := manifest.NewModuleIdentity(manifest.ModuleIdentitySpec{
		Type:  &moduleType,
		Path:  "arcoris.dev/control",
		Root:  &root,
		GoMod: &goMod,
	})
	if err != nil {
		t.Fatalf("NewModuleIdentity returned error: %v", err)
	}
	if identity.Root().String() != root || identity.GoMod().String() != goMod {
		t.Fatalf("unexpected identity overrides")
	}
}

func TestNewModuleIdentityRejectsInvalidFields(t *testing.T) {
	badType := "rust"
	for _, spec := range []manifest.ModuleIdentitySpec{
		{Type: &badType, Path: "arcoris.dev/control"},
		{Path: "control"},
		{Path: "arcoris.dev/control", Root: stringPtr("../control")},
		{Path: "arcoris.dev/control", GoMod: stringPtr(".")},
	} {
		if _, err := manifest.NewModuleIdentity(spec); err == nil {
			t.Fatalf("NewModuleIdentity(%#v) returned nil error", spec)
		}
	}
}

func TestNewModuleIdentityCollectsInvalidFields(t *testing.T) {
	badType := "rust"
	_, err := manifest.NewModuleIdentity(manifest.ModuleIdentitySpec{
		Type:  &badType,
		Path:  "control",
		Root:  stringPtr("../control"),
		GoMod: stringPtr("."),
	})

	requireValidationIssuePaths(t, err, "type", "path", "root", "goMod")
}
