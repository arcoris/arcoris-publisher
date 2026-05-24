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

package cli

import (
	"bytes"
	"strings"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/report"
)

func TestWriteUsageListsCommands(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	writeUsage(&buf)

	for _, want := range []string{
		"arcpub plan",
		"arcpub verify",
		"arcpub publish",
		"arcpub version",
		"arcpub help",
	} {
		if !strings.Contains(buf.String(), want) {
			t.Fatalf("usage missing %q:\n%s", want, buf.String())
		}
	}
}

func TestNewRendererUsesCLIReportOptions(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	err := newRenderer(report.Options{Format: report.FormatJSON, Pretty: true}).Plan(
		&buf,
		newFakeApplication(t).plan,
	)

	if err != nil {
		t.Fatalf("renderer.Plan() error = %v", err)
	}
	if !strings.Contains(buf.String(), `"kind": "plan"`) ||
		!strings.Contains(buf.String(), "\n  ") {
		t.Fatalf("renderer output = %q", buf.String())
	}
}
