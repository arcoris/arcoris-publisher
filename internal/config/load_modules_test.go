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
)

func TestLoadPublicationSetReportsMissingModuleManifest(t *testing.T) {
	reader := newMemoryReader()
	root := filepath.Join(t.TempDir(), "staging")
	stagingPath := topLevelManifestPath(root)

	reader.add(stagingPath, stagingYAML())
	reader.add(foundationModulePath(root), foundationModuleYAML())

	loader := NewLoader(LoaderOptions{Reader: reader})
	if _, err := loader.LoadPublicationSet(context.Background(), stagingPath); err == nil {
		t.Fatal("expected missing control module manifest error")
	}
}

func TestLoadPublicationSetReportsInvalidModuleManifest(t *testing.T) {
	reader := newMemoryReader()
	root := filepath.Join(t.TempDir(), "staging")
	stagingPath := topLevelManifestPath(root)

	reader.add(stagingPath, stagingYAML())
	reader.add(foundationModulePath(root), foundationModuleYAML())
	reader.add(controlModulePath(root), invalidModuleYAML())

	loader := NewLoader(LoaderOptions{Reader: reader})
	if _, err := loader.LoadPublicationSet(context.Background(), stagingPath); err == nil {
		t.Fatal("expected invalid module manifest error")
	}
}

func TestLoadPublicationSetReportsReadFailure(t *testing.T) {
	reader := newMemoryReader()
	root := filepath.Join(t.TempDir(), "staging")
	stagingPath := topLevelManifestPath(root)

	reader.add(stagingPath, stagingYAML())
	reader.fail(foundationModulePath(root), errors.New("permission denied"))
	reader.add(controlModulePath(root), controlModuleYAML())

	loader := NewLoader(LoaderOptions{Reader: reader})
	if _, err := loader.LoadPublicationSet(context.Background(), stagingPath); err == nil {
		t.Fatal("expected read failure")
	}
}
