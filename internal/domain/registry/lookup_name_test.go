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

func TestModuleByName(t *testing.T) {
	registry := mustRegistry(t, []manifest.ModuleSpec{moduleSpec("foundation"), moduleSpec("control")})

	module, ok := registry.ModuleByName(name("foundation"))
	if !ok || module.Name() != name("foundation") {
		t.Fatalf("ModuleByName() = %q, %v", module.Name(), ok)
	}
	if _, ok := registry.ModuleByName(name("missing")); ok {
		t.Fatalf("ModuleByName(missing) ok = true")
	}
	if !registry.ContainsName(name("foundation")) || registry.ContainsName(name("missing")) {
		t.Fatalf("ContainsName returned unexpected values")
	}
}
