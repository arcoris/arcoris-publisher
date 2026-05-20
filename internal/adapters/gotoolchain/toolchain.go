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

// Package gotoolchain implements the Go toolchain port by invoking the go
// binary.
//
// The adapter is responsible for turning typed options into stable Go command
// arguments and environment variables. It keeps command construction out of
// workflow code and keeps parsing of go list/go env output close to the command
// that produced it.
package gotoolchain

import (
	goport "arcoris.dev/arcoris-publisher/internal/ports/gotoolchain"
	processport "arcoris.dev/arcoris-publisher/internal/ports/process"
)

const defaultGoBinary = "go"

// Toolchain implements Go toolchain operations through the go command.
//
// Toolchain is safe to reuse. Per-call options are copied into process specs so
// callers can reuse option slices after a call without affecting future runs.
type Toolchain struct {
	runner processport.Runner
	goBin  string
	env    []string
}

// Options configures a Go toolchain adapter.
type Options struct {
	// GoBinary overrides the executable name or path. Empty means "go".
	GoBinary string
	// Env is added to every Go command invocation.
	Env []string
}

// New creates a Go toolchain adapter.
func New(runner processport.Runner, opts Options) *Toolchain {
	goBin := opts.GoBinary
	if goBin == "" {
		goBin = defaultGoBinary
	}
	return &Toolchain{runner: runner, goBin: goBin, env: append([]string(nil), opts.Env...)}
}

var _ goport.Toolchain = (*Toolchain)(nil)
