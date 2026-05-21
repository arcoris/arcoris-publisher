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

func TestParsePublishModeAcceptsExplicitProjection(t *testing.T) {
	got, err := manifest.ParsePublishMode(string(manifest.PublishModeExplicitProjection))
	if err != nil {
		t.Fatalf("ParsePublishMode returned error: %v", err)
	}
	if got != manifest.PublishModeExplicitProjection {
		t.Fatalf("unexpected publish mode: %q", got)
	}
}

func TestParseVersionPolicyAcceptsKnownPolicies(t *testing.T) {
	for _, want := range []manifest.VersionPolicy{manifest.VersionPolicyReleaseTrain, manifest.VersionPolicySnapshot} {
		got, err := manifest.ParseVersionPolicy(string(want))
		if err != nil {
			t.Fatalf("ParseVersionPolicy(%q) returned error: %v", want, err)
		}
		if got != want {
			t.Fatalf("ParseVersionPolicy(%q) = %q", want, got)
		}
	}
}

func TestParsePushPolicyAcceptsKnownPolicies(t *testing.T) {
	for _, want := range []manifest.PushPolicy{manifest.PushPolicyFastForwardOnly, manifest.PushPolicyCreateOnly, manifest.PushPolicyForceWithLease} {
		got, err := manifest.ParsePushPolicy(string(want))
		if err != nil {
			t.Fatalf("ParsePushPolicy(%q) returned error: %v", want, err)
		}
		if got != want {
			t.Fatalf("ParsePushPolicy(%q) = %q", want, got)
		}
	}
}

func TestPublishValueParsersRejectUnknownValues(t *testing.T) {
	checks := []func() error{
		func() error { _, err := manifest.ParsePublishMode("bad"); return err },
		func() error { _, err := manifest.ParseVersionPolicy("bad"); return err },
		func() error { _, err := manifest.ParsePushPolicy("bad"); return err },
	}
	for i, check := range checks {
		if err := check(); err == nil {
			t.Fatalf("check %d returned nil error", i)
		}
	}
}
