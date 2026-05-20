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

func TestNewReportsDuplicateIndexes(t *testing.T) {
	foundation := mustModule(t, moduleSpec("foundation"))
	duplicateName := mustModule(t, manifest.ModuleSpec{
		Name:       "foundation",
		ModulePath: "arcoris.dev/other",
		SourceDir:  "staging/src/arcoris.dev/other",
		Repository: "arcoris/other",
		Branches:   []manifest.BranchMappingSpec{{Source: "main", Target: "main"}},
	})
	duplicatePath := mustModule(t, manifest.ModuleSpec{
		Name:       "path-duplicate",
		ModulePath: "arcoris.dev/foundation",
		SourceDir:  "staging/src/arcoris.dev/path-duplicate",
		Repository: "arcoris/path-duplicate",
		Branches:   []manifest.BranchMappingSpec{{Source: "main", Target: "main"}},
	})
	duplicateDir := mustModule(t, manifest.ModuleSpec{
		Name:       "dir-duplicate",
		ModulePath: "arcoris.dev/dir-duplicate",
		SourceDir:  "staging/src/arcoris.dev/foundation",
		Repository: "arcoris/dir-duplicate",
		Branches:   []manifest.BranchMappingSpec{{Source: "main", Target: "main"}},
	})
	duplicateRepo := mustModule(t, manifest.ModuleSpec{
		Name:       "repo-duplicate",
		ModulePath: "arcoris.dev/repo-duplicate",
		SourceDir:  "staging/src/arcoris.dev/repo-duplicate",
		Repository: "arcoris/foundation",
		Branches:   []manifest.BranchMappingSpec{{Source: "main", Target: "main"}},
	})

	_, err := New([]manifest.Module{foundation, duplicateName, duplicatePath, duplicateDir, duplicateRepo})
	validationErr := mustValidationError(t, err)
	wantCodes := []IssueCode{IssueDuplicateModuleName, IssueDuplicateModulePath, IssueDuplicateSourceDir, IssueDuplicateRepository}
	for _, code := range wantCodes {
		if !hasIssueCode(validationErr.Issues, code) {
			t.Fatalf("issues = %#v, want code %s", validationErr.Issues, code)
		}
	}
}

func TestBuilderReportsDuplicateBranchMappings(t *testing.T) {
	builder := newBuilder(nil)
	builder.registry = newEmptyRegistry(nil)

	builder.indexBranchMappings(name("foundation"), []manifest.BranchMapping{
		mustBranchMapping(t, "main", "main"),
		mustBranchMapping(t, "main", "published-main"),
	})

	if !hasIssueCode(builder.issues, IssueDuplicateBranchMapping) {
		t.Fatalf("issues = %#v, want duplicate branch mapping", builder.issues)
	}
}
