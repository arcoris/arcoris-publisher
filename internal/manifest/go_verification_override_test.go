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

func TestNewGoVerificationOverrideAppliesOnlyDeclaredFields(t *testing.T) {
	mode := string(manifest.GoWorkspaceModeDefault)
	list := false
	override, err := manifest.NewGoVerificationOverride(manifest.GoVerifySpec{
		WorkspaceMode: &mode,
		List:          &list,
		Patterns:      []string{"./contracts/..."},
	})
	if err != nil {
		t.Fatalf("NewGoVerificationOverride returned error: %v", err)
	}
	merged := manifest.MergeGoVerification(
		manifest.BuiltInVerificationPolicy().Go(),
		override,
	)
	if merged.WorkspaceMode() != manifest.GoWorkspaceModeDefault ||
		merged.List() ||
		!merged.Test() ||
		!merged.Tidy() {
		t.Fatalf("unexpected merged Go verification policy")
	}
	if got := merged.Patterns(); len(got) != 1 || got[0] != "./contracts/..." {
		t.Fatalf("unexpected merged patterns: %#v", got)
	}
}

func TestNewGoVerificationOverrideRejectsInvalidFields(t *testing.T) {
	mode := "workspace"
	for _, spec := range []manifest.GoVerifySpec{
		{WorkspaceMode: &mode},
		{Patterns: []string{""}},
	} {
		if _, err := manifest.NewGoVerificationOverride(spec); err == nil {
			t.Fatalf("NewGoVerificationOverride(%#v) returned nil error", spec)
		}
	}
}
