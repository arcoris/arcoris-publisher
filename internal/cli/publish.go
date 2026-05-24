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

	"arcoris.dev/arcoris-publisher/internal/app"
)

func (c CLI) runPublish(ctx context.Context, args []string, stdout io.Writer, stderr io.Writer) int {
	var flags workflowFlags
	fs := newFlagSet("publish")
	addWorkflowFlags(fs, &flags, c.opts, true)
	if err := parseFlagSet(fs, args); err != nil {
		writeCLIError(stderr, err)
		return ExitUsage
	}

	version, err := parseVersion(flags.version)
	if err != nil {
		writeCLIError(stderr, err)
		return ExitUsage
	}

	reportOptions, err := parseReportOptions(flags.commonFlags)
	if err != nil {
		writeCLIError(stderr, err)
		return ExitUsage
	}

	appOptions := c.opts.App
	appOptions.Workflow.DryRun = flags.dryRun

	application, err := c.deps.application(appOptions)
	if err != nil {
		writeCLIError(stderr, err)
		return ExitError
	}

	result, err := application.Publish(ctx, app.Request{
		ManifestPath:        flags.manifest,
		Version:             version,
		SourceRepositoryDir: flags.sourceRepositoryDir,
		StagingDir:          flags.stagingDir,
		TargetRootDir:       flags.targetRootDir,
	})
	if err != nil {
		writeCLIError(stderr, &Error{Code: CodeUseCaseFailed, Message: "publish failed", Cause: err})
		return ExitError
	}

	workflowResult := result.Workflow()
	if err := newRenderer(reportOptions).Workflow(stdout, workflowResult); err != nil {
		writeCLIError(stderr, &Error{Code: CodeReportFailed, Message: "render workflow report failed", Cause: err})
		return ExitError
	}
	if workflowResult.Verify().Failed() {
		return ExitVerificationFailed
	}
	return ExitOK
}
