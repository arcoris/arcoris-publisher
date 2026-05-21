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
	"arcoris.dev/arcoris-publisher/internal/domain/manifest"
	"arcoris.dev/arcoris-publisher/internal/domain/registry"
)

// assignmentBuilder normalizes AssignmentSpec before building module assignments.
type assignmentBuilder struct {
	registry registry.Registry
	spec     AssignmentSpec
}

// newAssignmentBuilder captures the immutable inputs needed to build assignments.
func newAssignmentBuilder(registryValue registry.Registry, spec AssignmentSpec) assignmentBuilder {
	return assignmentBuilder{registry: registryValue, spec: spec}
}

// build resolves the policy-specific version and assigns it to publishable modules.
func (b assignmentBuilder) build() (Assignments, error) {
	policy := b.policy()
	version, err := b.version(policy)
	if err != nil {
		return Assignments{}, err
	}
	return assignVersion(b.registry.PublishableModules(), policy, version)
}

// policy returns the explicit policy or the release-train default.
func (b assignmentBuilder) policy() manifest.VersionPolicy {
	if b.spec.Policy == "" {
		return manifest.VersionPolicyReleaseTrain
	}
	return b.spec.Policy
}
