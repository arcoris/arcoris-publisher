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
	"context"
	"io"
)

func (c CLI) runPlan(ctx context.Context, args []string, stdout io.Writer, stderr io.Writer) int {
	var flags commonFlags
	fs := newFlagSet("plan")
	addCommonFlags(fs, &flags, c.opts, true)
	if err := parseFlagSet(fs, args); err != nil {
		writeCLIError(stderr, err)
		return ExitUsage
	}

	version, err := parseVersion(flags.version)
	if err != nil {
		writeCLIError(stderr, err)
		return ExitUsage
	}

	reportOptions, err := parseReportOptions(flags)
	if err != nil {
		writeCLIError(stderr, err)
		return ExitUsage
	}

	application, err := c.deps.application(c.opts.App)
	if err != nil {
		writeCLIError(stderr, err)
		return ExitError
	}

	plan, err := application.BuildPlan(ctx, flags.manifest, version)
	if err != nil {
		writeCLIError(stderr, &Error{Code: CodeUseCaseFailed, Message: "plan failed", Cause: err})
		return ExitError
	}

	if err := newRenderer(reportOptions).Plan(stdout, plan); err != nil {
		writeCLIError(stderr, &Error{Code: CodeReportFailed, Message: "render plan report failed", Cause: err})
		return ExitError
	}
	return ExitOK
}
