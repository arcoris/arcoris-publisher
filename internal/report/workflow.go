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

	"arcoris.dev/arcoris-publisher/internal/workflow"
)

// WorkflowReport is the stable aggregate DTO for a workflow run.
type WorkflowReport struct {
	Kind       string           `json:"kind"`
	Status     string           `json:"status"`
	Source     SourceReport     `json:"source"`
	Target     TargetReport     `json:"target"`
	Construct  ConstructReport  `json:"construct"`
	ModuleFile ModuleFileReport `json:"moduleFile"`
	Verify     VerifyReport     `json:"verify"`
	Publish    PublishReport    `json:"publish"`
}

// BuildWorkflowReport converts a completed or partial workflow result into a
// stable aggregate report DTO.
func BuildWorkflowReport(result workflow.Result, opts Options) WorkflowReport {
	sourceReport := BuildSourceReport(result.Source(), opts)
	targetReport := BuildTargetReport(result.Target(), opts)
	constructReport := BuildConstructReport(result.Construct(), opts)
	moduleFileReport := BuildModuleFileReport(result.ModuleFile(), opts)
	verifyReport := BuildVerifyReport(result.Verify(), opts)
	publishReport := BuildPublishReport(result.Publish(), opts)

	return WorkflowReport{
		Kind:       "workflow",
		Status:     workflowStatus(sourceReport, targetReport, constructReport, moduleFileReport, verifyReport, publishReport),
		Source:     sourceReport,
		Target:     targetReport,
		Construct:  constructReport,
		ModuleFile: moduleFileReport,
		Verify:     verifyReport,
		Publish:    publishReport,
	}
}

func workflowStatus(
	sourceReport SourceReport,
	targetReport TargetReport,
	constructReport ConstructReport,
	moduleFileReport ModuleFileReport,
	verifyReport VerifyReport,
	publishReport PublishReport,
) string {
	if !sourceReport.Present &&
		!targetReport.Present &&
		!constructReport.Present &&
		!moduleFileReport.Present &&
		!verifyReport.Present &&
		!publishReport.Present {
		return "empty"
	}
	if verifyReport.Status == "empty" {
		return "partial"
	}
	if verifyReport.Status == "failed" {
		return "verification_failed"
	}
	switch publishReport.Status {
	case "published":
		return "published"
	case "skipped":
		return "skipped"
	default:
		return "verified"
	}
}

func writeWorkflowText(w io.Writer, value any) error {
	report := value.(WorkflowReport)
	if err := writeLine(w, "Workflow"); err != nil {
		return err
	}
	if err := writeLine(w, "  Status: %s", report.Status); err != nil {
		return err
	}
	if err := writeLine(w, "  Source modules: %d", len(report.Source.Modules)); err != nil {
		return err
	}
	if err := writeLine(w, "  Target workspaces: %d", report.Target.WorkspaceCount); err != nil {
		return err
	}
	if err := writeLine(w, "  Construct changed: %t", report.Construct.Changed); err != nil {
		return err
	}
	if err := writeLine(w, "  Module files changed: %t", report.ModuleFile.Changed); err != nil {
		return err
	}
	if err := writeLine(w, "  Verification: %s", report.Verify.Status); err != nil {
		return err
	}
	if err := writeLine(w, "  Publication: %s", report.Publish.Status); err != nil {
		return err
	}
	return nil
}
