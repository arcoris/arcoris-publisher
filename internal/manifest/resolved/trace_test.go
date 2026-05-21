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

	"arcoris.dev/arcoris-publisher/internal/manifest/resolved"
)

func TestResolutionTraceRecordsFieldsInOrder(t *testing.T) {
	var trace resolved.ResolutionTrace

	trace.AddBuiltInDefault(
		"modules[0].visibility",
		"public",
		"visibility",
	)
	trace.AddStagingDefault(
		"modules[0].manifest",
		"arcpub.module.yaml",
		"defaults.moduleManifest.path",
	)

	fields := trace.Fields()
	if len(fields) != 2 ||
		fields[0].Path != "modules[0].visibility" ||
		fields[1].Source.Kind != resolved.SourceStagingDefault {
		t.Fatalf("unexpected trace fields: %#v", fields)
	}
}

func TestResolutionTraceFieldsReturnsDetachedSlice(t *testing.T) {
	var trace resolved.ResolutionTrace
	trace.AddBuiltInDefault("path", "value", "path")

	fields := trace.Fields()
	fields[0].Value = "mutated"

	if trace.Fields()[0].Value == "mutated" {
		t.Fatalf("Fields accessor leaked internal slice")
	}
}
