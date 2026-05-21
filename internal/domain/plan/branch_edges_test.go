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

package plan

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

func TestNewBranchPlanRejectsEmptyTarget(t *testing.T) {
	if _, err := newBranchPlan(branchName(t, "main"), ""); err == nil {
		t.Fatal("newBranchPlan returned nil error for empty target")
	}
}

func TestBranchPlansPropagatesInvalidMapping(t *testing.T) {
	if _, err := branchPlans([]manifest.BranchMapping{{}}); err == nil {
		t.Fatal("branchPlans returned nil error for invalid mapping")
	}
}
