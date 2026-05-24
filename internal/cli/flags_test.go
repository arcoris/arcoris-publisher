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
	"errors"
	"testing"

	"arcoris.dev/arcoris-publisher/internal/report"
)

func TestParseReportOptions(t *testing.T) {
	t.Parallel()

	opts, err := parseReportOptions(commonFlags{output: " JSON ", pretty: true, includeLocalPaths: true})
	if err != nil {
		t.Fatalf("parseReportOptions() error = %v", err)
	}
	if opts.Format != report.FormatJSON || !opts.Pretty || !opts.IncludeLocalPaths {
		t.Fatalf("parseReportOptions() = %+v", opts)
	}
}

func TestParseReportOptionsRejectsInvalidOutput(t *testing.T) {
	t.Parallel()

	_, err := parseReportOptions(commonFlags{output: "yaml"})
	if err == nil {
		t.Fatal("parseReportOptions() error = nil")
	}

	var cliErr *Error
	if !errors.As(err, &cliErr) || cliErr.Code != CodeInvalidFlags {
		t.Fatalf("parseReportOptions() error = %v", err)
	}
}

func TestParseReportOptionsCompactOverridesPretty(t *testing.T) {
	t.Parallel()

	opts, err := parseReportOptions(commonFlags{output: "json", pretty: true, compact: true})
	if err != nil {
		t.Fatalf("parseReportOptions() error = %v", err)
	}
	if opts.Pretty {
		t.Fatalf("parseReportOptions() Pretty = true")
	}
}

func TestParseVersion(t *testing.T) {
	t.Parallel()

	version, err := parseVersion("v1.2.3")
	if err != nil {
		t.Fatalf("parseVersion() error = %v", err)
	}
	if version.String() != "v1.2.3" {
		t.Fatalf("parseVersion() = %q", version)
	}
}

func TestParseVersionRequiresValue(t *testing.T) {
	t.Parallel()

	_, err := parseVersion("")
	if err == nil {
		t.Fatal("parseVersion() error = nil")
	}
	var cliErr *Error
	if !errors.As(err, &cliErr) || cliErr.Code != CodeInvalidVersion {
		t.Fatalf("parseVersion() error = %v", err)
	}
}

func TestAddWorkflowFlagsDryRunIsPublishOnly(t *testing.T) {
	t.Parallel()

	var verifyFlags workflowFlags
	verifySet := newFlagSet("verify")
	addWorkflowFlags(verifySet, &verifyFlags, DefaultOptions(), false)
	if err := verifySet.Parse([]string{"--dry-run"}); err == nil {
		t.Fatal("verify flags accepted --dry-run")
	}

	var publishFlags workflowFlags
	publishSet := newFlagSet("publish")
	addWorkflowFlags(publishSet, &publishFlags, DefaultOptions(), true)
	if err := publishSet.Parse([]string{"--dry-run"}); err != nil {
		t.Fatalf("publish flags rejected --dry-run: %v", err)
	}
	if !publishFlags.dryRun {
		t.Fatal("publish dry-run flag was not set")
	}
}
