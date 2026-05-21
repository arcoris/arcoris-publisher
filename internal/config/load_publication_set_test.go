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

func TestLoadPublicationSetWithTrace(t *testing.T) {
	reader, _, stagingPath := publicationSetReader(t)

	loader := NewLoader(LoaderOptions{Reader: reader})
	result, err := loader.LoadPublicationSetWithTrace(
		context.Background(),
		stagingPath,
	)
	if err != nil {
		t.Fatal(err)
	}

	modules := result.Set.Modules()
	if len(modules) != 2 {
		t.Fatalf("got %d modules", len(modules))
	}
	if modules[1].Name().String() != "control" {
		t.Fatalf("unexpected second module: %s", modules[1].Name())
	}
	if modules[1].Verification().Go().Test() {
		t.Fatal("expected module verification override to disable go test")
	}
	if len(result.Trace.Fields()) == 0 {
		t.Fatal("expected resolution trace")
	}

	assertControlDependsOnFoundation(t, modules[1].Dependencies())
}

func TestLoadPublicationSetSuccessWrapper(t *testing.T) {
	reader, _, stagingPath := publicationSetReader(t)

	loader := NewLoader(LoaderOptions{Reader: reader})
	set, err := loader.LoadPublicationSet(context.Background(), stagingPath)
	if err != nil {
		t.Fatal(err)
	}
	if set.Metadata().Name() != "arcoris" {
		t.Fatalf("unexpected metadata: %s", set.Metadata().Name())
	}
}

func TestDiscoverPublicationSet(t *testing.T) {
	reader, root, _ := publicationSetReader(t)

	loader := NewLoader(LoaderOptions{Reader: reader})
	set, err := loader.DiscoverPublicationSet(context.Background(), root)
	if err != nil {
		t.Fatal(err)
	}
	if len(set.Modules()) != 2 {
		t.Fatalf("got %d modules", len(set.Modules()))
	}
}

func TestDiscoverPublicationSetWithTraceReportsLocatorFailure(t *testing.T) {
	loader := NewLoader(LoaderOptions{Reader: newMemoryReader()})
	_, err := loader.DiscoverPublicationSetWithTrace(
		context.Background(),
		t.TempDir(),
	)
	if err == nil {
		t.Fatal("expected locator failure")
	}
}

func TestLoadPublicationSetWithTraceReportsResolveFailure(t *testing.T) {
	reader := newMemoryReader()
	root := filepath.Join(t.TempDir(), "staging")
	stagingPath := topLevelManifestPath(root)

	reader.add(stagingPath, stagingYAML())
	reader.add(foundationModulePath(root), foundationModuleYAML())
	reader.add(controlModulePath(root), foundationModuleYAML())

	loader := NewLoader(LoaderOptions{Reader: reader})
	if _, err := loader.LoadPublicationSetWithTrace(context.Background(), stagingPath); err == nil {
		t.Fatal("expected resolve failure")
	}
}

func publicationSetReader(t *testing.T) (*memoryReader, string, string) {
	t.Helper()

	reader := newMemoryReader()
	root := filepath.Join(t.TempDir(), "staging")
	stagingPath := topLevelManifestPath(root)

	reader.add(stagingPath, stagingYAML())
	reader.add(foundationModulePath(root), foundationModuleYAML())
	reader.add(controlModulePath(root), controlModuleYAML())

	return reader, root, stagingPath
}
