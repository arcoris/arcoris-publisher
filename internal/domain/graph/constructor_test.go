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

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

func TestNewBuildsDependencyGraph(t *testing.T) {
	graph := mustGraph(t, []manifest.ModuleSpec{
		moduleSpec("foundation"),
		moduleSpec("control", "foundation"),
		moduleSpec("scheduler", "foundation", "control"),
	})

	if graph.Len() != 3 {
		t.Fatalf("Len() = %d, want 3", graph.Len())
	}
	assertNames(t, graph.ModuleNames(), "foundation", "control", "scheduler")
	assertNames(t, graph.DependenciesOf(name("scheduler")), "control", "foundation")
	assertNames(t, graph.DependentsOf(name("foundation")), "control", "scheduler")
	if !graph.HasDependency(name("scheduler"), name("control")) {
		t.Fatalf("expected scheduler to depend on control")
	}
	if graph.HasDependency(name("control"), name("scheduler")) {
		t.Fatalf("did not expect control to depend on scheduler")
	}
}

func TestNewReportsUnknownDependencyWithoutManifestAggregate(t *testing.T) {
	module, err := manifest.NewModule(moduleSpec("control", "foundation"))
	if err != nil {
		t.Fatalf("NewModule() error = %v", err)
	}

	_, err = New([]manifest.Module{module})
	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("New() error = %T %v, want *ValidationError", err, err)
	}
	if len(validationErr.Issues) != 1 || validationErr.Issues[0].Code != IssueUnknownDependency {
		t.Fatalf("issues = %#v, want one unknown dependency", validationErr.Issues)
	}
}
