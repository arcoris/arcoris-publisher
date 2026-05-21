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

import "arcoris.dev/arcoris-publisher/internal/domain/registry"

// New builds version assignments for publishable modules in registry.
//
// The spec policy defaults to release-train because explicit release versions
// are the safest production path. Snapshot inputs are accepted only when the
// snapshot policy is selected.
func New(registryValue registry.Registry, spec AssignmentSpec) (Assignments, error) {
	return newAssignmentBuilder(registryValue, spec).build()
}

// Must constructs Assignments and panics on validation failure.
func Must(registryValue registry.Registry, spec AssignmentSpec) Assignments {
	assignments, err := New(registryValue, spec)
	if err != nil {
		panic(err)
	}
	return assignments
}
