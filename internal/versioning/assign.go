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

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/manifest/resolved"
)

// Assign assigns module versions according to the effective publication policy.
func Assign(req Request) (Assignments, error) {
	assigner := assigner{
		request: req,
		issues:  newIssueCollector(),
	}
	return assigner.assign()
}

// assigner owns the mutable state for one assignment pass.
type assigner struct {
	// request is the immutable input bundle for this assignment run.
	request Request

	// issues accumulates validation failures in deterministic discovery order.
	issues issueCollector

	// order is the graph publication order for assigned public modules.
	order []manifest.ModuleName

	// versions stores the assigned version value for every public module.
	versions map[manifest.ModuleName]ModuleVersion

	// requirements stores direct public dependency requirements by module name.
	requirements map[manifest.ModuleName][]Requirement
}

// assign validates inputs, assigns versions, and derives direct requirements.
func (a *assigner) assign() (Assignments, error) {
	if err := a.validateRequest(); err != nil {
		return Assignments{}, err
	}

	if err := a.assignVersions(); err != nil {
		return Assignments{}, err
	}

	a.assignRequirements()
	if err := a.issues.Err(); err != nil {
		return Assignments{}, err
	}

	return newAssignments(a.order, a.versions, a.requirements), nil
}

// validateRequest checks policy/version compatibility and cross-index
// consistency before assignment starts.
func (a *assigner) validateRequest() error {
	a.validatePolicyAndVersion()
	a.validatePublishableInputs()

	return a.issues.Err()
}

// validatePolicyAndVersion enforces the resolved version policy against the
// caller-supplied version.
func (a *assigner) validatePolicyAndVersion() {
	policy := a.request.Set.Publish().VersionPolicy()
	if a.request.Version.IsZero() {
		a.issues.Add(IssueInvalidVersion, "version", "version is required")
		return
	}

	switch policy {
	case manifest.VersionPolicyReleaseTrain:
		a.validateReleaseTrainVersion()
	case manifest.VersionPolicySnapshot:
		a.validateSnapshotVersion()
	default:
		a.addInvalidPolicyIssue(policy)
	}
}

// validateReleaseTrainVersion requires a canonical, non-pseudo SemVer value.
func (a *assigner) validateReleaseTrainVersion() {
	if a.request.Version.IsRelease() {
		return
	}

	a.issues.Add(
		IssueInvalidVersion,
		"version",
		"release-train requires a non-pseudo SemVer version",
	)
}

// validateSnapshotVersion requires a Go pseudo-version.
func (a *assigner) validateSnapshotVersion() {
	if a.request.Version.IsPseudo() {
		return
	}

	a.issues.Add(
		IssueInvalidVersion,
		"version",
		"snapshot policy requires a Go pseudo-version",
	)
}

// addInvalidPolicyIssue records a missing or unsupported resolved policy.
func (a *assigner) addInvalidPolicyIssue(policy manifest.VersionPolicy) {
	a.issues.Add(
		IssueInvalidRequest,
		"publish.versionPolicy",
		"unsupported or missing version policy %q",
		policy,
	)
}

// validatePublishableInputs rejects requests whose resolved set, registry, and
// graph no longer describe the same publishable modules.
func (a *assigner) validatePublishableInputs() {
	for i, module := range a.request.Set.Modules() {
		if module.Visibility() != manifest.VisibilityPublic {
			continue
		}

		path := fmt.Sprintf("modules[%d]", i)
		a.validateRegistryModule(path, module)
		a.validateGraphNode(path, module)
	}
}

// validateRegistryModule confirms the registry still exposes the same
// publishable module identity as the resolved set.
func (a *assigner) validateRegistryModule(
	path string,
	module resolved.PublicationModule,
) {
	registryModule, ok := a.request.Registry.ModuleByName(module.Name())
	if !ok {
		a.issues.Add(
			IssueUnknownModule,
			path+".name",
			"module %q is absent from registry",
			module.Name(),
		)
		return
	}

	if registryModule.ModulePath() != module.ModulePath() {
		a.issues.Add(
			IssueInvalidRequest,
			path+".module.path",
			"registry module %q has path %q, want %q",
			module.Name(),
			registryModule.ModulePath(),
			module.ModulePath(),
		)
	}
}

// validateGraphNode confirms the graph includes the publishable module.
func (a *assigner) validateGraphNode(path string, module resolved.PublicationModule) {
	if a.request.Graph.Contains(module.Name()) {
		return
	}

	a.issues.Add(
		IssueUnknownModule,
		path+".name",
		"module %q is absent from graph",
		module.Name(),
	)
}

// assignVersions walks graph publication order and creates one assignment per
// public module.
func (a *assigner) assignVersions() error {
	order, err := a.request.Graph.PublishOrder()
	if err != nil {
		return graphOrderError(err)
	}

	a.order = order
	a.versions = make(map[manifest.ModuleName]ModuleVersion, len(order))
	a.requirements = make(map[manifest.ModuleName][]Requirement, len(order))

	for i, name := range order {
		a.assignVersion(i, name)
	}

	return a.issues.Err()
}

// graphOrderError wraps graph ordering failures in versioning diagnostics.
func graphOrderError(err error) *ValidationError {
	return &ValidationError{
		Scope: "versioning",
		Issues: []Issue{
			{
				Code:    IssueGraphOrder,
				Path:    "graph.publishOrder",
				Message: err.Error(),
			},
		},
	}
}

// assignVersion stores the caller-supplied version for one public module.
func (a *assigner) assignVersion(index int, name manifest.ModuleName) {
	path := fmt.Sprintf("publishOrder[%d]", index)
	module, ok := a.request.Registry.ModuleByName(name)
	if !ok {
		a.issues.Add(IssueUnknownModule, path, "module %q is absent from registry", name)
		return
	}

	if module.Visibility() != manifest.VisibilityPublic {
		a.issues.Add(IssueInvalidRequest, path, "module %q is not public", name)
		return
	}

	a.versions[name] = newModuleVersion(name, module.ModulePath(), a.request.Version)
	a.requirements[name] = nil
}

// assignRequirements derives direct go.mod requirements for assigned modules.
func (a *assigner) assignRequirements() {
	for i, name := range a.order {
		a.assignModuleRequirements(i, name)
	}
}

// assignModuleRequirements derives all direct dependency requirements for one
// assigned module.
func (a *assigner) assignModuleRequirements(index int, name manifest.ModuleName) {
	dependencies, _ := a.request.Graph.DirectDependencies(name)
	requirements := make([]Requirement, 0, len(dependencies))

	for j, dependency := range dependencies {
		path := fmt.Sprintf("modules[%d].dependencies[%d]", index, j)
		requirement, ok := a.requirementForDependency(path, dependency)
		if !ok {
			continue
		}

		requirements = append(requirements, requirement)
	}

	a.requirements[name] = requirements
}

// requirementForDependency validates one dependency and converts its assignment
// into a go.mod requirement value.
func (a *assigner) requirementForDependency(
	path string,
	dependency manifest.ModuleName,
) (Requirement, bool) {
	depModule, ok := a.request.Registry.ModuleByName(dependency)
	if !ok {
		a.issues.Add(
			IssueUnknownModule,
			path,
			"dependency %q is absent from registry",
			dependency,
		)
		return Requirement{}, false
	}

	if depModule.Visibility() != manifest.VisibilityPublic {
		a.issues.Add(
			IssueNonPublishableDependency,
			path,
			"dependency %q is not publishable",
			dependency,
		)
		return Requirement{}, false
	}

	assigned, ok := a.versions[dependency]
	if !ok {
		a.issues.Add(
			IssueMissingAssignment,
			path,
			"dependency %q has no assigned version",
			dependency,
		)
		return Requirement{}, false
	}

	return newRequirement(
		dependency,
		assigned.ModulePath(),
		assigned.Version(),
	), true
}
