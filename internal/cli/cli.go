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
	"errors"
	"io"

	"arcoris.dev/arcoris-publisher/internal/buildinfo"
)

// CLI routes command-line invocations to application use cases and report
// renderers.
type CLI struct {
	deps Dependencies
	opts Options
}

// New creates a CLI router.
func New(deps Dependencies, opts Options) CLI {
	if deps.BuildInfo == nil {
		deps.BuildInfo = buildinfo.Current
	}
	return CLI{deps: deps, opts: normalizeOptions(opts)}
}

// Run executes one command and returns a process-style exit code.
func (c CLI) Run(ctx context.Context, args []string, stdout io.Writer, stderr io.Writer) int {
	if stdout == nil {
		stdout = io.Discard
	}
	if stderr == nil {
		stderr = io.Discard
	}

	root := c.newRootCommand(stdout, stderr)
	root.SetArgs(args)

	if err := root.ExecuteContext(ctx); err != nil {
		return exitCodeFor(normalizeCobraError(err), stderr)
	}
	return ExitOK
}

// Run executes one command with default CLI options.
func Run(ctx context.Context, args []string, stdout io.Writer, stderr io.Writer) int {
	return New(Dependencies{}, Options{}).Run(ctx, args, stdout, stderr)
}

func exitCodeFor(err error, stderr io.Writer) int {
	if errors.Is(err, errVerificationFailed) {
		return ExitVerificationFailed
	}

	writeCLIError(stderr, err)

	var cliErr *Error
	if !errors.As(err, &cliErr) {
		return ExitUsage
	}
	if cliErr.isUsage() {
		return ExitUsage
	}
	return ExitError
}
