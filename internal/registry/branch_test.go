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

package registry

import "testing"

func TestBranchMappingLookup(t *testing.T) {
	registry := testRegistry(t)

	mapping, ok := registry.BranchMapping("control", "release")
	if !ok {
		t.Fatal("expected control release mapping")
	}
	if mapping.Target() != "stable" {
		t.Fatalf("unexpected target branch: %s", mapping.Target())
	}

	if _, ok := registry.BranchMapping("missing", "main"); ok {
		t.Fatal("unexpected missing module mapping")
	}
	if _, ok := registry.BranchMapping("control", "main"); ok {
		t.Fatal("unexpected missing source mapping")
	}
}

func TestBranchListsAreDeduplicatedAndDeterministic(t *testing.T) {
	registry := testRegistry(t)

	sources := registry.SourceBranches()
	if len(sources) != 2 || sources[0] != "main" || sources[1] != "release" {
		t.Fatalf("unexpected source branches: %#v", sources)
	}

	targets := registry.TargetBranches()
	if len(targets) != 2 || targets[0] != "main" || targets[1] != "stable" {
		t.Fatalf("unexpected target branches: %#v", targets)
	}
}
