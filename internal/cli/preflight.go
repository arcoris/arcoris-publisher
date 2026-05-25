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

func (c CLI) newPreflightCommand(output *outputFlags) *cobra.Command {
	var flags workflowFlags
	cmd := &cobra.Command{
		Use:   "preflight",
		Short: "Check whether publish can safely start without mutating targets",
		Long: `Preflight validates current local and remote state without constructing target
worktrees, rewriting module files, writing transaction state, or pushing refs.
It does not reserve refs; publish can still fail later if state changes.`,
		Args: noArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.executePreflight(cmd.Context(), flags, outputForCommand(cmd, *output), cmd.OutOrStdout())
		},
	}
	addWorkflowFlags(cmd.Flags(), &flags, c.opts, false)
	cmd.Flags().StringVar(&flags.stateDir, "state-dir", flags.stateDir, "publish transaction state directory")
	return cmd
}

func (c CLI) executePreflight(ctx context.Context, flags workflowFlags, output outputFlags, stdout io.Writer) error {
	version, err := parseVersion(flags.version)
	if err != nil {
		return err
	}

	reportOptions, err := parseReportOptions(output)
	if err != nil {
		return err
	}

	appOptions := c.opts.App
	appOptions.Workflow.Preflight.StateDir = publishStateDir(flags)
	appOptions.Workflow.Preflight.RemoteName = appOptions.Workflow.Publish.RemoteName

	application, err := c.deps.application(appOptions)
	if err != nil {
		return err
	}

	result, err := application.Preflight(ctx, app.Request{
		ManifestPath:        flags.manifest,
		Version:             version,
		SourceRepositoryDir: flags.sourceRepositoryDir,
		StagingDir:          flags.stagingDir,
		TargetRootDir:       flags.targetRootDir,
	})
	if err != nil {
		if !result.Preflight().Empty() {
			_ = newRenderer(reportOptions).Preflight(stdout, result.Preflight())
		}
		return &Error{Code: CodeUseCaseFailed, Message: "preflight failed", Cause: err}
	}

	preflightResult := result.Preflight()
	if err := newRenderer(reportOptions).Preflight(stdout, preflightResult); err != nil {
		return &Error{Code: CodeReportFailed, Message: "render preflight report failed", Cause: err}
	}
	if preflightResult.Failed() {
		return &Error{Code: CodeUseCaseFailed, Message: "preflight failed"}
	}
	return nil
}
