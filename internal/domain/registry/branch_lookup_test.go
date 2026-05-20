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
	registry := mustRegistry(t, specs(
		moduleSpec("foundation", withBranch("main", "main"), withBranch("release-1.0", "release-1.0")),
	))

	mapping, ok := registry.BranchMapping(name("foundation"), branch("release-1.0"))
	if !ok {
		t.Fatalf("BranchMapping() ok = false")
	}
	if mapping.Source() != branch("release-1.0") || mapping.Target() != branch("release-1.0") {
		t.Fatalf("mapping = %q -> %q", mapping.Source(), mapping.Target())
	}
	if _, ok := registry.BranchMapping(name("foundation"), branch("missing")); ok {
		t.Fatalf("BranchMapping(missing branch) ok = true")
	}
	if _, ok := registry.BranchMapping(name("missing"), branch("main")); ok {
		t.Fatalf("BranchMapping(missing module) ok = true")
	}
}
