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

// resolveVisibility applies per-module visibility override before the public default.
func (r *resolver) resolveVisibility(path string, sm staging.Module) manifest.Visibility {
	tracePath := path + ".visibility"

	if override, ok := sm.VisibilityOverride(); ok {
		r.trace.AddStagingModule(tracePath, override.String())
		return override
	}

	r.trace.AddBuiltInDefault(
		tracePath,
		manifest.VisibilityPublic.String(),
		"visibility",
	)
	return manifest.VisibilityPublic
}
