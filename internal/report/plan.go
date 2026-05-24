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

package report

import (
	"io"

	"arcoris.dev/arcoris-publisher/internal/plan"
)

// PlanReport is the stable JSON/text DTO for a publication plan.
type PlanReport struct {
	Kind          string              `json:"kind"`
	ModuleCount   int                 `json:"moduleCount"`
	Source        SourcePolicyReport  `json:"source"`
	PublishPolicy PublishPolicyReport `json:"publishPolicy"`
	Modules       []PlanModuleReport  `json:"modules"`
}

// SourcePolicyReport describes top-level source repository policy values.
type SourcePolicyReport struct {
	Repository    string `json:"repository"`
	DefaultBranch string `json:"defaultBranch"`
	StagingRoot   string `json:"stagingRoot"`
	ModuleRoot    string `json:"moduleRoot"`
	DirtyPolicy   string `json:"dirtyPolicy"`
}

// PlanModuleReport is the report representation of one planned public module.
type PlanModuleReport struct {
	Name           string                        `json:"name"`
	ModulePath     string                        `json:"modulePath"`
	ModuleType     string                        `json:"moduleType"`
	Repository     string                        `json:"repository"`
	Version        string                        `json:"version"`
	SourceDir      string                        `json:"sourceDir"`
	ModuleRoot     string                        `json:"moduleRoot"`
	GoMod          string                        `json:"goMod"`
	Branches       []BranchReport                `json:"branches"`
	PublishEntries []PublishEntryReport          `json:"publishEntries"`
	Requirements   []DependencyRequirementReport `json:"requirements"`
	Verification   VerificationPolicyReport      `json:"verification"`
}

// BuildPlanReport converts a publication plan to a stable report DTO.
func BuildPlanReport(p plan.Plan, opts Options) PlanReport {
	modules := p.Modules()
	out := PlanReport{
		Kind:        "plan",
		ModuleCount: len(modules),
		Source: SourcePolicyReport{
			Repository:    p.Source().Repository().String(),
			DefaultBranch: p.Source().DefaultBranch().String(),
			StagingRoot:   p.Source().StagingRoot().String(),
			ModuleRoot:    p.Source().ModuleRoot().String(),
			DirtyPolicy:   string(p.Source().DirtyPolicy()),
		},
		PublishPolicy: publishPolicyReport(p, opts),
		Modules:       make([]PlanModuleReport, 0, len(modules)),
	}
	for _, mod := range modules {
		out.Modules = append(out.Modules, planModuleReport(mod))
	}
	return out
}

func publishPolicyReport(p plan.Plan, opts Options) PublishPolicyReport {
	policy := p.PublishPolicy()
	provenanceFile := ""
	if policy.Provenance().FileEnabled() {
		provenanceFile = includePath(policy.Provenance().File().String(), opts)
	}
	return PublishPolicyReport{
		Mode:                  string(policy.Mode()),
		PushPolicy:            string(policy.PushPolicy()),
		TagPolicy:             string(policy.Tags().Mode()),
		TagEnabled:            policy.Tags().Enabled(),
		ProvenanceFileEnabled: policy.Provenance().FileEnabled(),
		ProvenanceFile:        provenanceFile,
		CommitTrailers:        policy.Provenance().CommitTrailers(),
	}
}

func planModuleReport(mod plan.ModulePlan) PlanModuleReport {
	entries := mod.PublishEntries()
	requirements := mod.Requirements()
	branches := mod.Branches()

	out := PlanModuleReport{
		Name:           mod.Name().String(),
		ModulePath:     mod.ModulePath().String(),
		ModuleType:     string(mod.ModuleType()),
		Repository:     mod.Repository().String(),
		Version:        mod.Version().String(),
		SourceDir:      mod.SourceDir().String(),
		ModuleRoot:     mod.ModuleRoot().String(),
		GoMod:          mod.GoMod().String(),
		Branches:       make([]BranchReport, 0, len(branches)),
		PublishEntries: make([]PublishEntryReport, 0, len(entries)),
		Requirements:   make([]DependencyRequirementReport, 0, len(requirements)),
		Verification:   verificationPolicyReport(mod),
	}
	for _, branch := range branches {
		out.Branches = append(out.Branches, BranchReport{
			Source: branch.Source().String(),
			Target: branch.Target().String(),
		})
	}
	for _, entry := range entries {
		out.PublishEntries = append(out.PublishEntries, PublishEntryReport{
			Kind:      string(entry.Kind()),
			From:      entry.From().String(),
			To:        entry.To().String(),
			Optional:  entry.Optional(),
			Recursive: entry.Recursive(),
		})
	}
	for _, req := range requirements {
		out.Requirements = append(out.Requirements, DependencyRequirementReport{
			Module:     req.Module().String(),
			ModulePath: req.ModulePath().String(),
			Version:    req.Version().String(),
		})
	}
	return out
}

func verificationPolicyReport(mod plan.ModulePlan) VerificationPolicyReport {
	policy := mod.Verification()
	goPolicy := policy.Go()
	return VerificationPolicyReport{
		GoListEnabled:      goPolicy.List(),
		GoTestEnabled:      goPolicy.Test(),
		GoTidyEnabled:      goPolicy.Tidy(),
		GoPatterns:         goPolicy.Patterns(),
		WorkspaceMode:      string(goPolicy.WorkspaceMode()),
		LocalReplacePolicy: string(policy.LocalReplacePolicy()),
	}
}

func writePlanText(w io.Writer, report PlanReport) error {
	if err := writeLine(w, "Plan"); err != nil {
		return err
	}
	if err := writeLine(w, "  Modules: %d", report.ModuleCount); err != nil {
		return err
	}
	if err := writeLine(w, "  Source: %s", report.Source.Repository); err != nil {
		return err
	}
	if err := writeLine(w, "  Publish mode: %s", report.PublishPolicy.Mode); err != nil {
		return err
	}
	if err := writeLine(w, "  Push policy: %s", report.PublishPolicy.PushPolicy); err != nil {
		return err
	}
	if err := writeLine(w, "  Tag policy: %s", report.PublishPolicy.TagPolicy); err != nil {
		return err
	}
	if err := writeLine(w, "  Commit trailers: %t", report.PublishPolicy.CommitTrailers); err != nil {
		return err
	}
	if err := writeLine(w, "  Provenance file: %t", report.PublishPolicy.ProvenanceFileEnabled); err != nil {
		return err
	}
	if err := writeLine(w, ""); err != nil {
		return err
	}
	for _, mod := range report.Modules {
		if err := writeLine(w, "  %s", mod.Name); err != nil {
			return err
		}
		if err := writeLine(w, "    module path: %s", mod.ModulePath); err != nil {
			return err
		}
		if err := writeLine(w, "    repository:  %s", mod.Repository); err != nil {
			return err
		}
		if err := writeLine(w, "    version:     %s", mod.Version); err != nil {
			return err
		}
		if err := writeLine(w, "    source dir:  %s", mod.SourceDir); err != nil {
			return err
		}
		if err := writeLine(w, "    branches:    %s", branchList(mod.Branches)); err != nil {
			return err
		}
		if err := writeLine(w, "    entries:"); err != nil {
			return err
		}
		for _, entry := range mod.PublishEntries {
			if err := writeLine(w, "      %s %s -> %s", entry.Kind, entry.From, entry.To); err != nil {
				return err
			}
		}
		if err := writeLine(w, "    requires:    %s", requirementList(mod.Requirements)); err != nil {
			return err
		}
		if err := writeLine(w, "    verification: %s", verificationPolicyText(mod.Verification)); err != nil {
			return err
		}
		if err := writeLine(w, ""); err != nil {
			return err
		}
	}
	return nil
}

func verificationPolicyText(policy VerificationPolicyReport) string {
	return fmtBool("list", policy.GoListEnabled) + ", " +
		fmtBool("test", policy.GoTestEnabled) + ", " +
		fmtBool("tidy", policy.GoTidyEnabled) + ", " +
		"workspace=" + policy.WorkspaceMode + ", " +
		"localReplace=" + policy.LocalReplacePolicy
}

func branchList(branches []BranchReport) string {
	if len(branches) == 0 {
		return "-"
	}
	items := make([]string, 0, len(branches))
	for _, branch := range branches {
		items = append(items, branch.Source+" -> "+branch.Target)
	}
	return commaOrDash(items)
}

func requirementList(requirements []DependencyRequirementReport) string {
	if len(requirements) == 0 {
		return "-"
	}
	items := make([]string, 0, len(requirements))
	for _, req := range requirements {
		items = append(items, req.ModulePath+" "+req.Version)
	}
	return commaOrDash(items)
}
