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
	"path/filepath"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
)

func TestEnsureInsideAcceptsNestedPath(t *testing.T) {
	root := filepath.Join(t.TempDir(), "root")
	candidate := filepath.Join(root, "a", "b")
	if err := ensureInside(root, candidate); err != nil {
		t.Fatal(err)
	}
}

func TestEnsureInsideRejectsEscapes(t *testing.T) {
	root := filepath.Join(t.TempDir(), "root")
	candidate := filepath.Join(root, "..", "other")
	if err := ensureInside(root, candidate); err == nil {
		t.Fatal("expected escape error")
	}
}

func TestEnsureInsideAcceptsRoot(t *testing.T) {
	root := filepath.Join(t.TempDir(), "root")
	if err := ensureInside(root, root); err != nil {
		t.Fatal(err)
	}
}

func TestResolveModuleManifestPathRejectsEscapes(t *testing.T) {
	sourcePath := filepath.Join(t.TempDir(), "module")
	_, err := resolveModuleManifestPath(sourcePath, manifest.RelativePath("../x.yaml"))
	if err == nil {
		t.Fatal("expected escape error")
	}
}

func TestModuleManifestPathUsesOverride(t *testing.T) {
	spec := staging.ModuleSpec{
		Name:       "foundation",
		SourceDir:  "src/foundation",
		Repository: "arcoris/foundation",
		Manifest:   stringPtr("custom.yaml"),
	}
	mod, err := staging.NewModule(spec)
	if err != nil {
		t.Fatal(err)
	}

	got := moduleManifestPath(staging.Manifest{}, mod)
	if got.String() != "custom.yaml" {
		t.Fatalf("got %q", got)
	}
}
