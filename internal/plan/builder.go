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

package plan

import (
	"fmt"

	"arcoris.dev/arcoris-publisher/internal/graph"
	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/manifest/resolved"
	"arcoris.dev/arcoris-publisher/internal/registry"
	"arcoris.dev/arcoris-publisher/internal/versioning"
)

// Build constructs an executable publication plan from resolved manifests,
// registry indexes, graph topology, and version assignments.
func Build(req Request) (Plan, error) {
	builder := builder{
		request: req,
		issues:  newIssueCollector(),
	}
	return builder.build()
}

// FromPublicationSet builds registry, graph, version assignments, and then the
// final publication plan for set.
//
// The helper is convenient for tests and simple command paths. Code that already
// has registry, graph, or assignments should call Build directly to avoid
// rebuilding deterministic intermediate models.
func FromPublicationSet(set resolved.PublicationSet, version versioning.Version) (Plan, error) {
	reg, err := registry.New(set)
	if err != nil {
		return Plan{}, fmt.Errorf("registry: %w", err)
	}

	g, err := graph.New(reg)
	if err != nil {
		return Plan{}, fmt.Errorf("graph: %w", err)
	}

	assignments, err := versioning.Assign(versioning.Request{
		Set:      set,
		Registry: reg,
		Graph:    g,
		Version:  version,
	})
	if err != nil {
		return Plan{}, fmt.Errorf("versioning: %w", err)
	}

	return Build(Request{
		Set:         set,
		Registry:    reg,
		Graph:       g,
		Assignments: assignments,
	})
}

// builder owns mutable state for one plan-building pass.
type builder struct {
	// request is the immutable input bundle for the plan run.
	request Request

	// issues accumulates module-level diagnostics in deterministic order.
	issues issueCollector
}

// build creates the executable plan in dependency-before-dependent order.
func (b *builder) build() (Plan, error) {
	order, err := b.publishOrder()
	if err != nil {
		return Plan{}, err
	}

	modules := b.buildModulePlans(order)
	if err := b.issues.Err(); err != nil {
		return Plan{}, err
	}

	return b.finalize(modules)
}

// publishOrder obtains and validates the graph's public module order.
func (b *builder) publishOrder() ([]manifest.ModuleName, error) {
	order, err := b.request.Graph.PublishOrder()
	if err != nil {
		return nil, graphOrderError(err)
	}

	if len(order) == 0 {
		return nil, emptyPlanError()
	}

	return order, nil
}

// buildModulePlans converts ordered graph names into executable module plans.
func (b *builder) buildModulePlans(order []manifest.ModuleName) []ModulePlan {
	modules := make([]ModulePlan, 0, len(order))
	for i, name := range order {
		mod, ok := b.request.Registry.ModuleByName(name)
		if !ok {
			b.addUnknownModuleIssue(i, name)
			continue
		}

		modulePlan, ok := b.buildModulePlan(i, mod)
		if ok {
			modules = append(modules, modulePlan)
		}
	}

	return modules
}

// finalize detaches the completed module list, rebuilds lookup indexes, and
// runs final plan validation.
func (b *builder) finalize(modules []ModulePlan) (Plan, error) {
	out := Plan{
		metadata: b.request.Set.Metadata(),
		source:   b.request.Set.Source(),
		publish:  b.request.Set.Publish(),
		modules:  modules,
	}

	if err := out.rebuildIndexes(); err != nil {
		return Plan{}, err
	}

	if err := validate(out); err != nil {
		return Plan{}, err
	}

	return out, nil
}

// graphOrderError wraps graph ordering failures in plan diagnostics.
func graphOrderError(err error) *ValidationError {
	return &ValidationError{
		Scope: "plan",
		Issues: []Issue{
			{
				Code:    IssueGraphOrder,
				Path:    "graph.publishOrder",
				Message: err.Error(),
			},
		},
	}
}

// emptyPlanError reports that no public modules can be planned.
func emptyPlanError() *ValidationError {
	return &ValidationError{
		Scope: "plan",
		Issues: []Issue{
			{
				Code:    IssueEmptyPlan,
				Path:    "modules",
				Message: "publication plan has no public modules",
			},
		},
	}
}

// addUnknownModuleIssue records a publish order name absent from the registry.
func (b *builder) addUnknownModuleIssue(index int, name manifest.ModuleName) {
	b.issues.Add(
		IssueUnknownModule,
		fmt.Sprintf("publishOrder[%d]", index),
		"module %q is absent from registry",
		name,
	)
}

// buildModulePlan converts one resolved public module into execution data.
func (b *builder) buildModulePlan(index int, mod resolved.PublicationModule) (ModulePlan, bool) {
	path := fmt.Sprintf("modules[%d]", index)
	if mod.Visibility() != manifest.VisibilityPublic {
		b.addNonPublishableModuleIssue(path, mod)
		return ModulePlan{}, false
	}

	version := b.moduleVersion(path, mod)
	branches := newBranchPlans(mod.Branches())
	if len(branches) == 0 {
		b.addEmptyBranchesIssue(path, mod)
	}

	entries := mod.PublishEntries()
	if len(entries) == 0 {
		b.addEmptyPublishEntriesIssue(path, mod)
	}

	requirements, ok := b.request.Assignments.RequirementsFor(mod.Name())
	if !ok {
		b.addMissingRequirementsIssue(path, mod)
	}

	return ModulePlan{
		name:         mod.Name(),
		modulePath:   mod.ModulePath(),
		moduleType:   mod.ModuleType(),
		sourceDir:    mod.SourceDir(),
		moduleRoot:   mod.ModuleRoot(),
		goMod:        mod.GoMod(),
		repository:   mod.Repository(),
		visibility:   mod.Visibility(),
		version:      version,
		branches:     branches,
		entries:      manifest.ClonePublishEntries(entries),
		requirements: newDependencyRequirements(requirements),
		verification: mod.Verification(),
	}, true
}

// moduleVersion validates and returns the assigned version for mod.
func (b *builder) moduleVersion(
	path string,
	mod resolved.PublicationModule,
) versioning.Version {
	version, ok := b.request.Assignments.VersionOf(mod.Name())
	if !ok || version.IsZero() {
		b.addMissingVersionIssue(path, mod)
	}

	versioned, ok := b.request.Assignments.ModuleVersion(mod.Name())
	if !ok {
		b.addMissingModuleVersionIssue(path, mod)
		return version
	}

	if versioned.ModulePath() != mod.ModulePath() {
		b.addAssignmentPathMismatchIssue(path, mod, versioned)
	}

	return version
}

// addNonPublishableModuleIssue records an attempt to plan an internal or
// disabled module.
func (b *builder) addNonPublishableModuleIssue(
	path string,
	mod resolved.PublicationModule,
) {
	b.issues.Add(
		IssueNonPublishableModule,
		path+".visibility",
		"module %q is not public",
		mod.Name(),
	)
}

// addMissingVersionIssue records a missing scalar version assignment.
func (b *builder) addMissingVersionIssue(path string, mod resolved.PublicationModule) {
	b.issues.Add(
		IssueMissingAssignment,
		path+".version",
		"module %q has no assigned version",
		mod.Name(),
	)
}

// addMissingModuleVersionIssue records a missing full module-version record.
func (b *builder) addMissingModuleVersionIssue(
	path string,
	mod resolved.PublicationModule,
) {
	b.issues.Add(
		IssueMissingAssignment,
		path+".version",
		"module %q has no module version record",
		mod.Name(),
	)
}

// addAssignmentPathMismatchIssue records an assignment for the wrong module
// path.
func (b *builder) addAssignmentPathMismatchIssue(
	path string,
	mod resolved.PublicationModule,
	versioned versioning.ModuleVersion,
) {
	b.issues.Add(
		IssueMissingAssignment,
		path+".modulePath",
		"module %q assignment path %q does not match resolved module path %q",
		mod.Name(),
		versioned.ModulePath(),
		mod.ModulePath(),
	)
}

// addEmptyBranchesIssue records missing effective branch mappings.
func (b *builder) addEmptyBranchesIssue(path string, mod resolved.PublicationModule) {
	b.issues.Add(
		IssueEmptyBranches,
		path+".branches",
		"module %q has no branch mappings",
		mod.Name(),
	)
}

// addEmptyPublishEntriesIssue records missing explicit publication content.
func (b *builder) addEmptyPublishEntriesIssue(
	path string,
	mod resolved.PublicationModule,
) {
	b.issues.Add(
		IssueEmptyPublishEntries,
		path+".publish.entries",
		"module %q has no explicit publish entries",
		mod.Name(),
	)
}

// addMissingRequirementsIssue records missing dependency requirement data.
func (b *builder) addMissingRequirementsIssue(
	path string,
	mod resolved.PublicationModule,
) {
	b.issues.Add(
		IssueMissingRequirements,
		path+".requirements",
		"module %q has no dependency requirement set",
		mod.Name(),
	)
}

// newBranchPlans converts resolved branch mappings into detached plan values.
func newBranchPlans(mappings []manifest.BranchMapping) []BranchPlan {
	out := make([]BranchPlan, 0, len(mappings))
	for _, mapping := range mappings {
		out = append(out, newBranchPlan(mapping))
	}
	return out
}

// newDependencyRequirements converts version assignments into plan values.
func newDependencyRequirements(requirements []versioning.Requirement) []DependencyRequirement {
	out := make([]DependencyRequirement, 0, len(requirements))
	for _, requirement := range requirements {
		out = append(out, newDependencyRequirement(requirement))
	}
	return out
}
