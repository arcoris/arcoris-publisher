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
	"testing"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

func TestValidateSucceedsForBuiltRegistry(t *testing.T) {
	registry := mustRegistry(t, []manifest.ModuleSpec{moduleSpec("foundation")})
	if err := registry.Validate(); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestValidateZeroRegistrySucceeds(t *testing.T) {
	var registry Registry
	if err := registry.Validate(); err != nil {
		t.Fatalf("zero registry Validate() error = %v", err)
	}
}

func TestValidateReportsDuplicateModules(t *testing.T) {
	module := mustModule(t, moduleSpec("foundation"))
	registry := Registry{modules: []manifest.Module{module, module}}

	validationErr := mustValidationError(t, registry.Validate())
	if !hasIssueCode(validationErr.Issues, IssueDuplicateModuleName) {
		t.Fatalf("issues = %#v, want duplicate module name", validationErr.Issues)
	}
}
