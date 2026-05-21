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

	"arcoris.dev/arcoris-publisher/internal/domain/graph"
	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
	"arcoris.dev/arcoris-publisher/internal/domain/registry"
	"arcoris.dev/arcoris-publisher/internal/domain/versioning"
)

func TestDependencyPlansReportsUnassignedDependency(t *testing.T) {
	spec := testSpec()
	spec.Modules[1].Dependencies = []string{"internal-tools"}
	manifestValue, err := manifest.New(spec)
	if err != nil {
		t.Fatalf("manifest.New() error = %v", err)
	}
	registryValue, err := registry.FromManifest(manifestValue)
	if err != nil {
		t.Fatalf("registry.FromManifest() error = %v", err)
	}
	assignments, err := versioning.ReleaseTrain(registryValue, versioning.MustVersion("v0.3.0"))
	if err != nil {
		t.Fatalf("ReleaseTrain() error = %v", err)
	}
	module, ok := registryValue.ModuleByName(moduleName(t, "control"))
	if !ok {
		t.Fatal("control not found")
	}

	_, err = dependencyPlans(registryValue, assignments, module)
	validationErr := mustValidationError(t, err)
	if !hasIssueCode(validationErr.Issues, IssueInvalidDependency) {
		t.Fatalf("issues = %#v, want invalid dependency", validationErr.Issues)
	}
}

func TestPublishOrderRejectsInvalidPublishableGraph(t *testing.T) {
	spec := testSpec()
	spec.Modules[1].Dependencies = []string{"internal-tools"}
	manifestValue, err := manifest.New(spec)
	if err != nil {
		t.Fatalf("manifest.New() error = %v", err)
	}
	registryValue, err := registry.FromManifest(manifestValue)
	if err != nil {
		t.Fatalf("registry.FromManifest() error = %v", err)
	}
	graphValue, err := graph.FromManifest(manifestValue)
	if err != nil {
		t.Fatalf("graph.FromManifest() error = %v", err)
	}
	assignments, err := versioning.ReleaseTrain(registryValue, versioning.MustVersion("v0.3.0"))
	if err != nil {
		t.Fatalf("ReleaseTrain() error = %v", err)
	}
	builder := newBuilder(manifestValue, registryValue, graphValue, assignments)

	_, err = builder.publishOrder()
	validationErr := mustValidationError(t, err)
	if !hasIssueCode(validationErr.Issues, IssueInvalidGraph) {
		t.Fatalf("issues = %#v, want invalid graph", validationErr.Issues)
	}
}
