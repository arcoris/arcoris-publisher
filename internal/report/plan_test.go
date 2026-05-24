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

func TestBuildPlanReport(t *testing.T) {
	t.Parallel()

	p := reportPlan(t).Plan
	report := BuildPlanReport(p, Options{})

	if report.Kind != "plan" {
		t.Fatalf("Kind = %q", report.Kind)
	}
	if report.ModuleCount != 2 {
		t.Fatalf("ModuleCount = %d", report.ModuleCount)
	}
	if got := report.Modules[0].Name; got != "foundation" {
		t.Fatalf("first module = %q", got)
	}
	if got := report.Modules[1].Name; got != "control" {
		t.Fatalf("second module = %q", got)
	}
	if len(report.Modules[1].Requirements) != 1 {
		t.Fatalf("control requirements = %+v", report.Modules[1].Requirements)
	}
	if report.Modules[1].Requirements[0].ModulePath != "arcoris.dev/foundation" {
		t.Fatalf("control requirement = %+v", report.Modules[1].Requirements[0])
	}
	if len(report.Modules[0].Branches) == 0 {
		t.Fatalf("foundation branches = %+v", report.Modules[0].Branches)
	}
	if len(report.Modules[0].PublishEntries) != 2 {
		t.Fatalf("foundation entries = %+v", report.Modules[0].PublishEntries)
	}
}

func TestRendererPlanJSON(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	err := New(Options{Format: FormatJSON, Pretty: true}).Plan(&buf, reportPlan(t).Plan)
	if err != nil {
		t.Fatalf("Plan(JSON) error = %v", err)
	}
	if !strings.HasSuffix(buf.String(), "\n") {
		t.Fatalf("JSON missing trailing newline: %q", buf.String())
	}

	var decoded PlanReport
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if decoded.ModuleCount != 2 {
		t.Fatalf("decoded.ModuleCount = %d", decoded.ModuleCount)
	}
}

func TestRendererPlanText(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	err := New(Options{Format: FormatText}).Plan(&buf, reportPlan(t).Plan)
	if err != nil {
		t.Fatalf("Plan(text) error = %v", err)
	}
	text := buf.String()
	for _, want := range []string{"Plan", "foundation", "control", "arcoris.dev/foundation", "v0.3.0"} {
		if !strings.Contains(text, want) {
			t.Fatalf("text report missing %q:\n%s", want, text)
		}
	}
}
