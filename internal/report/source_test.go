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

func TestBuildSourceReportHidesLocalPathsByDefault(t *testing.T) {
	t.Parallel()

	report := BuildSourceReport(reportWorkflowResult(t, workflowReportFixture{}).Source(), Options{})

	if !report.Present || report.Status != StatusPresent {
		t.Fatalf("source report presence = %+v", report)
	}
	if report.Repository.Head == "" || report.Repository.Branch != "main" {
		t.Fatalf("source repository = %+v", report.Repository)
	}
	if report.Repository.RepositoryDir != "" || report.Repository.StagingDir != "" {
		t.Fatalf("source report leaked repository paths: %+v", report.Repository)
	}
	if report.ModuleCount != 2 ||
		len(report.Modules) != 2 ||
		report.Modules[0].Hash == "" ||
		report.Modules[0].EntryCount != 2 {
		t.Fatalf("source modules = %+v", report.Modules)
	}
	if report.Modules[0].SourceDir != "" || report.Modules[0].ModuleRootDir != "" {
		t.Fatalf("source report leaked module paths: %+v", report.Modules[0])
	}
}

func TestBuildSourceReportCanIncludeLocalPaths(t *testing.T) {
	t.Parallel()

	report := BuildSourceReport(
		reportWorkflowResult(t, workflowReportFixture{}).Source(),
		Options{IncludeLocalPaths: true},
	)

	if report.Repository.RepositoryDir != "/repo" || report.Repository.StagingDir != "/repo/staging" {
		t.Fatalf("source repository paths = %+v", report.Repository)
	}
	assertContains(t, report.Modules[0].SourceDir, "/repo/staging/src/arcoris.dev/foundation")
}
