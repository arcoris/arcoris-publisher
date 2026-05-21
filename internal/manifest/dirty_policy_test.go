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

func TestParseDirtyPolicyAcceptsKnownPolicies(t *testing.T) {
	for _, want := range []manifest.DirtyPolicy{
		manifest.DirtyPolicyFail,
		manifest.DirtyPolicyWarn,
		manifest.DirtyPolicyAllow,
	} {
		got, err := manifest.ParseDirtyPolicy(string(want))
		if err != nil {
			t.Fatalf("ParseDirtyPolicy(%q) returned error: %v", want, err)
		}
		if got != want {
			t.Fatalf("ParseDirtyPolicy(%q) = %q", want, got)
		}
	}
}

func TestParseDirtyPolicyRejectsUnknownPolicy(t *testing.T) {
	if _, err := manifest.ParseDirtyPolicy("ignore"); err == nil {
		t.Fatalf("expected unsupported dirty policy error")
	}
}
