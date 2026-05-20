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

// Package exec implements the process port by using the operating-system
// process execution facilities from the Go standard library.
//
// The adapter is deliberately narrow: it starts commands, captures bounded
// output, maps process lifecycle failures to process-port error codes, and
// redacts configured sensitive values from returned diagnostics.
package exec

import (
	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

// defaultOutputLimit caps captured stdout/stderr when the caller does not
// provide a more specific limit.
//
// Process output can become very large during failed builds. Keeping a bounded
// default protects the publisher process while still preserving enough context
// for useful diagnostics.
const defaultOutputLimit int64 = 8 << 20 // 8 MiB

// Runner executes external processes through os/exec.
//
// Runner is safe to reuse across calls. Its base environment is copied during
// construction, and each Run invocation builds a fresh process specification.
type Runner struct {
	baseEnv []string
}

// Options configures Runner.
type Options struct {
	// Env is a base environment overlaid on top of os.Environ for every command.
	//
	// Per-command Spec.Env values are applied after Options.Env, so a workflow can
	// override adapter defaults for one invocation without mutating the runner.
	Env []string
}

// New creates a process runner.
func New(opts Options) *Runner {
	return &Runner{baseEnv: append([]string(nil), opts.Env...)}
}

var _ processport.Runner = (*Runner)(nil)
