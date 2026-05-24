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

func TestBuildTargetReportHidesLocalPathsByDefault(t *testing.T) {
	t.Parallel()

	report := BuildTargetReport(reportWorkflowResult(t, workflowReportFixture{}).Target(), Options{})

	if !report.Present || report.WorkspaceCount != 2 {
		t.Fatalf("target report = %+v", report)
	}
	if report.Workspaces[0].Module != "foundation" || report.Workspaces[0].Repository != "arcoris/foundation" {
		t.Fatalf("first workspace = %+v", report.Workspaces[0])
	}
	if report.Workspaces[0].WorktreeDir != "" {
		t.Fatalf("target report leaked worktree path: %+v", report.Workspaces[0])
	}
	if len(report.Workspaces[0].Branches) == 0 {
		t.Fatalf("target branches = %+v", report.Workspaces[0].Branches)
	}
}

func TestBuildTargetReportCanIncludeLocalPaths(t *testing.T) {
	t.Parallel()

	report := BuildTargetReport(
		reportWorkflowResult(t, workflowReportFixture{}).Target(),
		Options{IncludeLocalPaths: true},
	)

	if report.Workspaces[0].WorktreeDir != "/target/arcoris__foundation" {
		t.Fatalf("worktree path = %+v", report.Workspaces[0])
	}
}
