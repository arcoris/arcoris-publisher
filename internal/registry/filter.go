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

import (
	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/manifest/resolved"
)

// PublishableModules returns public modules in declaration order.
func (r Registry) PublishableModules() []resolved.PublicationModule {
	return r.modulesByVisibility(manifest.VisibilityPublic)
}

// InternalModules returns internal modules in declaration order.
func (r Registry) InternalModules() []resolved.PublicationModule {
	return r.modulesByVisibility(manifest.VisibilityInternal)
}

// DisabledModules returns disabled modules in declaration order.
func (r Registry) DisabledModules() []resolved.PublicationModule {
	return r.modulesByVisibility(manifest.VisibilityDisabled)
}

func (r Registry) modulesByVisibility(
	visibility manifest.Visibility,
) []resolved.PublicationModule {
	out := make([]resolved.PublicationModule, 0)
	for _, module := range r.modules {
		if module.Visibility() == visibility {
			out = append(out, module)
		}
	}

	return out
}
