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

package porttest

import (
	"context"

	"arcoris.dev/arcoris-publisher/internal/ports/gotoolchain"
)

// GoToolchain is a deterministic fake Go toolchain port.
type GoToolchain struct {
	// ModTidyHook mutates test state during ModTidy when set.
	ModTidyHook func(context.Context, string) error

	// ModTidyError forces ModTidy to fail.
	ModTidyError error
}

// Env returns an empty environment result.
func (g GoToolchain) Env(context.Context, gotoolchain.EnvOptions) (gotoolchain.EnvResult, error) {
	return gotoolchain.EnvResult{}, nil
}

// List returns an empty list result.
func (g GoToolchain) List(
	context.Context,
	string,
	gotoolchain.ListOptions,
) (gotoolchain.ListResult, error) {
	return gotoolchain.ListResult{}, nil
}

// Test returns an empty test result.
func (g GoToolchain) Test(
	context.Context,
	string,
	gotoolchain.TestOptions,
) (gotoolchain.TestResult, error) {
	return gotoolchain.TestResult{}, nil
}

// ModTidy runs the configured hook or returns ModTidyError.
func (g GoToolchain) ModTidy(
	ctx context.Context,
	moduleDir string,
	_ gotoolchain.ModTidyOptions,
) (gotoolchain.ModTidyResult, error) {
	if g.ModTidyError != nil {
		return gotoolchain.ModTidyResult{}, g.ModTidyError
	}
	if g.ModTidyHook != nil {
		return gotoolchain.ModTidyResult{}, g.ModTidyHook(ctx, moduleDir)
	}

	return gotoolchain.ModTidyResult{}, nil
}
