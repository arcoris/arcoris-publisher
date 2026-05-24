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
	"fmt"
	"strings"

	"arcoris.dev/arcoris-publisher/internal/report"
	"arcoris.dev/arcoris-publisher/internal/versioning"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type outputFlags struct {
	output            string
	includeLocalPaths bool
	pretty            bool
	compact           bool
	// prettySet and compactSet distinguish explicit user intent from defaults.
	//
	// Pretty defaults to true for human-readable JSON, so treating the raw
	// boolean value as "set" would make --compact conflict with the default.
	prettySet  bool
	compactSet bool
}

type planFlags struct {
	manifest string
	version  string
}

type workflowFlags struct {
	manifest            string
	version             string
	sourceRepositoryDir string
	stagingDir          string
	targetRootDir       string
	dryRun              bool
}

func defaultOutputFlags(opts Options) outputFlags {
	return outputFlags{
		output:            opts.Report.Format.String(),
		includeLocalPaths: opts.Report.IncludeLocalPaths,
		pretty:            opts.Report.Pretty,
	}
}

func addOutputFlags(flags *pflag.FlagSet, values *outputFlags) {
	flags.StringVar(&values.output, "output", values.output, "output format: text or json")
	flags.BoolVar(&values.includeLocalPaths, "include-local-paths", values.includeLocalPaths, "include local absolute filesystem paths in reports")
	flags.BoolVar(&values.pretty, "pretty", values.pretty, "render pretty JSON output when --output=json")
	flags.BoolVar(&values.compact, "compact", values.compact, "render compact JSON output when --output=json")
}

func outputForCommand(cmd *cobra.Command, values outputFlags) outputFlags {
	if flag := cmd.Flag("pretty"); flag != nil {
		values.prettySet = flag.Changed
	}
	if flag := cmd.Flag("compact"); flag != nil {
		values.compactSet = flag.Changed
	}
	return values
}

func addPlanFlags(flags *pflag.FlagSet, values *planFlags, opts Options) {
	values.manifest = opts.ManifestPath
	flags.StringVar(&values.manifest, "manifest", values.manifest, "path to arcpub.yaml")
	flags.StringVar(&values.version, "version", "", "publication version, for example v0.3.0")
}

func addWorkflowFlags(flags *pflag.FlagSet, values *workflowFlags, opts Options, includeDryRun bool) {
	values.manifest = opts.ManifestPath
	values.sourceRepositoryDir = opts.SourceRepositoryDir
	values.stagingDir = opts.StagingDir
	values.targetRootDir = opts.TargetRootDir
	values.dryRun = opts.App.Workflow.DryRun

	flags.StringVar(&values.manifest, "manifest", values.manifest, "path to arcpub.yaml")
	flags.StringVar(&values.version, "version", "", "publication version, for example v0.3.0")
	flags.StringVar(&values.sourceRepositoryDir, "source-repo", values.sourceRepositoryDir, "source Git checkout root")
	flags.StringVar(&values.stagingDir, "staging-dir", values.stagingDir, "staging directory containing module sources")
	flags.StringVar(&values.targetRootDir, "target-root", values.targetRootDir, "directory containing target worktrees")
	if includeDryRun {
		flags.BoolVar(&values.dryRun, "dry-run", values.dryRun, "run through verification without publishing mutations")
	}
}

// parseReportOptions converts CLI output flags into report renderer options.
//
// Text reports deliberately ignore JSON indentation flags. Compact JSON may
// override the default pretty JSON, but an explicit --pretty --compact pair is a
// usage error because it hides user intent.
func parseReportOptions(flags outputFlags) (report.Options, error) {
	if flags.prettySet && flags.compactSet {
		return report.Options{}, &Error{
			Code:    CodeInvalidFlags,
			Message: "--pretty and --compact cannot be used together",
		}
	}

	format, err := report.ParseFormat(flags.output)
	if err != nil {
		return report.Options{}, &Error{Code: CodeInvalidFlags, Message: "invalid output format", Cause: err}
	}

	pretty := flags.pretty
	if flags.compact {
		pretty = false
	}
	if format == report.FormatText {
		pretty = false
	}

	return report.Options{
		Format:            format,
		Pretty:            pretty,
		IncludeLocalPaths: flags.includeLocalPaths,
	}, nil
}

// parseVersion keeps version validation in the CLI layer while leaving version
// syntax and normalization to the versioning package.
func parseVersion(value string) (versioning.Version, error) {
	if strings.TrimSpace(value) == "" {
		return "", &Error{Code: CodeInvalidVersion, Message: "--version is required"}
	}
	version, err := versioning.Parse(value)
	if err != nil {
		return "", &Error{Code: CodeInvalidVersion, Message: "invalid --version", Cause: err}
	}
	return version, nil
}

// flagError adapts pflag parse failures to the CLI error model so Cobra does
// not print duplicate usage/errors on its own.
func flagError(cmd *cobra.Command, err error) error {
	if err == nil {
		return nil
	}
	return &Error{
		Code:    CodeInvalidFlags,
		Message: fmt.Sprintf("invalid flags for %s", cmd.CommandPath()),
		Cause:   err,
	}
}

// noArgs rejects accidental positional arguments for commands whose full input
// surface is flags.
func noArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}
	return &Error{
		Code:    CodeInvalidFlags,
		Message: fmt.Sprintf("%s does not accept arguments", cmd.CommandPath()),
	}
}

func usageError(message string) error {
	return &Error{Code: CodeInvalidFlags, Message: message}
}

// normalizeCobraError turns raw Cobra command-routing errors into typed CLI
// errors while preserving sentinel outcomes produced by command handlers.
func normalizeCobraError(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, errVerificationFailed) {
		return err
	}

	var cliErr *Error
	if errors.As(err, &cliErr) {
		return err
	}
	return &Error{Code: CodeInvalidCommand, Message: strings.TrimSpace(err.Error())}
}
