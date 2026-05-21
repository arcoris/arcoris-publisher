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

func TestTopologicalOrderLinear(t *testing.T) {
	g := mustGraph(t,
		testModule{name: "foundation"},
		testModule{name: "control", dependencies: []string{"foundation"}},
		testModule{name: "scheduler", dependencies: []string{"control"}},
	)
	order, err := g.TopologicalOrder()
	if err != nil {
		t.Fatalf("TopologicalOrder() error = %v", err)
	}
	assertNames(t, order, "foundation", "control", "scheduler")
}

func TestTopologicalOrderDiamondIsDeterministic(t *testing.T) {
	g := mustGraph(t,
		testModule{name: "foundation"},
		testModule{name: "control", dependencies: []string{"foundation"}},
		testModule{name: "runtime", dependencies: []string{"foundation"}},
		testModule{name: "scheduler", dependencies: []string{"control", "runtime"}},
	)
	order, err := g.TopologicalOrder()
	if err != nil {
		t.Fatalf("TopologicalOrder() error = %v", err)
	}
	assertNames(t, order, "foundation", "control", "runtime", "scheduler")
}

func TestPublishOrderOmitsInternalModules(t *testing.T) {
	g := mustGraph(t,
		testModule{name: "foundation"},
		testModule{name: "tooling", visibility: string(manifest.VisibilityInternal)},
		testModule{name: "control", dependencies: []string{"foundation"}},
	)
	order, err := g.PublishOrder()
	if err != nil {
		t.Fatalf("PublishOrder() error = %v", err)
	}
	assertNames(t, order, "foundation", "control")
}

func TestTopologicalOrderReportsCycle(t *testing.T) {
	g := mustGraphWithCycle(t)
	_, err := g.TopologicalOrder()
	if err == nil {
		t.Fatalf("TopologicalOrder() error = nil, want cycle error")
	}
	var validation *ValidationError
	if !errors.As(err, &validation) || !validation.Has(IssueDependencyCycle) {
		t.Fatalf("error = %v, want dependency cycle validation error", err)
	}
}

func TestPublishOrderReportsCycle(t *testing.T) {
	g := mustGraphWithCycle(t)
	_, err := g.PublishOrder()
	if err == nil {
		t.Fatal("PublishOrder() error = nil, want cycle error")
	}
}
