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

// Package gitcli implements the Git port by invoking the system git binary.
//
// The adapter keeps Git command construction in one place and exposes the
// typed, workflow-facing port interfaces from internal/ports/git. It does not
// implement repository publishing policy; it only translates Git operations,
// output, and failures.
package gitcli

import (
	gitport "arcoris.dev/arcoris-publisher/internal/ports/git"
	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

const defaultGitBinary = "git"

// Client implements Git operations using the git command-line executable.
//
// Client is safe to reuse. The process runner and environment are supplied at
// construction time so tests can inject a fake runner while production code uses
// the exec adapter.
type Client struct {
	runner processport.Runner
	gitBin string
	env    []string
}

// Options configures a Git CLI client.
type Options struct {
	// GitBinary overrides the executable name or path. Empty means "git".
	GitBinary string
	// Env is added to every Git process invocation.
	Env []string
}

// New creates a Git CLI client.
func New(runner processport.Runner, opts Options) *Client {
	gitBin := opts.GitBinary
	if gitBin == "" {
		gitBin = defaultGitBinary
	}
	return &Client{runner: runner, gitBin: gitBin, env: append([]string(nil), opts.Env...)}
}

var _ gitport.Client = (*Client)(nil)
