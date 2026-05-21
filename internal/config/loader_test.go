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

package config

import "testing"

func TestNewLoaderAppliesDefaultDependencies(t *testing.T) {
	loader := NewLoader(LoaderOptions{})
	if loader.reader == nil || loader.decoder == nil || len(loader.locator.Names) == 0 {
		t.Fatalf("loader defaults not applied: %#v", loader)
	}
}

func TestNewLoaderPreservesInjectedDependencies(t *testing.T) {
	reader := newMemoryReader()
	decoder := StrictDecoder{}
	locator := Locator{Names: []string{"custom.yaml"}}

	loader := NewLoader(LoaderOptions{
		Reader:  reader,
		Decoder: decoder,
		Locator: locator,
	})

	if loader.reader != reader {
		t.Fatal("reader was not preserved")
	}
	if len(loader.locator.Names) != 1 || loader.locator.Names[0] != "custom.yaml" {
		t.Fatalf("locator was not preserved: %#v", loader.locator)
	}
}
