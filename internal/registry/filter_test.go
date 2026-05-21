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

import "testing"

func TestVisibilityFilters(t *testing.T) {
	registry := testRegistry(t)

	publishable := registry.PublishableModules()
	if len(publishable) != 2 || publishable[0].Name() != "foundation" {
		t.Fatalf("unexpected publishable modules: %#v", publishable)
	}

	internal := registry.InternalModules()
	if len(internal) != 1 || internal[0].Name() != "internal-tool" {
		t.Fatalf("unexpected internal modules: %#v", internal)
	}

	disabled := registry.DisabledModules()
	if len(disabled) != 1 || disabled[0].Name() != "disabled-tool" {
		t.Fatalf("unexpected disabled modules: %#v", disabled)
	}
}

func TestVisibilityFiltersReturnDetachedSlices(t *testing.T) {
	registry := testRegistry(t)
	publishable := registry.PublishableModules()
	publishable = publishable[:1]

	if len(registry.PublishableModules()) != 2 {
		t.Fatal("PublishableModules returned aliased slice")
	}
}
