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

package versioning

import "arcoris.dev/arcoris-publisher/internal/manifest"

// ModuleVersion binds one publishable module identity to an assigned version.
type ModuleVersion struct {
	// name is the resolved manifest module name.
	name manifest.ModuleName

	// modulePath is the Go module path that will receive the version.
	modulePath manifest.ModulePath

	// version is the concrete release or snapshot version assigned to the
	// module.
	version Version
}

// newModuleVersion creates the value object stored in assignment indexes.
func newModuleVersion(
	name manifest.ModuleName,
	modulePath manifest.ModulePath,
	version Version,
) ModuleVersion {
	return ModuleVersion{name: name, modulePath: modulePath, version: version}
}

// Module returns the assigned module name.
func (v ModuleVersion) Module() manifest.ModuleName { return v.name }

// ModulePath returns the assigned module path.
func (v ModuleVersion) ModulePath() manifest.ModulePath { return v.modulePath }

// Version returns the assigned publication version.
func (v ModuleVersion) Version() Version { return v.version }
