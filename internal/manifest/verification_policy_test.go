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

func TestBuiltInVerificationPolicyUsesSafeDefaults(t *testing.T) {
	policy := manifest.BuiltInVerificationPolicy()
	if policy.LocalReplacePolicy() != manifest.LocalReplacePolicyForbid {
		t.Fatalf("unexpected local replace policy: %q", policy.LocalReplacePolicy())
	}
	if policy.Go().WorkspaceMode() != manifest.GoWorkspaceModeOff || !policy.Go().List() || !policy.Go().Test() || !policy.Go().Tidy() {
		t.Fatalf("unexpected Go verification defaults")
	}
	if got := policy.Go().Patterns(); len(got) != 1 || got[0] != "./..." {
		t.Fatalf("unexpected default patterns: %#v", got)
	}
}

func TestGoVerificationPolicyPatternsAreDetached(t *testing.T) {
	policy := manifest.BuiltInVerificationPolicy()
	patterns := policy.Go().Patterns()
	patterns[0] = "mutated"
	if policy.Go().Patterns()[0] == "mutated" {
		t.Fatalf("patterns accessor leaked internal slice")
	}
}
