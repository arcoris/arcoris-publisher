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
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	modulemanifest "arcoris.dev/arcoris-publisher/internal/manifest/module"
)

// indexModuleManifests builds a deterministic lookup by metadata.name.
func (r *resolver) indexModuleManifests() (map[manifest.ModuleName]modulemanifest.Manifest, error) {
	var collector manifest.IssueCollector
	out := make(map[manifest.ModuleName]modulemanifest.Manifest, len(r.input.Modules))
	for i, mod := range r.input.Modules {
		name := manifest.ModuleName(mod.Metadata().Name())
		if prev, exists := out[name]; exists {
			collector.Add(manifest.IssueDuplicateValue, fmt.Sprintf("moduleManifests[%d].metadata.name", i), "duplicate module manifest name %q also declared by %q", name, prev.Metadata().Name())
			continue
		}
		out[name] = mod
	}
	if err := collector.Err(); err != nil {
		return nil, err
	}
	return out, nil
}
