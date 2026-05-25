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
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"arcoris.dev/arcoris-publisher/internal/app"
	"github.com/spf13/cobra"
)

type transactionFlags struct {
	targetRootDir string
	stateDir      string
}

type transactionPruneFlags struct {
	statuses  []string
	olderThan string
	dryRun    bool
}

func (c CLI) newTransactionsCommand(output *outputFlags) *cobra.Command {
	var flags transactionFlags
	cmd := &cobra.Command{
		Use:   "transactions",
		Short: "Inspect durable publish transactions",
	}
	cmd.AddCommand(c.newTransactionsListCommand(&flags, output))
	cmd.AddCommand(c.newTransactionsShowCommand(&flags, output))
	cmd.AddCommand(c.newTransactionsPruneCommand(&flags, output))
	addTransactionFlags(cmd.PersistentFlags(), &flags, c.opts)
	return cmd
}

func (c CLI) newTransactionsListCommand(flags *transactionFlags, output *outputFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List publish transactions",
		Args:  noArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.executeTransactionsList(cmd.Context(), *flags, outputForCommand(cmd, *output), cmd.OutOrStdout())
		},
	}
}

func (c CLI) newTransactionsShowCommand(flags *transactionFlags, output *outputFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "show <id>",
		Short: "Show one publish transaction",
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) != 1 {
				return usageError("transaction id is required")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.executeTransactionsShow(cmd.Context(), *flags, outputForCommand(cmd, *output), args[0], cmd.OutOrStdout())
		},
	}
}

func (c CLI) newTransactionsPruneCommand(flags *transactionFlags, output *outputFlags) *cobra.Command {
	var pruneFlags transactionPruneFlags
	cmd := &cobra.Command{
		Use:   "prune",
		Short: "Prune old terminal publish transaction journals",
		Long: "Prune removes only terminal publish transaction journals selected by explicit filters. " +
			"Pending, failed, rolling_back, and rollback_failed journals are preserved; publish locks are never removed.",
		Args: noArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.executeTransactionsPrune(cmd.Context(), *flags, pruneFlags, outputForCommand(cmd, *output), cmd.OutOrStdout())
		},
	}
	cmd.Flags().StringArrayVar(&pruneFlags.statuses, "status", nil, "terminal status to prune: committed or rolled_back; may be repeated or comma-separated")
	cmd.Flags().StringVar(&pruneFlags.olderThan, "older-than", "", "only prune journals older than this age, for example 720h or 30d")
	cmd.Flags().BoolVar(&pruneFlags.dryRun, "dry-run", false, "preview matching terminal journals without deleting them")
	return cmd
}

func (c CLI) newRollbackCommand(output *outputFlags) *cobra.Command {
	var flags transactionFlags
	var id string
	cmd := &cobra.Command{
		Use:   "rollback --transaction <id>",
		Short: "Roll back a failed publish transaction",
		Args:  noArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if strings.TrimSpace(id) == "" {
				return usageError("--transaction is required")
			}
			return c.executeRollback(cmd.Context(), flags, outputForCommand(cmd, *output), id, cmd.OutOrStdout())
		},
	}
	addTransactionFlags(cmd.Flags(), &flags, c.opts)
	cmd.Flags().StringVar(&id, "transaction", "", "publish transaction id to roll back")
	return cmd
}

func addTransactionFlags(flags flagSet, values *transactionFlags, opts Options) {
	values.targetRootDir = opts.TargetRootDir
	flags.StringVar(&values.targetRootDir, "target-root", values.targetRootDir, "directory containing target worktrees")
	flags.StringVar(&values.stateDir, "state-dir", values.stateDir, "publish transaction state directory")
}

type flagSet interface {
	StringVar(*string, string, string, string)
}

func (c CLI) executeTransactionsList(ctx context.Context, flags transactionFlags, output outputFlags, stdout io.Writer) error {
	reportOptions, err := parseReportOptions(output)
	if err != nil {
		return err
	}
	application, err := c.deps.application(c.opts.App)
	if err != nil {
		return err
	}
	result, err := application.ListTransactions(ctx, app.TransactionRequest{StateDir: transactionStateDir(flags)})
	if err != nil {
		return &Error{Code: CodeUseCaseFailed, Message: "list transactions failed", Cause: err}
	}
	return newRenderer(reportOptions).TransactionList(stdout, result.Summaries())
}

func (c CLI) executeTransactionsShow(ctx context.Context, flags transactionFlags, output outputFlags, id string, stdout io.Writer) error {
	reportOptions, err := parseReportOptions(output)
	if err != nil {
		return err
	}
	application, err := c.deps.application(c.opts.App)
	if err != nil {
		return err
	}
	result, err := application.ShowTransaction(ctx, app.TransactionRequest{
		StateDir:      transactionStateDir(flags),
		TransactionID: app.TransactionID(id),
	})
	if err != nil {
		return &Error{Code: CodeUseCaseFailed, Message: "show transaction failed", Cause: err}
	}
	return newRenderer(reportOptions).Transaction(stdout, result.Journal())
}

func (c CLI) executeRollback(ctx context.Context, flags transactionFlags, output outputFlags, id string, stdout io.Writer) error {
	reportOptions, err := parseReportOptions(output)
	if err != nil {
		return err
	}
	application, err := c.deps.application(c.opts.App)
	if err != nil {
		return err
	}
	result, err := application.RollbackTransaction(ctx, app.TransactionRequest{
		StateDir:      transactionStateDir(flags),
		TransactionID: app.TransactionID(id),
	})
	if renderErr := newRenderer(reportOptions).Transaction(stdout, result.Journal()); renderErr != nil {
		return &Error{Code: CodeReportFailed, Message: "render rollback report failed", Cause: renderErr}
	}
	if err != nil {
		return &Error{Code: CodeUseCaseFailed, Message: "rollback transaction failed", Cause: err}
	}
	return nil
}

func (c CLI) executeTransactionsPrune(
	ctx context.Context,
	flags transactionFlags,
	pruneFlags transactionPruneFlags,
	output outputFlags,
	stdout io.Writer,
) error {
	reportOptions, err := parseReportOptions(output)
	if err != nil {
		return err
	}
	statuses, err := parsePruneStatuses(pruneFlags.statuses)
	if err != nil {
		return err
	}
	olderThan, err := parsePruneDuration(pruneFlags.olderThan)
	if err != nil {
		return err
	}
	if !pruneFlags.dryRun && len(statuses) == 0 && olderThan <= 0 {
		return usageError("refusing to prune without an explicit --status or --older-than filter")
	}

	application, err := c.deps.application(c.opts.App)
	if err != nil {
		return err
	}
	result, err := application.PruneTransactions(ctx, app.TransactionPruneRequest{
		StateDir:  transactionStateDir(flags),
		Statuses:  statuses,
		OlderThan: olderThan,
		DryRun:    pruneFlags.dryRun,
	})
	if renderErr := newRenderer(reportOptions).TransactionPrune(stdout, result.Result()); renderErr != nil {
		return &Error{Code: CodeReportFailed, Message: "render transaction prune report failed", Cause: renderErr}
	}
	if err != nil {
		return &Error{Code: CodeUseCaseFailed, Message: "prune transactions failed", Cause: err}
	}
	return nil
}

func transactionStateDir(flags transactionFlags) string {
	if strings.TrimSpace(flags.stateDir) != "" {
		return flags.stateDir
	}
	return filepath.Join(flags.targetRootDir, ".arcpub", "state")
}

func parsePruneStatuses(values []string) ([]app.TransactionStatus, error) {
	var out []app.TransactionStatus
	for _, value := range values {
		for _, part := range strings.Split(value, ",") {
			status := strings.TrimSpace(strings.ToLower(part))
			if status == "" {
				continue
			}
			switch app.TransactionStatus(status) {
			case app.TransactionStatusCommitted, app.TransactionStatusRolledBack:
				out = append(out, app.TransactionStatus(status))
			default:
				return nil, usageError("invalid --status; expected committed or rolled_back")
			}
		}
	}
	return out, nil
}

func parsePruneDuration(value string) (time.Duration, error) {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return 0, nil
	}
	if strings.HasSuffix(value, "d") {
		days, err := strconv.ParseFloat(strings.TrimSuffix(value, "d"), 64)
		if err != nil || days < 0 {
			return 0, usageError("invalid --older-than; expected Go duration such as 720h or day duration such as 30d")
		}
		return time.Duration(days * float64(24*time.Hour)), nil
	}
	duration, err := time.ParseDuration(value)
	if err != nil || duration < 0 {
		return 0, usageError("invalid --older-than; expected Go duration such as 720h or day duration such as 30d")
	}
	return duration, nil
}
