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

import (
	"context"
	"path/filepath"
	"testing"
)

func TestLoadModuleJSON(t *testing.T) {
	reader := newMemoryReader()
	path := filepath.Join(t.TempDir(), "arcpub.module.json")
	reader.add(path, foundationModuleJSON())

	loader := NewLoader(LoaderOptions{Reader: reader})
	mod, err := loader.LoadModule(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	if mod.Metadata().Name() != "foundation" {
		t.Fatalf("unexpected module: %s", mod.Metadata().Name())
	}
}

func TestLoadModuleReportsDecodeFailure(t *testing.T) {
	reader := newMemoryReader()
	path := filepath.Join(t.TempDir(), "arcpub.module.yaml")
	reader.add(path, "apiVersion: [bad]\n")

	loader := NewLoader(LoaderOptions{Reader: reader})
	if _, err := loader.LoadModule(context.Background(), path); err == nil {
		t.Fatal("expected decode failure")
	}
}
