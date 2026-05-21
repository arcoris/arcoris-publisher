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

func TestNewBranchMappingRoundTripsSpec(t *testing.T) {
	mapping, err := manifest.NewBranchMapping(manifest.BranchMappingSpec{Source: "main", Target: "release"})
	if err != nil {
		t.Fatalf("NewBranchMapping returned error: %v", err)
	}
	if mapping.Source().String() != "main" || mapping.Target().String() != "release" {
		t.Fatalf("unexpected branch mapping: %#v", mapping)
	}
	if mapping.Spec().Target != "release" {
		t.Fatalf("unexpected branch mapping spec: %#v", mapping.Spec())
	}
}

func TestNewBranchMappingRejectsInvalidSourceOrTarget(t *testing.T) {
	for _, spec := range []manifest.BranchMappingSpec{
		{Source: "bad branch", Target: "main"},
		{Source: "main", Target: "bad branch"},
	} {
		if _, err := manifest.NewBranchMapping(spec); err == nil {
			t.Fatalf("NewBranchMapping(%#v) returned nil error", spec)
		}
	}
}
