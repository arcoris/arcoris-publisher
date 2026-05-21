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

func TestParseModuleTypeAcceptsGoAndRejectsUnknownTypes(t *testing.T) {
	got, err := manifest.ParseModuleType(string(manifest.ModuleTypeGo))
	if err != nil {
		t.Fatalf("ParseModuleType returned error: %v", err)
	}
	if got != manifest.ModuleTypeGo {
		t.Fatalf("unexpected module type: %q", got)
	}
	if _, err := manifest.ParseModuleType("rust"); err == nil {
		t.Fatalf("expected invalid module type error")
	}
}
