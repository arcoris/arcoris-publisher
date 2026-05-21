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

func TestReadManifestFileReturnsDataFormatAndPath(t *testing.T) {
	reader := newMemoryReader()
	path := filepath.Join(t.TempDir(), "arcpub.yaml")
	reader.add(path, "kind: StagingManifest\n")

	loader := NewLoader(LoaderOptions{Reader: reader})
	file, err := loader.readManifestFile(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}

	if file.Format != FormatYAML {
		t.Fatalf("unexpected format: %s", file.Format)
	}
	if string(file.Data) != "kind: StagingManifest\n" {
		t.Fatalf("unexpected data: %q", file.Data)
	}
	if !filepath.IsAbs(file.Path) {
		t.Fatalf("expected absolute path, got %q", file.Path)
	}
}

func TestReadManifestFileReportsReadFailure(t *testing.T) {
	loader := NewLoader(LoaderOptions{Reader: newMemoryReader()})
	_, err := loader.readManifestFile(
		context.Background(),
		filepath.Join(t.TempDir(), "missing.yaml"),
	)
	if err == nil {
		t.Fatal("expected read failure")
	}
}
