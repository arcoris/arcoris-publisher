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
	"arcoris.dev/arcoris-publisher/internal/domain/versioning"
)

func TestValidateZeroPlanSucceeds(t *testing.T) {
	var planValue Plan
	if err := planValue.Validate(); err != nil {
		t.Fatalf("zero plan Validate() error = %v", err)
	}
}

func TestValidateReportsInvalidModuleDetails(t *testing.T) {
	modulePlan := testPlan(t).Modules()[0]
	modulePlan.version = ""
	modulePlan.action = Action("archive")
	modulePlan.branches = []BranchPlan{{source: "", target: ""}}
	modulePlan.dependencies = []DependencyPlan{{}}
	modulePlan.requirements = map[manifest.ModulePath]versioning.Version{
		modulePath(t, "arcoris.dev/extra"): "v9.9.9",
	}
	manual := Plan{
		modules:  []ModulePlan{modulePlan},
		byModule: map[manifest.ModuleName]int{modulePlan.Name(): 0},
		byPath:   map[manifest.ModulePath]int{modulePlan.ModulePath(): 0},
		byRepo:   map[manifest.RepositoryRef]int{modulePlan.Repository(): 0},
	}

	validationErr := mustValidationError(t, manual.Validate())
	for _, code := range []IssueCode{IssueMissingVersion, IssueInvalidAction, IssueInvalidBranch, IssueInvalidDependency} {
		if !hasIssueCode(validationErr.Issues, code) {
			t.Fatalf("issues = %#v, want %s", validationErr.Issues, code)
		}
	}
}

func TestValidateReportsRequirementVersionMismatch(t *testing.T) {
	modulePlan := testPlan(t).Modules()[1]
	modulePlan.requirements[modulePath(t, "arcoris.dev/foundation")] = "v9.9.9"
	manual := Plan{
		modules:  []ModulePlan{modulePlan},
		byModule: map[manifest.ModuleName]int{modulePlan.Name(): 0},
		byPath:   map[manifest.ModulePath]int{modulePlan.ModulePath(): 0},
		byRepo:   map[manifest.RepositoryRef]int{modulePlan.Repository(): 0},
	}

	validationErr := mustValidationError(t, manual.Validate())
	if !hasIssueCode(validationErr.Issues, IssueInvalidDependency) {
		t.Fatalf("issues = %#v, want invalid dependency", validationErr.Issues)
	}
}

func TestValidateReportsDuplicateNaturalKeys(t *testing.T) {
	modulePlan := testPlan(t).Modules()[0]
	manual := Plan{
		modules:  []ModulePlan{modulePlan, modulePlan},
		byModule: map[manifest.ModuleName]int{modulePlan.Name(): 1},
		byPath:   map[manifest.ModulePath]int{modulePlan.ModulePath(): 1},
		byRepo:   map[manifest.RepositoryRef]int{modulePlan.Repository(): 1},
	}

	validationErr := mustValidationError(t, manual.Validate())
	for _, code := range []IssueCode{IssueDuplicateModule, IssueDuplicateModulePath, IssueDuplicateRepository} {
		if !hasIssueCode(validationErr.Issues, code) {
			t.Fatalf("issues = %#v, want %s", validationErr.Issues, code)
		}
	}
}

func TestValidateReportsBrokenIndexes(t *testing.T) {
	planValue := testPlan(t)
	delete(planValue.byModule, moduleName(t, "control"))
	planValue.byModule[moduleName(t, "foundation")] = 1
	planValue.byModule[moduleName(t, "unexpected")] = 0
	delete(planValue.byPath, modulePath(t, "arcoris.dev/control"))
	planValue.byPath[modulePath(t, "arcoris.dev/foundation")] = 1
	planValue.byPath[modulePath(t, "arcoris.dev/unexpected")] = 0
	planValue.byRepo[repositoryRef(t, "arcoris/control")] = 0
	planValue.byRepo[repositoryRef(t, "arcoris/unexpected")] = 0

	validationErr := mustValidationError(t, planValue.Validate())
	for _, index := range []string{"byModule", "byPath", "byRepo"} {
		if !hasIssueIndex(validationErr.Issues, index) {
			t.Fatalf("issues = %#v, want index %s", validationErr.Issues, index)
		}
	}
}
