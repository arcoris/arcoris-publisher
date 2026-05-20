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

package registry

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

func TestModulesByVisibilityPreservesDeclarationOrder(t *testing.T) {
	registry := mustRegistry(t, []manifest.ModuleSpec{
		moduleSpec("foundation"),
		moduleSpec("internal-runtime", withVisibility("internal")),
		moduleSpec("control"),
		moduleSpec("old-module", withVisibility("disabled")),
	})

	assertModules(t, registry.ModulesByVisibility(manifest.VisibilityPublic), "foundation", "control")
	assertModules(t, registry.ModulesByVisibility(manifest.VisibilityInternal), "internal-runtime")
	assertModules(t, registry.ModulesByVisibility(manifest.VisibilityDisabled), "old-module")
	if got := registry.ModulesByVisibility(manifest.Visibility("unknown")); len(got) != 0 {
		t.Fatalf("ModulesByVisibility(unknown) len = %d, want 0", len(got))
	}
}

func TestModulesByVisibilityReturnsDetachedSlice(t *testing.T) {
	registry := mustRegistry(t, []manifest.ModuleSpec{moduleSpec("foundation"), moduleSpec("control")})

	modules := registry.ModulesByVisibility(manifest.VisibilityPublic)
	modules[0] = modules[1]

	assertModules(t, registry.ModulesByVisibility(manifest.VisibilityPublic), "foundation", "control")
}
