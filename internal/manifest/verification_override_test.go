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

func TestNewVerificationOverrideAppliesLocalReplacePolicy(t *testing.T) {
	local := string(manifest.LocalReplacePolicyWarn)
	override, err := manifest.NewVerificationOverride(manifest.VerificationSpec{LocalReplacePolicy: &local})
	if err != nil {
		t.Fatalf("NewVerificationOverride returned error: %v", err)
	}
	merged := manifest.MergeVerification(manifest.BuiltInVerificationPolicy(), override)
	if merged.LocalReplacePolicy() != manifest.LocalReplacePolicyWarn {
		t.Fatalf("unexpected local replace policy: %q", merged.LocalReplacePolicy())
	}
}

func TestNewVerificationOverrideRejectsInvalidNestedPolicy(t *testing.T) {
	local := "ignore"
	if _, err := manifest.NewVerificationOverride(manifest.VerificationSpec{LocalReplacePolicy: &local}); err == nil {
		t.Fatalf("expected invalid local replace policy error")
	}
	mode := "workspace"
	if _, err := manifest.NewVerificationOverride(manifest.VerificationSpec{Go: manifest.GoVerifySpec{WorkspaceMode: &mode}}); err == nil {
		t.Fatalf("expected invalid Go verification override error")
	}
}
