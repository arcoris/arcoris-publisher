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
	"github.com/spf13/cobra"
)

func (c CLI) newVerifyCommand(output *outputFlags) *cobra.Command {
	var flags workflowFlags
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Construct and verify target module worktrees without publishing",
		Args:  noArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.executeVerify(cmd.Context(), flags, outputForCommand(cmd, *output), cmd.OutOrStdout())
		},
	}
	addWorkflowFlags(cmd.Flags(), &flags, c.opts, false)
	return cmd
}

func (c CLI) executeVerify(ctx context.Context, flags workflowFlags, output outputFlags, stdout io.Writer) error {
	version, err := parseVersion(flags.version)
	if err != nil {
		return err
	}

	reportOptions, err := parseReportOptions(output)
	if err != nil {
		return err
	}

	appOptions := c.opts.App

	application, err := c.deps.application(appOptions)
	if err != nil {
		return err
	}

	result, err := application.Verify(ctx, app.Request{
		ManifestPath:        flags.manifest,
		Version:             version,
		SourceRepositoryDir: flags.sourceRepositoryDir,
		StagingDir:          flags.stagingDir,
		TargetRootDir:       flags.targetRootDir,
	})
	if err != nil {
		return &Error{Code: CodeUseCaseFailed, Message: "verify failed", Cause: err}
	}

	workflowResult := result.Workflow()
	if err := newRenderer(reportOptions).Workflow(stdout, workflowResult); err != nil {
		return &Error{Code: CodeReportFailed, Message: "render workflow report failed", Cause: err}
	}
	if workflowResult.Verify().Failed() {
		return errVerificationFailed
	}
	return nil
}
