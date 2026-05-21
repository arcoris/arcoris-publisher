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

package resolved_test

import (
	"testing"

	modulemanifest "arcoris.dev/arcoris-publisher/internal/manifest/module"
	"arcoris.dev/arcoris-publisher/internal/manifest/resolved"
)

func TestResolveRejectsDuplicateModulePaths(t *testing.T) {
	_, err := resolved.Resolve(resolved.ResolveInput{
		Staging: stagingManifest(t, baseStagingSpec()),
		Modules: []modulemanifest.Manifest{
			moduleManifest(t, "foundation", "arcoris.dev/shared", nil),
			moduleManifest(t, "control", "arcoris.dev/shared", nil),
		},
	})

	if err == nil {
		t.Fatalf("expected duplicate module path error")
	}
}

func TestResolveRejectsDuplicatePublicRepositories(t *testing.T) {
	spec := baseStagingSpec()
	spec.Modules[1].Repository = spec.Modules[0].Repository

	_, err := resolved.Resolve(resolved.ResolveInput{
		Staging: stagingManifest(t, spec),
		Modules: []modulemanifest.Manifest{
			moduleManifest(t, "foundation", "arcoris.dev/foundation", nil),
			moduleManifest(t, "control", "arcoris.dev/control", nil),
		},
	})

	if err == nil {
		t.Fatalf("expected duplicate public repository error")
	}
}
