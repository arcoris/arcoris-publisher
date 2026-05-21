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
	"strings"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/domain/graph"
	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
	"arcoris.dev/arcoris-publisher/internal/domain/registry"
	"arcoris.dev/arcoris-publisher/internal/domain/versioning"
)

func TestNewRejectsMissingVersion(t *testing.T) {
	manifestValue, registryValue, graphValue, _ := testInputs(t)
	if _, err := New(manifestValue, registryValue, graphValue, versioning.Assignments{}); err == nil {
		t.Fatal("New returned nil error, want missing version error")
	}
}

func TestNewRejectsPublicDependencyOnInternalModule(t *testing.T) {
	spec := testSpec()
	spec.Modules[1].Dependencies = []string{"internal-tools"}
	manifestValue, err := manifest.New(spec)
	if err != nil {
		t.Fatalf("manifest.New returned error: %v", err)
	}
	registryValue, err := registry.FromManifest(manifestValue)
	if err != nil {
		t.Fatalf("registry.FromManifest returned error: %v", err)
	}
	graphValue, err := graph.FromManifest(manifestValue)
	if err != nil {
		t.Fatalf("graph.FromManifest returned error: %v", err)
	}
	assignments, err := versioning.ReleaseTrain(registryValue, versioning.MustVersion("v0.3.0"))
	if err != nil {
		t.Fatalf("versioning.ReleaseTrain returned error: %v", err)
	}
	_, err = New(manifestValue, registryValue, graphValue, assignments)
	if err == nil {
		t.Fatal("New returned nil error, want error")
	}
	if !strings.Contains(err.Error(), "publishable graph") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateReportsManualPlanIssues(t *testing.T) {
	planValue := Plan{assignments: versioning.Assignments{}, modules: []ModulePlan{{action: ActionPublish}}}
	if err := planValue.Validate(); err == nil {
		t.Fatal("Validate returned nil error, want error")
	}
}

func TestValidationErrorFormatting(t *testing.T) {
	if got := (*ValidationError)(nil).Error(); got != "publish plan validation failed" {
		t.Fatalf("unexpected nil error string: %q", got)
	}
	one := (&ValidationError{Issues: []Issue{{Message: "one"}}}).Error()
	if one != "one" {
		t.Fatalf("unexpected one issue string: %q", one)
	}
	many := (&ValidationError{Issues: []Issue{{Message: "one"}, {Message: "two"}}}).Error()
	if !strings.Contains(many, "2 issues") {
		t.Fatalf("unexpected many issue string: %q", many)
	}
}

func TestValidateReportsExtraAssignment(t *testing.T) {
	planValue := testPlan(t)
	modules := planValue.Modules()[:1]
	partial := Plan{
		modules:     modules,
		assignments: planValue.Versions(),
		byModule:    map[manifest.ModuleName]int{modules[0].Name(): 0},
		byPath:      map[manifest.ModulePath]int{modules[0].ModulePath(): 0},
		byRepo:      map[manifest.RepositoryRef]int{modules[0].Repository(): 0},
	}
	err := partial.Validate()
	if err == nil {
		t.Fatal("Validate returned nil error, want extra assignment error")
	}
	validation, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("error type = %T, want *ValidationError", err)
	}
	found := false
	for _, issue := range validation.Issues {
		if issue.Code == IssueExtraVersion {
			found = true
		}
	}
	if !found {
		t.Fatalf("extra version issue not found: %#v", validation.Issues)
	}
}
