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

// Repositories returns target repositories in manifest declaration order.
//
// The returned repositories are suitable for planning remote operations; the
// registry keeps uniqueness checks centralized so callers can treat each value
// as belonging to exactly one module.
func (r Registry) Repositories() []manifest.RepositoryRef {
	values := make([]manifest.RepositoryRef, 0, len(r.modules))
	for _, module := range r.modules {
		values = append(values, module.Repository())
	}
	return values
}
