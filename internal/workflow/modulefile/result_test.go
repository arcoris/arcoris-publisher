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

package modulefile

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestResultAccessorsDetach(t *testing.T) {
	name := manifest.ModuleName("control")
	result := Result{modules: []ModuleResult{{
		module:    name,
		goModPath: "/target/go.mod",
		changed:   true,
		requirements: []RequirementUpdate{{
			modulePath: manifest.ModulePath("arcoris.dev/foundation"),
			version:    "v0.1.0",
		}},
	}}}

	modules := result.Modules()
	modules[0].module = "mutated"
	if got, ok := result.ModuleByName(name); !ok || got.Module() != name {
		t.Fatal("Modules() returned attached slice")
	}
	if !result.Changed() {
		t.Fatal("Changed() = false")
	}

	requirements := result.Modules()[0].Requirements()
	requirements[0].version = "mutated"
	if result.Modules()[0].Requirements()[0].Version() != "v0.1.0" {
		t.Fatal("Requirements() returned attached slice")
	}
}
