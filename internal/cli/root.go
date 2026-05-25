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
	"io"

	"github.com/spf13/cobra"
)

// newRootCommand creates a fresh Cobra tree for one Run invocation.
//
// Cobra commands retain flag values and output writers after execution, so the
// CLI never stores a package-level root command. Rebuilding the tree keeps tests
// and concurrent embedders isolated from previous invocations.
func (c CLI) newRootCommand(stdout io.Writer, stderr io.Writer) *cobra.Command {
	output := defaultOutputFlags(c.opts)
	root := &cobra.Command{
		Use:           "arcpub",
		Short:         "Publish ARCORIS staging modules",
		Long:          rootDescription,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	configureCobra(root, stdout, stderr)
	root.CompletionOptions.DisableDefaultCmd = true
	addOutputFlags(root.PersistentFlags(), &output)

	root.AddCommand(c.newPlanCommand(&output))
	root.AddCommand(c.newPreflightCommand(&output))
	root.AddCommand(c.newVerifyCommand(&output))
	root.AddCommand(c.newPublishCommand(&output))
	root.AddCommand(c.newTransactionsCommand(&output))
	root.AddCommand(c.newRollbackCommand(&output))
	root.AddCommand(c.newVersionCommand(&output))
	root.AddCommand(c.newCompletionCommand())

	return root
}

// configureCobra routes all Cobra output through caller-owned writers and lets
// CLI.Run decide how errors become process exit codes.
func configureCobra(cmd *cobra.Command, stdout io.Writer, stderr io.Writer) {
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		return flagError(cmd, err)
	})
}

const rootDescription = `ARCORIS Publisher builds explicit publication plans, verifies target
repositories, and publishes resolved module projections.`
