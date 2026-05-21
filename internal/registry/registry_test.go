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

func TestNewIndexesPublicationSetInDeclarationOrder(t *testing.T) {
	registry, err := New(testPublicationSet(t))
	if err != nil {
		t.Fatal(err)
	}

	modules := registry.Modules()
	if len(modules) != 4 {
		t.Fatalf("got %d modules", len(modules))
	}
	if modules[0].Name().String() != "foundation" {
		t.Fatalf("unexpected first module: %s", modules[0].Name())
	}
	if modules[1].Name().String() != "control" {
		t.Fatalf("unexpected second module: %s", modules[1].Name())
	}
}

func TestMustPanicsOnInvalidPublicationSet(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()

	Must(duplicatePublicationSet(t))
}

func TestModulesReturnsDetachedSlice(t *testing.T) {
	registry := testRegistry(t)
	modules := registry.Modules()
	modules = modules[:1]

	if len(registry.Modules()) != 4 {
		t.Fatal("Modules returned aliased slice")
	}
}
