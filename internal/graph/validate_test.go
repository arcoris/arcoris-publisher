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

package graph

import (
	"errors"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestValidationErrorHelpers(t *testing.T) {
	err := &ValidationError{Issues: []Issue{{Code: IssueUnknownNode, Path: "x", Message: "missing"}}}
	if err.Error() == "" {
		t.Fatalf("Error() is empty")
	}
	if !err.Has(IssueUnknownNode) {
		t.Fatalf("Has(unknown_node) = false")
	}
	if err.Has(IssueDependencyCycle) {
		t.Fatalf("Has(dependency_cycle) = true")
	}
	if err.Empty() {
		t.Fatalf("Empty() = true")
	}
	if (&ValidationError{}).Error() == "" {
		t.Fatalf("empty Error() is empty")
	}
	if !((*ValidationError)(nil)).Empty() {
		t.Fatalf("nil Empty() = false")
	}
	if ((*ValidationError)(nil)).Has(IssueUnknownNode) {
		t.Fatalf("nil Has() = true")
	}
}

func TestIssueError(t *testing.T) {
	issue := Issue{
		Code:    IssueUnknownDependency,
		Path:    "modules[0].dependencies[0]",
		Message: "unknown dependency",
	}
	if issue.Error() == "" {
		t.Fatalf("Issue.Error() is empty")
	}
	issue.Path = ""
	if issue.Error() == "" {
		t.Fatalf("Issue.Error() without path is empty")
	}
}

func TestNewRejectsSelfDependencyWhenInputBypassesResolvedValidation(t *testing.T) {
	g := Graph{
		order: []manifest.ModuleName{"foundation"},
		nodes: map[manifest.ModuleName]Node{
			"foundation": {
				name:       "foundation",
				modulePath: "arcoris.dev/foundation",
				visibility: manifest.VisibilityPublic,
			},
		},
		dependencies: map[manifest.ModuleName][]manifest.ModuleName{"foundation": {"foundation"}},
		dependents:   map[manifest.ModuleName][]manifest.ModuleName{"foundation": {"foundation"}},
	}
	if !g.HasCycle() {
		t.Fatalf("manual self-cycle HasCycle() = false")
	}
}

func TestNewRejectsDuplicateNodesWhenRegistryIsCorrupt(t *testing.T) {
	_, err := New(registryWithDuplicateModule(t))
	if err == nil {
		t.Fatal("expected duplicate node error")
	}

	var validation *ValidationError
	if !errors.As(err, &validation) || !validation.Has(IssueDuplicateNode) {
		t.Fatalf("error = %v, want duplicate node validation error", err)
	}
}

func TestBuilderAddDependencyDefensiveErrors(t *testing.T) {
	reg := mustRegistry(t,
		testModule{name: "foundation"},
		testModule{name: "control"},
		testModule{name: "disabled", visibility: string(manifest.VisibilityDisabled)},
	)
	builder := builder{
		registry:     reg,
		nodes:        map[manifest.ModuleName]Node{"foundation": {}},
		dependencies: map[manifest.ModuleName][]manifest.ModuleName{"foundation": nil},
		dependents:   map[manifest.ModuleName][]manifest.ModuleName{"foundation": nil},
	}

	builder.addDependency("self", "foundation", "foundation")
	builder.addDependency("unknown", "foundation", "missing")
	builder.addDependency("disabled", "foundation", "disabled")
	builder.addDependency("absent", "foundation", "control")

	err := builder.issues.Err()
	if err == nil {
		t.Fatal("expected validation errors")
	}

	var validation *ValidationError
	if !errors.As(err, &validation) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	for _, code := range []IssueCode{
		IssueSelfDependency,
		IssueUnknownDependency,
		IssueDisabledDependency,
	} {
		if !validation.Has(code) {
			t.Fatalf("missing issue %s in %#v", code, validation.Issues)
		}
	}
}

func TestFromPublicationSetAndMust(t *testing.T) {
	set := mustPublicationSet(t, testModule{name: "foundation"})

	g, err := FromPublicationSet(set)
	if err != nil {
		t.Fatalf("FromPublicationSet() error = %v", err)
	}
	if g.Len() != 1 {
		t.Fatalf("Len() = %d, want 1", g.Len())
	}

	reg := mustRegistry(t, testModule{name: "foundation"})
	if Must(reg).Len() != 1 {
		t.Fatal("Must() returned unexpected graph")
	}
}

func TestFromPublicationSetReportsRegistryError(t *testing.T) {
	_, err := FromPublicationSet(publicationSetWithDuplicateModule(t))
	if err == nil {
		t.Fatal("expected registry error")
	}
}

func TestMustPanicsOnInvalidRegistry(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()

	Must(registryWithDuplicateModule(t))
}
