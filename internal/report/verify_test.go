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

	"arcoris.dev/arcoris-publisher/internal/workflow/verify"
)

func TestBuildVerifyReportZeroValue(t *testing.T) {
	t.Parallel()

	report := BuildVerifyReport(verify.Result{}, Options{})
	if report.Kind != "verify" {
		t.Fatalf("Kind = %q", report.Kind)
	}
	if report.Status != "empty" {
		t.Fatalf("Status = %q", report.Status)
	}
	if report.ModuleCount != 0 || report.FailedCount != 0 {
		t.Fatalf("unexpected counts: %+v", report)
	}
}

func TestBuildVerifyReportPassedAndFailedStatuses(t *testing.T) {
	t.Parallel()

	passed := BuildVerifyReport(
		reportWorkflowResult(t, workflowReportFixture{}).Verify(),
		Options{},
	)
	if passed.Status != "passed" || passed.FailedCount != 0 {
		t.Fatalf("passed report = %+v", passed)
	}

	failed := BuildVerifyReport(
		reportWorkflowResult(t, workflowReportFixture{verifyFails: true}).Verify(),
		Options{},
	)
	if failed.Status != "failed" || failed.FailedCount == 0 {
		t.Fatalf("failed report = %+v", failed)
	}
}

func TestRendererVerifyFailureIsReportContent(t *testing.T) {
	t.Parallel()

	result := reportWorkflowResult(t, workflowReportFixture{verifyFails: true}).Verify()
	text := renderReport(t, func(buf *bytes.Buffer) error {
		return New(Options{Format: FormatText}).Verify(buf, result)
	})
	assertContains(t, text, "Verification", "Status: failed", "tidy failed")
}

func TestRendererVerifyTextZeroValue(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := New(Options{Format: FormatText}).Verify(&buf, verify.Result{}); err != nil {
		t.Fatalf("Verify(text) error = %v", err)
	}
	if !strings.Contains(buf.String(), "Verification") {
		t.Fatalf("text report missing header:\n%s", buf.String())
	}
}
