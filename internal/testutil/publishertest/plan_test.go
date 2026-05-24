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

package publishertest

import "testing"

func TestPlanBuildsDependencyOrderedFixture(t *testing.T) {
	plan, err := Plan(
		PlanOptions{},
		Module{Name: "foundation"},
		Module{Name: "control", Dependencies: []string{"foundation"}},
	)
	if err != nil {
		t.Fatalf("Plan() error = %v", err)
	}

	got := plan.ModuleNames()
	if len(got) != 2 || got[0] != "foundation" || got[1] != "control" {
		t.Fatalf("ModuleNames() = %v", got)
	}
}
