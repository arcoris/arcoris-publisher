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

	"github.com/spf13/cobra"
)

func (c CLI) newPlanCommand(output *outputFlags) *cobra.Command {
	var flags planFlags
	cmd := &cobra.Command{
		Use:   "plan",
		Short: "Render the executable publication plan",
		Args:  noArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.executePlan(cmd.Context(), flags, outputForCommand(cmd, *output), cmd.OutOrStdout())
		},
	}
	addPlanFlags(cmd.Flags(), &flags, c.opts)
	return cmd
}

func (c CLI) executePlan(ctx context.Context, flags planFlags, output outputFlags, stdout io.Writer) error {
	version, err := parseVersion(flags.version)
	if err != nil {
		return err
	}

	reportOptions, err := parseReportOptions(output)
	if err != nil {
		return err
	}

	application, err := c.deps.application(c.opts.App)
	if err != nil {
		return err
	}

	plan, err := application.BuildPlan(ctx, flags.manifest, version)
	if err != nil {
		return &Error{Code: CodeUseCaseFailed, Message: "plan failed", Cause: err}
	}

	if err := newRenderer(reportOptions).Plan(stdout, plan); err != nil {
		return &Error{Code: CodeReportFailed, Message: "render plan report failed", Cause: err}
	}
	return nil
}
