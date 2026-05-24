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
	cmd, rest, err := parseCommand(args)
	if err != nil {
		writeCLIError(stderr, err)
		writeUsage(stderr)
		return ExitUsage
	}

	switch cmd {
	case commandHelp:
		writeUsage(stdout)
		return ExitOK
	default:
		if isHelpRequest(rest) {
			writeUsage(stdout)
			return ExitOK
		}
	}

	switch cmd {
	case commandPlan:
		return c.runPlan(ctx, rest, stdout, stderr)
	case commandVerify:
		return c.runVerify(ctx, rest, stdout, stderr)
	case commandPublish:
		return c.runPublish(ctx, rest, stdout, stderr)
	case commandVersion:
		return c.runVersion(ctx, rest, stdout, stderr)
	default:
		writeCLIError(stderr, &Error{Code: CodeInvalidCommand, Message: "unknown command"})
		writeUsage(stderr)
		return ExitUsage
	}
}

// Run executes one command with default CLI options.
func Run(ctx context.Context, args []string, stdout io.Writer, stderr io.Writer) int {
	return New(Dependencies{}, Options{}).Run(ctx, args, stdout, stderr)
}
