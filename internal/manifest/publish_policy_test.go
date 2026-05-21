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

func TestNewPublishPolicyAppliesBuiltInDefaults(t *testing.T) {
	policy, err := manifest.NewPublishPolicy(manifest.PublishSpec{})
	if err != nil {
		t.Fatalf("NewPublishPolicy returned error: %v", err)
	}
	if policy.Mode() != manifest.PublishModeExplicitProjection || policy.VersionPolicy() != manifest.VersionPolicyReleaseTrain || policy.PushPolicy() != manifest.PushPolicyFastForwardOnly {
		t.Fatalf("unexpected publish defaults")
	}
	if !policy.Tags().Enabled() || policy.Tags().Mode() != manifest.TagModeSemver {
		t.Fatalf("unexpected tag defaults")
	}
	if !policy.Provenance().CommitTrailers() || policy.Provenance().FileEnabled() {
		t.Fatalf("unexpected provenance defaults")
	}
}

func TestNewPublishPolicyAcceptsExplicitValues(t *testing.T) {
	versionPolicy := string(manifest.VersionPolicySnapshot)
	pushPolicy := string(manifest.PushPolicyCreateOnly)
	policy, err := manifest.NewPublishPolicy(manifest.PublishSpec{VersionPolicy: &versionPolicy, PushPolicy: &pushPolicy})
	if err != nil {
		t.Fatalf("NewPublishPolicy returned error: %v", err)
	}
	if policy.VersionPolicy() != manifest.VersionPolicySnapshot || policy.PushPolicy() != manifest.PushPolicyCreateOnly {
		t.Fatalf("explicit publish values were not applied")
	}
}

func TestNewPublishPolicyRejectsInvalidNestedPolicy(t *testing.T) {
	for _, spec := range []manifest.PublishSpec{
		{Mode: stringPtr("bad")},
		{VersionPolicy: stringPtr("bad")},
		{PushPolicy: stringPtr("bad")},
		{Tags: manifest.TagPolicySpec{Mode: stringPtr("bad")}},
		{Provenance: manifest.ProvenanceSpec{File: stringPtr("../provenance.json")}},
	} {
		if _, err := manifest.NewPublishPolicy(spec); err == nil {
			t.Fatalf("NewPublishPolicy(%#v) returned nil error", spec)
		}
	}
}
