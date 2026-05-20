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

package gotoolchain

import (
	"context"

	goport "arcoris.dev/arcoris-publisher/internal/ports/gotoolchain"
)

// Test runs go test with typed test options.
func (t *Toolchain) Test(ctx context.Context, moduleDir string, opts goport.TestOptions) (goport.TestResult, error) {
	result, err := t.runner.Run(ctx, t.command(moduleDir, testArgs(opts), opts.CommonOptions))
	out := goport.TestResult{Stdout: result.Stdout, Stderr: result.Stderr}
	if err != nil {
		return out, wrapGoError(goport.CodeTestFailed, "go test failed", result, err)
	}
	return out, nil
}
