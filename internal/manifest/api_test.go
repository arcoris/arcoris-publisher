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

func TestValidateAPIVersionAcceptsSupportedVersion(t *testing.T) {
	got, err := manifest.ValidateAPIVersion(string(manifest.APIVersionV1Alpha1))
	if err != nil {
		t.Fatalf("ValidateAPIVersion returned error: %v", err)
	}
	if got != manifest.APIVersionV1Alpha1 {
		t.Fatalf("unexpected version: %q", got)
	}
}

func TestValidateAPIVersionRejectsMissingAndUnsupportedVersions(t *testing.T) {
	for _, value := range []string{"", "v1"} {
		if _, err := manifest.ValidateAPIVersion(value); err == nil {
			t.Fatalf("ValidateAPIVersion(%q) returned nil error", value)
		}
	}
}

func TestValidateKindAcceptsExpectedKind(t *testing.T) {
	got, err := manifest.ValidateKind(string(manifest.KindStagingManifest), manifest.KindStagingManifest)
	if err != nil {
		t.Fatalf("ValidateKind returned error: %v", err)
	}
	if got != manifest.KindStagingManifest {
		t.Fatalf("unexpected kind: %q", got)
	}
}

func TestValidateKindRejectsMissingAndUnexpectedKinds(t *testing.T) {
	for _, value := range []string{"", string(manifest.KindModuleManifest)} {
		if _, err := manifest.ValidateKind(value, manifest.KindStagingManifest); err == nil {
			t.Fatalf("ValidateKind(%q) returned nil error", value)
		}
	}
}
