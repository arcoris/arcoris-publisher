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

package manifest_test

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestNewMetadataRoundTripsValidatedName(t *testing.T) {
	metadata, err := manifest.NewMetadata(manifest.MetadataSpec{Name: "arcoris-core"})
	if err != nil {
		t.Fatalf("NewMetadata returned error: %v", err)
	}
	if metadata.Name() != "arcoris-core" {
		t.Fatalf("unexpected metadata name: %q", metadata.Name())
	}
	if metadata.Spec().Name != "arcoris-core" {
		t.Fatalf("unexpected metadata spec: %#v", metadata.Spec())
	}
}

func TestNewMetadataRejectsInvalidName(t *testing.T) {
	if _, err := manifest.NewMetadata(manifest.MetadataSpec{Name: "bad/name"}); err == nil {
		t.Fatalf("expected invalid metadata name error")
	}
}
