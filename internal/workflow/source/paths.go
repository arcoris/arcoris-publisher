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

package source

import (
	"path/filepath"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/workflow/pathutil"
)

// resolveModuleDir resolves a module SourceDir under the staging root.
func resolveModuleDir(stagingDir string, sourceDir manifest.SourceDir) string {
	return filepath.Clean(filepath.Join(stagingDir, filepath.FromSlash(sourceDir.String())))
}

// resolveModuleRootDir resolves a module root under the module source dir.
func resolveModuleRootDir(moduleDir string, moduleRoot manifest.RelativePath) string {
	return pathutil.JoinRelative(moduleDir, moduleRoot)
}

// resolveEntrySource resolves one explicit publish entry source path.
func resolveEntrySource(moduleRootDir string, entry manifest.PublishEntry) string {
	return pathutil.JoinRelative(moduleRootDir, entry.From())
}
