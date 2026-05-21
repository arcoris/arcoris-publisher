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

func TestAffectedBy(t *testing.T) {
	g := mustGraph(t,
		testModule{name: "foundation"},
		testModule{name: "control", dependencies: []string{"foundation"}},
		testModule{name: "runtime", dependencies: []string{"foundation"}},
		testModule{name: "scheduler", dependencies: []string{"control", "runtime"}},
	)
	affected, err := g.AffectedBy(manifest.ModuleName("foundation"))
	if err != nil {
		t.Fatalf("AffectedBy() error = %v", err)
	}
	assertNames(t, affected, "foundation", "control", "runtime", "scheduler")
}

func TestPublishClosure(t *testing.T) {
	g := mustGraph(t,
		testModule{name: "foundation"},
		testModule{name: "control", dependencies: []string{"foundation"}},
		testModule{name: "runtime", dependencies: []string{"foundation"}},
		testModule{name: "scheduler", dependencies: []string{"control", "runtime"}},
	)
	closure, err := g.PublishClosure(manifest.ModuleName("scheduler"))
	if err != nil {
		t.Fatalf("PublishClosure() error = %v", err)
	}
	assertNames(t, closure, "foundation", "control", "runtime", "scheduler")
}

func TestPublishClosureDeduplicatesInput(t *testing.T) {
	g := mustGraph(t,
		testModule{name: "foundation"},
		testModule{name: "control", dependencies: []string{"foundation"}},
	)
	closure, err := g.PublishClosure(
		manifest.ModuleName("control"),
		manifest.ModuleName("control"),
	)
	if err != nil {
		t.Fatalf("PublishClosure() error = %v", err)
	}
	assertNames(t, closure, "foundation", "control")
}

func TestImpactUnknownNode(t *testing.T) {
	g := mustGraph(t, testModule{name: "foundation"})
	_, err := g.AffectedBy(manifest.ModuleName("missing"))
	if err == nil {
		t.Fatalf("AffectedBy(missing) error = nil")
	}
	var validation *ValidationError
	if !errors.As(err, &validation) || !validation.Has(IssueUnknownNode) {
		t.Fatalf("error = %v, want unknown node validation error", err)
	}
}

func TestImpactUnknownNodesUseCollectionPath(t *testing.T) {
	g := mustGraph(t, testModule{name: "foundation"})
	_, err := g.AffectedBy(
		manifest.ModuleName("missing-a"),
		manifest.ModuleName("missing-b"),
	)
	if err == nil {
		t.Fatal("AffectedBy(missing-a, missing-b) error = nil")
	}

	var validation *ValidationError
	if !errors.As(err, &validation) {
		t.Fatalf("error = %v, want validation error", err)
	}
	for _, issue := range validation.Issues {
		if issue.Path != "changed[]" {
			t.Fatalf("issue path = %q, want changed[]", issue.Path)
		}
	}
}

func TestPublishClosureUnknownNode(t *testing.T) {
	g := mustGraph(t, testModule{name: "foundation"})
	_, err := g.PublishClosure(manifest.ModuleName("missing"))
	if err == nil {
		t.Fatal("PublishClosure(missing) error = nil")
	}
}
