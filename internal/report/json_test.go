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
	"encoding/json"
	"strings"
	"testing"
)

func TestWriteJSONProducesValidNewlineTerminatedOutput(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := writeJSON(&buf, BuildPlanReport(reportPlan(t).Plan, Options{}), true); err != nil {
		t.Fatalf("writeJSON() error = %v", err)
	}
	if !strings.HasSuffix(buf.String(), "\n") {
		t.Fatalf("JSON missing trailing newline: %q", buf.String())
	}
	if !json.Valid(buf.Bytes()) {
		t.Fatalf("invalid JSON:\n%s", buf.String())
	}
}

func TestWriteJSONIsDeterministic(t *testing.T) {
	t.Parallel()

	value := BuildPlanReport(reportPlan(t).Plan, Options{})
	var first bytes.Buffer
	var second bytes.Buffer

	if err := writeJSON(&first, value, true); err != nil {
		t.Fatalf("first writeJSON() error = %v", err)
	}
	if err := writeJSON(&second, value, true); err != nil {
		t.Fatalf("second writeJSON() error = %v", err)
	}
	if first.String() != second.String() {
		t.Fatalf("JSON output differs:\nfirst:\n%s\nsecond:\n%s", first.String(), second.String())
	}
}

func TestWriteJSONSupportsCompactOutput(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	if err := writeJSON(&buf, BuildPlanReport(reportPlan(t).Plan, Options{}), false); err != nil {
		t.Fatalf("writeJSON() error = %v", err)
	}
	if strings.Contains(buf.String(), "\n  ") {
		t.Fatalf("compact JSON contains indentation:\n%s", buf.String())
	}
}

func TestRendererJSONKindsAndDeterminism(t *testing.T) {
	t.Parallel()

	result := reportWorkflowResult(t, workflowReportFixture{publish: true})
	renderer := New(Options{Format: FormatJSON, Pretty: true})

	tests := []struct {
		name string
		kind string
		run  func(*bytes.Buffer) error
	}{
		{name: "plan", kind: "plan", run: func(buf *bytes.Buffer) error {
			return renderer.Plan(buf, reportPlan(t).Plan)
		}},
		{name: "workflow", kind: "workflow", run: func(buf *bytes.Buffer) error {
			return renderer.Workflow(buf, result)
		}},
		{name: "verify", kind: "verify", run: func(buf *bytes.Buffer) error {
			return renderer.Verify(buf, result.Verify())
		}},
		{name: "publish", kind: "publish", run: func(buf *bytes.Buffer) error {
			return renderer.Publish(buf, result.Publish())
		}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			first := renderReport(t, tt.run)
			second := renderReport(t, tt.run)
			if first != second {
				t.Fatalf("JSON output differs:\nfirst:\n%s\nsecond:\n%s", first, second)
			}
			assertTrailingNewline(t, first)
			assertNoLocalPaths(t, first)

			var decoded struct {
				Kind string `json:"kind"`
			}
			if err := json.Unmarshal([]byte(first), &decoded); err != nil {
				t.Fatalf("json.Unmarshal() error = %v", err)
			}
			if decoded.Kind != tt.kind {
				t.Fatalf("kind = %q, want %q", decoded.Kind, tt.kind)
			}
		})
	}
}

func TestRendererJSONPrettyOption(t *testing.T) {
	t.Parallel()

	compact := renderReport(t, func(buf *bytes.Buffer) error {
		return New(Options{Format: FormatJSON}).Plan(buf, reportPlan(t).Plan)
	})
	pretty := renderReport(t, func(buf *bytes.Buffer) error {
		return New(Options{Format: FormatJSON, Pretty: true}).Plan(buf, reportPlan(t).Plan)
	})

	if strings.Contains(compact, "\n  ") {
		t.Fatalf("compact JSON contains indentation:\n%s", compact)
	}
	if !strings.Contains(pretty, "\n  ") {
		t.Fatalf("pretty JSON missing indentation:\n%s", pretty)
	}
}
