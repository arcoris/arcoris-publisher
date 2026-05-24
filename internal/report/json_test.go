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
