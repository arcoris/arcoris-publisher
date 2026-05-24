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

import "testing"

func TestBuildConstructReportHidesOperationPathsByDefault(t *testing.T) {
	t.Parallel()

	report := BuildConstructReport(reportWorkflowResult(t, workflowReportFixture{}).Construct(), Options{})

	if !report.Present || !report.Changed || report.ModuleCount != 2 || report.OperationCount == 0 {
		t.Fatalf("construct report = %+v", report)
	}
	first := report.Modules[0]
	if first.WorktreeDir != "" {
		t.Fatalf("construct report leaked worktree path: %+v", first)
	}
	for _, op := range first.Operations {
		if op.SourcePath != "" || op.TargetPath != "" {
			t.Fatalf("construct report leaked operation paths: %+v", op)
		}
	}
}

func TestBuildConstructReportCanIncludeLocalPaths(t *testing.T) {
	t.Parallel()

	report := BuildConstructReport(
		reportWorkflowResult(t, workflowReportFixture{}).Construct(),
		Options{IncludeLocalPaths: true},
	)

	first := report.Modules[0]
	if first.WorktreeDir != "/target/arcoris__foundation" {
		t.Fatalf("worktree path = %+v", first)
	}
	if len(first.Operations) == 0 || first.Operations[0].TargetPath == "" {
		t.Fatalf("construct operations = %+v", first.Operations)
	}
}
