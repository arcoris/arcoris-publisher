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

package config

import (
	"context"

	modulemanifest "arcoris.dev/arcoris-publisher/internal/manifest/module"
	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
)

// loadModuleManifests resolves every module declaration from the staging
// manifest to a concrete module manifest file, preserving the loaded locations
// for diagnostics and callers that want to display provenance.
func (l Loader) loadModuleManifests(
	ctx context.Context,
	stagingPath string,
	stg staging.Manifest,
) ([]modulemanifest.Manifest, []ModuleManifestLocation, error) {
	root, err := stagingRootPath(stagingPath, stg)
	if err != nil {
		return nil, nil, stagingRootError(stagingPath, err)
	}

	modules := stg.Modules()
	manifests := make([]modulemanifest.Manifest, 0, len(modules))
	locations := make([]ModuleManifestLocation, 0, len(modules))

	for _, mod := range modules {
		manifest, location, err := l.loadModuleManifest(ctx, root, stg, mod)
		if err != nil {
			return nil, nil, err
		}

		manifests = append(manifests, manifest)
		locations = append(locations, location)
	}

	return manifests, locations, nil
}

// loadModuleManifest computes the expected module manifest location first, then
// delegates actual parsing and validation to LoadModule.
func (l Loader) loadModuleManifest(
	ctx context.Context,
	root string,
	stg staging.Manifest,
	mod staging.Module,
) (modulemanifest.Manifest, ModuleManifestLocation, error) {
	location, err := moduleManifestLocation(root, stg, mod)
	if err != nil {
		return modulemanifest.Manifest{}, ModuleManifestLocation{}, err
	}

	manifest, err := l.LoadModule(ctx, location.Path)
	if err != nil {
		return modulemanifest.Manifest{}, ModuleManifestLocation{}, err
	}

	return manifest, location, nil
}
