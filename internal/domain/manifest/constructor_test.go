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

func TestNewAcceptsValidSpec(t *testing.T) {
	manifest, err := New(validSpec())
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if got := manifest.Version(); got != VersionV1 {
		t.Fatalf("Version() = %q, want %q", got, VersionV1)
	}
	if got := len(manifest.Modules()); got != 2 {
		t.Fatalf("len(Modules()) = %d, want 2", got)
	}
	control, ok := manifest.ModuleByName(ModuleName("control"))
	if !ok {
		t.Fatalf("ModuleByName(control) not found")
	}
	if deps := control.Dependencies(); len(deps) != 1 || deps[0].Module() != ModuleName("foundation") {
		t.Fatalf("control dependencies = %#v", deps)
	}
}

func TestMustPanicsOnInvalidSpec(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("Must() did not panic")
		}
	}()
	spec := validSpec()
	spec.Version = "v9"
	_ = Must(spec)
}

func TestNewRejectsInvalidVersion(t *testing.T) {
	spec := validSpec()
	spec.Version = "v2"
	if _, err := New(spec); err == nil {
		t.Fatalf("New() error = nil, want error")
	}
}

func TestNewRejectsInvalidSource(t *testing.T) {
	spec := validSpec()
	spec.Source.Repository = "arcoris"
	if _, err := New(spec); err == nil {
		t.Fatalf("New() error = nil, want error")
	}
}

func TestNewRejectsInvalidPolicy(t *testing.T) {
	spec := validSpec()
	spec.Policy.PushPolicy = "force"
	if _, err := New(spec); err == nil {
		t.Fatalf("New() error = nil, want error")
	}
}
