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

func TestTopologicalOrderReturnsDependencyFirstOrder(t *testing.T) {
	graph := mustGraph(t, []manifest.ModuleSpec{
		moduleSpec("scheduler", "control", "foundation"),
		moduleSpec("control", "foundation"),
		moduleSpec("foundation"),
	})

	order, err := graph.TopologicalOrder()
	if err != nil {
		t.Fatalf("TopologicalOrder() error = %v", err)
	}
	assertNames(t, order, "foundation", "control", "scheduler")
}

func TestTopologicalOrderUsesDeclarationOrderAsTieBreaker(t *testing.T) {
	graph := mustGraph(t, []manifest.ModuleSpec{
		moduleSpec("runtime"),
		moduleSpec("foundation"),
		moduleSpec("control", "foundation", "runtime"),
	})

	order, err := graph.TopologicalOrder()
	if err != nil {
		t.Fatalf("TopologicalOrder() error = %v", err)
	}
	assertNames(t, order, "runtime", "foundation", "control")
}

func TestPublishOrderAliasesTopologicalOrder(t *testing.T) {
	graph := mustGraph(t, []manifest.ModuleSpec{
		moduleSpec("foundation"),
		moduleSpec("control", "foundation"),
	})

	order, err := graph.PublishOrder()
	if err != nil {
		t.Fatalf("PublishOrder() error = %v", err)
	}
	assertNames(t, order, "foundation", "control")
}

func TestTopologicalOrderReturnsCycleError(t *testing.T) {
	graph := graphWithCycle()

	_, err := graph.TopologicalOrder()
	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("TopologicalOrder() error = %T %v, want *ValidationError", err, err)
	}
	if len(validationErr.Issues) == 0 || validationErr.Issues[0].Code != IssueCycle {
		t.Fatalf("issues = %#v, want cycle", validationErr.Issues)
	}
}
