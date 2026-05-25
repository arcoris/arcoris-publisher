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
	"arcoris.dev/arcoris-publisher/internal/manifest"
	"github.com/spf13/cobra"
)

type targetPrepareFlags struct {
	manifest       string
	version        string
	targetRootDir  string
	remoteTemplate string
}

func (c CLI) newTargetCommand(output *outputFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "target",
		Short: "Manage target publication worktrees",
	}
	cmd.AddCommand(c.newTargetPrepareCommand(output))
	return cmd
}

func (c CLI) newTargetPrepareCommand(output *outputFlags) *cobra.Command {
	var flags targetPrepareFlags
	cmd := &cobra.Command{
		Use:   "prepare",
		Short: "Prepare target Git worktrees for a publication plan",
		Long: `Prepare target Git worktrees for a publication plan.

The command may create the target root, clone missing worktrees, add a missing
remote, fetch remote state, and checkout target branches. It does not construct
publication files, rewrite go.mod, write provenance, commit, tag, push, or write
publish transaction journals.`,
		Args: noArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.executeTargetPrepare(cmd.Context(), flags, outputForCommand(cmd, *output), cmd.OutOrStdout())
		},
	}
	flags.manifest = c.opts.ManifestPath
	flags.targetRootDir = c.opts.TargetRootDir
	cmd.Flags().StringVar(&flags.manifest, "manifest", flags.manifest, "path to arcpub.yaml")
	cmd.Flags().StringVar(&flags.version, "version", "", "publication version, for example v0.3.0")
	cmd.Flags().StringVar(&flags.targetRootDir, "target-root", flags.targetRootDir, "directory containing target worktrees")
	cmd.Flags().StringVar(&flags.remoteTemplate, "remote-template", "", "Git remote URL template using {repository}, {owner}, or {name}")
	return cmd
}

func (c CLI) executeTargetPrepare(ctx context.Context, flags targetPrepareFlags, output outputFlags, stdout io.Writer) error {
	version, err := parseVersion(flags.version)
	if err != nil {
		return err
	}
	if flags.remoteTemplate != "" {
		if _, err := manifest.ParseRemoteTemplate(flags.remoteTemplate); err != nil {
			return &Error{Code: CodeInvalidFlags, Message: "invalid --remote-template", Cause: err}
		}
	}

	reportOptions, err := parseReportOptions(output)
	if err != nil {
		return err
	}

	application, err := c.deps.application(c.opts.App)
	if err != nil {
		return err
	}
	result, err := application.PrepareTargets(ctx, app.Request{
		ManifestPath:         flags.manifest,
		Version:              version,
		TargetRootDir:        flags.targetRootDir,
		TargetRemoteTemplate: flags.remoteTemplate,
	})
	prepareResult := result.TargetPrepare()
	if !prepareResult.Empty() {
		if renderErr := newRenderer(reportOptions).TargetPrepare(stdout, prepareResult); renderErr != nil {
			return &Error{Code: CodeReportFailed, Message: "render target prepare report failed", Cause: renderErr}
		}
	}
	if err != nil {
		return &Error{Code: CodeUseCaseFailed, Message: "target prepare failed", Cause: err}
	}
	if prepareResult.Failed() {
		return &Error{Code: CodeUseCaseFailed, Message: "target prepare failed"}
	}
	return nil
}
