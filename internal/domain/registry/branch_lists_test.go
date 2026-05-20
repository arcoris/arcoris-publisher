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

func TestSourceAndTargetBranches(t *testing.T) {
	registry := mustRegistry(t, specs(
		moduleSpec("foundation", withBranch("main", "published-main"), withBranch("release-1.0", "published-release-1.0")),
	))

	assertBranches(t, registry.SourceBranches(name("foundation")), "main", "release-1.0")
	assertBranches(t, registry.TargetBranches(name("foundation")), "published-main", "published-release-1.0")
	if got := registry.SourceBranches(name("missing")); got != nil {
		t.Fatalf("SourceBranches(missing) = %#v, want nil", got)
	}
	if got := registry.TargetBranches(name("missing")); got != nil {
		t.Fatalf("TargetBranches(missing) = %#v, want nil", got)
	}
}
