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
	"errors"
	"path/filepath"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

type memoryReader struct {
	files map[string][]byte
	errs  map[string]error
}

func newMemoryReader() *memoryReader {
	return &memoryReader{
		files: map[string][]byte{},
		errs:  map[string]error{},
	}
}

func (r *memoryReader) add(path string, data string) {
	abs, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	r.files[filepath.Clean(abs)] = []byte(data)
}

func (r *memoryReader) fail(path string, err error) {
	abs, absErr := filepath.Abs(path)
	if absErr != nil {
		panic(absErr)
	}
	r.errs[filepath.Clean(abs)] = err
}

func (r *memoryReader) ReadFile(ctx context.Context, path string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	path = filepath.Clean(abs)
	if err := r.errs[path]; err != nil {
		return nil, err
	}
	data, ok := r.files[path]
	if !ok {
		return nil, errors.New("not found")
	}
	out := make([]byte, len(data))
	copy(out, data)
	return out, nil
}

func (r *memoryReader) Exists(ctx context.Context, path string) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, err
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return false, err
	}
	path = filepath.Clean(abs)
	if err := r.errs[path]; err != nil {
		return false, err
	}
	_, ok := r.files[path]
	return ok, nil
}

func topLevelManifestPath(root string) string {
	return filepath.Join(root, "arcpub.yaml")
}

func foundationModulePath(root string) string {
	return filepath.Join(
		root,
		"src/arcoris.dev/foundation/arcpub.module.yaml",
	)
}

func controlModulePath(root string) string {
	return filepath.Join(
		root,
		"src/arcoris.dev/control/arcpub.module.yaml",
	)
}

func stringPtr(value string) *string {
	return &value
}

func assertControlDependsOnFoundation(
	t *testing.T,
	dependencies []manifest.ModuleName,
) {
	t.Helper()

	if len(dependencies) != 1 || dependencies[0].String() != "foundation" {
		t.Fatalf("unexpected dependencies: %#v", dependencies)
	}
}

func invalidModuleYAML() string {
	return `apiVersion: arcpub.arcoris.dev/v1alpha1
kind: ModuleManifest
`
}

func foundationModuleJSON() string {
	return `{
	"apiVersion": "arcpub.arcoris.dev/v1alpha1",
	"kind": "ModuleManifest",
	"metadata": {
		"name": "foundation"
	},
	"module": {
		"type": "go",
		"path": "arcoris.dev/foundation"
	},
	"publish": {
		"entries": [
			{
				"type": "file",
				"from": "go.mod",
				"to": "go.mod"
			}
		]
	}
}`
}

func stagingJSONWithUnknownField() []byte {
	return []byte(`{
	"apiVersion": "arcpub.arcoris.dev/v1alpha1",
	"kind": "StagingManifest",
	"metadata": {
		"name": "x"
	},
	"source": {
		"repository": "a/b",
		"defaultBranch": "main"
	},
	"modules": [],
	"unknown": true
}`)
}

func stagingJSONWithTrailingData() []byte {
	return []byte(`{
	"apiVersion": "arcpub.arcoris.dev/v1alpha1",
	"kind": "StagingManifest",
	"metadata": {
		"name": "x"
	},
	"source": {
		"repository": "a/b",
		"defaultBranch": "main"
	},
	"modules": []
} {}`)
}

func stagingYAML() string {
	return `apiVersion: arcpub.arcoris.dev/v1alpha1
kind: StagingManifest
metadata:
  name: arcoris
source:
  repository: arcoris/arcoris
  defaultBranch: main
  stagingRoot: .
  moduleRoot: src/arcoris.dev
publish:
  mode: explicit-projection
  versionPolicy: release-train
  pushPolicy: fast-forward-only
defaults:
  moduleManifest:
    path: arcpub.module.yaml
  verification:
    go:
      test: true
modules:
  - name: foundation
    sourceDir: src/arcoris.dev/foundation
    repository: arcoris/foundation
  - name: control
    sourceDir: src/arcoris.dev/control
    repository: arcoris/control
`
}

func foundationModuleYAML() string {
	return `apiVersion: arcpub.arcoris.dev/v1alpha1
kind: ModuleManifest
metadata:
  name: foundation
module:
  type: go
  path: arcoris.dev/foundation
publish:
  entries:
    - type: file
      from: go.mod
      to: go.mod
    - type: directory
      from: contracts
      to: contracts
`
}

func controlModuleYAML() string {
	return `apiVersion: arcpub.arcoris.dev/v1alpha1
kind: ModuleManifest
metadata:
  name: control
module:
  type: go
  path: arcoris.dev/control
dependencies:
  internal:
    - foundation
publish:
  entries:
    - type: file
      from: go.mod
      to: go.mod
    - type: directory
      from: runtime
      to: runtime
verification:
  go:
    test: false
`
}
