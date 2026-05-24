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
	execadapter "arcoris.dev/arcoris-publisher/internal/adapters/exec"
	"arcoris.dev/arcoris-publisher/internal/adapters/filesystem"
	"arcoris.dev/arcoris-publisher/internal/adapters/gitcli"
	"arcoris.dev/arcoris-publisher/internal/adapters/gotoolchain"
	"arcoris.dev/arcoris-publisher/internal/buildinfo"
	"arcoris.dev/arcoris-publisher/internal/cli"
	"arcoris.dev/arcoris-publisher/internal/config"
	fsport "arcoris.dev/arcoris-publisher/internal/ports/filesystem"
	gitport "arcoris.dev/arcoris-publisher/internal/ports/git"
	goport "arcoris.dev/arcoris-publisher/internal/ports/gotoolchain"
	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

// Dependencies contains concrete collaborators used to build app and CLI
// boundaries.
//
// Tests may provide selected fields and let Runtime fill the rest with default
// production adapters. Production callers usually use NewDependencies through
// New or NewWithDependencies.
type Dependencies struct {
	// Process executes external processes for Git and Go adapters.
	Process processport.Runner

	// FileSystem provides local filesystem operations.
	FileSystem fsport.FileSystem

	// Git provides all Git capabilities required by source, target, verify, and
	// publish workflow stages.
	Git gitport.Client

	// Go provides Go toolchain operations used by verification.
	Go goport.Toolchain

	// Loader reads top-level and module-level manifests.
	Loader *config.Loader

	// BuildInfo supplies build metadata to the CLI version command.
	BuildInfo cli.BuildInfoFunc
}

// NewDependencies constructs production adapters from Options.
func NewDependencies(opts Options) Dependencies {
	return normalizeDependencies(Dependencies{}, opts)
}

func normalizeDependencies(deps Dependencies, opts Options) Dependencies {
	opts = normalizeOptions(opts)
	if deps.Process == nil {
		deps.Process = execadapter.New(execadapter.Options{Env: copyStrings(opts.Env)})
	}
	if deps.FileSystem == nil {
		deps.FileSystem = filesystem.New()
	}
	if deps.Git == nil {
		deps.Git = gitcli.New(deps.Process, gitcli.Options{
			GitBinary: opts.GitBinary,
			Env:       copyStrings(opts.Env),
		})
	}
	if deps.Go == nil {
		deps.Go = gotoolchain.New(deps.Process, gotoolchain.Options{
			GoBinary: opts.GoBinary,
			Env:      copyStrings(opts.Env),
		})
	}
	if deps.Loader == nil {
		loader := config.NewLoader(opts.Loader)
		deps.Loader = &loader
	}
	if deps.BuildInfo == nil {
		deps.BuildInfo = buildinfo.Current
	}
	return deps
}
