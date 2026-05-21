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

package staging_test

import (
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
)

func TestValidateAcceptsUniqueModuleNamesAndSourceDirs(t *testing.T) {
	m, err := staging.New(validSpec())
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}
	if err := m.Validate(); err != nil {
		t.Fatalf("Validate returned error: %v", err)
	}
}

func TestValidateRejectsDuplicateModuleNames(t *testing.T) {
	spec := validSpec()
	spec.Modules[1].Name = "foundation"
	if _, err := staging.New(spec); err == nil {
		t.Fatalf("expected duplicate module name error")
	}
}

func TestValidateRejectsDuplicateSourceDirs(t *testing.T) {
	spec := validSpec()
	spec.Modules[1].SourceDir = spec.Modules[0].SourceDir
	if _, err := staging.New(spec); err == nil {
		t.Fatalf("expected duplicate sourceDir error")
	}
}
