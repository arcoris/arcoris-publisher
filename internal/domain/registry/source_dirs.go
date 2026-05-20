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

// SourceDirs returns source directories in manifest declaration order.
//
// Source directories are the staged module roots used by filesystem adapters.
// Returning them as values keeps callers from depending on the internal
// source-directory index representation.
func (r Registry) SourceDirs() []manifest.SourceDir {
	values := make([]manifest.SourceDir, 0, len(r.modules))
	for _, module := range r.modules {
		values = append(values, module.SourceDir())
	}
	return values
}
