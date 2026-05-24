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

package runtime

import (
	"arcoris.dev/arcoris-publisher/internal/app"
	"arcoris.dev/arcoris-publisher/internal/cli"
	"arcoris.dev/arcoris-publisher/internal/config"
)

// Options configures production runtime wiring.
//
// Options is intentionally limited to process/tool selection and high-level app
// and CLI defaults. Workflow-stage behavior remains inside app.Options and
// workflow.Options, while command-line overrides remain inside cli.Options.
type Options struct {
	// Env is added to every external process launched by production adapters.
	//
	// The slice is copied before it is stored in adapters, so callers may reuse or
	// mutate their original slice after constructing a Runtime.
	Env []string

	// GitBinary overrides the Git executable name or path. Empty means "git".
	GitBinary string

	// GoBinary overrides the Go executable name or path. Empty means "go".
	GoBinary string

	// Loader configures manifest loading. Empty LoaderOptions use config defaults.
	Loader config.LoaderOptions

	// App configures application use cases and workflow stages.
	//
	// Runtime.CLI copies this value into cli.Options.App before constructing the
	// CLI so command-level overrides such as --dry-run still flow through the CLI
	// application factory path.
	App app.Options

	// CLI configures command defaults such as manifest path and report output.
	//
	// Runtime.CLI treats Options.App as authoritative for app defaults and writes
	// it into this value before constructing the CLI.
	CLI cli.Options
}

// DefaultOptions returns conservative production runtime defaults.
func DefaultOptions() Options {
	return Options{
		CLI: cli.DefaultOptions(),
	}
}

func normalizeOptions(opts Options) Options {
	opts.Env = copyStrings(opts.Env)
	opts.CLI = normalizeCLIOptions(opts.CLI)
	return opts
}

func normalizeCLIOptions(opts cli.Options) cli.Options {
	defaults := cli.DefaultOptions()
	if opts.ManifestPath == "" {
		opts.ManifestPath = defaults.ManifestPath
	}
	if opts.SourceRepositoryDir == "" {
		opts.SourceRepositoryDir = defaults.SourceRepositoryDir
	}
	if opts.StagingDir == "" {
		opts.StagingDir = defaults.StagingDir
	}
	if opts.TargetRootDir == "" {
		opts.TargetRootDir = defaults.TargetRootDir
	}
	if opts.Report.Format == "" {
		includeLocalPaths := opts.Report.IncludeLocalPaths
		opts.Report = defaults.Report
		opts.Report.IncludeLocalPaths = includeLocalPaths
	}
	return opts
}

func copyStrings(values []string) []string {
	if len(values) == 0 {
		return nil
	}
	return append([]string(nil), values...)
}
