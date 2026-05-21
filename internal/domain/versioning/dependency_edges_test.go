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

package versioning

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

func TestDependencyVersionsFailForUnknownDependency(t *testing.T) {
	module := mustStandaloneModule(t, manifest.ModuleSpec{
		Name:         "public",
		ModulePath:   "arcoris.dev/public",
		SourceDir:    "staging/src/arcoris.dev/public",
		Repository:   "arcoris/public",
		Branches:     []manifest.BranchMappingSpec{{Source: "main", Target: "main"}},
		Dependencies: []string{"missing"},
	})
	assignments, err := ReleaseTrain(testRegistry(t), MustVersion("v0.2.0"))
	if err != nil {
		t.Fatalf("ReleaseTrain() error = %v", err)
	}

	_, err = assignments.DependencyVersions(testRegistry(t), module)
	validationErr := mustValidationError(t, err)
	if !hasIssueCode(validationErr.Issues, IssueUnknownDependency) {
		t.Fatalf("issues = %#v, want unknown dependency", validationErr.Issues)
	}
}

func TestRequirementMapPropagatesDependencyError(t *testing.T) {
	module := mustStandaloneModule(t, manifest.ModuleSpec{
		Name:         "public",
		ModulePath:   "arcoris.dev/public",
		SourceDir:    "staging/src/arcoris.dev/public",
		Repository:   "arcoris/public",
		Branches:     []manifest.BranchMappingSpec{{Source: "main", Target: "main"}},
		Dependencies: []string{"missing"},
	})
	assignments, err := ReleaseTrain(testRegistry(t), MustVersion("v0.2.0"))
	if err != nil {
		t.Fatalf("ReleaseTrain() error = %v", err)
	}

	if _, err := assignments.RequirementMap(testRegistry(t), module); err == nil {
		t.Fatalf("RequirementMap() succeeded, expected error")
	}
}

func mustStandaloneModule(t *testing.T, spec manifest.ModuleSpec) manifest.Module {
	t.Helper()
	module, err := manifest.NewModule(spec)
	if err != nil {
		t.Fatalf("NewModule() error = %v", err)
	}
	return module
}
