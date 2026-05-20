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
package gotoolchain

import "context"

// Toolchain groups Go toolchain capabilities used by publisher verification and
// construction workflows.
type Toolchain interface {
	Env(ctx context.Context, opts EnvOptions) (EnvResult, error)
	List(ctx context.Context, moduleDir string, opts ListOptions) (ListResult, error)
	Test(ctx context.Context, moduleDir string, opts TestOptions) (TestResult, error)
	ModTidy(ctx context.Context, moduleDir string, opts ModTidyOptions) (ModTidyResult, error)
}
