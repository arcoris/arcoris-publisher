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
	"testing"

	"arcoris.dev/arcoris-publisher/internal/graph"
	"arcoris.dev/arcoris-publisher/internal/registry"
	"arcoris.dev/arcoris-publisher/internal/versioning"
)

func TestBuildRejectsMissingAssignment(t *testing.T) {
	req := mustRequest(t, "v0.3.0", testModule{name: "foundation"})
	req.Assignments = versioning.Assignments{}
	_, err := Build(req)
	assertPlanError(t, err, IssueMissingAssignment)
}

func TestBuildRejectsGraphCycle(t *testing.T) {
	set := mustPublicationSet(t,
		testModule{name: "a", dependencies: []string{"c"}},
		testModule{name: "b", dependencies: []string{"a"}},
		testModule{name: "c", dependencies: []string{"b"}},
	)
	reg, err := registry.New(set)
	if err != nil {
		t.Fatalf("registry.New() error = %v", err)
	}
	g, err := graph.New(reg)
	if err != nil {
		t.Fatalf("graph.New() error = %v", err)
	}
	_, err = Build(Request{Set: set, Registry: reg, Graph: g, Assignments: versioning.Assignments{}})
	assertPlanError(t, err, IssueGraphOrder)
}

func TestBuildRejectsEmptyPlan(t *testing.T) {
	set := mustPublicationSet(t, testModule{
		name:       "helper",
		visibility: "internal",
	})
	_, err := FromPublicationSet(set, versioning.Must("v0.3.0"))
	assertPlanError(t, err, IssueEmptyPlan)
}

// assertPlanError verifies that err is a planning ValidationError containing
// code.
func assertPlanError(t *testing.T, err error, code IssueCode) {
	t.Helper()
	if err == nil {
		t.Fatalf("error = nil, want %s", code)
	}
	validation, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("error type = %T, want *ValidationError: %v", err, err)
	}
	if !validation.Has(code) {
		t.Fatalf("error %v does not contain code %s", err, code)
	}
}
