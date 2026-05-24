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

func TestBuildModuleFileReportHidesGoModPathByDefault(t *testing.T) {
	t.Parallel()

	report := BuildModuleFileReport(reportWorkflowResult(t, workflowReportFixture{}).ModuleFile(), Options{})

	if !report.Present || !report.Changed || report.ModuleCount != 2 {
		t.Fatalf("modulefile report = %+v", report)
	}
	if report.Modules[0].GoModPath != "" {
		t.Fatalf("modulefile report leaked go.mod path: %+v", report.Modules[0])
	}
	if len(report.Modules[1].Requirements) != 1 {
		t.Fatalf("control requirements = %+v", report.Modules[1].Requirements)
	}
}

func TestBuildModuleFileReportCanIncludeLocalPaths(t *testing.T) {
	t.Parallel()

	report := BuildModuleFileReport(
		reportWorkflowResult(t, workflowReportFixture{}).ModuleFile(),
		Options{IncludeLocalPaths: true},
	)

	if report.Modules[0].GoModPath != "/target/arcoris__foundation/go.mod" {
		t.Fatalf("go.mod path = %+v", report.Modules[0])
	}
}
