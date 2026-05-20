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

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

func TestValidateReportsBrokenIndexes(t *testing.T) {
	registry := mustRegistry(t, []manifest.ModuleSpec{moduleSpec("foundation"), moduleSpec("control")})
	registry.byName[name("foundation")] = 1
	registry.byName[name("unexpected")] = 0
	registry.byModulePath[modulePath("arcoris.dev/control")] = 0
	delete(registry.bySourceDir, sourceDir("staging/src/arcoris.dev/control"))
	registry.byRepository[repository("arcoris/unexpected")] = 0
	delete(registry.byVisibility, manifest.VisibilityPublic)
	registry.byVisibility[manifest.Visibility("unexpected")] = []int{0}
	registry.byBranch[name("foundation")][branch("main")] = mustBranchMapping(t, "main", "published-main")
	registry.byBranch[name("foundation")][branch("unexpected")] = mustBranchMapping(t, "unexpected", "unexpected")
	delete(registry.byBranch, name("control"))
	registry.byBranch[name("unexpected")] = map[manifest.BranchName]manifest.BranchMapping{}

	validationErr := mustValidationError(t, registry.Validate())
	if !hasIssueCode(validationErr.Issues, IssueInvalidIndex) {
		t.Fatalf("issues = %#v, want invalid index issue", validationErr.Issues)
	}
	for _, index := range []string{"byName", "byModulePath", "bySourceDir", "byRepository", "byVisibility", "byBranch[foundation]", "byBranch"} {
		if !hasIssueIndex(validationErr.Issues, index) {
			t.Fatalf("issues = %#v, want index %s", validationErr.Issues, index)
		}
	}
}

func TestValidateReportsWrongVisibilityBucket(t *testing.T) {
	registry := mustRegistry(t, []manifest.ModuleSpec{moduleSpec("foundation"), moduleSpec("control")})
	registry.byVisibility[manifest.VisibilityPublic] = []int{1, 0}

	validationErr := mustValidationError(t, registry.Validate())
	if !hasIssueIndex(validationErr.Issues, "byVisibility") {
		t.Fatalf("issues = %#v, want byVisibility issue", validationErr.Issues)
	}
}

func TestValidateReportsMissingVisibilityBucket(t *testing.T) {
	registry := mustRegistry(t, []manifest.ModuleSpec{moduleSpec("foundation")})
	delete(registry.byVisibility, manifest.VisibilityPublic)

	validationErr := mustValidationError(t, registry.Validate())
	if !hasIssueIndex(validationErr.Issues, "byVisibility") {
		t.Fatalf("issues = %#v, want byVisibility issue", validationErr.Issues)
	}
}

func TestValidateReportsMissingBranchModule(t *testing.T) {
	registry := mustRegistry(t, []manifest.ModuleSpec{moduleSpec("foundation")})
	delete(registry.byBranch, name("foundation"))

	validationErr := mustValidationError(t, registry.Validate())
	if !hasIssueIndex(validationErr.Issues, "byBranch") {
		t.Fatalf("issues = %#v, want byBranch issue", validationErr.Issues)
	}
}

func TestValidateReportsMissingBranchMapping(t *testing.T) {
	registry := mustRegistry(t, []manifest.ModuleSpec{moduleSpec("foundation")})
	delete(registry.byBranch[name("foundation")], branch("main"))

	validationErr := mustValidationError(t, registry.Validate())
	if !hasIssueIndex(validationErr.Issues, "byBranch[foundation]") {
		t.Fatalf("issues = %#v, want byBranch[foundation] issue", validationErr.Issues)
	}
}
