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

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/manifest/resolved"
)

func TestResolveBuildsPublicationSet(t *testing.T) {
	spec := baseStagingSpec()
	spec.Defaults.Branches = []manifest.BranchMappingSpec{{Source: "main", Target: "main"}}
	set, err := resolved.Resolve(resolved.ResolveInput{Staging: stagingManifest(t, spec), Modules: standardModules(t)})
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	mods := set.Modules()
	if len(mods) != 2 {
		t.Fatalf("expected two modules, got %d", len(mods))
	}
	control := mods[1]
	if control.Name() != "control" || control.ModulePath() != "arcoris.dev/control" {
		t.Fatalf("unexpected control module: %#v", control)
	}
	if len(control.Dependencies()) != 1 || control.Dependencies()[0] != "foundation" {
		t.Fatalf("unexpected dependencies")
	}
}

func TestResolveWithTraceRecordsValueSources(t *testing.T) {
	result, err := resolved.ResolveWithTrace(resolved.ResolveInput{Staging: stagingManifest(t, baseStagingSpec()), Modules: standardModules(t)})
	if err != nil {
		t.Fatalf("ResolveWithTrace returned error: %v", err)
	}
	if len(result.Trace.Fields()) == 0 {
		t.Fatalf("expected resolution trace")
	}
	if len(result.Set.Modules()) != 2 {
		t.Fatalf("expected resolved set in result")
	}
}
