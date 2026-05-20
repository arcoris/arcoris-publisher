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

func TestPublishableSubgraphRejectsDependencyOnInternalModule(t *testing.T) {
	graph := mustGraph(t, []manifest.ModuleSpec{
		moduleSpecWithVisibility("foundation", string(manifest.VisibilityInternal)),
		moduleSpec("control", "foundation"),
	})

	_, err := graph.PublishableSubgraph()
	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("PublishableSubgraph() error = %T %v, want *ValidationError", err, err)
	}
	if validationErr.Issues[0].Code != IssueUnknownDependency {
		t.Fatalf("issue code = %s, want %s", validationErr.Issues[0].Code, IssueUnknownDependency)
	}
}

func TestPublishableSubgraphContainsOnlyPublicModules(t *testing.T) {
	graph := mustGraph(t, []manifest.ModuleSpec{
		moduleSpec("foundation"),
		moduleSpecWithVisibility("internal-tools", string(manifest.VisibilityInternal)),
		moduleSpec("control", "foundation"),
	})

	subgraph, err := graph.PublishableSubgraph()
	if err != nil {
		t.Fatalf("PublishableSubgraph() error = %v", err)
	}
	assertNames(t, subgraph.ModuleNames(), "foundation", "control")
}
