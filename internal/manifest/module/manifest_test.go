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

package module_test

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	modulemanifest "arcoris.dev/arcoris-publisher/internal/manifest/module"
)

func TestNewBuildsValidatedModuleManifest(t *testing.T) {
	m, err := modulemanifest.New(validSpec())
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	if m.APIVersion() != manifest.APIVersionV1Alpha1 || m.Kind() != manifest.KindModuleManifest {
		t.Fatalf("unexpected manifest api metadata")
	}
	if m.Metadata().Name() != "control" {
		t.Fatalf("unexpected metadata: %#v", m.Metadata())
	}
	if m.Module().Type() != manifest.ModuleTypeGo ||
		m.Module().Root().String() != "." ||
		m.Module().GoMod().String() != "go.mod" {
		t.Fatalf("unexpected module identity defaults")
	}
	if len(m.Dependencies().Internal()) != 1 || len(m.Publish().Entries()) != 2 {
		t.Fatalf("unexpected dependencies or publish entries")
	}
	verification := manifest.MergeVerification(manifest.BuiltInVerificationPolicy(), m.Verification())
	if verification.LocalReplacePolicy() != manifest.LocalReplacePolicyForbid {
		t.Fatalf("unexpected verification override")
	}
}

func TestNewCollectsValidationErrorsAcrossSections(t *testing.T) {
	spec := validSpec()
	spec.APIVersion = "v1"
	spec.Kind = string(manifest.KindStagingManifest)
	spec.Metadata.Name = "Control"
	if _, err := modulemanifest.New(spec); err == nil {
		t.Fatalf("expected validation error")
	}
}
