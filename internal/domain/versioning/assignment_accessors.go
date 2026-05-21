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

import "arcoris.dev/arcoris-publisher/internal/domain/manifest"

// Policy returns the assignment policy.
func (a Assignments) Policy() manifest.VersionPolicy { return a.policy }

// Len returns the number of module version assignments.
func (a Assignments) Len() int { return len(a.items) }

// Empty reports whether no publishable modules received a version.
func (a Assignments) Empty() bool { return len(a.items) == 0 }

// Modules returns assigned module names in deterministic declaration order.
func (a Assignments) Modules() []manifest.ModuleName {
	modules := make([]manifest.ModuleName, 0, len(a.items))
	for _, item := range a.items {
		modules = append(modules, item.Module())
	}
	return modules
}

// Items returns detached module version assignments in deterministic declaration order.
func (a Assignments) Items() []ModuleVersion {
	return cloneModuleVersions(a.items)
}
