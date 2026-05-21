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

func TestNewTagPolicyAppliesDefaultsAndExplicitEnabled(t *testing.T) {
	policy, err := manifest.NewTagPolicy(manifest.TagPolicySpec{Enabled: boolPtr(false)})
	if err != nil {
		t.Fatalf("NewTagPolicy returned error: %v", err)
	}
	if policy.Enabled() || policy.Mode() != manifest.TagModeSemver {
		t.Fatalf("unexpected tag policy: %#v", policy)
	}
}

func TestParseTagModeAcceptsSemverAndRejectsUnknownModes(t *testing.T) {
	got, err := manifest.ParseTagMode(string(manifest.TagModeSemver))
	if err != nil {
		t.Fatalf("ParseTagMode returned error: %v", err)
	}
	if got != manifest.TagModeSemver {
		t.Fatalf("unexpected tag mode: %q", got)
	}
	if _, err := manifest.ParseTagMode("calendar"); err == nil {
		t.Fatalf("expected invalid tag mode error")
	}
}

func TestNewTagPolicyRejectsInvalidMode(t *testing.T) {
	if _, err := manifest.NewTagPolicy(
		manifest.TagPolicySpec{
			Mode: stringPtr("calendar"),
		},
	); err == nil {
		t.Fatalf("expected invalid tag policy error")
	}
}
