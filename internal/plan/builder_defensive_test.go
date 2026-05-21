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
	"testing"

	"arcoris.dev/arcoris-publisher/internal/manifest"
	"arcoris.dev/arcoris-publisher/internal/registry"
	"arcoris.dev/arcoris-publisher/internal/versioning"
)

func TestBuildRejectsRegistryMissingPublishOrderModule(t *testing.T) {
	req := mustRequest(t, "v0.3.0", testModule{name: "foundation"})
	req.Registry = registry.Registry{}

	_, err := Build(req)

	assertPlanError(t, err, IssueUnknownModule)
}

func TestBuildRejectsAssignmentPathMismatch(t *testing.T) {
	req := mustRequest(t, "v0.3.0", testModule{name: "foundation"})

	alternate := mustRequest(t, "v0.3.0", testModule{
		name:       "foundation",
		modulePath: "arcoris.dev/foundation-alt",
	})
	req.Assignments = alternate.Assignments

	_, err := Build(req)

	assertPlanError(t, err, IssueMissingAssignment)
}

func TestBuilderRejectsNonPublishableModule(t *testing.T) {
	set := mustPublicationSet(t, testModule{
		name:       "tooling",
		visibility: string(manifest.VisibilityInternal),
	})
	module := set.Modules()[0]

	var builder builder
	_, ok := builder.buildModulePlan(0, module)

	if ok {
		t.Fatal("buildModulePlan(internal) succeeded")
	}
	assertPlanError(t, builder.issues.Err(), IssueNonPublishableModule)
}

func TestBuilderRecordsMissingExecutionData(t *testing.T) {
	var builder builder
	set := mustPublicationSet(t, testModule{name: "foundation"})
	module := set.Modules()[0]

	builder.addEmptyBranchesIssue("modules[0]", module)
	builder.addEmptyPublishEntriesIssue("modules[0]", module)

	err := builder.issues.Err()
	assertPlanError(t, err, IssueEmptyBranches)
	assertPlanError(t, err, IssueEmptyPublishEntries)
}

func TestFinalizeRejectsDuplicateIndexes(t *testing.T) {
	module := validModulePlan("foundation")
	builder := builder{request: Request{
		Set: mustPublicationSet(t, testModule{name: "foundation"}),
	}}

	_, err := builder.finalize([]ModulePlan{module, module})

	assertPlanError(t, err, IssueDuplicateModuleName)
	assertPlanError(t, err, IssueDuplicateModulePath)
	assertPlanError(t, err, IssueDuplicateRepository)
}

func TestValidateRejectsInvalidModulePlan(t *testing.T) {
	p := Plan{modules: []ModulePlan{{
		name:       manifest.ModuleName("tooling"),
		modulePath: manifest.ModulePath("arcoris.dev/tooling"),
		repository: manifest.RepositoryRef("arcoris/tooling"),
		visibility: manifest.VisibilityInternal,
	}}}

	err := validate(p)

	assertPlanError(t, err, IssueNonPublishableModule)
	assertPlanError(t, err, IssueMissingAssignment)
	assertPlanError(t, err, IssueEmptyBranches)
	assertPlanError(t, err, IssueEmptyPublishEntries)
}

func TestFinalizeRejectsInvalidPlan(t *testing.T) {
	builder := builder{request: Request{
		Set: mustPublicationSet(t, testModule{name: "foundation"}),
	}}

	_, err := builder.finalize([]ModulePlan{{
		name:       manifest.ModuleName("foundation"),
		modulePath: manifest.ModulePath("arcoris.dev/foundation"),
		repository: manifest.RepositoryRef("arcoris/foundation"),
		visibility: manifest.VisibilityPublic,
	}})

	assertPlanError(t, err, IssueMissingAssignment)
	assertPlanError(t, err, IssueEmptyBranches)
	assertPlanError(t, err, IssueEmptyPublishEntries)
}

func TestValidateRejectsZeroPlan(t *testing.T) {
	assertPlanError(t, validate(Plan{}), IssueEmptyPlan)
}

// validModulePlan returns a minimal valid executable module plan.
func validModulePlan(name string) ModulePlan {
	entry, err := manifest.NewPublishEntry(manifest.PublishEntrySpec{
		Type: string(manifest.PublishEntryFile),
		From: "go.mod",
		To:   "go.mod",
	})
	if err != nil {
		panic(err)
	}

	branch := BranchPlan{
		source: manifest.BranchName("main"),
		target: manifest.BranchName("main"),
	}

	return ModulePlan{
		name:       manifest.ModuleName(name),
		modulePath: manifest.ModulePath("arcoris.dev/" + name),
		sourceDir:  manifest.SourceDir("src/arcoris.dev/" + name),
		repository: manifest.RepositoryRef("arcoris/" + name),
		visibility: manifest.VisibilityPublic,
		version:    versioning.Must("v0.3.0"),
		branches:   []BranchPlan{branch},
		entries:    []manifest.PublishEntry{entry},
	}
}
