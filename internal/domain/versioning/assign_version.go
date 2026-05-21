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

import (
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
)

// assignVersion creates assignment items for already-filtered publishable modules.
func assignVersion(modules []manifest.Module, policy manifest.VersionPolicy, version Version) (Assignments, error) {
	if version == "" {
		return Assignments{}, validationErrorf(IssueInvalidAssignment, "version", "version is required")
	}
	items := make([]ModuleVersion, 0, len(modules))
	for i, module := range modules {
		item, err := NewModuleVersion(module.Name(), module.ModulePath(), version)
		if err != nil {
			return Assignments{}, validationErrorf(IssueInvalidAssignment, fmt.Sprintf("modules[%d]", i), "invalid assignment: %v", err)
		}
		items = append(items, item)
	}
	assignments := Assignments{policy: policy, items: cloneModuleVersions(items)}
	if err := assignments.Validate(); err != nil {
		return Assignments{}, err
	}
	return assignments, nil
}
