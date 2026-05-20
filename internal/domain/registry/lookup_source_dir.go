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

package registry

import "arcoris.dev/arcoris-publisher/internal/domain/manifest"

// ModuleBySourceDir returns the module with source directory and whether it was found.
//
// Source directories are unique to keep staged filesystem operations
// unambiguous. A copy, clean, or hash operation can therefore map a path root
// back to a single module.
func (r Registry) ModuleBySourceDir(sourceDir manifest.SourceDir) (manifest.Module, bool) {
	index, ok := r.bySourceDir[sourceDir]
	if !ok {
		return manifest.Module{}, false
	}
	return r.modules[index], true
}

// ContainsSourceDir reports whether a module with source directory exists.
func (r Registry) ContainsSourceDir(sourceDir manifest.SourceDir) bool {
	_, ok := r.bySourceDir[sourceDir]
	return ok
}
