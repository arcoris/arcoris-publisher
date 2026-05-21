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

func TestResolveRejectsUnknownDependency(t *testing.T) {
	_, err := resolved.Resolve(resolved.ResolveInput{
		Staging: stagingManifest(t, baseStagingSpec()),
		Modules: []modulemanifest.Manifest{
			moduleManifest(t, "foundation", "arcoris.dev/foundation", nil),
			moduleManifest(t, "control", "arcoris.dev/control", []string{"missing"}),
		},
	})

	if err == nil {
		t.Fatalf("expected unknown dependency error")
	}
}

func TestResolveRejectsSelfDependency(t *testing.T) {
	_, err := resolved.Resolve(resolved.ResolveInput{
		Staging: stagingManifest(t, baseStagingSpec()),
		Modules: []modulemanifest.Manifest{
			moduleManifest(t, "foundation", "arcoris.dev/foundation", nil),
			moduleManifest(t, "control", "arcoris.dev/control", []string{"control"}),
		},
	})

	if err == nil {
		t.Fatalf("expected self dependency error")
	}
}

func TestResolveRejectsPublicDependencyOnInternal(t *testing.T) {
	spec := baseStagingSpec()
	spec.Modules[0].Visibility = stringPtr("internal")

	_, err := resolved.Resolve(resolved.ResolveInput{
		Staging: stagingManifest(t, spec),
		Modules: standardModules(t),
	})

	if err == nil {
		t.Fatalf("expected public dependency on internal module error")
	}
}
