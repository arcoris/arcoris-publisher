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

func TestManifestAccessorsReturnDetachedSlices(t *testing.T) {
	manifest := Must(validSpec())
	modules := manifest.Modules()
	modules[0] = Module{}
	modulesAgain := manifest.Modules()
	if modulesAgain[0].Name() != ModuleName("foundation") {
		t.Fatalf("Modules() returned attached slice")
	}
}

func TestManifestSourceAndPolicyAccessors(t *testing.T) {
	manifest := Must(validSpec())
	if manifest.Source().Repository() != RepositoryRef("arcoris/arcoris") {
		t.Fatalf("Source().Repository() = %q", manifest.Source().Repository())
	}
	if manifest.Policy().PushPolicy() != PushPolicyFastForwardOnly {
		t.Fatalf("Policy().PushPolicy() = %q", manifest.Policy().PushPolicy())
	}
	if _, ok := manifest.ModuleByName(ModuleName("missing")); ok {
		t.Fatalf("ModuleByName(missing) found unexpected module")
	}
}
