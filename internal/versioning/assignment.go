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

import "arcoris.dev/arcoris-publisher/internal/manifest"

// Assignments is an immutable-by-convention snapshot of module versions and
// direct internal dependency requirements.
type Assignments struct {
	// order preserves publication order for every assigned public module.
	order []manifest.ModuleName

	// byName indexes assigned versions by resolved module name.
	byName map[manifest.ModuleName]ModuleVersion

	// byPath indexes assigned versions by published Go module path.
	byPath map[manifest.ModulePath]ModuleVersion

	// requirements stores direct publishable dependency requirements by module
	// name. Values are cloned on construction and access.
	requirements map[manifest.ModuleName][]Requirement
}

// Empty reports whether no modules received a version assignment.
func (a Assignments) Empty() bool { return len(a.order) == 0 }

// Len returns the number of assigned modules.
func (a Assignments) Len() int { return len(a.order) }

// Modules returns assigned module versions in deterministic publication order.
func (a Assignments) Modules() []ModuleVersion {
	out := make([]ModuleVersion, 0, len(a.order))
	for _, name := range a.order {
		out = append(out, a.byName[name])
	}
	return out
}

// ModuleNames returns assigned module names in deterministic publication order.
func (a Assignments) ModuleNames() []manifest.ModuleName {
	out := make([]manifest.ModuleName, len(a.order))
	copy(out, a.order)
	return out
}

// VersionOf returns the version assigned to name.
func (a Assignments) VersionOf(name manifest.ModuleName) (Version, bool) {
	value, ok := a.byName[name]
	if !ok {
		return "", false
	}
	return value.Version(), true
}

// ModuleVersion returns the module version assigned to name.
func (a Assignments) ModuleVersion(name manifest.ModuleName) (ModuleVersion, bool) {
	value, ok := a.byName[name]
	return value, ok
}

// ModuleVersionByPath returns the module version assigned to module path.
func (a Assignments) ModuleVersionByPath(path manifest.ModulePath) (ModuleVersion, bool) {
	value, ok := a.byPath[path]
	return value, ok
}

// RequirementsFor returns direct internal dependency requirements for name.
func (a Assignments) RequirementsFor(name manifest.ModuleName) ([]Requirement, bool) {
	reqs, ok := a.requirements[name]
	if !ok {
		return nil, false
	}
	return cloneRequirements(reqs), true
}

// newAssignments detaches mutable builder maps before returning the public
// assignment snapshot.
func newAssignments(
	order []manifest.ModuleName,
	versions map[manifest.ModuleName]ModuleVersion,
	requirements map[manifest.ModuleName][]Requirement,
) Assignments {
	byName := make(map[manifest.ModuleName]ModuleVersion, len(versions))
	byPath := make(map[manifest.ModulePath]ModuleVersion, len(versions))
	for name, value := range versions {
		byName[name] = value
		byPath[value.ModulePath()] = value
	}

	reqs := make(map[manifest.ModuleName][]Requirement, len(requirements))
	for name, value := range requirements {
		reqs[name] = cloneRequirements(value)
	}

	outOrder := make([]manifest.ModuleName, len(order))
	copy(outOrder, order)

	return Assignments{
		order:        outOrder,
		byName:       byName,
		byPath:       byPath,
		requirements: reqs,
	}
}

// RequirementMapFor returns direct internal dependency requirements for name as
// a detached module-path-to-version map.
func (a Assignments) RequirementMapFor(
	name manifest.ModuleName,
) (map[manifest.ModulePath]Version, bool) {
	reqs, ok := a.requirements[name]
	if !ok {
		return nil, false
	}
	out := make(map[manifest.ModulePath]Version, len(reqs))
	for _, req := range reqs {
		out[req.ModulePath()] = req.Version()
	}

	return out, true
}
