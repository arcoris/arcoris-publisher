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

package internal_test

import (
	"context"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/config"
	"arcoris.dev/arcoris-publisher/internal/graph"
	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/plan"
	"arcoris.dev/arcoris-publisher/internal/registry"
	"arcoris.dev/arcoris-publisher/internal/versioning"
)

func TestStaticPipelineBuildsDependencyOrderedPlan(t *testing.T) {
	set, err := config.NewLoader(config.LoaderOptions{}).LoadPublicationSet(
		context.Background(),
		"config/testdata/minimal/arcpub.yaml",
	)
	if err != nil {
		t.Fatalf("LoadPublicationSet() error = %v", err)
	}

	reg, err := registry.New(set)
	if err != nil {
		t.Fatalf("registry.New() error = %v", err)
	}
	if !reg.ContainsName("foundation") || !reg.ContainsName("control") {
		t.Fatalf("registry missing expected modules: %v", reg.Modules())
	}

	deps, err := graph.New(reg)
	if err != nil {
		t.Fatalf("graph.New() error = %v", err)
	}
	order, err := deps.PublishOrder()
	if err != nil {
		t.Fatalf("PublishOrder() error = %v", err)
	}
	assertNames(t, order, "foundation", "control")

	assignments, err := versioning.Assign(versioning.Request{
		Set:      set,
		Registry: reg,
		Graph:    deps,
		Version:  versioning.Must("v0.3.0"),
	})
	if err != nil {
		t.Fatalf("versioning.Assign() error = %v", err)
	}

	publicationPlan, err := plan.Build(plan.Request{
		Set:         set,
		Registry:    reg,
		Graph:       deps,
		Assignments: assignments,
	})
	if err != nil {
		t.Fatalf("plan.Build() error = %v", err)
	}
	assertNames(t, publicationPlan.ModuleNames(), "foundation", "control")

	control, ok := publicationPlan.ModuleByName("control")
	if !ok {
		t.Fatal("control plan missing")
	}
	requirements := control.Requirements()
	if len(requirements) != 1 || requirements[0].ModulePath() != "arcoris.dev/foundation" {
		t.Fatalf("control requirements = %#v", requirements)
	}

	for _, module := range publicationPlan.Modules() {
		for _, entry := range module.PublishEntries() {
			if entry.Kind() != manifest.PublishEntryFile &&
				entry.Kind() != manifest.PublishEntryDirectory {
				t.Fatalf("unexpected publish entry kind %q", entry.Kind())
			}
		}
	}
}

func assertNames(t *testing.T, got []manifest.ModuleName, want ...manifest.ModuleName) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("names = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("names[%d] = %q, want %q; all = %v", i, got[i], want[i], got)
		}
	}
}
