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

func TestNewProvenancePolicyAppliesDefaults(t *testing.T) {
	policy, err := manifest.NewProvenancePolicy(manifest.ProvenanceSpec{})
	if err != nil {
		t.Fatalf("NewProvenancePolicy returned error: %v", err)
	}
	if !policy.CommitTrailers() || policy.FileEnabled() {
		t.Fatalf("unexpected provenance defaults")
	}
}

func TestNewProvenancePolicyAcceptsFileAndExplicitTrailers(t *testing.T) {
	file := ".arcpub/provenance.json"
	policy, err := manifest.NewProvenancePolicy(manifest.ProvenanceSpec{
		CommitTrailers: boolPtr(false),
		File:           &file,
	})
	if err != nil {
		t.Fatalf("NewProvenancePolicy returned error: %v", err)
	}
	if policy.CommitTrailers() || !policy.FileEnabled() || policy.File().String() != file {
		t.Fatalf("unexpected provenance policy")
	}
}

func TestNewProvenancePolicyRejectsUnsafeFilePath(t *testing.T) {
	file := "../provenance.json"
	if _, err := manifest.NewProvenancePolicy(manifest.ProvenanceSpec{File: &file}); err == nil {
		t.Fatalf("expected invalid provenance file path")
	}
}
