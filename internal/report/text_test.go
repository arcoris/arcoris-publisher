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
)

func TestTextRenderingIsNewlineTerminatedAndDeterministic(t *testing.T) {
	t.Parallel()

	p := reportPlan(t).Plan
	first := renderReport(t, func(buf *bytes.Buffer) error {
		return New(Options{Format: FormatText}).Plan(buf, p)
	})
	second := renderReport(t, func(buf *bytes.Buffer) error {
		return New(Options{Format: FormatText}).Plan(buf, p)
	})

	if !strings.HasSuffix(first, "\n") {
		t.Fatalf("text missing trailing newline: %q", first)
	}
	if first != second {
		t.Fatalf("text output differs:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

func TestTextRenderingHasNoTerminalDecorationsOrLocalPaths(t *testing.T) {
	t.Parallel()

	text := renderReport(t, func(buf *bytes.Buffer) error {
		return New(Options{Format: FormatText}).Workflow(
			buf,
			reportWorkflowResult(t, workflowReportFixture{}),
		)
	})

	if strings.Contains(text, "\x1b[") {
		t.Fatalf("text contains ANSI escape sequence:\n%s", text)
	}
	assertNoLocalPaths(t, text)
}
