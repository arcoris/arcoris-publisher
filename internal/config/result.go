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

import "arcoris.dev/arcoris-publisher/internal/manifest/resolved"

// LoadResult contains a resolved publication set and the module manifest files
// that were read to build it.
type LoadResult struct {
	Set     resolved.PublicationSet
	Trace   resolved.ResolutionTrace
	Staging string
	Modules []ModuleManifestLocation
}

// ModuleLocations returns detached module manifest locations.
func (r LoadResult) ModuleLocations() []ModuleManifestLocation {
	out := make([]ModuleManifestLocation, len(r.Modules))
	copy(out, r.Modules)
	return out
}
