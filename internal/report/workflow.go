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
	Status     Status           `json:"status"`
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
) Status {
	if !sourceReport.Present &&
		!targetReport.Present &&
		!constructReport.Present &&
		!moduleFileReport.Present &&
		!verifyReport.Present &&
		!publishReport.Present {
		return StatusEmpty
	}
	if verifyReport.Status == StatusEmpty {
		return StatusPartial
	}
	if verifyReport.Status == StatusFailed {
		return StatusVerificationFailed
	}
	switch publishReport.Status {
	case StatusPublished:
		return StatusPublished
	case StatusSkipped:
		return StatusSkipped
	default:
		return StatusVerified
	}
}

func writeWorkflowText(w io.Writer, report WorkflowReport) error {
	if err := writeLine(w, "Workflow"); err != nil {
		return err
	}
	if err := writeLine(w, "  Status: %s", report.Status); err != nil {
		return err
	}
	if err := writeLine(
		w,
		"  Source: %s, modules=%d",
		report.Source.Status,
		report.Source.ModuleCount,
	); err != nil {
		return err
	}
	if err := writeLine(
		w,
		"  Target: %s, workspaces=%d",
		report.Target.Status,
		report.Target.WorkspaceCount,
	); err != nil {
		return err
	}
	if err := writeLine(
		w,
		"  Construct: %s, changed=%t, operations=%d",
		report.Construct.Status,
		report.Construct.Changed,
		report.Construct.OperationCount,
	); err != nil {
		return err
	}
	if err := writeLine(
		w,
		"  Module files: %s, changed=%t, modules=%d",
		report.ModuleFile.Status,
		report.ModuleFile.Changed,
		report.ModuleFile.ModuleCount,
	); err != nil {
		return err
	}
	if err := writeLine(
		w,
		"  Verification: %s, failedChecks=%d",
		report.Verify.Status,
		report.Verify.FailedCount,
	); err != nil {
		return err
	}
	if err := writeLine(
		w,
		"  Publication: %s, published=%d, skipped=%d",
		report.Publish.Status,
		report.Publish.PublishedCount,
		report.Publish.SkippedCount,
	); err != nil {
		return err
	}
	return nil
}
