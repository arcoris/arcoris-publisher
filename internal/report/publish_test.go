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
	"bytes"
	"strings"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/workflow/publish"
)

func TestBuildPublishReportZeroValue(t *testing.T) {
	t.Parallel()

	report := BuildPublishReport(publish.Result{}, Options{})
	if report.Kind != "publish" {
		t.Fatalf("Kind = %q", report.Kind)
	}
	if report.Status != "empty" {
		t.Fatalf("Status = %q", report.Status)
	}
	if report.ModuleCount != 0 || report.PublishedCount != 0 || report.SkippedCount != 0 {
		t.Fatalf("unexpected counts: %+v", report)
	}
}

func TestBuildPublishReportPublishedAndSkippedStatuses(t *testing.T) {
	t.Parallel()

	published := BuildPublishReport(
		reportWorkflowResult(
			t,
			workflowReportFixture{
				publish:      true,
				dirtyModules: []string{"foundation", "control"},
			},
		).Publish(),
		Options{},
	)
	if published.Status != "published" || published.PublishedCount != 2 || published.SkippedCount != 0 {
		t.Fatalf("published report = %+v", published)
	}
	if len(published.Modules[0].Tags) == 0 {
		t.Fatalf("published tags = %+v", published.Modules[0].Tags)
	}

	skipped := BuildPublishReport(
		reportWorkflowResult(t, workflowReportFixture{publish: true}).Publish(),
		Options{},
	)
	if skipped.Status != "skipped" || skipped.PublishedCount != 0 || skipped.SkippedCount != 2 {
		t.Fatalf("skipped report = %+v", skipped)
	}
}

func TestRendererPublishSkippedIsReportContent(t *testing.T) {
	t.Parallel()

	result := reportWorkflowResult(t, workflowReportFixture{publish: true}).Publish()
	text := renderReport(t, func(buf *bytes.Buffer) error {
		return New(Options{Format: FormatText}).Publish(buf, result)
	})
	assertContains(t, text, "Publication", "Status: skipped", "foundation: skipped")
}

func TestRendererPublishJSONZeroValue(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := New(Options{Format: FormatJSON, Pretty: true}).Publish(&buf, publish.Result{}); err != nil {
		t.Fatalf("Publish(JSON) error = %v", err)
	}
	if !strings.Contains(buf.String(), `"kind": "publish"`) {
		t.Fatalf("json report missing kind:\n%s", buf.String())
	}
}
