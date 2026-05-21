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

	"arcoris.dev/arcoris-publisher/internal/manifest/resolved"
)

// DiscoverPublicationSet discovers a top-level manifest in dir and resolves it.
func (l Loader) DiscoverPublicationSet(
	ctx context.Context,
	dir string,
) (resolved.PublicationSet, error) {
	result, err := l.DiscoverPublicationSetWithTrace(ctx, dir)
	if err != nil {
		return resolved.PublicationSet{}, err
	}

	return result.Set, nil
}

// DiscoverPublicationSetWithTrace discovers a top-level manifest in dir and
// returns the resolved publication set with value-origin trace.
func (l Loader) DiscoverPublicationSetWithTrace(
	ctx context.Context,
	dir string,
) (LoadResult, error) {
	path, err := l.locator.Find(ctx, l.reader, dir)
	if err != nil {
		return LoadResult{}, err
	}

	return l.LoadPublicationSetWithTrace(ctx, path)
}

// LoadPublicationSet loads a top-level manifest and every referenced module
// manifest, then resolves them into an effective publication set.
func (l Loader) LoadPublicationSet(
	ctx context.Context,
	path string,
) (resolved.PublicationSet, error) {
	result, err := l.LoadPublicationSetWithTrace(ctx, path)
	if err != nil {
		return resolved.PublicationSet{}, err
	}

	return result.Set, nil
}

// LoadPublicationSetWithTrace loads all manifests and returns resolution trace.
func (l Loader) LoadPublicationSetWithTrace(
	ctx context.Context,
	path string,
) (LoadResult, error) {
	path, err := normalizeInputPath(path)
	if err != nil {
		return LoadResult{}, err
	}

	stagingManifest, err := l.LoadStaging(ctx, path)
	if err != nil {
		return LoadResult{}, err
	}

	manifests, locations, err := l.loadModuleManifests(
		ctx,
		path,
		stagingManifest,
	)
	if err != nil {
		return LoadResult{}, err
	}

	resolvedResult, err := resolved.ResolveWithTrace(resolved.ResolveInput{
		Staging: stagingManifest,
		Modules: manifests,
	})
	if err != nil {
		return LoadResult{}, resolveError(path, err)
	}

	return newLoadResult(path, locations, resolvedResult), nil
}

// newLoadResult keeps the result construction in one place so the resolution
// output, staging path, and module locations stay synchronized.
func newLoadResult(
	stagingPath string,
	locations []ModuleManifestLocation,
	result resolved.ResolveResult,
) LoadResult {
	return LoadResult{
		Set:     result.Set,
		Trace:   result.Trace,
		Staging: stagingPath,
		Modules: locations,
	}
}
