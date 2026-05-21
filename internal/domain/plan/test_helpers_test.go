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
	"errors"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/domain/graph"
	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
	"arcoris.dev/arcoris-publisher/internal/domain/registry"
	"arcoris.dev/arcoris-publisher/internal/domain/versioning"
)

func testSpec() manifest.Spec {
	return manifest.Spec{
		Version: "v1",
		Source:  manifest.SourceSpec{Repository: "arcoris/arcoris", DefaultBranch: "main"},
		Policy:  manifest.PolicySpec{VersionPolicy: "release-train", PushPolicy: "fast-forward-only"},
		Modules: []manifest.ModuleSpec{
			{
				Name:       "foundation",
				ModulePath: "arcoris.dev/foundation",
				SourceDir:  "staging/src/arcoris.dev/foundation",
				Repository: "arcoris/foundation",
				Branches:   []manifest.BranchMappingSpec{{Source: "main", Target: "main"}},
			},
			{
				Name:         "control",
				ModulePath:   "arcoris.dev/control",
				SourceDir:    "staging/src/arcoris.dev/control",
				Repository:   "arcoris/control",
				Branches:     []manifest.BranchMappingSpec{{Source: "main", Target: "main"}, {Source: "release/v1", Target: "release/v1"}},
				Dependencies: []string{"foundation"},
			},
			{
				Name:       "internal-tools",
				ModulePath: "arcoris.dev/internal-tools",
				SourceDir:  "staging/src/arcoris.dev/internal-tools",
				Repository: "arcoris/internal-tools",
				Branches:   []manifest.BranchMappingSpec{{Source: "main", Target: "main"}},
				Visibility: "internal",
			},
			{
				Name:       "old-module",
				ModulePath: "arcoris.dev/old-module",
				SourceDir:  "staging/src/arcoris.dev/old-module",
				Repository: "arcoris/old-module",
				Branches:   []manifest.BranchMappingSpec{{Source: "main", Target: "main"}},
				Visibility: "disabled",
			},
		},
	}
}

func testManifest(t *testing.T) manifest.Manifest {
	t.Helper()
	value, err := manifest.New(testSpec())
	if err != nil {
		t.Fatalf("manifest.New returned error: %v", err)
	}
	return value
}

func testInputs(t *testing.T) (manifest.Manifest, registry.Registry, graph.Graph, versioning.Assignments) {
	t.Helper()
	manifestValue := testManifest(t)
	return testInputsFromManifest(t, manifestValue)
}

func testInputsFromManifest(t *testing.T, manifestValue manifest.Manifest) (manifest.Manifest, registry.Registry, graph.Graph, versioning.Assignments) {
	t.Helper()
	registryValue, err := registry.FromManifest(manifestValue)
	if err != nil {
		t.Fatalf("registry.FromManifest returned error: %v", err)
	}
	graphValue, err := graph.FromManifest(manifestValue)
	if err != nil {
		t.Fatalf("graph.FromManifest returned error: %v", err)
	}
	assignments, err := versioning.ReleaseTrain(registryValue, versioning.MustVersion("v0.3.0"))
	if err != nil {
		t.Fatalf("versioning.ReleaseTrain returned error: %v", err)
	}
	return manifestValue, registryValue, graphValue, assignments
}

func testInputsFromSpec(t *testing.T, spec manifest.Spec) (manifest.Manifest, registry.Registry, graph.Graph, versioning.Assignments) {
	t.Helper()
	manifestValue, err := manifest.New(spec)
	if err != nil {
		t.Fatalf("manifest.New returned error: %v", err)
	}
	return testInputsFromManifest(t, manifestValue)
}

func testPlan(t *testing.T) Plan {
	t.Helper()
	manifestValue, registryValue, graphValue, assignments := testInputs(t)
	planValue, err := New(manifestValue, registryValue, graphValue, assignments)
	if err != nil {
		t.Fatalf("plan.New returned error: %v", err)
	}
	return planValue
}

func moduleName(t *testing.T, value string) manifest.ModuleName {
	t.Helper()
	name, err := manifest.ParseModuleName(value)
	if err != nil {
		t.Fatalf("ParseModuleName returned error: %v", err)
	}
	return name
}

func modulePath(t *testing.T, value string) manifest.ModulePath {
	t.Helper()
	path, err := manifest.ParseModulePath(value)
	if err != nil {
		t.Fatalf("ParseModulePath returned error: %v", err)
	}
	return path
}

func repositoryRef(t *testing.T, value string) manifest.RepositoryRef {
	t.Helper()
	repo, err := manifest.ParseRepositoryRef(value)
	if err != nil {
		t.Fatalf("ParseRepositoryRef returned error: %v", err)
	}
	return repo
}

func branchName(t *testing.T, value string) manifest.BranchName {
	t.Helper()
	branch, err := manifest.ParseBranchName(value)
	if err != nil {
		t.Fatalf("ParseBranchName returned error: %v", err)
	}
	return branch
}

func mustValidationError(t *testing.T, err error) *ValidationError {
	t.Helper()
	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %T %v, want *ValidationError", err, err)
	}
	return validationErr
}

func hasIssueCode(issues []Issue, code IssueCode) bool {
	for _, issue := range issues {
		if issue.Code == code {
			return true
		}
	}
	return false
}

func hasIssueIndex(issues []Issue, index string) bool {
	for _, issue := range issues {
		if issue.Index == index {
			return true
		}
	}
	return false
}
