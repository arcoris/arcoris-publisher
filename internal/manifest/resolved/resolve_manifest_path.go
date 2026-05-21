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
	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/manifest/staging"
)

// resolveManifestPath applies per-module manifest path override before defaults.
func (r *resolver) resolveManifestPath(path string, sm staging.Module) manifest.RelativePath {
	tracePath := path + ".manifest"

	if override, ok := sm.ManifestPathOverride(); ok {
		r.trace.AddStagingModule(tracePath, override.String())
		return override
	}

	value := r.input.Staging.Defaults().ModuleManifest().Path()
	r.trace.AddStagingDefault(
		tracePath,
		value.String(),
		"defaults.moduleManifest.path",
	)
	return value
}
