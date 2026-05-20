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
	"testing"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

func TestTransitiveDependenciesOfReturnsDependencyFirstOrder(t *testing.T) {
	graph := mustGraph(t, []manifest.ModuleSpec{
		moduleSpec("foundation"),
		moduleSpec("control", "foundation"),
		moduleSpec("scheduler", "control", "foundation"),
	})

	assertNames(t, graph.TransitiveDependenciesOf(name("scheduler")), "foundation", "control")
}

func TestTransitiveDependentsOfReturnsDeclarationOrder(t *testing.T) {
	graph := mustGraph(t, []manifest.ModuleSpec{
		moduleSpec("foundation"),
		moduleSpec("control", "foundation"),
		moduleSpec("scheduler", "control"),
		moduleSpec("runtime", "foundation"),
	})

	assertNames(t, graph.TransitiveDependentsOf(name("foundation")), "control", "scheduler", "runtime")
}
