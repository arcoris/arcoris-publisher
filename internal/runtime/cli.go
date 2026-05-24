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

import "arcoris.dev/arcoris-publisher/internal/cli"

// AppFactory builds app instances from command-specific app options.
func (r Runtime) AppFactory() cli.AppFactory {
	return cli.ApplicationFactoryFromApp(r.AppDependencies())
}

// CLIDependencies builds the CLI dependency graph.
func (r Runtime) CLIDependencies() cli.Dependencies {
	return cli.Dependencies{
		AppFactory: r.AppFactory(),
		BuildInfo:  r.deps.BuildInfo,
	}
}

// CLIOptions returns CLI defaults with runtime app options applied.
//
// Runtime Options.App is authoritative for app defaults. Command-level flags
// handled by internal/cli may still override selected fields, such as dry-run,
// before calling the application factory.
func (r Runtime) CLIOptions() cli.Options {
	opts := r.opts.CLI
	opts.App = r.opts.App
	return opts
}

// CLI creates a command router wired to the production application graph.
func (r Runtime) CLI() cli.CLI {
	return cli.New(r.CLIDependencies(), r.CLIOptions())
}
