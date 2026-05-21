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

package resolved

import (
	modulemanifest "arcoris.dev/arcoris-publisher/internal/manifest/module"
	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
)

// resolveModule combines one staging module and its matching module manifest.
func (r *resolver) resolveModule(
	path string,
	sm staging.Module,
	mm modulemanifest.Manifest,
) PublicationModule {
	manifestPath := r.resolveManifestPath(path, sm)
	visibility := r.resolveVisibility(path, sm)
	branches := r.resolveBranches(path, sm)
	verification := r.resolveVerification(path, mm)
	identity := mm.Module()
	return PublicationModule{
		name:         sm.Name(),
		sourceDir:    sm.SourceDir(),
		manifestPath: manifestPath,
		repository:   sm.Repository(),
		visibility:   visibility,
		branches:     branches,
		moduleType:   identity.Type(),
		modulePath:   identity.Path(),
		moduleRoot:   identity.Root(),
		goMod:        identity.GoMod(),
		dependencies: mm.Dependencies().Internal(),
		entries:      mm.Publish().Entries(),
		verification: verification,
	}
}
