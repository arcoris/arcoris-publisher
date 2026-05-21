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
	"fmt"
	"path/filepath"
	"strings"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
)

// ModuleManifestLocation ties a top-level module declaration to the concrete
// module manifest file loaded from disk.
type ModuleManifestLocation struct {
	Name       manifest.ModuleName
	SourceDir  string
	Path       string
	Relative   string
	Repository manifest.RepositoryRef
}

// stagingRootPath resolves source.stagingRoot relative to the loaded top-level
// manifest file.
func stagingRootPath(stagingManifestPath string, stg staging.Manifest) (string, error) {
	base := filepath.Dir(stagingManifestPath)

	root := filepath.Join(
		base,
		filepath.FromSlash(stg.Source().StagingRoot().String()),
	)

	root, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}

	return filepath.Clean(root), nil
}

// moduleManifestLocation converts one staging module declaration into the exact
// manifest file path that should be loaded from disk.
func moduleManifestLocation(
	stagingRoot string,
	stg staging.Manifest,
	mod staging.Module,
) (ModuleManifestLocation, error) {
	sourcePath, err := moduleSourcePath(stagingRoot, mod)
	if err != nil {
		return ModuleManifestLocation{}, err
	}

	manifestPath := moduleManifestPath(stg, mod)
	modulePath, err := resolveModuleManifestPath(sourcePath, manifestPath)
	if err != nil {
		return ModuleManifestLocation{}, err
	}

	relativePath, err := filepath.Rel(stagingRoot, modulePath)
	if err != nil {
		return ModuleManifestLocation{}, err
	}

	return ModuleManifestLocation{
		Name:       mod.Name(),
		SourceDir:  sourcePath,
		Path:       modulePath,
		Relative:   filepath.ToSlash(relativePath),
		Repository: mod.Repository(),
	}, nil
}

// moduleSourcePath resolves module.sourceDir and ensures it remains inside the
// staging root boundary.
func moduleSourcePath(
	stagingRoot string,
	mod staging.Module,
) (string, error) {
	sourcePath, err := absoluteJoinedPath(
		stagingRoot,
		mod.SourceDir().String(),
	)
	if err != nil {
		return "", err
	}

	if err := ensureInside(stagingRoot, sourcePath); err != nil {
		return "", sourceDirEscapeError(sourcePath, err)
	}

	return sourcePath, nil
}

// moduleManifestPath applies the per-module manifest override, falling back to
// defaults.moduleManifest.path when no override was declared.
func moduleManifestPath(
	stg staging.Manifest,
	mod staging.Module,
) manifest.RelativePath {
	manifestPath, ok := mod.ManifestPathOverride()
	if !ok {
		manifestPath = stg.Defaults().ModuleManifest().Path()
	}

	return manifestPath
}

// resolveModuleManifestPath joins a module source directory with the relative
// manifest path and enforces that the result cannot escape the module directory.
func resolveModuleManifestPath(
	sourcePath string,
	manifestPath manifest.RelativePath,
) (string, error) {
	modulePath, err := absoluteJoinedPath(sourcePath, manifestPath.String())
	if err != nil {
		return "", err
	}

	if err := ensureInside(sourcePath, modulePath); err != nil {
		return "", moduleManifestEscapeError(modulePath, err)
	}

	return modulePath, nil
}

// absoluteJoinedPath joins a slash-style manifest path to a native filesystem
// root and returns its clean absolute form.
func absoluteJoinedPath(root string, path string) (string, error) {
	joined := filepath.Join(root, filepath.FromSlash(path))

	absolutePath, err := filepath.Abs(joined)
	if err != nil {
		return "", err
	}

	return filepath.Clean(absolutePath), nil
}

// ensureInside verifies that candidate is root itself or a descendant of root.
func ensureInside(root string, candidate string) error {
	root = filepath.Clean(root)
	candidate = filepath.Clean(candidate)
	rel, err := filepath.Rel(root, candidate)
	if err != nil {
		return err
	}
	if rel == "." {
		return nil
	}
	if strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
		return fmt.Errorf("%q is outside %q", candidate, root)
	}
	return nil
}
