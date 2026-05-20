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

package registry

import (
	"errors"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

type moduleSpecOption func(*manifest.ModuleSpec)

func specs(values ...manifest.ModuleSpec) []manifest.ModuleSpec {
	return append([]manifest.ModuleSpec(nil), values...)
}

func moduleSpec(moduleName string, options ...moduleSpecOption) manifest.ModuleSpec {
	spec := manifest.ModuleSpec{
		Name:       moduleName,
		ModulePath: "arcoris.dev/" + moduleName,
		SourceDir:  "staging/src/arcoris.dev/" + moduleName,
		Repository: "arcoris/" + moduleName,
		Branches: []manifest.BranchMappingSpec{
			{Source: "main", Target: "main"},
		},
	}
	for _, option := range options {
		option(&spec)
	}
	return spec
}

func withDependency(dependency string) moduleSpecOption {
	return func(spec *manifest.ModuleSpec) {
		spec.Dependencies = append(spec.Dependencies, dependency)
	}
}

func withVisibility(visibility string) moduleSpecOption {
	return func(spec *manifest.ModuleSpec) {
		spec.Visibility = visibility
	}
}

func withBranch(source string, target string) moduleSpecOption {
	return func(spec *manifest.ModuleSpec) {
		if len(spec.Branches) == 1 && spec.Branches[0].Source == "main" && spec.Branches[0].Target == "main" {
			spec.Branches = nil
		}
		spec.Branches = append(spec.Branches, manifest.BranchMappingSpec{Source: source, Target: target})
	}
}

func mustModule(t *testing.T, spec manifest.ModuleSpec) manifest.Module {
	t.Helper()
	module, err := manifest.NewModule(spec)
	if err != nil {
		t.Fatalf("NewModule(%#v) error = %v", spec, err)
	}
	return module
}

func mustBranchMapping(t *testing.T, source string, target string) manifest.BranchMapping {
	t.Helper()
	mapping, err := manifest.NewBranchMapping(manifest.BranchMappingSpec{Source: source, Target: target})
	if err != nil {
		t.Fatalf("NewBranchMapping(%q, %q) error = %v", source, target, err)
	}
	return mapping
}

func mustRegistry(t *testing.T, specs []manifest.ModuleSpec) Registry {
	t.Helper()
	modules := make([]manifest.Module, 0, len(specs))
	for _, spec := range specs {
		modules = append(modules, mustModule(t, spec))
	}
	registry, err := New(modules)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return registry
}

func mustValidationError(t *testing.T, err error) *ValidationError {
	t.Helper()
	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error = %T %v, want *ValidationError", err, err)
	}
	return validationErr
}

func name(value string) manifest.ModuleName {
	name, err := manifest.ParseModuleName(value)
	if err != nil {
		panic(err)
	}
	return name
}

func modulePath(value string) manifest.ModulePath {
	modulePath, err := manifest.ParseModulePath(value)
	if err != nil {
		panic(err)
	}
	return modulePath
}

func sourceDir(value string) manifest.SourceDir {
	sourceDir, err := manifest.ParseSourceDir(value)
	if err != nil {
		panic(err)
	}
	return sourceDir
}

func repository(value string) manifest.RepositoryRef {
	repository, err := manifest.ParseRepositoryRef(value)
	if err != nil {
		panic(err)
	}
	return repository
}

func branch(value string) manifest.BranchName {
	branch, err := manifest.ParseBranchName(value)
	if err != nil {
		panic(err)
	}
	return branch
}

func assertModules(t *testing.T, got []manifest.Module, want ...string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("modules len = %d, want %d: got %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i].Name() != name(want[i]) {
			t.Fatalf("modules[%d] = %q, want %q; full = %v", i, got[i].Name(), want[i], got)
		}
	}
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

func assertNames(t *testing.T, got []manifest.ModuleName, want ...string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("names len = %d, want %d: got %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != name(want[i]) {
			t.Fatalf("names[%d] = %q, want %q; full = %v", i, got[i], want[i], got)
		}
	}
}

func assertModulePaths(t *testing.T, got []manifest.ModulePath, want ...string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("module paths len = %d, want %d: got %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != modulePath(want[i]) {
			t.Fatalf("module paths[%d] = %q, want %q; full = %v", i, got[i], want[i], got)
		}
	}
}

func assertSourceDirs(t *testing.T, got []manifest.SourceDir, want ...string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("source dirs len = %d, want %d: got %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != sourceDir(want[i]) {
			t.Fatalf("source dirs[%d] = %q, want %q; full = %v", i, got[i], want[i], got)
		}
	}
}

func assertRepositories(t *testing.T, got []manifest.RepositoryRef, want ...string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("repositories len = %d, want %d: got %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != repository(want[i]) {
			t.Fatalf("repositories[%d] = %q, want %q; full = %v", i, got[i], want[i], got)
		}
	}
}

func assertBranches(t *testing.T, got []manifest.BranchName, want ...string) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("branches len = %d, want %d: got %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != branch(want[i]) {
			t.Fatalf("branches[%d] = %q, want %q; full = %v", i, got[i], want[i], got)
		}
	}
}
