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
	"arcoris.dev/arcoris-publisher/internal/app"
	"arcoris.dev/arcoris-publisher/internal/report"
)

const (
	defaultManifestPath        = "arcpub.yaml"
	defaultSourceRepositoryDir = "."
	defaultStagingDir          = "."
	defaultTargetRootDir       = ".arcpub/targets"
)

// Options configures CLI defaults that are independent from one specific
// command invocation.
type Options struct {
	// ManifestPath is the default arcpub.yaml path.
	ManifestPath string

	// SourceRepositoryDir is the default source repository root for workflow
	// commands.
	SourceRepositoryDir string

	// StagingDir is the default staging root for workflow commands.
	StagingDir string

	// TargetRootDir is the default target worktree root for workflow commands.
	TargetRootDir string

	// App configures high-level app behavior, including workflow options.
	App app.Options

	// Report configures output rendering defaults.
	Report report.Options
}

// DefaultOptions returns safe CLI defaults.
func DefaultOptions() Options {
	return Options{
		ManifestPath:        defaultManifestPath,
		SourceRepositoryDir: defaultSourceRepositoryDir,
		StagingDir:          defaultStagingDir,
		TargetRootDir:       defaultTargetRootDir,
		Report:              report.DefaultOptions(),
	}
}

func normalizeOptions(opts Options) Options {
	defaults := DefaultOptions()
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
		opts.Report = defaults.Report
	}
	return opts
}
