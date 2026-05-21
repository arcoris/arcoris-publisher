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
	"errors"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
	"arcoris.dev/arcoris-publisher/internal/domain/registry"
	"arcoris.dev/arcoris-publisher/internal/domain/versioning"
)

func TestNewRejectsModuleMissingFromRegistry(t *testing.T) {
	manifestValue, registryValue, graphValue, _ := testInputs(t)
	foundation, ok := registryValue.ModuleByName(moduleName(t, "foundation"))
	if !ok {
		t.Fatal("foundation not found")
	}
	subsetRegistry := registry.Must([]manifest.Module{foundation})
	assignments, err := versioning.ReleaseTrain(subsetRegistry, versioning.MustVersion("v0.3.0"))
	if err != nil {
		t.Fatalf("ReleaseTrain() error = %v", err)
	}

	_, err = New(manifestValue, subsetRegistry, graphValue, assignments)
	validationErr := mustValidationError(t, err)
	if !hasIssueCode(validationErr.Issues, IssueMissingModule) {
		t.Fatalf("issues = %#v, want missing module", validationErr.Issues)
	}
}

func TestNewRejectsModuleMissingVersion(t *testing.T) {
	manifestValue, registryValue, graphValue, _ := testInputs(t)
	foundation, ok := registryValue.ModuleByName(moduleName(t, "foundation"))
	if !ok {
		t.Fatal("foundation not found")
	}
	subsetRegistry := registry.Must([]manifest.Module{foundation})
	assignments, err := versioning.ReleaseTrain(subsetRegistry, versioning.MustVersion("v0.3.0"))
	if err != nil {
		t.Fatalf("ReleaseTrain() error = %v", err)
	}

	_, err = New(manifestValue, registryValue, graphValue, assignments)
	validationErr := mustValidationError(t, err)
	if !hasIssueCode(validationErr.Issues, IssueMissingVersion) {
		t.Fatalf("issues = %#v, want missing version", validationErr.Issues)
	}
}

func TestNewRejectsInvalidAssignments(t *testing.T) {
	manifestValue, registryValue, graphValue, _ := testInputs(t)

	_, err := New(manifestValue, registryValue, graphValue, versioning.Assignments{})
	validationErr := mustValidationError(t, err)
	if !hasIssueCode(validationErr.Issues, IssueMissingVersion) {
		t.Fatalf("issues = %#v, want missing version", validationErr.Issues)
	}
}

func TestBuilderRejectsInvalidPublishOrderGraph(t *testing.T) {
	spec := testSpec()
	spec.Modules[0].Dependencies = []string{"control"}
	manifestValue, err := manifest.New(spec)
	if err != nil {
		t.Fatalf("manifest.New() error = %v", err)
	}
	registryValue, err := registry.FromManifest(manifestValue)
	if err != nil {
		t.Fatalf("registry.FromManifest() error = %v", err)
	}
	assignments, err := versioning.ReleaseTrain(registryValue, versioning.MustVersion("v0.3.0"))
	if err != nil {
		t.Fatalf("ReleaseTrain() error = %v", err)
	}

	_, err = FromManifest(manifestValue, assignments)
	validationErr := mustValidationError(t, err)
	if !hasIssueCode(validationErr.Issues, IssueInvalidGraph) {
		t.Fatalf("issues = %#v, want invalid graph", validationErr.Issues)
	}
}

func TestBuilderRejectsExtraAssignmentDuringFinalValidation(t *testing.T) {
	manifestValue, _, graphValue, _ := testInputs(t)
	extraSpec := testSpec()
	extraSpec.Modules = append(extraSpec.Modules, manifest.ModuleSpec{
		Name:       "extra",
		ModulePath: "arcoris.dev/extra",
		SourceDir:  "staging/src/arcoris.dev/extra",
		Repository: "arcoris/extra",
		Branches:   []manifest.BranchMappingSpec{{Source: "main", Target: "main"}},
	})
	_, registryWithExtra, _, assignmentsWithExtra := testInputsFromSpec(t, extraSpec)

	_, err := New(manifestValue, registryWithExtra, graphValue, assignmentsWithExtra)
	validationErr := mustValidationError(t, err)
	if !hasIssueCode(validationErr.Issues, IssueExtraVersion) {
		t.Fatalf("issues = %#v, want extra version", validationErr.Issues)
	}
}

func TestDependencyIssuesWrapsPlainError(t *testing.T) {
	issues := dependencyIssues(moduleName(t, "control"), errors.New("plain failure"))
	if len(issues) != 1 || issues[0].Code != IssueInvalidDependency {
		t.Fatalf("issues = %#v, want invalid dependency", issues)
	}
}

func TestDependencyIssuesPassesValidationIssuesThrough(t *testing.T) {
	err := &ValidationError{Issues: []Issue{{Code: IssueInvalidDependency, Module: moduleName(t, "control"), Message: "bad"}}}
	issues := dependencyIssues(moduleName(t, "control"), err)
	if len(issues) != 1 || issues[0].Message != "bad" {
		t.Fatalf("issues = %#v", issues)
	}
}

func TestValidateReportsDuplicateNaturalKeysAfterIndexing(t *testing.T) {
	planValue := testPlan(t)
	modules := planValue.Modules()
	manual := Plan{modules: []ModulePlan{modules[0], modules[0]}}

	indexPlan(&manual)
	validationErr := mustValidationError(t, manual.Validate())
	for _, code := range []IssueCode{IssueDuplicateModule, IssueDuplicateModulePath, IssueDuplicateRepository} {
		if !hasIssueCode(validationErr.Issues, code) {
			t.Fatalf("issues = %#v, want %s", validationErr.Issues, code)
		}
	}
}

func TestBuilderModulePlanReportsInvalidBranch(t *testing.T) {
	issue := invalidBranchIssue(moduleName(t, "foundation"), errors.New("bad branch"))
	if issue.Code != IssueInvalidBranch || issue.Module != moduleName(t, "foundation") {
		t.Fatalf("issue = %#v", issue)
	}
}
