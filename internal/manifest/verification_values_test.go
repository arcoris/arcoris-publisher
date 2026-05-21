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

func TestParseLocalReplacePolicyAcceptsKnownPolicies(t *testing.T) {
	for _, want := range []manifest.LocalReplacePolicy{manifest.LocalReplacePolicyForbid, manifest.LocalReplacePolicyWarn, manifest.LocalReplacePolicyAllow} {
		got, err := manifest.ParseLocalReplacePolicy(string(want))
		if err != nil {
			t.Fatalf("ParseLocalReplacePolicy(%q) returned error: %v", want, err)
		}
		if got != want {
			t.Fatalf("ParseLocalReplacePolicy(%q) = %q", want, got)
		}
	}
}

func TestParseGoWorkspaceModeAcceptsKnownModes(t *testing.T) {
	for _, want := range []manifest.GoWorkspaceMode{manifest.GoWorkspaceModeOff, manifest.GoWorkspaceModeDefault} {
		got, err := manifest.ParseGoWorkspaceMode(string(want))
		if err != nil {
			t.Fatalf("ParseGoWorkspaceMode(%q) returned error: %v", want, err)
		}
		if got != want {
			t.Fatalf("ParseGoWorkspaceMode(%q) = %q", want, got)
		}
	}
}

func TestVerificationValueParsersRejectUnknownValues(t *testing.T) {
	checks := []func() error{
		func() error { _, err := manifest.ParseLocalReplacePolicy("bad"); return err },
		func() error { _, err := manifest.ParseGoWorkspaceMode("bad"); return err },
	}
	for i, check := range checks {
		if err := check(); err == nil {
			t.Fatalf("check %d returned nil error", i)
		}
	}
}
