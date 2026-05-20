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

package manifest

import "testing"

func validModuleSpec() ModuleSpec {
	return ModuleSpec{
		Name:         "control",
		ModulePath:   "arcoris.dev/control",
		SourceDir:    "staging/src/arcoris.dev/control",
		Repository:   "arcoris/control",
		Branches:     []BranchMappingSpec{{Source: "main", Target: "main"}},
		Dependencies: []string{"foundation"},
		Visibility:   "public",
	}
}

func TestNewModuleAcceptsValidSpec(t *testing.T) {
	module, err := NewModule(validModuleSpec())
	if err != nil {
		t.Fatalf("NewModule() error = %v", err)
	}
	if module.Name() != ModuleName("control") {
		t.Fatalf("Name() = %q", module.Name())
	}
	if !module.Publishable() {
		t.Fatalf("Publishable() = false, want true")
	}
	if module.Repository().Owner() != "arcoris" || module.Repository().Name() != "control" {
		t.Fatalf("repository owner/name = %q/%q", module.Repository().Owner(), module.Repository().Name())
	}
}

func TestNewModuleDefaultsVisibilityToPublic(t *testing.T) {
	spec := validModuleSpec()
	spec.Visibility = ""
	module, err := NewModule(spec)
	if err != nil {
		t.Fatalf("NewModule() error = %v", err)
	}
	if module.Visibility() != VisibilityPublic {
		t.Fatalf("Visibility() = %q, want %q", module.Visibility(), VisibilityPublic)
	}
}

func TestNewModuleRejectsSelfDependency(t *testing.T) {
	spec := validModuleSpec()
	spec.Dependencies = []string{"control"}
	if _, err := NewModule(spec); err == nil {
		t.Fatalf("NewModule() error = nil, want error")
	}
}

func TestNewModuleRejectsDuplicateDependency(t *testing.T) {
	spec := validModuleSpec()
	spec.Dependencies = []string{"foundation", "foundation"}
	if _, err := NewModule(spec); err == nil {
		t.Fatalf("NewModule() error = nil, want error")
	}
}

func TestModuleAccessorsReturnDetachedSlices(t *testing.T) {
	module := Must(validSpec()).Modules()[1]
	branches := module.Branches()
	branches[0] = BranchMapping{}
	if got := module.Branches()[0].Source(); got != BranchName("main") {
		t.Fatalf("Branches() returned attached slice: %q", got)
	}
	deps := module.Dependencies()
	deps[0] = Dependency{}
	if got := module.Dependencies()[0].Module(); got != ModuleName("foundation") {
		t.Fatalf("Dependencies() returned attached slice: %q", got)
	}
}

func TestModuleValueAccessors(t *testing.T) {
	module := Must(validSpec()).Modules()[1]
	if module.ModulePath() != ModulePath("arcoris.dev/control") {
		t.Fatalf("ModulePath() = %q", module.ModulePath())
	}
	if module.SourceDir() != SourceDir("staging/src/arcoris.dev/control") {
		t.Fatalf("SourceDir() = %q", module.SourceDir())
	}
}

func TestNewModuleRejectsInvalidFields(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*ModuleSpec)
	}{
		{name: "name", mutate: func(s *ModuleSpec) { s.Name = "Control" }},
		{name: "module path", mutate: func(s *ModuleSpec) { s.ModulePath = "../control" }},
		{name: "source dir", mutate: func(s *ModuleSpec) { s.SourceDir = "../control" }},
		{name: "repository", mutate: func(s *ModuleSpec) { s.Repository = "arcoris" }},
		{name: "branch", mutate: func(s *ModuleSpec) { s.Branches = []BranchMappingSpec{{Source: "main branch", Target: "main"}} }},
		{name: "dependency", mutate: func(s *ModuleSpec) { s.Dependencies = []string{"Bad"} }},
		{name: "visibility", mutate: func(s *ModuleSpec) { s.Visibility = "private" }},
		{name: "empty branches", mutate: func(s *ModuleSpec) { s.Branches = nil }},
		{name: "duplicate branch", mutate: func(s *ModuleSpec) {
			s.Branches = []BranchMappingSpec{{Source: "main", Target: "main"}, {Source: "main", Target: "other"}}
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := validModuleSpec()
			tt.mutate(&spec)
			if _, err := NewModule(spec); err == nil {
				t.Fatalf("NewModule() error = nil, want error")
			}
		})
	}
}
