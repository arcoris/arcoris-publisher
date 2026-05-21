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

package versioning

import (
	"errors"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
	"arcoris.dev/arcoris-publisher/internal/domain/registry"
)

func testRegistry(t *testing.T) registry.Registry {
	t.Helper()
	manifestValue := manifest.Must(manifest.Spec{
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
				Branches:     []manifest.BranchMappingSpec{{Source: "main", Target: "main"}},
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
				Name:       "disabled",
				ModulePath: "arcoris.dev/disabled",
				SourceDir:  "staging/src/arcoris.dev/disabled",
				Repository: "arcoris/disabled",
				Branches:   []manifest.BranchMappingSpec{{Source: "main", Target: "main"}},
				Visibility: "disabled",
			},
		},
	})
	return registry.Must(manifestValue.Modules())
}

func testModuleName(t *testing.T, name string) manifest.ModuleName {
	t.Helper()
	moduleName, err := manifest.ParseModuleName(name)
	if err != nil {
		t.Fatal(err)
	}
	return moduleName
}

func testModulePath(t *testing.T, value string) manifest.ModulePath {
	t.Helper()
	modulePath, err := manifest.ParseModulePath(value)
	if err != nil {
		t.Fatal(err)
	}
	return modulePath
}

func mustModule(t *testing.T, registryValue registry.Registry, name string) manifest.Module {
	t.Helper()
	module, ok := registryValue.ModuleByName(testModuleName(t, name))
	if !ok {
		t.Fatalf("module %q not found", name)
	}
	return module
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
