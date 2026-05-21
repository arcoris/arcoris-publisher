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

	"arcoris.dev/arcoris-publisher/internal/domain/versioning"
)

func TestNewBuildsPublishPlan(t *testing.T) {
	planValue := testPlan(t)
	if planValue.Len() != 2 {
		t.Fatalf("Len() = %d, want 2", planValue.Len())
	}
	if planValue.Empty() {
		t.Fatal("plan unexpectedly empty")
	}
	order := planValue.PublishOrder()
	if len(order) != 2 || order[0] != moduleName(t, "foundation") || order[1] != moduleName(t, "control") {
		t.Fatalf("unexpected publish order: %#v", order)
	}
	if planValue.Source().DefaultBranch() == "" {
		t.Fatal("missing source")
	}
	if planValue.Policy().VersionPolicy() == "" {
		t.Fatal("missing policy")
	}
	if planValue.Versions().Len() != 2 {
		t.Fatalf("Versions().Len() = %d, want 2", planValue.Versions().Len())
	}
}

func TestFromManifestBuildsPlan(t *testing.T) {
	manifestValue, registryValue, _, assignments := testInputs(t)
	planValue, err := FromManifest(manifestValue, assignments)
	if err != nil {
		t.Fatalf("FromManifest returned error: %v", err)
	}
	if planValue.Registry().Len() != registryValue.Len() {
		t.Fatalf("registry length = %d, want %d", planValue.Registry().Len(), registryValue.Len())
	}
}

func TestMustPanicsOnInvalidPlan(t *testing.T) {
	manifestValue, registryValue, graphValue, assignments := testInputs(t)
	emptyAssignments := versioning.Assignments{}
	defer func() {
		if recover() == nil {
			t.Fatal("Must did not panic")
		}
	}()
	_ = Must(manifestValue, registryValue, graphValue, emptyAssignments)
	_ = assignments
}

func TestMustReturnsPlanOnValidInput(t *testing.T) {
	manifestValue, registryValue, graphValue, assignments := testInputs(t)
	planValue := Must(manifestValue, registryValue, graphValue, assignments)
	if planValue.Len() != 2 {
		t.Fatalf("Len() = %d, want 2", planValue.Len())
	}
}

func TestPlanDomainAccessors(t *testing.T) {
	manifestValue, registryValue, graphValue, assignments := testInputs(t)
	planValue := Must(manifestValue, registryValue, graphValue, assignments)
	if planValue.Manifest().Version() != manifestValue.Version() {
		t.Fatalf("unexpected manifest accessor")
	}
	if planValue.Graph().Len() != graphValue.Len() {
		t.Fatalf("unexpected graph accessor")
	}
}
