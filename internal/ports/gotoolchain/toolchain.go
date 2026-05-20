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

// Package gotoolchain defines the infrastructure port for executing Go
// toolchain operations such as go list, go test, and go mod tidy.
//
// The port keeps workflow code independent from command construction details:
// adapters decide how to map options to flags, environment variables, timeouts,
// and output parsing. Workflow code consumes typed requests and results.
//
// Implementations should isolate workspace state when requested, redact
// sensitive values from diagnostics, preserve raw stdout/stderr for debugging,
// and return stable porterr.Error codes for failures reported by the Go command.
package gotoolchain

import "context"

// Toolchain groups Go toolchain capabilities used by publisher verification and
// construction workflows.
type Toolchain interface {
	// Env returns selected Go environment values.
	//
	// Implementations may return a provider-defined subset, but Values should
	// contain every variable the caller explicitly requested through adapter
	// configuration. Missing values can be distinguished with EnvResult.HasValue.
	Env(ctx context.Context, opts EnvOptions) (EnvResult, error)
	// List runs go list for moduleDir and returns raw and optionally parsed data.
	//
	// moduleDir is the module root or a directory inside the module. Adapters
	// should apply WorkspaceMode before invoking the Go command so workspace files
	// outside moduleDir do not accidentally affect publisher decisions.
	List(ctx context.Context, moduleDir string, opts ListOptions) (ListResult, error)
	// Test runs go test for moduleDir according to opts.
	//
	// The result carries raw stdout and stderr even when tests fail so callers can
	// include useful diagnostics in publish reports.
	Test(ctx context.Context, moduleDir string, opts TestOptions) (TestResult, error)
	// ModTidy runs go mod tidy for moduleDir.
	//
	// This operation may modify go.mod and go.sum. Callers are expected to run it
	// only inside prepared worktrees or staging directories.
	ModTidy(ctx context.Context, moduleDir string, opts ModTidyOptions) (ModTidyResult, error)
}
