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

func TestParseModuleNameAcceptsLowerKebabName(t *testing.T) {
	got, err := manifest.ParseModuleName("control-plane")
	if err != nil {
		t.Fatalf("ParseModuleName returned error: %v", err)
	}
	if got.String() != "control-plane" {
		t.Fatalf("unexpected module name: %q", got.String())
	}
}

func TestParseModuleNameRejectsUnsafeNames(t *testing.T) {
	for _, value := range []string{"", "Control", "control_plane", "control "} {
		if _, err := manifest.ParseModuleName(value); err == nil {
			t.Fatalf("ParseModuleName(%q) returned nil error", value)
		}
	}
}
