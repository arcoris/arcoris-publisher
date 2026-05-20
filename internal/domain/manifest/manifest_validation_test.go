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

import (
	"errors"
	"testing"
)

func TestNewRejectsDuplicateModuleName(t *testing.T) {
	spec := validSpec()
	spec.Modules[1].Name = "foundation"
	spec.Modules[1].Dependencies = nil
	_, err := New(spec)
	validationErr := mustValidationError(t, err)
	if validationErr.Issues[0].Code != IssueDuplicateModule {
		t.Fatalf("issue code = %q, want %q", validationErr.Issues[0].Code, IssueDuplicateModule)
	}
}

func TestNewRejectsUnknownDependency(t *testing.T) {
	spec := validSpec()
	spec.Modules[1].Dependencies = []string{"missing"}
	_, err := New(spec)
	validationErr := mustValidationError(t, err)
	if validationErr.Issues[0].Code != IssueUnknownDependency {
		t.Fatalf("issue code = %q, want %q", validationErr.Issues[0].Code, IssueUnknownDependency)
	}
}

func TestNewRejectsEmptyModules(t *testing.T) {
	spec := validSpec()
	spec.Modules = nil
	_, err := New(spec)
	validationErr := mustValidationError(t, err)
	if validationErr.Issues[0].Path != "modules" {
		t.Fatalf("issue path = %q, want modules", validationErr.Issues[0].Path)
	}
}

func TestNewRejectsDuplicateModulePathSourceDirAndRepository(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Spec)
	}{
		{name: "module path", mutate: func(s *Spec) { s.Modules[1].ModulePath = s.Modules[0].ModulePath }},
		{name: "source dir", mutate: func(s *Spec) { s.Modules[1].SourceDir = s.Modules[0].SourceDir }},
		{name: "repository", mutate: func(s *Spec) { s.Modules[1].Repository = s.Modules[0].Repository }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := validSpec()
			tt.mutate(&spec)
			if _, err := New(spec); err == nil {
				t.Fatalf("New() error = nil, want error")
			}
		})
	}
}

func mustValidationError(t *testing.T, err error) *ValidationError {
	t.Helper()
	if err == nil {
		t.Fatalf("New() error = nil, want validation error")
	}
	var validationErr *ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("error type = %T, want *ValidationError", err)
	}
	return validationErr
}
