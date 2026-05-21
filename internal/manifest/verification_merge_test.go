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

func TestMergeVerificationKeepsInheritedValuesWhenOverrideIsEmpty(t *testing.T) {
	base := manifest.BuiltInVerificationPolicy()
	merged := manifest.MergeVerification(base, manifest.VerificationOverride{})
	if merged.LocalReplacePolicy() != base.LocalReplacePolicy() ||
		merged.Go().WorkspaceMode() != base.Go().WorkspaceMode() {
		t.Fatalf("empty override changed base policy")
	}
}

func TestMergeVerificationAppliesNestedOverride(t *testing.T) {
	list := false
	test := false
	tidy := false
	patterns := []string{"./contracts/..."}
	override, err := manifest.NewVerificationOverride(manifest.VerificationSpec{
		Go: manifest.GoVerifySpec{
			List:     &list,
			Test:     &test,
			Tidy:     &tidy,
			Patterns: patterns,
		},
	})
	if err != nil {
		t.Fatalf("NewVerificationOverride returned error: %v", err)
	}
	merged := manifest.MergeVerification(manifest.BuiltInVerificationPolicy(), override)
	if merged.Go().List() {
		t.Fatalf("expected list override")
	}
	if merged.Go().Test() || merged.Go().Tidy() {
		t.Fatalf("expected test/tidy overrides")
	}
	got := merged.Go().Patterns()
	if len(got) != 1 || got[0] != "./contracts/..." {
		t.Fatalf("unexpected patterns: %#v", got)
	}
}
