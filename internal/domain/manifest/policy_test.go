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

package manifest

import "testing"

func TestNewPolicyDefaults(t *testing.T) {
	policy, err := NewPolicy(PolicySpec{})
	if err != nil {
		t.Fatalf("NewPolicy() error = %v", err)
	}
	if policy.VersionPolicy() != VersionPolicyReleaseTrain {
		t.Fatalf("VersionPolicy() = %q", policy.VersionPolicy())
	}
	if policy.PushPolicy() != PushPolicyFastForwardOnly {
		t.Fatalf("PushPolicy() = %q", policy.PushPolicy())
	}
}

func TestParseVersionPolicyRejectsUnknown(t *testing.T) {
	if _, err := ParseVersionPolicy("floating"); err == nil {
		t.Fatalf("ParseVersionPolicy() error = nil, want error")
	}
}

func TestParsePushPolicyRejectsUnknown(t *testing.T) {
	if _, err := ParsePushPolicy("force"); err == nil {
		t.Fatalf("ParsePushPolicy() error = nil, want error")
	}
}

func TestParseVisibilityRejectsUnknown(t *testing.T) {
	if _, err := ParseVisibility("private"); err == nil {
		t.Fatalf("ParseVisibility() error = nil, want error")
	}
}

func TestNewSource(t *testing.T) {
	source, err := NewSource(SourceSpec{Repository: "arcoris/arcoris", DefaultBranch: "main"})
	if err != nil {
		t.Fatalf("NewSource() error = %v", err)
	}
	if source.Repository() != RepositoryRef("arcoris/arcoris") || source.DefaultBranch() != BranchName("main") {
		t.Fatalf("source = %#v", source)
	}
}

func TestParseVersion(t *testing.T) {
	if version, err := ParseVersion("v1"); err != nil || version != VersionV1 {
		t.Fatalf("ParseVersion(v1) = %q, %v", version, err)
	}
	if _, err := ParseVersion(""); err == nil {
		t.Fatalf("ParseVersion(empty) error = nil, want error")
	}
}

func TestSourceAndPolicySpecRoundTrip(t *testing.T) {
	source := Must(validSpec()).Source()
	if got := source.Spec(); got.Repository != "arcoris/arcoris" || got.DefaultBranch != "main" {
		t.Fatalf("Source.Spec() = %#v", got)
	}
	policy := Must(validSpec()).Policy()
	if got := policy.Spec(); got.VersionPolicy != "release-train" || got.PushPolicy != "fast-forward-only" {
		t.Fatalf("Policy.Spec() = %#v", got)
	}
}
