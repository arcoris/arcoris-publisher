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

package manifest_test

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
)

func TestParseModulePathAcceptsGoModulePath(t *testing.T) {
	got, err := manifest.ParseModulePath("arcoris.dev/control")
	if err != nil {
		t.Fatalf("ParseModulePath returned error: %v", err)
	}
	if got.String() != "arcoris.dev/control" {
		t.Fatalf("unexpected module path: %q", got.String())
	}
}

func TestParseModulePathRejectsUnsafeForms(t *testing.T) {
	for _, value := range []string{
		"",
		"control",
		" arcoris.dev/control",
		"arcoris.dev//control",
		"/arcoris.dev/control",
		"arcoris.dev/control/",
		"arcoris.dev/../control",
		"arcoris.dev\\.control",
	} {
		if _, err := manifest.ParseModulePath(value); err == nil {
			t.Fatalf("ParseModulePath(%q) returned nil error", value)
		}
	}
}
