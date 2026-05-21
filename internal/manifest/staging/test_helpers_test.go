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
	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
)

func validSpec() staging.Spec {
	return staging.Spec{
		APIVersion: string(manifest.APIVersionV1Alpha1),
		Kind:       string(manifest.KindStagingManifest),
		Metadata:   manifest.MetadataSpec{Name: "arcoris"},
		Source:     manifest.SourceSpec{Repository: "arcoris/arcoris", DefaultBranch: "main"},
		Modules: []staging.ModuleSpec{
			{Name: "foundation", SourceDir: "src/arcoris.dev/foundation", Repository: "arcoris/foundation"},
			{Name: "control", SourceDir: "src/arcoris.dev/control", Repository: "arcoris/control"},
		},
	}
}

func stringPtr(value string) *string { return &value }

func boolPtr(value bool) *bool { return &value }
