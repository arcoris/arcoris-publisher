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
	"flag"
	"fmt"
	"io"
	"strings"

	"arcoris.dev/arcoris-publisher/internal/report"
	"arcoris.dev/arcoris-publisher/internal/versioning"
)

type commonFlags struct {
	manifest          string
	version           string
	output            string
	includeLocalPaths bool
	pretty            bool
	compact           bool
}

type workflowFlags struct {
	commonFlags
	sourceRepositoryDir string
	stagingDir          string
	targetRootDir       string
	dryRun              bool
}

func newFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet("arcpub "+name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	return fs
}

func addCommonFlags(fs *flag.FlagSet, flags *commonFlags, opts Options, requireVersion bool) {
	fs.StringVar(&flags.manifest, "manifest", opts.ManifestPath, "path to arcpub.yaml")
	if requireVersion {
		fs.StringVar(&flags.version, "version", "", "publication version, for example v0.3.0")
	}
	fs.StringVar(&flags.output, "output", opts.Report.Format.String(), "output format: text or json")
	fs.BoolVar(&flags.includeLocalPaths, "include-local-paths", opts.Report.IncludeLocalPaths, "include local absolute filesystem paths in reports")
	fs.BoolVar(&flags.pretty, "pretty", opts.Report.Pretty, "render pretty JSON output when --output=json")
	fs.BoolVar(&flags.compact, "compact", false, "render compact JSON output when --output=json")
}

func addWorkflowFlags(fs *flag.FlagSet, flags *workflowFlags, opts Options, includeDryRun bool) {
	addCommonFlags(fs, &flags.commonFlags, opts, true)
	fs.StringVar(&flags.sourceRepositoryDir, "source-repo", opts.SourceRepositoryDir, "source Git checkout root")
	fs.StringVar(&flags.stagingDir, "staging-dir", opts.StagingDir, "staging directory containing module sources")
	fs.StringVar(&flags.targetRootDir, "target-root", opts.TargetRootDir, "directory containing target worktrees")
	if includeDryRun {
		fs.BoolVar(&flags.dryRun, "dry-run", opts.App.Workflow.DryRun, "run through verification without publishing mutations")
	}
}

func parseReportOptions(flags commonFlags) (report.Options, error) {
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

func parseFlagSet(fs *flag.FlagSet, args []string) error {
	if err := fs.Parse(args); err != nil {
		return &Error{Code: CodeInvalidFlags, Message: fmt.Sprintf("invalid flags for %s", fs.Name()), Cause: err}
	}
	return nil
}
