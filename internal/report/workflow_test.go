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
	"testing"

	"arcoris.dev/arcoris-publisher/internal/workflow"
)

func TestBuildWorkflowReportZeroValueIsSafe(t *testing.T) {
	t.Parallel()

	report := BuildWorkflowReport(workflow.Result{}, Options{})

	if report.Kind != "workflow" || report.Status != "empty" {
		t.Fatalf("workflow report = %+v", report)
	}
}

func TestBuildWorkflowReportStatuses(t *testing.T) {
	t.Parallel()

	verified := BuildWorkflowReport(reportWorkflowResult(t, workflowReportFixture{}), Options{})
	if verified.Status != "verified" {
		t.Fatalf("verified status = %+v", verified)
	}

	failed := BuildWorkflowReport(
		reportWorkflowResult(t, workflowReportFixture{verifyFails: true}),
		Options{},
	)
	if failed.Status != "verification_failed" {
		t.Fatalf("failed status = %+v", failed)
	}

	published := BuildWorkflowReport(
		reportWorkflowResult(
			t,
			workflowReportFixture{
				publish:      true,
				dirtyModules: []string{"foundation", "control"},
			},
		),
		Options{},
	)
	if published.Status != "published" {
		t.Fatalf("published status = %+v", published)
	}

	skipped := BuildWorkflowReport(
		reportWorkflowResult(t, workflowReportFixture{publish: true}),
		Options{},
	)
	if skipped.Status != "skipped" {
		t.Fatalf("skipped status = %+v", skipped)
	}
}

func TestWorkflowStatusReportsPartialBeforeVerification(t *testing.T) {
	t.Parallel()

	got := workflowStatus(
		SourceReport{Present: true},
		TargetReport{Present: true},
		ConstructReport{Present: true},
		ModuleFileReport{Present: true},
		VerifyReport{Status: "empty"},
		PublishReport{Status: "empty"},
	)
	if got != "partial" {
		t.Fatalf("workflowStatus() = %q", got)
	}
}
